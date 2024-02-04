package service

import (
	"context"
	"fmt"
	"time"

	"github.com/parviz-yu/digital-wallet/internal/config"
	"github.com/parviz-yu/digital-wallet/internal/models"
	"github.com/parviz-yu/digital-wallet/internal/storage"
	customerrors "github.com/parviz-yu/digital-wallet/pkg/custom-errors"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
)

type ServiceI interface {
	DoesWalletExists(ctx context.Context, userID string) (int, error)
	PutFunds(ctx context.Context, payment *models.PaymentReq) error
	GetWalletStats(ctx context.Context, userID string) (*models.WalletStatResp, error)
	GetWalletBalance(ctx context.Context, userID string) (*models.WalletResp, error)
}

type service struct {
	cfg  config.Config
	log  logger.LoggerI
	strg storage.StorageI
}

func NewService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) ServiceI {
	return &service{
		cfg:  cfg,
		log:  log,
		strg: strg,
	}
}

func (s *service) DoesWalletExists(ctx context.Context, userID string) (int, error) {
	const fn = "service.DoesWalletExists"

	walletID, err := s.strg.Wallet().GetWallet(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return walletID, nil
}

func (s *service) GetWalletBalance(ctx context.Context, userID string) (*models.WalletResp, error) {
	const fn = "service.GetWalletBalance"

	wllt, err := s.strg.Wallet().CheckBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	res := &models.WalletResp{
		Balance: float64(wllt.Balance) / 100,
	}

	return res, nil
}

func (s *service) GetWalletStats(ctx context.Context, userID string) (*models.WalletStatResp, error) {
	const fn = "service.GetWalletStats"

	walledID, err := s.DoesWalletExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	today := time.Now()
	statRange := &models.WalletStatsRange{
		WalletID:  walledID,
		DateBegin: monthStart(today),
		DateEnd:   monthEnd(today),
	}

	monthlyStats, err := s.strg.Transaction().GetMonthlyStats(ctx, statRange)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	res := &models.WalletStatResp{
		Number: monthlyStats.Number,
		Amount: float64(monthlyStats.Amount) / 100,
	}

	return res, nil
}

func (s *service) PutFunds(ctx context.Context, payment *models.PaymentReq) error {
	const fn = "service.PutFunds"

	wallet, err := s.strg.Wallet().CheckBalance(ctx, payment.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	limit, err := s.strg.Wallet().GetLimit(ctx, wallet.Type)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	smallestUnit := int(payment.Amount * 100)
	if smallestUnit+wallet.Balance > limit.MaxAmount {
		err := customerrors.ErrLimitExceeded{WalletType: limit.Name, MaxAmount: limit.MaxAmount}
		return fmt.Errorf("%s: %w", fn, err)
	}

	tx, err := s.strg.Transaction().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer tx.Rollback()

	pay := &models.Payment{
		Amount:   smallestUnit,
		WalletID: wallet.ID,
	}
	_, err = s.strg.Transaction().PutFunds(ctx, tx, pay)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if err := s.strg.Wallet().UpdateBalance(ctx, tx, pay); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	tx.Commit()

	return nil
}

func monthStart(now time.Time) time.Time {
	year, month, location := now.Year(), now.Month(), now.Location()
	return time.Date(year, month, 1, 0, 0, 0, 0, location)
}

func monthEnd(now time.Time) time.Time {
	var lastDay int
	year, month, location := now.Year(), now.Month(), now.Location()

	if month == 2 && (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
		lastDay = 29
	} else if month == 2 {
		lastDay = 28
	} else if month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12 {
		lastDay = 31
	} else {
		lastDay = 30
	}

	return time.Date(year, month, lastDay, 23, 59, 59, int(time.Second-time.Nanosecond), location)
}

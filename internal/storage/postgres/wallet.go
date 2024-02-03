package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/parviz-yu/digital-wallet/internal/models"
	customerrors "github.com/parviz-yu/digital-wallet/pkg/custom-errors"
)

type walletRepo struct {
	db *sql.DB
}

func newWalletRepo(db *sql.DB) *walletRepo {
	return &walletRepo{
		db: db,
	}
}

// GetWallet return wallet's id if wallet exists
func (r *walletRepo) GetWallet(ctx context.Context, userID string) (int, error) {
	const fn = "storage.postgres.GetWallet"

	var id int
	query := "SELECT id FROM wallets WHERE user_id = $1"

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", fn, customerrors.ErrWalletNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

// CheckBalance return wallet's balance and wallet's type
func (r *walletRepo) CheckBalance(ctx context.Context, userID string) (*models.Wallet, error) {
	const fn = "storage.postgres.CheckBalance"

	wllt := &models.Wallet{}
	query := "SELECT balance, type FROM wallets WHERE user_id = $1"

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&wllt.Balance, &wllt.Type)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", fn, customerrors.ErrWalletNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return wllt, nil
}

// UpdateBalance updates wallet's balance
func (r *walletRepo) UpdateBalance(ctx context.Context, tx *sql.Tx, payment *models.Payment) error {
	const fn = "storage.postgres.UpdateBalance"

	query := "UPDATE wallets SET balance = balance + $2 WHERE id = $1"
	_, err := tx.ExecContext(ctx, query, payment.WalletID, payment.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *walletRepo) GetLimit(ctx context.Context, id int) (*models.Limit, error) {
	const fn = "storage.postgres.GetLimits"

	limit := &models.Limit{}
	query := `SELECT name, max_amount FROM limits WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&limit.Name, &limit.MaxAmount); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return limit, nil
}

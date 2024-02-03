package storage

import (
	"context"
	"database/sql"

	"github.com/parviz-yu/digital-wallet/internal/models"
)

type StorageI interface {
	CloseDB()
	Wallet() WalletRepoI
	Transaction() TxRepoI
}

type WalletRepoI interface {
	GetWallet(ctx context.Context, userID string) (int, error)
	CheckBalance(ctx context.Context, userID string) (*models.Wallet, error)
	UpdateBalance(ctx context.Context, tx *sql.Tx, payment *models.Payment) error
	GetLimit(ctx context.Context, id int) (*models.Limit, error)
}

type TxRepoI interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetMonthlyStats(ctx context.Context, statRange *models.WalletStatsRange) (*models.WalletStat, error)
	PutFunds(ctx context.Context, tx *sql.Tx, payment *models.Payment) (int, error)
}

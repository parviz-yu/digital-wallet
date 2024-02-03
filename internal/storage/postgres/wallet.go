package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *walletRepo) IsExists(ctx context.Context, userID string) error {
	const fn = "storage.postgres.IsExists"

	var id int
	query := "SELECT id FROM wallets WHERE user_id = $1"

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: %w", fn, customerrors.ErrWalletNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (r *walletRepo) CheckBalance(ctx context.Context, userID string) (int, error) {
	const fn = "storage.postgres.CheckBalance"

	var balance int
	query := "SELECT balance FROM wallets WHERE user_id = $1"

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", fn, customerrors.ErrWalletNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return balance, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/parviz-yu/digital-wallet/internal/models"
)

type txRepo struct {
	db *sql.DB
}

func newReplanishmentRepo(db *sql.DB) *txRepo {
	return &txRepo{
		db: db,
	}
}

func (r *txRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	const fn = "storage.postgres.BeginTx"

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return tx, nil
}

// PutFunds adds info of the new refill
func (r *txRepo) PutFunds(ctx context.Context, tx *sql.Tx, payment *models.Payment) (int, error) {
	const fn = "storage.postgres.PutFunds"

	var id int
	query := `INSERT INTO transactions(wallet_id, amount) VALUES ($1, $2) RETURNING id`
	err := tx.QueryRowContext(ctx, query, payment.WalletID, payment.Amount).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

// GetMonthlyStats calculates refills' stats of the speciefic month
func (r *txRepo) GetMonthlyStats(ctx context.Context, statRange *models.WalletStatsRange) (*models.WalletStat, error) {
	const fn = "storage.postgres.MonthlyStats"

	var (
		number sql.NullInt64
		amount sql.NullInt64
	)

	result := &models.WalletStat{}
	query := `SELECT COUNT(amount) AS number, SUM(amount) AS total FROM transactions
	WHERE wallet_id = $1 AND created_at BETWEEN $2 AND $3`

	err := r.db.QueryRowContext(
		ctx,
		query,
		statRange.WalletID,
		statRange.DateBegin,
		statRange.DateEnd,
	).Scan(&number, &amount)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if number.Valid {
		result.Number = int(number.Int64)
	}
	if amount.Valid {
		result.Amount = int(amount.Int64)
	}

	return result, nil
}

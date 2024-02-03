package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/parviz-yu/digital-wallet/internal/config"
	"github.com/parviz-yu/digital-wallet/internal/storage"

	_ "github.com/lib/pq"
)

type store struct {
	db         *sql.DB
	walletRepo *walletRepo
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	const fn = "storage.postgres.NewStorage"

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return newStorage(db), nil
}

func newStorage(db *sql.DB) *store {
	return &store{
		db:         db,
		walletRepo: newWalletRepo(db),
	}
}

func (s *store) CloseDB() {
	s.db.Close()
}

func (s *store) Wallet() storage.WalletRepoI {
	return s.walletRepo
}

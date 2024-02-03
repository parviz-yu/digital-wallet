package storage

import "context"

type StorageI interface {
	CloseDB()
	Wallet() WalletRepoI
}

type WalletRepoI interface {
	IsExists(ctx context.Context, userID string) error
	CheckBalance(ctx context.Context, userID string) (int, error)
}

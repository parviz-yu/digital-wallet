package customerrors

import (
	"errors"
	"fmt"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type ErrLimitExceeded struct {
	WalletType string
	MaxAmount  int
}

func (e ErrLimitExceeded) Error() string {
	return fmt.Sprintf("limit exceeded for %s — %d TJS", e.WalletType, e.MaxAmount/100)
}

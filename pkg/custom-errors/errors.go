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
	return fmt.Sprintf("limit exceeded %d TJS", e.MaxAmount/100)
}

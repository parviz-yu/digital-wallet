package models

import "time"

type Wallet struct {
	ID      int
	Balance int // smalles unit (diram)
	Type    int
}

type Payment struct {
	Amount   int // smalles unit (diram)
	WalletID int
}

type WalletStatsRange struct {
	DateBegin time.Time
	DateEnd   time.Time
	WalletID  int
}

type WalletStatResult struct {
	Number int
	Amount int
}

type Limit struct {
	Name      string
	MaxAmount int
}

type PaymentReq struct {
	UserID string
	Amount float64
}

type WalletResp struct {
	Balance float64 `json:"balance"`
}

type WalletStatResp struct {
	Number int     `json:"number"`
	Amount float64 `json:"amount"`
}

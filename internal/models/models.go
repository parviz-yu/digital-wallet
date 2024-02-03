package models

import "time"

type Wallet struct {
	Balance int // smalles unit (diram)
	Type    int
}

type Payment struct {
	WalletID int
	Amount   int // smalles unit (diram)
}

type WalletStatsRange struct {
	WalletID  int
	DateBegin time.Time
	DateEnd   time.Time
}

type WalletStat struct {
	Number int
	Amount int
}

type Limit struct {
	Name      string
	MaxAmount int
}

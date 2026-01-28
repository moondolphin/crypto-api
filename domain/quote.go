package domain

import "time"

type Quote struct {
	CoinID   int64
	Symbol   string
	Provider string
	Currency string
	Price    string
	QuotedAt time.Time
}

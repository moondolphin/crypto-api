package domain

import "time"

type Quote struct {
	ID        int64
	CoinID    int64
	Symbol    string
	Provider  string
	Currency  string
	Price     string
	QuotedAt  time.Time
	CreatedAt time.Time
}

package domain

import "time"

type QuoteFilter struct {
	Symbol   string
	Provider string
	Currency string

	MinPrice *float64
	MaxPrice *float64

	From *time.Time
	To   *time.Time

	Page     int
	PageSize int
}

package domain

//go:generate echo Generating mocks for quote_filter.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=quote_filter.go -destination=../test/mocks/quote_filter_mock.go -package=mocks

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

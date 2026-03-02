package domain

//go:generate echo Generating mocks for quote_filters.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=quote_filters.go -destination=../test/mocks/quote_filters_mock.go -package=mocks

import "time"

type QuoteFilters struct {
	Symbols    []string
	Providers  []string
	Currencies []string

	MinPrice string
	MaxPrice string

	From *time.Time
	To   *time.Time
}

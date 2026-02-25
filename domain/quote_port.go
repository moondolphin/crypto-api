package domain

//go:generate echo Generating mocks for quote_port.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=quote_port.go -destination=../test/mocks/quote_port_mock.go -package=mocks

import "context"

type QuoteRepository interface {
	Insert(ctx context.Context, q Quote) error

	GetLatest(ctx context.Context, symbol, provider, currency string) (*PriceQuote, error)

	ListFilter(ctx context.Context, f QuoteFilter) ([]Quote, int, error)
}

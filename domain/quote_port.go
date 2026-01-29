package domain

import "context"

type QuoteRepository interface {
	Insert(ctx context.Context, q Quote) error

	GetLatest(ctx context.Context, symbol, provider, currency string) (*PriceQuote, error)
}

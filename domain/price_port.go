package domain

import "context"

type PriceProvider interface {
	Name() string
	GetCurrentPrice(ctx context.Context, coin Coin, currency string) (PriceQuote, error)
}

type PriceProviderRegistry interface {
	Get(name string) (PriceProvider, bool)
}

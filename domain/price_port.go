package domain

//go:generate echo Generating mocks for price_port.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=price_port.go -destination=../test/mocks/price_port_mock.go -package=mocks

import "context"

type PriceProvider interface {
	Name() string
	GetCurrentPrice(ctx context.Context, coin Coin, currency string) (PriceQuote, error)
}

type PriceProviderRegistry interface {
	Get(name string) (PriceProvider, bool)
}

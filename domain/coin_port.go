package domain

//go:generate echo Generating mocks for coin_port.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=coin_port.go -destination=../test/mocks/coin_port_mock.go -package=mocks

import "context"

type CoinRepository interface {
	// GetEnabledBySymbol devuelve la coin si existe y esta habilitada. Si no existe, retorna nil, nil.
	GetEnabledBySymbol(ctx context.Context, symbol string) (*Coin, error)

	// ListEnabled devuelve todas las coins habilitadas como de interes
	ListEnabled(ctx context.Context) ([]Coin, error)

	GetBySymbol(ctx context.Context, symbol string) (*Coin, error)

	Upsert(ctx context.Context, c Coin) (*Coin, error)
}

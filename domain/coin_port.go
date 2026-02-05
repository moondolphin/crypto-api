package domain

import "context"

type CoinRepository interface {
	// GetEnabledBySymbol devuelve la coin si existe y esta habilitada. Si no existe, retorna nil, nil.
	GetEnabledBySymbol(ctx context.Context, symbol string) (*Coin, error)

	// ListEnabled devuelve todas las coins habilitadas como de interes
	ListEnabled(ctx context.Context) ([]Coin, error)

	GetBySymbol(ctx context.Context, symbol string) (*Coin, error)

	Upsert(ctx context.Context, c Coin) (*Coin, error)
}

package domain

import "context"

type CoinRepository interface {
	// GetEnabledBySymbol devuelve la coin habilitada o nil si no existe / no está enabled
	GetEnabledBySymbol(ctx context.Context, symbol string) (*Coin, error)
}

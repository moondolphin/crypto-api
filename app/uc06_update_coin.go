package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrInvalidCoinUpdate = errors.New("invalid_coin_update")
	ErrCoinNotFound      = errors.New("coin_not_found")
)

type UpdateCoinInput struct {
	Symbol        string `json:"-"` // viene por path
	Enabled       *bool  `json:"enabled,omitempty"`
	CoinGeckoID   string `json:"coingecko_id,omitempty"`
	BinanceSymbol string `json:"binance_symbol,omitempty"`
}

type UpdateCoinUseCase struct {
	CoinRepo domain.CoinRepository
	Now      func() time.Time // opcional, por consistencia/test
}

func (in UpdateCoinInput) IsEmpty() bool {
	return in.Enabled == nil &&
		in.CoinGeckoID == "" &&
		in.BinanceSymbol == ""
}

func (uc UpdateCoinUseCase) Execute(ctx context.Context, in UpdateCoinInput) (*domain.Coin, error) {
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	if symbol == "" {
		return nil, ErrInvalidCoinUpdate
	}
	if in.IsEmpty() {
		return nil, ErrInvalidCoinUpdate
	}

	// Traemos la coin existente (si no existe -> nil)
	existing, err := uc.CoinRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrCoinNotFound // si no la tenés creada, la creamos en el paso 3/handler o la agregás ahora
	}

	// Aplicamos cambios solo si vinieron
	if in.Enabled != nil {
		existing.Enabled = *in.Enabled
	}
	if strings.TrimSpace(in.CoinGeckoID) != "" {
		existing.CoinGeckoID = strings.TrimSpace(in.CoinGeckoID)
	}
	if strings.TrimSpace(in.BinanceSymbol) != "" {
		existing.BinanceSymbol = strings.TrimSpace(in.BinanceSymbol)
	}

	// Persistimos (Upsert)
	updated, err := uc.CoinRepo.Upsert(ctx, *existing)
	if err != nil {
		return nil, err
	}

	_ = uc.Now // por ahora no lo usamos; queda para futuro/auditoría/tests
	return updated, nil
}

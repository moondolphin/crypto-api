package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrBadRequest           = errors.New("bad_request")
	ErrCoinNotEnabled       = errors.New("coin_not_enabled")
	ErrProviderNotSupported = errors.New("provider_not_supported")
	ErrExternalService      = errors.New("external_service_error")
)

type GetCurrentPriceInput struct {
	Symbol   string
	Currency string
	Provider string
}

type GetCurrentPriceUseCase struct {
	CoinRepo  domain.CoinRepository
	Providers domain.PriceProviderRegistry
	Now       func() time.Time
}

func (uc GetCurrentPriceUseCase) Execute(ctx context.Context, in GetCurrentPriceInput) (domain.PriceQuote, error) {
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	currency := strings.ToUpper(strings.TrimSpace(in.Currency))
	providerName := strings.ToLower(strings.TrimSpace(in.Provider))

	if symbol == "" || currency == "" || providerName == "" {
		return domain.PriceQuote{}, ErrBadRequest
	}

	coin, err := uc.CoinRepo.GetEnabledBySymbol(ctx, symbol)
	if err != nil {
		return domain.PriceQuote{}, err
	}
	if coin == nil {
		return domain.PriceQuote{}, ErrCoinNotEnabled
	}

	provider, ok := uc.Providers.Get(providerName)
	if !ok {
		return domain.PriceQuote{}, ErrProviderNotSupported
	}

	quote, err := provider.GetCurrentPrice(ctx, *coin, currency)
	if err != nil {
		return domain.PriceQuote{}, ErrExternalService
	}

	now := uc.Now
	if now == nil {
		now = time.Now
	}
	if quote.Timestamp == "" {
		quote.Timestamp = now().UTC().Format(time.RFC3339)
	}

	quote.Symbol = symbol
	quote.Currency = currency
	quote.Provider = provider.Name()

	return quote, nil
}

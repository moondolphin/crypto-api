package app

import (
	"context"
	"errors"
	"strings"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrQuoteNotFound = errors.New("quote_not_found")
)

type GetLastPriceInput struct {
	Symbol   string
	Currency string // opcional
	Provider string // opcional
}

type GetLastPriceUseCase struct {
	CoinRepo  domain.CoinRepository
	QuoteRepo domain.QuoteRepository
}

func (uc GetLastPriceUseCase) Execute(ctx context.Context, in GetLastPriceInput) (domain.PriceQuote, error) {
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	currency := strings.ToUpper(strings.TrimSpace(in.Currency))
	provider := strings.ToLower(strings.TrimSpace(in.Provider))

	if symbol == "" {
		return domain.PriceQuote{}, ErrBadRequest
	}

	coin, err := uc.CoinRepo.GetEnabledBySymbol(ctx, symbol)
	if err != nil {
		return domain.PriceQuote{}, err
	}
	if coin == nil {
		return domain.PriceQuote{}, ErrCoinNotEnabled
	}

	q, err := uc.QuoteRepo.GetLatest(ctx, symbol, provider, currency)
	if err != nil {
		return domain.PriceQuote{}, err
	}
	if q == nil {
		return domain.PriceQuote{}, ErrQuoteNotFound
	}

	return *q, nil
}

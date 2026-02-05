package app

import (
	"context"
	"errors"
	"strings"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrInvalidCoinInput = errors.New("invalid_coin_input")
)

type CreateCoinInput struct {
	Symbol        string `json:"symbol"`
	Enabled       *bool  `json:"enabled,omitempty"`
	CoinGeckoID   string `json:"coingecko_id,omitempty"`
	BinanceSymbol string `json:"binance_symbol,omitempty"`
}

type CreateCoinOutput struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Enabled       bool   `json:"enabled"`
	CoinGeckoID   string `json:"coingecko_id"`
	BinanceSymbol string `json:"binance_symbol"`
}

type CreateCoinUseCase struct {
	CoinRepo domain.CoinRepository
}

func (uc CreateCoinUseCase) Execute(ctx context.Context, in CreateCoinInput) (CreateCoinOutput, error) {
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	if symbol == "" {
		return CreateCoinOutput{}, ErrInvalidCoinInput
	}

	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}

	coin := domain.Coin{
		Symbol:        symbol,
		Enabled:       enabled,
		CoinGeckoID:   strings.TrimSpace(in.CoinGeckoID),
		BinanceSymbol: strings.TrimSpace(in.BinanceSymbol),
	}

	out, err := uc.CoinRepo.Upsert(ctx, coin)
	if err != nil {
		return CreateCoinOutput{}, err
	}

	return CreateCoinOutput{
		ID:            out.ID,
		Symbol:        out.Symbol,
		Enabled:       out.Enabled,
		CoinGeckoID:   out.CoinGeckoID,
		BinanceSymbol: out.BinanceSymbol,
	}, nil
}

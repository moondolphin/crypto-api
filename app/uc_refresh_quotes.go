package app

import (
	"context"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

type RefreshQuotesOutput struct {
	CoinsProcessed int `json:"coins_processed"`
	QuotesSaved    int `json:"quotes_saved"`
	Failed         int `json:"failed"`
}

type RefreshQuotesUseCase struct {
	CoinRepo   domain.CoinRepository
	QuoteRepo  domain.QuoteRepository
	Providers  domain.PriceProviderRegistry
	Now        func() time.Time
	ProviderFX map[string]string // provider -> currency (ej: binance->USDT, coingecko->USD)
}

func (uc RefreshQuotesUseCase) Execute(ctx context.Context) (RefreshQuotesOutput, error) {
	now := uc.Now
	if now == nil {
		now = time.Now
	}

	coins, err := uc.CoinRepo.ListEnabled(ctx)
	if err != nil {
		return RefreshQuotesOutput{}, err
	}

	out := RefreshQuotesOutput{CoinsProcessed: len(coins)}

	for _, coin := range coins {
		for providerName, currency := range uc.ProviderFX {
			if providerName == "binance" && coin.BinanceSymbol == "" {
				continue
			}
			if providerName == "coingecko" && coin.CoinGeckoID == "" {
				continue
			}

			p, ok := uc.Providers.Get(providerName)
			if !ok {
				out.Failed++
				continue
			}

			q, err := p.GetCurrentPrice(ctx, coin, currency)
			if err != nil {
				out.Failed++
				continue
			}

			quotedAt := now().UTC()
			if q.Timestamp != "" {
				if t, parseErr := time.Parse(time.RFC3339, q.Timestamp); parseErr == nil {
					quotedAt = t.UTC()
				}
			}

			err = uc.QuoteRepo.Insert(ctx, domain.Quote{
				CoinID:   coin.ID,
				Symbol:   coin.Symbol,
				Provider: p.Name(),
				Currency: currency,
				Price:    q.Price,
				QuotedAt: quotedAt,
			})
			if err != nil {
				out.Failed++
				continue
			}

			out.QuotesSaved++
		}
	}

	return out, nil
}

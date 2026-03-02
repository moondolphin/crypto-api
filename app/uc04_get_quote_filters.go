package app

import (
	"context"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

type QuoteFiltersDTO struct {
	Symbols    []string `json:"symbols"`
	Providers  []string `json:"providers"`
	Currencies []string `json:"currencies"`

	MinPrice string `json:"min_price"`
	MaxPrice string `json:"max_price"`

	From string `json:"from"`
	To   string `json:"to"`
}

type GetQuoteFiltersOutput struct {
	Filters QuoteFiltersDTO `json:"filters"`
}

type GetQuoteFiltersUseCase struct {
	Repo domain.QuoteRepository
}

func (uc GetQuoteFiltersUseCase) Execute(ctx context.Context, in SearchQuotesInput) (GetQuoteFiltersOutput, error) {
	// Misma normalización que SearchQuotesUseCase (pero sin paginado)
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	currency := strings.ToUpper(strings.TrimSpace(in.Currency))

	f := domain.QuoteFilter{
		Symbol:   symbol,
		Provider: provider,
		Currency: currency,
		MinPrice: in.MinPrice,
		MaxPrice: in.MaxPrice,
		From:     in.From,
		To:       in.To,
		// Page/PageSize no aplican acá
	}

	facets, err := uc.Repo.ListAvailableFilters(ctx, f)
	if err != nil {
		return GetQuoteFiltersOutput{}, err
	}

	dto := QuoteFiltersDTO{
		Symbols:    facets.Symbols,
		Providers:  facets.Providers,
		Currencies: facets.Currencies,
		MinPrice:   facets.MinPrice,
		MaxPrice:   facets.MaxPrice,
	}

	if facets.From != nil {
		dto.From = facets.From.UTC().Format(time.RFC3339)
	}
	if facets.To != nil {
		dto.To = facets.To.UTC().Format(time.RFC3339)
	}

	return GetQuoteFiltersOutput{Filters: dto}, nil
}

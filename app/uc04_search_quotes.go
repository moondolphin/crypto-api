package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrInvalidFilters = errors.New("invalid_filters")
)

type SearchQuotesInput struct {
	Symbol   string
	Provider string
	Currency string

	MinPrice *float64
	MaxPrice *float64

	From *time.Time
	To   *time.Time

	Page     int // 1..10
	PageSize int // default 50 (por ej), máximo 100
}

type QuoteItem struct {
	Symbol   string    `json:"symbol"`
	Provider string    `json:"provider"`
	Currency string    `json:"currency"`
	Price    string    `json:"price"`
	QuotedAt time.Time `json:"quoted_at"`
}

type QuotesSummary struct {
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
}

type SearchQuotesOutput struct {
	Items   []QuoteItem   `json:"items"`
	Summary QuotesSummary `json:"summary"`
}

type SearchQuotesUseCase struct {
	Repo domain.QuoteRepository
}

func (uc SearchQuotesUseCase) Execute(ctx context.Context, in SearchQuotesInput) (SearchQuotesOutput, error) {
	// normalizacion basica
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	provider := strings.ToLower(strings.TrimSpace(in.Provider))
	currency := strings.ToUpper(strings.TrimSpace(in.Currency))

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page > 10 {
		return SearchQuotesOutput{}, ErrInvalidFilters
	}

	// armar filtro de dominio
	f := domain.QuoteFilter{
		Symbol:   symbol,
		Provider: provider,
		Currency: currency,
		MinPrice: in.MinPrice,
		MaxPrice: in.MaxPrice,
		From:     in.From,
		To:       in.To,
		Page:     page,
		PageSize: pageSize,
	}

	quotes, total, err := uc.Repo.ListFilter(ctx, f)
	if err != nil {
		return SearchQuotesOutput{}, err
	}

	items := make([]QuoteItem, 0, len(quotes))
	for _, q := range quotes {
		items = append(items, QuoteItem{
			Symbol:   q.Symbol,
			Provider: q.Provider,
			Currency: q.Currency,
			Price:    q.Price,
			QuotedAt: q.QuotedAt,
		})
	}

	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	// maximo 10 paginas
	if totalPages > 10 {
		totalPages = 10
	}

	return SearchQuotesOutput{
		Items: items,
		Summary: QuotesSummary{
			TotalItems: total,
			TotalPages: totalPages,
			Page:       page,
			PageSize:   pageSize,
		},
	}, nil
}

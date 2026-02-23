package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/domain"
	"github.com/moondolphin/crypto-api/test/mocks"
)

func TestUC04SearchQuotes_InvalidPage_WhenPageExceedsMax(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		Page:     11,
		PageSize: 50,
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidFilters)
}

func TestUC04SearchQuotes_DefaultPageSize_WhenNotProvided(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	quotes := []domain.Quote{
		{Symbol: "BTC", Provider: "binance", Currency: "USD", Price: "45000", QuotedAt: time.Now()},
	}

	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
			// Verify default page size was applied (50)
			require.Equal(t, 50, f.PageSize)
			return quotes, 1, nil
		})

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		Page:     1,
		PageSize: 0, // Should default to 50
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, 50, result.Summary.PageSize)
	require.Equal(t, 1, result.Summary.TotalItems)
}

func TestUC04SearchQuotes_MaxPageSize_WhenExceeds100(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	quotes := []domain.Quote{
		{Symbol: "BTC", Provider: "binance", Currency: "USD", Price: "45000", QuotedAt: time.Now()},
	}

	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
			// Verify max page size was capped at 100
			require.Equal(t, 100, f.PageSize)
			return quotes, 1, nil
		})

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		Page:     1,
		PageSize: 200, // Should be capped to 100
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, 100, result.Summary.PageSize)
}

func TestUC04SearchQuotes_DefaultPage_WhenZeroOrNegative(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	quotes := []domain.Quote{
		{Symbol: "BTC", Provider: "binance", Currency: "USD", Price: "45000", QuotedAt: time.Now()},
	}

	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
			// Verify default page was applied (1)
			require.Equal(t, 1, f.Page)
			return quotes, 1, nil
		})

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		Page:     -1, // Should default to 1
		PageSize: 50,
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.Summary.Page)
}

func TestUC04SearchQuotes_RepoError_WhenQueryFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		Return(nil, 0, errors.New("db_error"))

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		Page:     1,
		PageSize: 50,
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "db_error")
}

func TestUC04SearchQuotes_Success_ReturnsQuotesWithCorrectSummary(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	quotedTime := time.Date(2026, 1, 22, 15, 30, 0, 0, time.UTC)

	quotes := []domain.Quote{
		{
			Symbol:   "BTC",
			Provider: "binance",
			Currency: "USD",
			Price:    "45000",
			QuotedAt: quotedTime,
		},
		{
			Symbol:   "BTC",
			Provider: "coingecko",
			Currency: "USD",
			Price:    "45100",
			QuotedAt: quotedTime,
		},
	}

	// Assume 150 total items, so with page size 50, we have 3 pages
	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		Return(quotes, 150, nil)

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		Provider: "binance",
		Currency: "USD",
		Page:     1,
		PageSize: 50,
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, 2, len(result.Items))
	require.Equal(t, "BTC", result.Items[0].Symbol)
	require.Equal(t, "45000", result.Items[0].Price)
	require.Equal(t, quotedTime, result.Items[0].QuotedAt)
	require.Equal(t, 150, result.Summary.TotalItems)
	require.Equal(t, 3, result.Summary.TotalPages)
	require.Equal(t, 1, result.Summary.Page)
	require.Equal(t, 50, result.Summary.PageSize)
}

func TestUC04SearchQuotes_Success_WithFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	minPrice := 40000.0
	maxPrice := 50000.0
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
			// Verify filters are applied
			require.Equal(t, "BTC", f.Symbol)
			require.Equal(t, &minPrice, f.MinPrice)
			require.Equal(t, &maxPrice, f.MaxPrice)
			require.Equal(t, &from, f.From)
			require.Equal(t, &to, f.To)
			return []domain.Quote{}, 0, nil
		})

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   "BTC",
		MinPrice: &minPrice,
		MaxPrice: &maxPrice,
		From:     &from,
		To:       &to,
		Page:     1,
		PageSize: 50,
	})

	// Assert
	require.NoError(t, err)
}

func TestUC04SearchQuotes_Success_NormalizesInputs(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	repo.EXPECT().
		ListFilter(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
			// Verify normalization
			require.Equal(t, "BTC", f.Symbol)       // uppercase
			require.Equal(t, "binance", f.Provider) // lowercase
			require.Equal(t, "USD", f.Currency)     // uppercase
			return []domain.Quote{}, 0, nil
		})

	uc := app.SearchQuotesUseCase{
		Repo: repo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol:   " btc ",
		Provider: " BINANCE ",
		Currency: " usd ",
		Page:     1,
		PageSize: 50,
	})

	// Assert
	require.NoError(t, err)
}

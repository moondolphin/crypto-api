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

func TestUC04GetQuoteFilters_Success_WithoutFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	filters := domain.QuoteFilters{
		Symbols:    []string{"BTC", "ETH"},
		Providers:  []string{"binance", "coingecko"},
		Currencies: []string{"USD", "EUR"},
		MinPrice:   "100.5",
		MaxPrice:   "50000.99",
	}

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Return(filters, nil)

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{})

	// Assert
	require.NoError(t, err)
	require.Equal(t, []string{"BTC", "ETH"}, result.Filters.Symbols)
	require.Equal(t, []string{"binance", "coingecko"}, result.Filters.Providers)
	require.Equal(t, []string{"USD", "EUR"}, result.Filters.Currencies)
	require.Equal(t, "100.5", result.Filters.MinPrice)
	require.Equal(t, "50000.99", result.Filters.MaxPrice)
}

func TestUC04GetQuoteFilters_Success_WithSymbolFilter(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	filters := domain.QuoteFilters{
		Symbols:    []string{"BTC"},
		Providers:  []string{"binance"},
		Currencies: []string{"USD"},
		MinPrice:   "40000",
		MaxPrice:   "50000",
	}

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, f domain.QuoteFilter) {
			// Verify symbol was normalized to uppercase
			require.Equal(t, "BTC", f.Symbol)
			require.Equal(t, "", f.Provider)
			require.Equal(t, "", f.Currency)
		}).
		Return(filters, nil)

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol: "  btc  ", // Should be normalized to uppercase
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, []string{"BTC"}, result.Filters.Symbols)
}

func TestUC04GetQuoteFilters_Success_WithProviderAndCurrencyFilter(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	filters := domain.QuoteFilters{
		Symbols:    []string{"BTC", "ETH"},
		Providers:  []string{"binance"},
		Currencies: []string{"USD"},
		MinPrice:   "100",
		MaxPrice:   "60000",
	}

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, f domain.QuoteFilter) {
			// Verify provider was normalized to lowercase
			// Verify currency was normalized to uppercase
			require.Equal(t, "binance", f.Provider)
			require.Equal(t, "USD", f.Currency)
		}).
		Return(filters, nil)

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Provider: "  BINANCE  ", // Should be normalized to lowercase
		Currency: "  usd  ",     // Should be normalized to uppercase
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, []string{"binance"}, result.Filters.Providers)
	require.Equal(t, []string{"USD"}, result.Filters.Currencies)
}

func TestUC04GetQuoteFilters_Success_WithPriceFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	minPrice := 1000.0
	maxPrice := 50000.0

	filters := domain.QuoteFilters{
		Symbols:    []string{"BTC"},
		Providers:  []string{"binance"},
		Currencies: []string{"USD"},
		MinPrice:   "1000.00",
		MaxPrice:   "50000.00",
	}

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, f domain.QuoteFilter) {
			require.NotNil(t, f.MinPrice)
			require.NotNil(t, f.MaxPrice)
			require.Equal(t, 1000.0, *f.MinPrice)
			require.Equal(t, 50000.0, *f.MaxPrice)
		}).
		Return(filters, nil)

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		MinPrice: &minPrice,
		MaxPrice: &maxPrice,
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, "1000.00", result.Filters.MinPrice)
	require.Equal(t, "50000.00", result.Filters.MaxPrice)
}

func TestUC04GetQuoteFilters_Success_WithDateFilters(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	fromTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	filters := domain.QuoteFilters{
		Symbols:    []string{"BTC"},
		Providers:  []string{"binance"},
		Currencies: []string{"USD"},
		MinPrice:   "40000",
		MaxPrice:   "50000",
		From:       &fromTime,
		To:         &toTime,
	}

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, f domain.QuoteFilter) {
			require.NotNil(t, f.From)
			require.NotNil(t, f.To)
			require.Equal(t, fromTime, *f.From)
			require.Equal(t, toTime, *f.To)
		}).
		Return(filters, nil)

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		From: &fromTime,
		To:   &toTime,
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, "2025-01-01T00:00:00Z", result.Filters.From)
	require.Equal(t, "2025-12-31T23:59:59Z", result.Filters.To)
}

func TestUC04GetQuoteFilters_Error_WhenRepositoryFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Return(domain.QuoteFilters{}, errors.New("database error"))

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.SearchQuotesInput{})

	// Assert
	require.Error(t, err)
	require.Equal(t, "database error", err.Error())
}

func TestUC04GetQuoteFilters_Success_EmptyFiltersResult(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockQuoteRepository(ctrl)

	filters := domain.QuoteFilters{
		Symbols:    []string{},
		Providers:  []string{},
		Currencies: []string{},
		MinPrice:   "",
		MaxPrice:   "",
	}

	repo.EXPECT().
		ListAvailableFilters(gomock.Any(), gomock.Any()).
		Return(filters, nil)

	uc := app.GetQuoteFiltersUseCase{
		Repo: repo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.SearchQuotesInput{
		Symbol: "NONEXISTENT",
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, len(result.Filters.Symbols))
	require.Equal(t, 0, len(result.Filters.Providers))
	require.Equal(t, 0, len(result.Filters.Currencies))
	require.Equal(t, "", result.Filters.MinPrice)
	require.Equal(t, "", result.Filters.MaxPrice)
}

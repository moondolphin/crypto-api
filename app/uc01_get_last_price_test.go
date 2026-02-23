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

func TestUC01LastPrice_BadRequest_WhenMissingSymbol(t *testing.T) {
	// Arrange
	uc := app.GetLastPriceUseCase{}

	// Act
	_, err := uc.Execute(context.Background(), app.GetLastPriceInput{
		Symbol:   "",
		Currency: "USD",
		Provider: "binance",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrBadRequest)
}

func TestUC01LastPrice_CoinNotEnabled_WhenRepoReturnsNilCoin(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	uc := app.GetLastPriceUseCase{
		CoinRepo:  coinRepo,
		QuoteRepo: quoteRepo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.GetLastPriceInput{
		Symbol:   " btc ",
		Currency: "usd",
		Provider: "binance",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrCoinNotEnabled)
}

func TestUC01LastPrice_QuoteNotFound_WhenRepoReturnsNil(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)

	coin := domain.Coin{ID: 1, Symbol: "BTC", Enabled: true}

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(&coin, nil)

	quoteRepo.EXPECT().
		GetLatest(gomock.Any(), "BTC", "binance", "USD").
		Return(nil, nil)

	uc := app.GetLastPriceUseCase{
		CoinRepo:  coinRepo,
		QuoteRepo: quoteRepo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.GetLastPriceInput{
		Symbol:   "BTC",
		Currency: "USD",
		Provider: "binance",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrQuoteNotFound)
}

func TestUC01LastPrice_RepoError_WhenQuoteRepoFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)

	coin := domain.Coin{ID: 1, Symbol: "BTC", Enabled: true}

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(&coin, nil)

	quoteRepo.EXPECT().
		GetLatest(gomock.Any(), "BTC", "binance", "USD").
		Return(nil, errors.New("db_error"))

	uc := app.GetLastPriceUseCase{
		CoinRepo:  coinRepo,
		QuoteRepo: quoteRepo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.GetLastPriceInput{
		Symbol:   "BTC",
		Currency: "USD",
		Provider: "binance",
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "db_error")
}

func TestUC01LastPrice_Success_ReturnsQuoteWithNormalization(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)

	coin := domain.Coin{ID: 1, Symbol: "BTC", Enabled: true}
	quotedTime := time.Date(2026, 1, 22, 15, 30, 0, 0, time.UTC)
	expectedQuote := &domain.PriceQuote{
		Symbol:    "BTC",
		Currency:  "USD",
		Provider:  "binance",
		Price:     "45000.50",
		Timestamp: quotedTime.Format(time.RFC3339),
	}

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(&coin, nil)

	quoteRepo.EXPECT().
		GetLatest(gomock.Any(), "BTC", "binance", "USD").
		Return(expectedQuote, nil)

	uc := app.GetLastPriceUseCase{
		CoinRepo:  coinRepo,
		QuoteRepo: quoteRepo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.GetLastPriceInput{
		Symbol:   " btc ",
		Currency: " usd ",
		Provider: " binance ",
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, "BTC", result.Symbol)
	require.Equal(t, "USD", result.Currency)
	require.Equal(t, "binance", result.Provider)
	require.Equal(t, "45000.50", result.Price)
	require.Equal(t, quotedTime.Format(time.RFC3339), result.Timestamp)
}

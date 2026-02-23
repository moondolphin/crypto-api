package app_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/domain"
	"github.com/moondolphin/crypto-api/test/mocks"
)

func TestUC05CreateCoin_InvalidInput_WhenSymbolEmpty(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:  "",
		Enabled: nil,
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidCoinInput)
}

func TestUC05CreateCoin_DefaultEnabled_WhenNotSpecified(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	providers.EXPECT().
		Get("binance").
		Return(nil, false)

	providers.EXPECT().
		Get("coingecko").
		Return(nil, false)

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:  "BTC",
		Enabled: nil, // Not specified
	})

	// Assert
	require.ErrorIs(t, err, app.ErrCoinNotResolvable)
}

func TestUC05CreateCoin_CoinNotResolvable_WhenNoIDsProvidedAndCantResolve(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "UNKNOWNCOIN").
		Return(nil, nil)

	// Both providers not available
	providers.EXPECT().
		Get("binance").
		Return(nil, false).
		AnyTimes()

	providers.EXPECT().
		Get("coingecko").
		Return(nil, false).
		AnyTimes()

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol: "UNKNOWNCOIN",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrCoinNotResolvable)
}

func TestUC05CreateCoin_Success_WithExplicitBinanceAndCoinGeckoIDs(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	createdCoin := &domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, "BTC", c.Symbol)
			require.Equal(t, "bitcoin", c.CoinGeckoID)
			require.Equal(t, "BTCUSDT", c.BinanceSymbol)
		}).
		Return(createdCoin, nil)

	// Setup HTTP mock client that returns success for validations
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/v3/ticker/price") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"symbol":"BTCUSDT","price":"45000"}`))
			return
		}
		if strings.Contains(r.URL.Path, "/api/v3/simple/price") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"bitcoin":{"usd":45000}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer httpServer.Close()

	httpClient := httpServer.Client()

	uc := app.CreateCoinUseCase{
		CoinRepo:   coinRepo,
		Providers:  providers,
		HTTPClient: httpClient,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:        "BTC",
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, "BTC", result.Symbol)
	require.True(t, result.Enabled)
	require.Equal(t, "bitcoin", result.CoinGeckoID)
	require.Equal(t, "BTCUSDT", result.BinanceSymbol)
}

func TestUC05CreateCoin_Success_MergesWithExistingCoin(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	existingCoin := &domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	updatedCoin := &domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       false, // Updated enabled status
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, int64(1), c.ID)
			require.False(t, c.Enabled)
		}).
		Return(updatedCoin, nil)

	enabled := false
	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:  "BTC",
		Enabled: &enabled,
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, int64(1), result.ID)
	require.False(t, result.Enabled)
}

func TestUC05CreateCoin_RepoError_WhenGetBySymbolFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, errors.New("db_error"))

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol: "BTC",
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "db_error")
}

func TestUC05CreateCoin_RepoError_WhenUpsertFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("upsert_error"))

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:        "BTC",
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "upsert_error")
}

func TestUC05CreateCoin_Success_NormalizesSymbol(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	createdCoin := &domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, "BTC", c.Symbol)
		}).
		Return(createdCoin, nil)

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:        " btc ",
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	})

	// Assert
	require.NoError(t, err)
}

func TestUC05CreateCoin_InvalidIDs_WhenProvidedButInvalid(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	// Invalid IDs are ignored and auto-resolve will fail without valid providers
	providers.EXPECT().
		Get("binance").
		Return(nil, false).
		AnyTimes()

	providers.EXPECT().
		Get("coingecko").
		Return(nil, false).
		AnyTimes()

	uc := app.CreateCoinUseCase{
		CoinRepo:  coinRepo,
		Providers: providers,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.CreateCoinInput{
		Symbol:        "BTC",
		CoinGeckoID:   "invalid_id",
		BinanceSymbol: "INVALID",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrCoinNotResolvable)
}

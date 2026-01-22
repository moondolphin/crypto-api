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

func TestUC01_BadRequest_WhenMissingParams(t *testing.T) {
	// Arrange
	uc := app.GetCurrentPriceUseCase{}

	// Act
	_, err := uc.Execute(context.Background(), app.GetCurrentPriceInput{
		Symbol:   "",
		Currency: "USD",
		Provider: "binance",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrBadRequest)
}

func TestUC01_CoinNotEnabled_WhenRepoReturnsNilCoin(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	reg := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	uc := app.GetCurrentPriceUseCase{
		CoinRepo:  coinRepo,
		Providers: reg,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.GetCurrentPriceInput{
		Symbol:   " btc ",
		Currency: "usd",
		Provider: "binance",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrCoinNotEnabled)
}

func TestUC01_ProviderNotSupported_WhenRegistryMiss(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	reg := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(&domain.Coin{ID: 1, Symbol: "BTC", Enabled: true}, nil)

	reg.EXPECT().
		Get("nope").
		Return(nil, false)

	uc := app.GetCurrentPriceUseCase{
		CoinRepo:  coinRepo,
		Providers: reg,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.GetCurrentPriceInput{
		Symbol:   "BTC",
		Currency: "USD",
		Provider: "nope",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrProviderNotSupported)
}

func TestUC01_ExternalServiceError_WhenProviderFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	reg := mocks.NewMockPriceProviderRegistry(ctrl)
	provider := mocks.NewMockPriceProvider(ctrl)

	coin := domain.Coin{ID: 1, Symbol: "BTC", Enabled: true}

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(&coin, nil)

	reg.EXPECT().
		Get("binance").
		Return(provider, true)

	provider.EXPECT().
		GetCurrentPrice(gomock.Any(), coin, "USD").
		Return(domain.PriceQuote{}, errors.New("boom"))

	uc := app.GetCurrentPriceUseCase{
		CoinRepo:  coinRepo,
		Providers: reg,
		Now:       func() time.Time { return time.Date(2026, 1, 22, 2, 0, 0, 0, time.UTC) },
	}

	// Act
	_, err := uc.Execute(context.Background(), app.GetCurrentPriceInput{
		Symbol:   "BTC",
		Currency: "USD",
		Provider: "binance",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrExternalService)
}

func TestUC01_Success_NormalizesAndSetsTimestampAndProviderName(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	reg := mocks.NewMockPriceProviderRegistry(ctrl)
	provider := mocks.NewMockPriceProvider(ctrl)

	coin := domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		CoinGeckoID:   "bitcoin",
		BinanceSymbol: "BTCUSDT",
	}

	fixedNow := time.Date(2026, 1, 22, 2, 1, 9, 0, time.UTC)

	coinRepo.EXPECT().
		GetEnabledBySymbol(gomock.Any(), "BTC").
		Return(&coin, nil)

	reg.EXPECT().
		Get("coingecko").
		Return(provider, true)

	provider.EXPECT().
		Name().
		Return("coingecko")

	// provider devuelve quote sin timestamp → el UC lo completa
	provider.EXPECT().
		GetCurrentPrice(gomock.Any(), coin, "USD").
		Return(domain.PriceQuote{
			Price: "90019",
		}, nil)

	uc := app.GetCurrentPriceUseCase{
		CoinRepo:  coinRepo,
		Providers: reg,
		Now:       func() time.Time { return fixedNow },
	}

	// Act
	out, err := uc.Execute(context.Background(), app.GetCurrentPriceInput{
		Symbol:   " btc ",
		Currency: " usd ",
		Provider: " COINGECKO ",
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, "BTC", out.Symbol)
	require.Equal(t, "USD", out.Currency)
	require.Equal(t, "90019", out.Price)
	require.Equal(t, "coingecko", out.Provider)
	require.Equal(t, fixedNow.UTC().Format(time.RFC3339), out.Timestamp)
}

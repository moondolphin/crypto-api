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

func TestUC12ManualRefresh_CooldownActive_WhenRefreshNotAllowed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshUC := &app.RefreshQuotesUseCase{}
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	lastRefresh := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)
	currentTime := time.Date(2026, 1, 22, 10, 10, 0, 0, time.UTC) // 10 minutes later
	cooldown := 20 * time.Minute

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(lastRefresh, true, nil)

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   *refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    cooldown,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.ErrorIs(t, err, app.ErrCooldownActive)
	require.NotEqual(t, 0, result.RetryAfterSeconds)
	// 20 minutes cooldown - 10 minutes elapsed = 10 minutes remaining (600 seconds)
	require.Equal(t, 600, result.RetryAfterSeconds)
}

func TestUC12ManualRefresh_CooldownActive_RetryAfterSeconds(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshUC := &app.RefreshQuotesUseCase{}
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	lastRefresh := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)
	currentTime := time.Date(2026, 1, 22, 10, 19, 59, 0, time.UTC) // 19 minutes 59 seconds later
	cooldown := 20 * time.Minute

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(lastRefresh, true, nil)

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   *refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    cooldown,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.ErrorIs(t, err, app.ErrCooldownActive)
	// 1 second remaining, should be rounded up to 1
	require.Equal(t, 1, result.RetryAfterSeconds)
}

func TestUC12ManualRefresh_Success_WhenCooldownExpired(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	lastRefresh := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)
	currentTime := time.Date(2026, 1, 22, 10, 25, 0, 0, time.UTC) // 25 minutes later
	cooldown := 20 * time.Minute

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return([]domain.Coin{}, nil)

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(lastRefresh, true, nil)

	controlRepo.EXPECT().
		SetLastManualRefresh(gomock.Any(), currentTime).
		Return(nil)

	refreshUC := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return currentTime },
		ProviderFX: map[string]string{},
	}

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    cooldown,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, result.RetryAfterSeconds)
	require.Equal(t, 0, result.CoinsProcessed)
}

func TestUC12ManualRefresh_Success_WhenNoLastRefresh(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	currentTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return([]domain.Coin{}, nil)

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(time.Time{}, false, nil) // No last refresh

	controlRepo.EXPECT().
		SetLastManualRefresh(gomock.Any(), currentTime).
		Return(nil)

	refreshUC := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return currentTime },
		ProviderFX: map[string]string{},
	}

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    20 * time.Minute,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, result.RetryAfterSeconds)
}

func TestUC12ManualRefresh_RepoError_WhenGetLastRefreshFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshUC := &app.RefreshQuotesUseCase{}
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	currentTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(time.Time{}, false, errors.New("db_error"))

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   *refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    20 * time.Minute,
	}

	// Act
	_, err := uc.Execute(context.Background())

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "db_error")
}

func TestUC12ManualRefresh_DefaultCooldown_WhenNotSet(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshUC := &app.RefreshQuotesUseCase{}
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	lastRefresh := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)
	currentTime := time.Date(2026, 1, 22, 10, 10, 0, 0, time.UTC) // 10 minutes later

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(lastRefresh, true, nil)

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   *refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    0, // Should default to 20 minutes
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.ErrorIs(t, err, app.ErrCooldownActive)
	require.Equal(t, 600, result.RetryAfterSeconds) // 10 minutes remaining with default 20 min cooldown
}

func TestUC12ManualRefresh_Success_ReturnsRefreshOutput(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)

	currentTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	coin := domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		BinanceSymbol: "BTCUSDT",
		CoinGeckoID:   "",
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return([]domain.Coin{coin}, nil)

	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true)

	binanceProvider.EXPECT().
		Name().
		Return("binance")

	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coin, "USDT").
		Return(domain.PriceQuote{Price: "45000", Timestamp: currentTime.Format(time.RFC3339)}, nil)

	quoteRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(nil)

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(time.Time{}, false, nil)

	controlRepo.EXPECT().
		SetLastManualRefresh(gomock.Any(), currentTime).
		Return(nil)

	refreshUC := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return currentTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    20 * time.Minute,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.CoinsProcessed)
	require.Equal(t, 1, result.QuotesSaved)
	require.Equal(t, 0, result.Failed)
	require.Equal(t, 0, result.RetryAfterSeconds)
}

func TestUC12ManualRefresh_Success_SetLastRefreshIgnoresError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	currentTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return([]domain.Coin{}, nil)

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(time.Time{}, false, nil)

	// SetLastManualRefresh fails but should be ignored
	controlRepo.EXPECT().
		SetLastManualRefresh(gomock.Any(), currentTime).
		Return(errors.New("store_error"))

	refreshUC := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return currentTime },
		ProviderFX: map[string]string{},
	}

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    20 * time.Minute,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err) // Should not fail despite SetLastManualRefresh error
	require.Equal(t, 0, result.CoinsProcessed)
}

func TestUC12ManualRefresh_CooldownMinimum_WhenLessThan1Second(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refreshUC := &app.RefreshQuotesUseCase{}
	controlRepo := mocks.NewMockRefreshControlRepository(ctrl)

	lastRefresh := time.Date(2026, 1, 22, 10, 0, 0, 500000000, time.UTC) // + 0.5 seconds
	currentTime := time.Date(2026, 1, 22, 10, 0, 0, 600000000, time.UTC) // + 0.6 seconds (0.1 seconds later)
	cooldown := 1 * time.Second

	controlRepo.EXPECT().
		GetLastManualRefresh(gomock.Any()).
		Return(lastRefresh, true, nil)

	uc := app.ManualRefreshWithCooldownUseCase{
		RefreshUC:   *refreshUC,
		ControlRepo: controlRepo,
		Now:         func() time.Time { return currentTime },
		Cooldown:    cooldown,
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.ErrorIs(t, err, app.ErrCooldownActive)
	// Should be minimum 1 second
	require.Equal(t, 1, result.RetryAfterSeconds)
}

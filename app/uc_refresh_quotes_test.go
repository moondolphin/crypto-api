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

func TestUCRefreshQuotes_Success_NoEnabledCoins(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return([]domain.Coin{}, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, result.CoinsProcessed)
	require.Equal(t, 0, result.QuotesSaved)
	require.Equal(t, 0, result.Failed)
}

func TestUCRefreshQuotes_RepoError_WhenListEnabledFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(nil, errors.New("db_error"))

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		ProviderFX: map[string]string{"binance": "USDT"},
	}

	// Act
	_, err := uc.Execute(context.Background())

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "db_error")
}

func TestUCRefreshQuotes_Success_ProcessesMultipleCoinsAndProviders(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)
	coingeckoProvider := mocks.NewMockPriceProvider(ctrl)

	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: "bitcoin"},
		{ID: 2, Symbol: "ETH", Enabled: true, BinanceSymbol: "ETHUSDT", CoinGeckoID: "ethereum"},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	// Setup provider registry
	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true).
		Times(2)

	providers.EXPECT().
		Get("coingecko").
		Return(coingeckoProvider, true).
		Times(2)

	binanceProvider.EXPECT().
		Name().
		Return("binance").
		Times(2)

	coingeckoProvider.EXPECT().
		Name().
		Return("coingecko").
		Times(2)

	// BTC - Binance
	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USDT").
		Return(domain.PriceQuote{Price: "45000", Timestamp: fixedTime.Format(time.RFC3339)}, nil)

	// BTC - CoinGecko
	coingeckoProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USD").
		Return(domain.PriceQuote{Price: "45050", Timestamp: fixedTime.Format(time.RFC3339)}, nil)

	// ETH - Binance
	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[1], "USDT").
		Return(domain.PriceQuote{Price: "2500", Timestamp: fixedTime.Format(time.RFC3339)}, nil)

	// ETH - CoinGecko
	coingeckoProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[1], "USD").
		Return(domain.PriceQuote{Price: "2550", Timestamp: fixedTime.Format(time.RFC3339)}, nil)

	// Quote repo inserts: 4 quotes (2 coins * 2 providers)
	quoteRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(4)

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 2, result.CoinsProcessed)
	require.Equal(t, 4, result.QuotesSaved)
	require.Equal(t, 0, result.Failed)
}

func TestUCRefreshQuotes_Success_SkipsMissingProviderIDs(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)

	// Coin missing CoinGecko ID
	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: ""},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	// Only binance provider called
	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true)

	binanceProvider.EXPECT().
		Name().
		Return("binance")

	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USDT").
		Return(domain.PriceQuote{Price: "45000", Timestamp: fixedTime.Format(time.RFC3339)}, nil)

	quoteRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.CoinsProcessed)
	require.Equal(t, 1, result.QuotesSaved)
	require.Equal(t, 0, result.Failed)
}

func TestUCRefreshQuotes_Success_CountsFailures_WhenProviderUnavailable(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)

	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: "bitcoin"},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	// Provider not available
	providers.EXPECT().
		Get("binance").
		Return(nil, false)

	providers.EXPECT().
		Get("coingecko").
		Return(nil, false)

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.CoinsProcessed)
	require.Equal(t, 0, result.QuotesSaved)
	require.Equal(t, 2, result.Failed) // 2 providers for 1 coin
}

func TestUCRefreshQuotes_Success_CountsFailures_WhenProviderFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)

	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: "bitcoin"},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true)

	providers.EXPECT().
		Get("coingecko").
		Return(nil, false)

	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USDT").
		Return(domain.PriceQuote{}, errors.New("api_error"))

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.CoinsProcessed)
	require.Equal(t, 0, result.QuotesSaved)
	require.Equal(t, 2, result.Failed)
}

func TestUCRefreshQuotes_Success_CountsFailures_WhenQuoteInsertFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)

	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: ""},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true)

	binanceProvider.EXPECT().
		Name().
		Return("binance")

	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USDT").
		Return(domain.PriceQuote{Price: "45000", Timestamp: fixedTime.Format(time.RFC3339)}, nil)

	quoteRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(errors.New("insert_error"))

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT", "coingecko": "USD"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.CoinsProcessed)
	require.Equal(t, 0, result.QuotesSaved)
	require.Equal(t, 1, result.Failed)
}

func TestUCRefreshQuotes_Success_UsesProviderTimestamp_WhenAvailable(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)

	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: ""},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)
	providerTime := time.Date(2026, 1, 22, 9, 30, 0, 0, time.UTC)

	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true)

	binanceProvider.EXPECT().
		Name().
		Return("binance")

	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USDT").
		Return(domain.PriceQuote{
			Price:     "45000",
			Timestamp: providerTime.Format(time.RFC3339),
		}, nil)

	quoteRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, q domain.Quote) {
			// Verify provider's timestamp was used
			require.Equal(t, providerTime, q.QuotedAt)
		}).
		Return(nil)

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.QuotesSaved)
}

func TestUCRefreshQuotes_Success_UsesCurrentTime_WhenProviderTimestampInvalid(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)
	quoteRepo := mocks.NewMockQuoteRepository(ctrl)
	providers := mocks.NewMockPriceProviderRegistry(ctrl)
	binanceProvider := mocks.NewMockPriceProvider(ctrl)

	coins := []domain.Coin{
		{ID: 1, Symbol: "BTC", Enabled: true, BinanceSymbol: "BTCUSDT", CoinGeckoID: ""},
	}

	coinRepo.EXPECT().
		ListEnabled(gomock.Any()).
		Return(coins, nil)

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	providers.EXPECT().
		Get("binance").
		Return(binanceProvider, true)

	binanceProvider.EXPECT().
		Name().
		Return("binance")

	binanceProvider.EXPECT().
		GetCurrentPrice(gomock.Any(), coins[0], "USDT").
		Return(domain.PriceQuote{
			Price:     "45000",
			Timestamp: "invalid-timestamp",
		}, nil)

	quoteRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, q domain.Quote) {
			// Verify current time was used when timestamp is invalid
			require.Equal(t, fixedTime, q.QuotedAt)
		}).
		Return(nil)

	uc := app.RefreshQuotesUseCase{
		CoinRepo:   coinRepo,
		QuoteRepo:  quoteRepo,
		Providers:  providers,
		Now:        func() time.Time { return fixedTime },
		ProviderFX: map[string]string{"binance": "USDT"},
	}

	// Act
	result, err := uc.Execute(context.Background())

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, result.QuotesSaved)
}

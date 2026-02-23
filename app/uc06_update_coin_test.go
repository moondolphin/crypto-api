package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/domain"
	"github.com/moondolphin/crypto-api/test/mocks"
)

func TestUC06UpdateCoin_InvalidInput_WhenSymbolEmpty(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := true
	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:  "",
		Enabled: &enabled,
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidCoinUpdate)
}

func TestUC06UpdateCoin_InvalidInput_WhenNoFieldsToUpdate(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:        "BTC",
		Enabled:       nil,
		CoinGeckoID:   "",
		BinanceSymbol: "",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidCoinUpdate)
}

func TestUC06UpdateCoin_CoinNotFound_WhenDoesNotExist(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := false
	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:  "BTC",
		Enabled: &enabled,
	})

	// Assert
	require.ErrorIs(t, err, app.ErrCoinNotFound)
}

func TestUC06UpdateCoin_RepoError_WhenGetBySymbolFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(nil, errors.New("db_error"))

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := true
	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:  "BTC",
		Enabled: &enabled,
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "db_error")
}

func TestUC06UpdateCoin_RepoError_WhenUpsertFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	existingCoin := &domain.Coin{
		ID:      1,
		Symbol:  "BTC",
		Enabled: true,
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("upsert_error"))

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := false
	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:  "BTC",
		Enabled: &enabled,
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "upsert_error")
}

func TestUC06UpdateCoin_Success_UpdatesEnabledField(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	existingCoin := &domain.Coin{
		ID:      1,
		Symbol:  "BTC",
		Enabled: true,
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	updatedCoin := &domain.Coin{
		ID:      1,
		Symbol:  "BTC",
		Enabled: false,
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, int64(1), c.ID)
			require.False(t, c.Enabled)
		}).
		Return(updatedCoin, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := false
	// Act
	result, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:  "BTC",
		Enabled: &enabled,
	})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.ID)
	require.False(t, result.Enabled)
}

func TestUC06UpdateCoin_Success_UpdatesCoinGeckoID(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	existingCoin := &domain.Coin{
		ID:          1,
		Symbol:      "BTC",
		Enabled:     true,
		CoinGeckoID: "bitcoin",
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	updatedCoin := &domain.Coin{
		ID:          1,
		Symbol:      "BTC",
		Enabled:     true,
		CoinGeckoID: "new-bitcoin-id",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, "new-bitcoin-id", c.CoinGeckoID)
		}).
		Return(updatedCoin, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:      "BTC",
		CoinGeckoID: "new-bitcoin-id",
	})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "new-bitcoin-id", result.CoinGeckoID)
}

func TestUC06UpdateCoin_Success_UpdatesBinanceSymbol(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	existingCoin := &domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		BinanceSymbol: "BTCUSDT",
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	updatedCoin := &domain.Coin{
		ID:            1,
		Symbol:        "BTC",
		Enabled:       true,
		BinanceSymbol: "BTCBUSD",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, "BTCBUSD", c.BinanceSymbol)
		}).
		Return(updatedCoin, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:        "BTC",
		BinanceSymbol: "BTCBUSD",
	})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "BTCBUSD", result.BinanceSymbol)
}

func TestUC06UpdateCoin_Success_UpdatesMultipleFields(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

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
		Enabled:       false,
		CoinGeckoID:   "new-bitcoin",
		BinanceSymbol: "BTCBUSD",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.False(t, c.Enabled)
			require.Equal(t, "new-bitcoin", c.CoinGeckoID)
			require.Equal(t, "BTCBUSD", c.BinanceSymbol)
		}).
		Return(updatedCoin, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := false
	// Act
	result, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:        "BTC",
		Enabled:       &enabled,
		CoinGeckoID:   "new-bitcoin",
		BinanceSymbol: "BTCBUSD",
	})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Enabled)
	require.Equal(t, "new-bitcoin", result.CoinGeckoID)
	require.Equal(t, "BTCBUSD", result.BinanceSymbol)
}

func TestUC06UpdateCoin_Success_NormalizesSymbol(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	// Should normalize " btc " to "BTC"
	existingCoin := &domain.Coin{
		ID:      1,
		Symbol:  "BTC",
		Enabled: true,
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	updatedCoin := &domain.Coin{
		ID:      1,
		Symbol:  "BTC",
		Enabled: false,
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Return(updatedCoin, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	enabled := false
	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:  " btc ",
		Enabled: &enabled,
	})

	// Assert
	require.NoError(t, err)
}

func TestUC06UpdateCoin_Success_TrimsWhitespacesInOptionalFields(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinRepo := mocks.NewMockCoinRepository(ctrl)

	existingCoin := &domain.Coin{
		ID:          1,
		Symbol:      "BTC",
		Enabled:     true,
		CoinGeckoID: "bitcoin",
	}

	coinRepo.EXPECT().
		GetBySymbol(gomock.Any(), "BTC").
		Return(existingCoin, nil)

	updatedCoin := &domain.Coin{
		ID:          1,
		Symbol:      "BTC",
		Enabled:     true,
		CoinGeckoID: "new-id",
	}

	coinRepo.EXPECT().
		Upsert(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, c domain.Coin) {
			require.Equal(t, "new-id", c.CoinGeckoID)
		}).
		Return(updatedCoin, nil)

	uc := app.UpdateCoinUseCase{
		CoinRepo: coinRepo,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.UpdateCoinInput{
		Symbol:      "BTC",
		CoinGeckoID: "  new-id  ",
	})

	// Assert
	require.NoError(t, err)
}

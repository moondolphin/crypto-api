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

func TestUC02RegisterUser_InvalidEmail_WhenMissingAtSign(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    "invalidgmail.com",
		Password: "SecurePassword123",
		Name:     "John Doe",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidEmail)
}

func TestUC02RegisterUser_InvalidEmail_WhenEmailTooShort(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    "a@b",
		Password: "SecurePassword123",
		Name:     "John Doe",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidEmail)
}

func TestUC02RegisterUser_InvalidPassword_WhenTooShort(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    "john@example.com",
		Password: "short",
		Name:     "John Doe",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidPassword)
}

func TestUC02RegisterUser_InvalidName_WhenEmpty(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
		Name:     "",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidName)
}

func TestUC02RegisterUser_EmailAlreadyRegistered_WhenExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	userRepo.EXPECT().
		ExistsByEmail(gomock.Any(), "john@example.com").
		Return(true, nil)

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
		Name:     "John Doe",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrEmailAlreadyRegistered)
}

func TestUC02RegisterUser_HasherError_WhenHashingFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	userRepo.EXPECT().
		ExistsByEmail(gomock.Any(), "john@example.com").
		Return(false, nil)

	hasher.EXPECT().
		Hash("SecurePassword123").
		Return("", errors.New("hashing_failed"))

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
		Name:     "John Doe",
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "hashing_failed")
}

func TestUC02RegisterUser_Success_CreatesUserAndReturnsIt(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)

	fixedTime := time.Date(2026, 1, 22, 10, 30, 0, 0, time.UTC)

	userRepo.EXPECT().
		ExistsByEmail(gomock.Any(), "john@example.com").
		Return(false, nil)

	hasher.EXPECT().
		Hash("SecurePassword123").
		Return("$2a$10$encrypted_hash", nil)

	createdUser := domain.User{
		ID:           1,
		Email:        "john@example.com",
		Name:         "John Doe",
		PasswordHash: "$2a$10$encrypted_hash",
		CreatedAt:    fixedTime,
	}

	userRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(createdUser, nil)

	uc := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Now:      func() time.Time { return fixedTime },
	}

	// Act
	result, err := uc.Execute(context.Background(), app.RegisterInput{
		Email:    " john@example.com ",
		Password: "SecurePassword123",
		Name:     " John Doe ",
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, "john@example.com", result.Email)
	require.Equal(t, "John Doe", result.Name)
	require.Equal(t, fixedTime, result.CreatedAt)
}

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

func TestUC03Login_BadRequest_WhenMissingEmail(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "",
		Password: "SecurePassword123",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrBadRequest)
}

func TestUC03Login_BadRequest_WhenMissingPassword(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrBadRequest)
}

func TestUC03Login_InvalidEmail_WhenMissingAtSign(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "invalidemail",
		Password: "SecurePassword123",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidEmailL)
}

func TestUC03Login_InvalidEmail_WhenEmailTooShort(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "a@b",
		Password: "SecurePassword123",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidEmailL)
}

func TestUC03Login_InvalidPassword_WhenTooShort(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "short",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidPasswordL)
}

func TestUC03Login_InvalidCredentials_WhenUserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	userRepo.EXPECT().
		FindByEmail(gomock.Any(), "john@example.com").
		Return(nil, nil)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidCredentials)
}

func TestUC03Login_InvalidCredentials_WhenPasswordMismatch(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	user := domain.User{
		ID:           1,
		Email:        "john@example.com",
		Name:         "John Doe",
		PasswordHash: "$2a$10$encrypted_hash",
	}

	userRepo.EXPECT().
		FindByEmail(gomock.Any(), "john@example.com").
		Return(&user, nil)

	hasher.EXPECT().
		Compare("$2a$10$encrypted_hash", "WrongPassword123").
		Return(false, nil)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "WrongPassword123",
	})

	// Assert
	require.ErrorIs(t, err, app.ErrInvalidCredentials)
}

func TestUC03Login_HasherError_WhenCompareFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	user := domain.User{
		ID:           1,
		Email:        "john@example.com",
		Name:         "John Doe",
		PasswordHash: "$2a$10$encrypted_hash",
	}

	userRepo.EXPECT().
		FindByEmail(gomock.Any(), "john@example.com").
		Return(&user, nil)

	hasher.EXPECT().
		Compare("$2a$10$encrypted_hash", "SecurePassword123").
		Return(false, errors.New("hasher_error"))

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "hasher_error")
}

func TestUC03Login_TokenServiceError_WhenGenerateTokenFails(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	user := domain.User{
		ID:           1,
		Email:        "john@example.com",
		Name:         "John Doe",
		PasswordHash: "$2a$10$encrypted_hash",
	}

	userRepo.EXPECT().
		FindByEmail(gomock.Any(), "john@example.com").
		Return(&user, nil)

	hasher.EXPECT().
		Compare("$2a$10$encrypted_hash", "SecurePassword123").
		Return(true, nil)

	tokens.EXPECT().
		Generate(int64(1), "john@example.com").
		Return("", errors.New("token_error"))

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
	}

	// Act
	_, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
	})

	// Assert
	require.Error(t, err)
	require.Equal(t, err.Error(), "token_error")
}

func TestUC03Login_Success_ReturnsAccessToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	user := domain.User{
		ID:           1,
		Email:        "john@example.com",
		Name:         "John Doe",
		PasswordHash: "$2a$10$encrypted_hash",
	}

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)
	ttl := 60 * time.Minute

	userRepo.EXPECT().
		FindByEmail(gomock.Any(), "john@example.com").
		Return(&user, nil)

	hasher.EXPECT().
		Compare("$2a$10$encrypted_hash", "SecurePassword123").
		Return(true, nil)

	tokens.EXPECT().
		Generate(int64(1), "john@example.com").
		Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", nil)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
		Now:      func() time.Time { return fixedTime },
		TTL:      ttl,
	}

	// Act
	result, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    " john@example.com ",
		Password: "SecurePassword123",
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, "Bearer", result.TokenType)
	require.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", result.AccessToken)
	require.Equal(t, fixedTime.Add(ttl), result.ExpiresAt)
}

func TestUC03Login_Success_UsesDefaultTTLWhenNotSet(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockPasswordHasher(ctrl)
	tokens := mocks.NewMockTokenService(ctrl)

	user := domain.User{
		ID:           1,
		Email:        "john@example.com",
		Name:         "John Doe",
		PasswordHash: "$2a$10$encrypted_hash",
	}

	fixedTime := time.Date(2026, 1, 22, 10, 0, 0, 0, time.UTC)

	userRepo.EXPECT().
		FindByEmail(gomock.Any(), "john@example.com").
		Return(&user, nil)

	hasher.EXPECT().
		Compare("$2a$10$encrypted_hash", "SecurePassword123").
		Return(true, nil)

	tokens.EXPECT().
		Generate(int64(1), "john@example.com").
		Return("token123", nil)

	uc := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   tokens,
		Now:      func() time.Time { return fixedTime },
		TTL:      0, // Use default
	}

	// Act
	result, err := uc.Execute(context.Background(), app.LoginInput{
		Email:    "john@example.com",
		Password: "SecurePassword123",
	})

	// Assert
	require.NoError(t, err)
	expectedTTL := 60 * time.Minute
	require.Equal(t, fixedTime.Add(expectedTTL), result.ExpiresAt)
}

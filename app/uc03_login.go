package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrInvalidEmailL      = errors.New("invalid_email")
	ErrInvalidPasswordL   = errors.New("invalid_password")
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type LoginUseCase struct {
	UserRepo domain.UserRepository
	Hasher   domain.PasswordHasher
	Tokens   domain.TokenService
	Now      func() time.Time
	TTL      time.Duration
}

func (uc LoginUseCase) Execute(ctx context.Context, in LoginInput) (LoginOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	password := in.Password

	if email == "" || password == "" {
		return LoginOutput{}, ErrBadRequest
	}

	if !strings.Contains(email, "@") || len(email) < 5 {
		return LoginOutput{}, ErrInvalidEmailL
	}
	if len(password) < 8 {
		return LoginOutput{}, ErrInvalidPasswordL
	}

	u, err := uc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return LoginOutput{}, err
	}
	if u == nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	ok, err := uc.Hasher.Compare(u.PasswordHash, password)
	if err != nil {
		return LoginOutput{}, err
	}
	if !ok {
		return LoginOutput{}, ErrInvalidCredentials
	}

	token, err := uc.Tokens.Generate(u.ID, u.Email)
	if err != nil {
		return LoginOutput{}, err
	}

	now := uc.Now
	if now == nil {
		now = time.Now
	}

	ttl := uc.TTL
	if ttl <= 0 {
		ttl = 60 * time.Minute
	}

	return LoginOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   now().UTC().Add(ttl),
	}, nil
}

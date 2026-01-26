package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email_already_registered")
	ErrInvalidEmail           = errors.New("invalid_email")
	ErrInvalidPassword        = errors.New("invalid_password")
	ErrInvalidName            = errors.New("invalid_name")
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserOutput struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUserUseCase struct {
	UserRepo domain.UserRepository
	Hasher   domain.PasswordHasher
	Now      func() time.Time
}

func (uc RegisterUserUseCase) Execute(ctx context.Context, in RegisterInput) (UserOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	name := strings.TrimSpace(in.Name)
	password := in.Password

	if !strings.Contains(email, "@") || len(email) < 5 {
		return UserOutput{}, ErrInvalidEmail
	}
	if len(password) < 8 {
		return UserOutput{}, ErrInvalidPassword
	}
	if name == "" {
		return UserOutput{}, ErrInvalidName
	}

	exists, err := uc.UserRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return UserOutput{}, err
	}
	if exists {
		return UserOutput{}, ErrEmailAlreadyRegistered
	}

	hash, err := uc.Hasher.Hash(password)
	if err != nil {
		return UserOutput{}, err
	}

	now := uc.Now
	if now == nil {
		now = time.Now
	}

	u := domain.User{
		Email:        email,
		Name:         name,
		PasswordHash: hash,
		CreatedAt:    now().UTC().Truncate(time.Second),
	}

	created, err := uc.UserRepo.Create(ctx, u)
	if err != nil {
		return UserOutput{}, err
	}

	return UserOutput{
		ID:        created.ID,
		Email:     created.Email,
		Name:      created.Name,
		CreatedAt: created.CreatedAt,
	}, nil
}

package domain

//go:generate echo Generating mocks for user_port.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=user_port.go -destination=../test/mocks/user_port_mock.go -package=mocks

import "context"

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, u User) (User, error)

	FindByEmail(ctx context.Context, email string) (*User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) (bool, error)
}

type TokenService interface {
	Generate(userID int64, email string) (string, error)
}

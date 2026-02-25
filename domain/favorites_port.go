package domain

//go:generate echo Generating mocks for favorites_port.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=favorites_port.go -destination=../test/mocks/favorites_port_mock.go -package=mocks

import "context"

type FavoritesRepository interface {
	AddFavoriteCoinToUser(ctx context.Context, userID, coinID int64) error
	RemoveFavoriteCoinFromUser(ctx context.Context, userID, coinID int64) error
	ListFavoriteCoinIDsByUser(ctx context.Context, userID int64) ([]Coin, error)
}

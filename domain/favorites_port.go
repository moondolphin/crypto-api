package domain

import "context"

type FavoritesRepository interface {
	AddFavoriteCoinToUser(ctx context.Context, userID, coinID int64) error
	RemoveFavoriteCoinFromUser(ctx context.Context, userID, coinID int64) error
	ListFavoriteCoinIDsByUser(ctx context.Context, userID int64) ([]Coin, error)
}

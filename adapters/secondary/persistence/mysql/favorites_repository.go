package mysql

import (
	"context"
	"database/sql"

	"github.com/moondolphin/crypto-api/domain"
)

type MySQLFavoritesRepository struct {
	DB *sql.DB
}

func NewMySQLFavoritesRepository(db *sql.DB) *MySQLFavoritesRepository {
	return &MySQLFavoritesRepository{DB: db}
}

// Idempotente: si ya existe (user_id, coin_id)
func (r *MySQLFavoritesRepository) AddFavoriteCoinToUser(ctx context.Context, userID, coinID int64) error {
	const q = `
		INSERT IGNORE INTO user_favorites (user_id, coin_id)
		VALUES (?, ?)
	`
	_, err := r.DB.ExecContext(ctx, q, userID, coinID)
	return err
}

// Idempotente: si no existe
func (r *MySQLFavoritesRepository) RemoveFavoriteCoinFromUser(ctx context.Context, userID, coinID int64) error {
	const q = `
		DELETE FROM user_favorites
		WHERE user_id = ? AND coin_id = ?
	`
	_, err := r.DB.ExecContext(ctx, q, userID, coinID)
	return err
}

func (r *MySQLFavoritesRepository) ListFavoriteCoinIDsByUser(ctx context.Context, userID int64) ([]domain.Coin, error) {
	const q = `
		SELECT c.id, c.symbol, c.enabled, c.coingecko_id, c.binance_symbol
		FROM user_favorites uf
		JOIN coins c ON c.id = uf.coin_id
		WHERE uf.user_id = ?
		ORDER BY c.symbol ASC
	`

	rows, err := r.DB.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Coin, 0, 16)
	for rows.Next() {
		var c domain.Coin
		if err := rows.Scan(&c.ID, &c.Symbol, &c.Enabled, &c.CoinGeckoID, &c.BinanceSymbol); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

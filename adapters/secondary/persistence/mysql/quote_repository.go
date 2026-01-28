package mysql

import (
	"context"
	"database/sql"

	"github.com/moondolphin/crypto-api/domain"
)

type MySQLQuoteRepository struct {
	DB *sql.DB
}

func NewMySQLQuoteRepository(db *sql.DB) *MySQLQuoteRepository {
	return &MySQLQuoteRepository{DB: db}
}

func (r *MySQLQuoteRepository) Insert(ctx context.Context, q domain.Quote) error {
	const stmt = `
		INSERT INTO quotes (coin_id, symbol, provider, currency, price, quoted_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.DB.ExecContext(ctx, stmt, q.CoinID, q.Symbol, q.Provider, q.Currency, q.Price, q.QuotedAt)
	return err
}

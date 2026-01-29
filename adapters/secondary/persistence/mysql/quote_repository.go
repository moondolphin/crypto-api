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

func (r MySQLQuoteRepository) GetLatest(ctx context.Context, symbol, provider, currency string) (*domain.PriceQuote, error) {
	const base = `
SELECT symbol, provider, currency, price, quoted_at
FROM quotes
WHERE symbol = ?
`

	q := base
	args := []any{symbol}

	if provider != "" {
		q += " AND provider = ?"
		args = append(args, provider)
	}
	if currency != "" {
		q += " AND currency = ?"
		args = append(args, currency)
	}

	q += " ORDER BY quoted_at DESC LIMIT 1"

	var out domain.PriceQuote
	if err := r.DB.QueryRowContext(ctx, q, args...).Scan(
		&out.Symbol,
		&out.Provider,
		&out.Currency,
		&out.Price,
		&out.Timestamp,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &out, nil
}

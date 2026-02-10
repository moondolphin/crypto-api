package mysql

import (
	"context"
	"database/sql"

	"github.com/moondolphin/crypto-api/domain"
)

type MySQLCoinRepository struct {
	DB *sql.DB
}

func NewMySQLCoinRepository(db *sql.DB) *MySQLCoinRepository {
	return &MySQLCoinRepository{DB: db}
}

func (r *MySQLCoinRepository) GetEnabledBySymbol(ctx context.Context, symbol string) (*domain.Coin, error) {
	const q = `
		SELECT id, symbol, enabled, coingecko_id, binance_symbol
		FROM coins
		WHERE symbol = ? AND enabled = true
		LIMIT 1
	`

	row := r.DB.QueryRowContext(ctx, q, symbol)

	var c domain.Coin
	err := row.Scan(&c.ID, &c.Symbol, &c.Enabled, &c.CoinGeckoID, &c.BinanceSymbol)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MySQLCoinRepository) ListEnabled(ctx context.Context) ([]domain.Coin, error) {
	const q = `
		SELECT id, symbol, enabled, coingecko_id, binance_symbol
		FROM coins
		WHERE enabled = true
		ORDER BY symbol ASC
	`

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Coin, 0, 64)
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

func (r *MySQLCoinRepository) GetBySymbol(ctx context.Context, symbol string) (*domain.Coin, error) {
	const q = `
		SELECT id, symbol, enabled, coingecko_id, binance_symbol
		FROM coins
		WHERE symbol = ?
		LIMIT 1
	`

	var c domain.Coin
	err := r.DB.QueryRowContext(ctx, q, symbol).Scan(
		&c.ID, &c.Symbol, &c.Enabled, &c.CoinGeckoID, &c.BinanceSymbol,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MySQLCoinRepository) Upsert(ctx context.Context, c domain.Coin) (*domain.Coin, error) {
	const stmt = `
	INSERT INTO coins (symbol, enabled, coingecko_id, binance_symbol)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		enabled = VALUES(enabled),
		coingecko_id = COALESCE(NULLIF(VALUES(coingecko_id), ''), coingecko_id),
		binance_symbol = COALESCE(NULLIF(VALUES(binance_symbol), ''), binance_symbol)
`

	_, err := r.DB.ExecContext(ctx, stmt, c.Symbol, c.Enabled, c.CoinGeckoID, c.BinanceSymbol)
	if err != nil {
		return nil, err
	}
	return r.GetBySymbol(ctx, c.Symbol)
}

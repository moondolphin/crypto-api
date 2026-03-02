package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *MySQLQuoteRepository) GetLatest(ctx context.Context, symbol, provider, currency string) (*domain.PriceQuote, error) {
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

// buildQuoteWhere arma el WHERE dinámico y sus args para reutilizarlo
// tanto en listados como en "faceted filters".
func buildQuoteWhere(f domain.QuoteFilter) (string, []any) {
	where := " WHERE 1=1"
	args := make([]any, 0, 12)

	if f.Symbol != "" {
		where += " AND symbol = ?"
		args = append(args, f.Symbol)
	}
	if f.Provider != "" {
		where += " AND provider = ?"
		args = append(args, f.Provider)
	}
	if f.Currency != "" {
		where += " AND currency = ?"
		args = append(args, f.Currency)
	}
	if f.From != nil {
		where += " AND quoted_at >= ?"
		args = append(args, *f.From)
	}
	if f.To != nil {
		where += " AND quoted_at <= ?"
		args = append(args, *f.To)
	}
	if f.MinPrice != nil {
		where += " AND price >= ?"
		args = append(args, *f.MinPrice)
	}
	if f.MaxPrice != nil {
		where += " AND price <= ?"
		args = append(args, *f.MaxPrice)
	}

	return where, args
}

func (r *MySQLQuoteRepository) ListFilter(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
	// defaults defensivos
	page := f.Page
	if page <= 0 {
		page = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	where, args := buildQuoteWhere(f)

	// 1) COUNT total (para summary)
	countSQL := "SELECT COUNT(*) FROM quotes" + where
	var total int
	if err := r.DB.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2) SELECT paginado
	listSQL := `
SELECT id, coin_id, symbol, provider, currency, price, quoted_at, created_at
FROM quotes
` + where + `
ORDER BY quoted_at DESC
LIMIT ? OFFSET ?
`

	listArgs := make([]any, 0, len(args)+2)
	listArgs = append(listArgs, args...)
	listArgs = append(listArgs, pageSize, offset)

	rows, err := r.DB.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]domain.Quote, 0, pageSize)

	for rows.Next() {
		var q domain.Quote
		var priceStr string // DECIMAL -> lo escaneamos como string para no perder precisión

		if err := rows.Scan(
			&q.ID,
			&q.CoinID,
			&q.Symbol,
			&q.Provider,
			&q.Currency,
			&priceStr,
			&q.QuotedAt,
			&q.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		q.Price = priceStr
		out = append(out, q)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return out, total, nil
}

// ListAvailableFilters devuelve "faceted filters":
// combos (symbol/provider/currency) y rangos (price/quoted_at) recalculados
// aplicando el mismo tamiz (WHERE dinámico) de QuoteFilter.
func (r *MySQLQuoteRepository) ListAvailableFilters(ctx context.Context, f domain.QuoteFilter) (domain.QuoteFilters, error) {
	where, args := buildQuoteWhere(f)

	distinctList := func(col string) ([]string, error) {
		q := fmt.Sprintf(`SELECT DISTINCT %s FROM quotes %s ORDER BY %s ASC`, col, where, col)
		rows, err := r.DB.QueryContext(ctx, q, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		out := make([]string, 0, 16)
		for rows.Next() {
			var v sql.NullString
			if err := rows.Scan(&v); err != nil {
				return nil, err
			}
			if v.Valid && v.String != "" {
				out = append(out, v.String)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return out, nil
	}

	symbols, err := distinctList("symbol")
	if err != nil {
		return domain.QuoteFilters{}, err
	}
	providers, err := distinctList("provider")
	if err != nil {
		return domain.QuoteFilters{}, err
	}
	currencies, err := distinctList("currency")
	if err != nil {
		return domain.QuoteFilters{}, err
	}

	// rangos de price (DECIMAL) -> string sin perder precisión
	var minPrice, maxPrice sql.NullString
	priceQ := `SELECT CAST(MIN(price) AS CHAR), CAST(MAX(price) AS CHAR) FROM quotes` + where
	if err := r.DB.QueryRowContext(ctx, priceQ, args...).Scan(&minPrice, &maxPrice); err != nil {
		return domain.QuoteFilters{}, err
	}

	// rangos de quoted_at
	var minTime, maxTime sql.NullTime
	timeQ := `SELECT MIN(quoted_at), MAX(quoted_at) FROM quotes` + where
	if err := r.DB.QueryRowContext(ctx, timeQ, args...).Scan(&minTime, &maxTime); err != nil {
		return domain.QuoteFilters{}, err
	}

	out := domain.QuoteFilters{
		Symbols:    symbols,
		Providers:  providers,
		Currencies: currencies,
	}

	if minPrice.Valid {
		out.MinPrice = minPrice.String
	}
	if maxPrice.Valid {
		out.MaxPrice = maxPrice.String
	}

	if minTime.Valid {
		t := minTime.Time.UTC()
		out.From = &t
	}
	if maxTime.Valid {
		t := maxTime.Time.UTC()
		out.To = &t
	}

	return out, nil
}

var _ = time.RFC3339

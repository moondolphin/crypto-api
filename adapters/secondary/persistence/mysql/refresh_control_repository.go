package mysql

import (
	"context"
	"database/sql"
	"time"
)

type MySQLRefreshControlRepository struct {
	DB *sql.DB
}

func NewMySQLRefreshControlRepository(db *sql.DB) *MySQLRefreshControlRepository {
	return &MySQLRefreshControlRepository{DB: db}
}

const keyLastManualRefresh = "last_manual_refresh_rfc3339"

func (r *MySQLRefreshControlRepository) GetLastManualRefresh(ctx context.Context) (time.Time, bool, error) {
	const q = `SELECT value FROM refresh_control WHERE ` + "`key`" + ` = ? LIMIT 1`

	var v string
	err := r.DB.QueryRowContext(ctx, q, keyLastManualRefresh).Scan(&v)
	if err == sql.ErrNoRows {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		// si quedó basura, lo tratamos como "no hay dato"
		return time.Time{}, false, nil
	}

	return t.UTC(), true, nil
}

func (r *MySQLRefreshControlRepository) SetLastManualRefresh(ctx context.Context, t time.Time) error {
	const stmt = `
		INSERT INTO refresh_control (` + "`key`" + `, value)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
			value = VALUES(value)
	`
	_, err := r.DB.ExecContext(ctx, stmt, keyLastManualRefresh, t.UTC().Format(time.RFC3339))
	return err
}

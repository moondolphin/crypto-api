package mysql

import (
	"context"
	"database/sql"

	"github.com/moondolphin/crypto-api/domain"
)

type MySQLUserRepository struct {
	DB *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{DB: db}
}

func (r *MySQLUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const q = `SELECT 1 FROM users WHERE email = ? LIMIT 1`
	var one int
	err := r.DB.QueryRowContext(ctx, q, email).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MySQLUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const q = `
		INSERT INTO users (email, name, password_hash, created_at)
		VALUES (?, ?, ?, ?)
	`
	res, err := r.DB.ExecContext(ctx, q, u.Email, u.Name, u.PasswordHash, u.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}

	u.ID = id
	return u, nil
}

func (r *MySQLUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, name, password_hash, created_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`
	row := r.DB.QueryRowContext(ctx, q, email)

	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

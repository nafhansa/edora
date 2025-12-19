package repository

import (
	"context"
	"database/sql"
	"errors"

	"edora/backend/internal/models"
)

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	q := `SELECT id, username, password, role, created_at FROM users WHERE username = $1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, username)
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

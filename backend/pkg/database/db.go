package database

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect opens a *sql.DB using the pgx stdlib driver and returns it as
// an interface{} to preserve the existing function signature used by the
// rest of the application. Callers should type-assert to *sql.DB.
func Connect(ctx context.Context, url string) (interface{}, error) {
	_ = ctx

	// Use the pgx stdlib driver registered under the name "pgx".
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	// Optionally verify connection now
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

package database

import "context"

// Connect is a light stub that returns nil. In production replace with real
// pgxpool connection. Returning nil allows the app to run in environments
// without Postgres for smoke tests.
func Connect(ctx context.Context, url string) (interface{}, error) {
    _ = ctx
    _ = url
    return nil, nil
}

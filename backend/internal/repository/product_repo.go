package repository

import (
    "context"
    "log"

    "edora/backend/internal/models"
)

type ProductRepository struct{
    db interface{}
}

func NewProductRepository(db interface{}) *ProductRepository {
    return &ProductRepository{db: db}
}

// List returns products. If a real DB connection is not provided, it returns
// an empty slice to allow smoke tests without Postgres.
func (r *ProductRepository) List(ctx context.Context, limit int) ([]models.Product, error) {
    if r.db == nil {
        log.Println("product repo: db is nil, returning empty slice")
        return []models.Product{}, nil
    }
    // In a fully wired environment, implement DB query here.
    return []models.Product{}, nil
}

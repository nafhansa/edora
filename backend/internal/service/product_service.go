package service

import (
    "context"
    "edora/backend/internal/models"
    "edora/backend/internal/repository"
)

type ProductService struct{
    repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository, _ interface{}) *ProductService {
    return &ProductService{repo: repo}
}

func (s *ProductService) List(ctx context.Context, limit int) ([]models.Product, error) {
    // Directly read from repository. Caching disabled in this smoke-test friendly build.
    return s.repo.List(ctx, limit)
}

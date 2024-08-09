package service

import (
	"chi-sqlx/database/entity"
	"chi-sqlx/database/repository"
	"context"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	return s.repo.CreateProduct(ctx, p)
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*entity.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]entity.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	return s.repo.UpdateProduct(ctx, p)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

package repository

import (
	"chi-sqlx/src/database/entity"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (repo *ProductRepository) CreateProduct(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	res, err := repo.db.NamedExecContext(ctx, "INSERT INTO product (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)
	if err != nil {
		return nil, fmt.Errorf("error inserting product: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	p.ID = id

	return p, nil
}

func (repo *ProductRepository) GetProduct(ctx context.Context, id int64) (*entity.Product, error) {
	var p entity.Product

	err := repo.db.GetContext(ctx, &p, "SELECT * FROM product WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %v", err)
	}

	return &p, nil
}

func (repo *ProductRepository) ListProducts(ctx context.Context) ([]*entity.Product, error) {
	var products []*entity.Product

	err := repo.db.SelectContext(ctx, &products, "SELECT * FROM product")
	if err != nil {
		return nil, fmt.Errorf("error listing products: %v", err)
	}

	return products, nil
}

func (repo *ProductRepository) UpdateProduct(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	_, err := repo.db.NamedExecContext(ctx, "UPDATE product SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock, updated_at=:updated_at WHERE id=:id", p)
	if err != nil {
		return nil, fmt.Errorf("error updating product: %v", err)
	}

	return p, nil
}

func (repo *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM product WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting product: %v", err)
	}

	return nil
}

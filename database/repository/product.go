package repository

import (
	"chi-sqlx/database/entity"
	"context"
	"fmt"
	"time"

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

const insertProduct = `
	INSERT INTO product (name, image, category, description, rating, num_reviews, price, count_in_stock) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at, updated_at
`

func (repo *ProductRepository) CreateProduct(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	var lastInsertID int64
	var createdAt time.Time
	var updatedAt time.Time

	err := repo.db.QueryRowContext(ctx, insertProduct,
		p.Name,
		p.Image,
		p.Category,
		p.Description,
		p.Rating,
		p.NumReviews,
		p.Price,
		p.CountInStock).
		Scan(&lastInsertID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", err)
	}

	p.ID = lastInsertID
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt

	return p, nil
}

func (repo *ProductRepository) GetProduct(ctx context.Context, id int64) (*entity.Product, error) {
	var p entity.Product

	err := repo.db.GetContext(ctx, &p, "SELECT * FROM product WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %v", err)
	}

	return &p, nil
}

func (repo *ProductRepository) ListProducts(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product

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

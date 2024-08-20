package repository

import (
	"chi-sqlx/database/entity"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	fn(db, mock)
}

func TestCreateProduct(t *testing.T) {
	p := &entity.Product{
		Name:         "test product",
		Image:        "test.png",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        1000.0,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *ProductRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				// Mock the expected query and result
				expectedID := int64(1)
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery("INSERT INTO product (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at").
					WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				record, err := repo.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.NotNil(t, record)
				require.Equal(t, expectedID, record.ID)
				require.Equal(t, expectedCreatedAt, record.CreatedAt)
				require.Equal(t, expectedUpdatedAt, record.UpdatedAt)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting product",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO product (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at").
					WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock).
					WillReturnError(fmt.Errorf("error inserting product"))

				_, err := repo.CreateProduct(context.Background(), p)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProductRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestGetProduct(t *testing.T) {
	p := &entity.Product{
		Name:         "test product",
		Image:        "test.png",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        1000.0,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *ProductRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt, p.DeletedAt)
				mock.ExpectQuery("SELECT * FROM product WHERE id=$1").WithArgs(1).WillReturnRows(rows)

				record, err := repo.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), record.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM product WHERE id=$1").WithArgs(1).WillReturnError(fmt.Errorf("error getting product"))

				_, err := repo.GetProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProductRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestListProduct(t *testing.T) {
	p := &entity.Product{
		Name:         "test product",
		Image:        "test.png",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        1000.0,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *ProductRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt, p.DeletedAt)
				mock.ExpectQuery("SELECT * FROM product").WillReturnRows(rows)

				records, err := repo.ListProducts(context.Background())
				require.NoError(t, err)
				require.Len(t, records, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed querying product",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM product").WillReturnError(fmt.Errorf("error querying product"))

				_, err := repo.ListProducts(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProductRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	p := &entity.Product{
		ID:           1,
		Name:         "test product",
		Image:        "test.png",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        1000.0,
		CountInStock: 10,
	}

	np := &entity.Product{
		ID:           1,
		Name:         "new test product",
		Image:        "test.png",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        1000.0,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *ProductRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				// Mock the expected query and result
				expectedID := int64(1)
				expectedCreatedAt := time.Now()
				expectedUpdatedAt := time.Now()

				mock.ExpectQuery("INSERT INTO product (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at").
					WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt))

				cp, err := repo.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, expectedID, cp.ID)
				require.Equal(t, expectedCreatedAt, cp.CreatedAt)
				require.Equal(t, expectedUpdatedAt, cp.UpdatedAt)

				mock.ExpectExec("UPDATE product SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").
					WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := repo.UpdateProduct(context.Background(), np)
				require.NoError(t, err)
				require.Equal(t, int64(1), up.ID)
				require.Equal(t, np.Name, up.Name)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating product",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE product SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").
					WillReturnError(fmt.Errorf("error updating product"))

				_, err := repo.UpdateProduct(context.Background(), p)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProductRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *ProductRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM product WHERE id=?").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

				err := repo.DeleteProduct(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T, repo *ProductRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM product WHERE id=?").
					WithArgs(1).WillReturnError(fmt.Errorf("error deleting product"))

				err := repo.DeleteProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewProductRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

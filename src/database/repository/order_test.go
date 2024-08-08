package repository

import (
	"chi-sqlx/src/database/entity"
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	ois := []entity.OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.png",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test.png",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &entity.Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *OrderRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO order (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_item (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_item (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectCommit()

				co, err := repo.CreateOrder(context.Background(), o)
				require.NoError(t, err)
				require.Equal(t, int64(1), co.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO order (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WillReturnError(fmt.Errorf("error creating order"))
				mock.ExpectRollback()

				_, err := repo.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order item",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO order (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_item (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnError(fmt.Errorf("error creating order item"))
				mock.ExpectRollback()

				_, err := repo.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed committing transaction",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO order (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_item (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO order_item (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("error committing transaction"))

				_, err := repo.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewOrderRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestGetOrder(t *testing.T) {
	ois := []entity.OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.png",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test.png",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &entity.Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *OrderRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt, o.DeletedAt)

				mock.ExpectQuery("SELECT * FROM order WHERE id=?").WithArgs(1).WillReturnRows(orows)

				oirows := sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(1, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, ois[0].OrderID).
					AddRow(2, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, ois[1].OrderID)

				mock.ExpectQuery("SELECT * FROM order_item WHERE order_id=?").WithArgs(1).WillReturnRows(oirows)

				mo, err := repo.GetOrder(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), mo.ID)

				for i, oi := range mo.Items {
					require.Equal(t, ois[i].Name, oi.Name)
					require.Equal(t, ois[i].Quantity, oi.Quantity)
					require.Equal(t, ois[i].Image, oi.Image)
					require.Equal(t, ois[i].Price, oi.Price)
					require.Equal(t, ois[i].ProductID, oi.ProductID)
				}

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting order",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM order WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting order"))

				_, err := repo.GetOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting order item",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt, o.DeletedAt)

				mock.ExpectQuery("SELECT * FROM order WHERE id=?").WithArgs(1).WillReturnRows(orows)

				mock.ExpectQuery("SELECT * FROM order_item WHERE order_id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting order item"))

				_, err := repo.GetOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewOrderRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestListOrders(t *testing.T) {
	ois := []entity.OrderItem{
		{
			Name:      "test product",
			Quantity:  1,
			Image:     "test.png",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "test product 2",
			Quantity:  2,
			Image:     "test.png",
			Price:     199.99,
			ProductID: 2,
		},
	}

	o := &entity.Order{
		PaymentMethod: "test payment method",
		TaxPrice:      10.0,
		ShippingPrice: 20.0,
		TotalPrice:    129.99,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *OrderRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt, o.DeletedAt)

				mock.ExpectQuery("SELECT * FROM order").WillReturnRows(orows)

				oirows := sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(1, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, ois[0].OrderID).
					AddRow(2, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, ois[1].OrderID)

				mock.ExpectQuery("SELECT * FROM order_item WHERE order_id=?").WithArgs(1).WillReturnRows(oirows)

				mo, err := repo.ListOrders(context.Background())
				require.NoError(t, err)
				require.Len(t, mo, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting order",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM order").WillReturnError(fmt.Errorf("error querying order"))

				_, err := repo.ListOrders(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting order item",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				orows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt, o.DeletedAt)

				mock.ExpectQuery("SELECT * FROM order").WillReturnRows(orows)

				mock.ExpectQuery("SELECT * FROM order_item WHERE order_id=?").WithArgs(1).WillReturnError(fmt.Errorf("error querying order item"))

				_, err := repo.ListOrders(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewOrderRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

func TestDeleteOrder(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *OrderRepository, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_item WHERE order_id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM order WHERE id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repo.DeleteOrder(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order item",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_item WHERE order_id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order item"))
				mock.ExpectRollback()

				err := repo.DeleteOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting order",
			test: func(t *testing.T, repo *OrderRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM order_item WHERE order_id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM order WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order"))
				mock.ExpectRollback()

				err := repo.DeleteOrder(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				repo := NewOrderRepository(db)
				tc.test(t, repo, mock)
			})
		})
	}
}

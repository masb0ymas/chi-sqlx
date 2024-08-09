package repository

import (
	"chi-sqlx/database/entity"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v", rbErr)
		}

		return fmt.Errorf("error in transaction: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (repo *OrderRepository) CreateOrder(ctx context.Context, o *entity.Order) (*entity.Order, error) {
	err := repo.execTx(ctx, func(tx *sqlx.Tx) error {
		// insert into order
		order, err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("error creating order: %v", err)
		}

		for _, oi := range o.Items {
			oi.OrderID = order.ID
			// insert into order item
			err = createOrderItem(ctx, tx, oi)
			if err != nil {
				return fmt.Errorf("error creating order items: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}

	return o, nil
}

func createOrder(ctx context.Context, tx *sqlx.Tx, o *entity.Order) (*entity.Order, error) {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO order (payment_method, tax_price, shipping_price, total_price) VALUES (:payment_method, :tax_price, :shipping_price, :total_price)", o)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}
	o.ID = id

	return o, nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi entity.OrderItem) error {
	res, err := tx.NamedExecContext(ctx, "INSERT INTO order_item (name, quantity, image, price, product_id, order_id) VALUES (:name, :quantity, :image, :price, :product_id, :order_id)", oi)
	if err != nil {
		return fmt.Errorf("error inserting order item: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}
	oi.ID = id

	return nil
}

func (repo *OrderRepository) GetOrder(ctx context.Context, id int64) (*entity.Order, error) {
	var o entity.Order
	err := repo.db.GetContext(ctx, &o, "SELECT * FROM order WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %v", err)
	}

	var oi []entity.OrderItem
	err = repo.db.SelectContext(ctx, &oi, "SELECT * FROM order_item WHERE order_id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %v", err)
	}

	o.Items = oi

	return &o, nil
}

func (repo *OrderRepository) ListOrders(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	err := repo.db.SelectContext(ctx, &orders, "SELECT * FROM order")
	if err != nil {
		return nil, fmt.Errorf("error listing order: %v", err)
	}

	for i := range orders {
		var items []entity.OrderItem
		err := repo.db.SelectContext(ctx, &items, "SELECT * FROM order_item WHERE order_id=?", orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting order items: %v", err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (repo *OrderRepository) DeleteOrder(ctx context.Context, id int64) error {
	err := repo.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM order_item WHERE order_id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order items: %v", err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM order WHERE id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error deleting order: %v", err)
	}

	return nil
}

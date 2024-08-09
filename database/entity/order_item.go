package entity

import "time"

type OrderItem struct {
	ID        int64      `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Name      string     `json:"name" db:"name"`
	Quantity  int64      `json:"quantity" db:"quantity"`
	Image     string     `json:"image" db:"image"`
	Price     float64    `json:"price" db:"price"`
	ProductID int64      `json:"product_id" db:"product_id"`
	OrderID   int64      `json:"order_id" db:"order_id"`
}

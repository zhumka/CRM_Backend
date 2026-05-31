package model

import "time"

// Product — продукт (товар) на складе.
type Product struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	CategoryID  *int      `db:"category_id" json:"category_id"`
	SupplierID  *int      `db:"supplier_id" json:"supplier_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ProductInput — данные для создания/обновления продукта.
type ProductInput struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price" binding:"gte=0"`
	CategoryID  *int    `json:"category_id" binding:"omitempty"`
	SupplierID  *int    `json:"supplier_id" binding:"omitempty"`
}

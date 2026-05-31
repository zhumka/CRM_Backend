package model

import "time"

// Supplier — поставщик продукции.
type Supplier struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	ContactName string    `db:"contact_name" json:"contact_name"`
	Phone       string    `db:"phone" json:"phone"`
	Email       string    `db:"email" json:"email"`
	Address     string    `db:"address" json:"address"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// SupplierInput — данные для создания/обновления поставщика.
type SupplierInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	ContactName string `json:"contact_name" binding:"omitempty,max=100"`
	Phone       string `json:"phone" binding:"omitempty,max=20"`
	Email       string `json:"email" binding:"omitempty,email,max=100"`
	Address     string `json:"address" binding:"omitempty"`
}

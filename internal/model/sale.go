package model

import "time"

// Статусы установки по продаже.
const (
	InstallationNotRequired = "not_required"
	InstallationScheduled   = "scheduled"
	InstallationCompleted   = "completed"
)

// Sale — выполненная продажа (и при необходимости установка).
type Sale struct {
	ID                 int       `db:"id" json:"id"`
	InvoiceID          *int      `db:"invoice_id" json:"invoice_id"`
	ProductName        string    `db:"product_name" json:"product_name"`
	Quantity           int       `db:"quantity" json:"quantity"`
	Amount             float64   `db:"amount" json:"amount"`
	SoldAt             time.Time `db:"sold_at" json:"sold_at"`
	InstallationStatus string    `db:"installation_status" json:"installation_status"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// SaleInput — данные для создания/обновления продажи.
type SaleInput struct {
	InvoiceID          *int    `json:"invoice_id" binding:"omitempty"`
	ProductName        string  `json:"product_name" binding:"required,max=150"`
	Quantity           int     `json:"quantity" binding:"gte=0"`
	Amount             float64 `json:"amount" binding:"gte=0"`
	InstallationStatus string  `json:"installation_status" binding:"omitempty,oneof=not_required scheduled completed"`
}

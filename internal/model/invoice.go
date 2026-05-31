package model

import "time"

// Статусы оплаты счёта-фактуры.
const (
	InvoiceStatusUnpaid = "unpaid"
	InvoiceStatusPaid   = "paid"
	InvoiceStatusVoid   = "void"
)

// Invoice — счёт-фактура, привязанный к заявке/продаже.
type Invoice struct {
	ID                int       `db:"id" json:"id"`
	Number            string    `db:"number" json:"number"`
	PurchaseRequestID *int      `db:"purchase_request_id" json:"purchase_request_id"`
	Amount            float64   `db:"amount" json:"amount"`
	Status            string    `db:"status" json:"status"`
	IssuedDate        time.Time `db:"issued_date" json:"issued_date"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// InvoiceInput — данные для создания/обновления счёта-фактуры.
type InvoiceInput struct {
	Number            string  `json:"number" binding:"required,max=50"`
	PurchaseRequestID *int    `json:"purchase_request_id" binding:"omitempty"`
	Amount            float64 `json:"amount" binding:"gte=0"`
	Status            string  `json:"status" binding:"omitempty,oneof=unpaid paid void"`
}

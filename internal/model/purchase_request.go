package model

import "time"

// Статусы заявки на закупку.
const (
	PurchaseStatusNew        = "new"
	PurchaseStatusInProgress = "in_progress"
	PurchaseStatusApproved   = "approved"
	PurchaseStatusRejected   = "rejected"
	PurchaseStatusCompleted  = "completed"
)

// PurchaseRequest — заявка на закупку продукта (подсистема обработки заявок).
type PurchaseRequest struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	ProductID int       `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Status    string    `db:"status" json:"status"`
	Comment   string    `db:"comment" json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// PurchaseRequestInput — данные для создания заявки на закупку.
type PurchaseRequestInput struct {
	ProductID int    `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gte=1"`
	Comment   string `json:"comment" binding:"omitempty"`
}

// PurchaseStatusInput — смена статуса заявки.
type PurchaseStatusInput struct {
	Status string `json:"status" binding:"required,oneof=new in_progress approved rejected completed"`
}

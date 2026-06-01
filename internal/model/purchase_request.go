package model

import "time"

// Статусы заявки на закупку (выровнены с фронтендом).
const (
	PurchaseStatusDraft     = "draft"
	PurchaseStatusPending   = "pending"
	PurchaseStatusChecking  = "checking"
	PurchaseStatusApproved  = "approved"
	PurchaseStatusOrdered   = "ordered"
	PurchaseStatusCompleted = "completed"
	PurchaseStatusRejected  = "rejected"
)

// PurchaseRequest — заявка на закупку продукта (подсистема обработки заявок).
type PurchaseRequest struct {
	ID         int       `db:"id" json:"id"`
	Number     string    `db:"number" json:"number"`
	UserID     int       `db:"user_id" json:"user_id"`
	ClientName string    `db:"client_name" json:"client_name"`
	ProductID  int       `db:"product_id" json:"product_id"`
	Quantity   int       `db:"quantity" json:"quantity"`
	Status     string    `db:"status" json:"status"`
	Comment    string    `db:"comment" json:"comment"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// PurchaseRequestInput — данные для создания заявки на закупку.
type PurchaseRequestInput struct {
	ProductID  int    `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,gte=1"`
	ClientName string `json:"client_name" binding:"omitempty,max=150"`
	Comment    string `json:"comment" binding:"omitempty"`
}

// PurchaseStatusInput — смена статуса заявки.
type PurchaseStatusInput struct {
	Status string `json:"status" binding:"required,oneof=draft pending checking approved ordered completed rejected"`
}

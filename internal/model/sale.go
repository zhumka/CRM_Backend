package model

import "time"

// Статусы установки по продаже.
const (
	InstallationNotRequired = "not_required"
	InstallationScheduled   = "scheduled"
	InstallationCompleted   = "completed"
)

// Sale — выполненная продажа (и при необходимости установка).
// Amount — база без налога; TaxAmount и Total вычисляются по TaxRate.
type Sale struct {
	ID                 int       `db:"id" json:"id"`
	InvoiceID          *int      `db:"invoice_id" json:"invoice_id"`
	ProductName        string    `db:"product_name" json:"product_name"`
	Quantity           int       `db:"quantity" json:"quantity"`
	Amount             float64   `db:"amount" json:"amount"`
	TaxRate            float64   `db:"tax_rate" json:"tax_rate"`
	TaxAmount          float64   `db:"-" json:"tax_amount"`
	Total              float64   `db:"-" json:"total"`
	SoldAt             time.Time `db:"sold_at" json:"sold_at"`
	InstallationStatus string    `db:"installation_status" json:"installation_status"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// ApplyTax заполняет вычисляемые поля TaxAmount и Total из Amount и TaxRate.
func (s *Sale) ApplyTax() {
	s.TaxAmount = s.Amount * s.TaxRate / 100
	s.Total = s.Amount + s.TaxAmount
}

// SaleInput — данные для создания/обновления продажи.
// Ставка налога: TaxRate (явное значение) > TaxID (из справочника) > ставка продукта.
type SaleInput struct {
	InvoiceID          *int     `json:"invoice_id" binding:"omitempty"`
	ProductName        string   `json:"product_name" binding:"required,max=150"`
	Quantity           int      `json:"quantity" binding:"gte=0"`
	Amount             float64  `json:"amount" binding:"gte=0"`
	TaxRate            *float64 `json:"tax_rate" binding:"omitempty,gte=0"`
	TaxID              *int     `json:"tax_id" binding:"omitempty"`
	InstallationStatus string   `json:"installation_status" binding:"omitempty,oneof=not_required scheduled completed"`
}

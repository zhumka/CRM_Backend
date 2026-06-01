package model

import "time"

// Tax — налоговая ставка для цен, счетов и отчётности.
type Tax struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Rate      float64   `db:"rate" json:"rate"`
	Active    bool      `db:"active" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TaxInput — данные для создания/обновления налоговой ставки.
type TaxInput struct {
	Name   string  `json:"name" binding:"required,max=100"`
	Rate   float64 `json:"rate" binding:"gte=0"`
	Active *bool   `json:"active" binding:"omitempty"`
}

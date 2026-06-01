package model

import "time"

// Unit — единица измерения для складского учёта.
type Unit struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	ShortName   string    `db:"short_name" json:"short_name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// UnitInput — данные для создания/обновления единицы измерения.
type UnitInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	ShortName   string `json:"short_name" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty"`
}

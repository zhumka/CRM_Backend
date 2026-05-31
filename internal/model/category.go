package model

import "time"

// Category — категория продукции для классификации и удобства поиска.
type Category struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// CategoryInput — данные для создания/обновления категории.
type CategoryInput struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty"`
}

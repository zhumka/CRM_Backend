package model

import "time"

// Report — сохранённый отчёт (формирование, хранение, экспорт).
type Report struct {
	ID          int       `db:"id" json:"id"`
	CategoryID  *int      `db:"category_id" json:"category_id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	CreatedDate time.Time `db:"created_date" json:"created_date"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ReportInput — данные для ручного создания/обновления отчёта.
type ReportInput struct {
	Title      string `json:"title" binding:"required,max=150"`
	Content    string `json:"content" binding:"omitempty"`
	CategoryID *int   `json:"category_id" binding:"omitempty"`
}

// GenerateSalesReportInput — параметры генерации отчёта о продажах за период.
// Даты в формате YYYY-MM-DD; обе опциональны.
type GenerateSalesReportInput struct {
	From  string `json:"from" binding:"omitempty,datetime=2006-01-02"`
	To    string `json:"to" binding:"omitempty,datetime=2006-01-02"`
	Title string `json:"title" binding:"omitempty,max=150"`
}

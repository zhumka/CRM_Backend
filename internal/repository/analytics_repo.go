package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// AnalyticsRepository — агрегированные выборки для аналитики и отчётности.
type AnalyticsRepository struct {
	db *sqlx.DB
}

func NewAnalyticsRepository(db *sqlx.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// Summary возвращает сводные KPI по всем сущностям.
func (r *AnalyticsRepository) Summary(ctx context.Context) (*model.Summary, error) {
	const q = `
		SELECT
			(SELECT count(*) FROM users)                                  AS total_users,
			(SELECT count(*) FROM categories)                             AS total_categories,
			(SELECT count(*) FROM suppliers)                              AS total_suppliers,
			(SELECT count(*) FROM products)                               AS total_products,
			(SELECT count(*) FROM purchase_requests)                      AS total_purchase_requests,
			(SELECT count(*) FROM invoices)                               AS total_invoices,
			(SELECT coalesce(sum(amount), 0) FROM invoices WHERE status = 'paid') AS total_revenue`
	var s model.Summary
	if err := r.db.GetContext(ctx, &s, q); err != nil {
		return nil, err
	}
	return &s, nil
}

// SalesAnalytics возвращает финансовую аналитику по счетам за период (границы опциональны).
func (r *AnalyticsRepository) SalesAnalytics(ctx context.Context, from, to *time.Time) (*model.SalesAnalytics, error) {
	const q = `
		SELECT
			count(*)                                                          AS invoice_count,
			coalesce(sum(amount), 0)                                          AS total_amount,
			coalesce(sum(amount) FILTER (WHERE status = 'paid'), 0)           AS paid_amount,
			coalesce(sum(amount) FILTER (WHERE status = 'unpaid'), 0)         AS unpaid_amount
		FROM invoices
		WHERE ($1::date IS NULL OR issued_date >= $1)
		  AND ($2::date IS NULL OR issued_date <= $2)`
	var sa model.SalesAnalytics
	if err := r.db.GetContext(ctx, &sa, q, from, to); err != nil {
		return nil, err
	}
	sa.From, sa.To = from, to
	return &sa, nil
}

// RequestsByStatus возвращает количество заявок на закупку в разрезе статусов.
func (r *AnalyticsRepository) RequestsByStatus(ctx context.Context) ([]model.StatusCount, error) {
	const q = `SELECT status, count(*) AS count FROM purchase_requests GROUP BY status ORDER BY status`
	items := []model.StatusCount{}
	if err := r.db.SelectContext(ctx, &items, q); err != nil {
		return nil, err
	}
	return items, nil
}

// TopProducts возвращает самые востребованные продукты по заявкам на закупку.
func (r *AnalyticsRepository) TopProducts(ctx context.Context, limit int) ([]model.TopProduct, error) {
	const q = `
		SELECT p.id AS product_id, p.name,
		       coalesce(sum(pr.quantity), 0) AS total_quantity,
		       count(*)                       AS request_count
		FROM purchase_requests pr
		JOIN products p ON p.id = pr.product_id
		GROUP BY p.id, p.name
		ORDER BY total_quantity DESC, request_count DESC
		LIMIT $1`
	items := []model.TopProduct{}
	if err := r.db.SelectContext(ctx, &items, q, limit); err != nil {
		return nil, err
	}
	return items, nil
}

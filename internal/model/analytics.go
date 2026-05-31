package model

import "time"

// Summary — сводные KPI системы для дашборда.
type Summary struct {
	TotalUsers            int     `db:"total_users" json:"total_users"`
	TotalCategories       int     `db:"total_categories" json:"total_categories"`
	TotalSuppliers        int     `db:"total_suppliers" json:"total_suppliers"`
	TotalProducts         int     `db:"total_products" json:"total_products"`
	TotalPurchaseRequests int     `db:"total_purchase_requests" json:"total_purchase_requests"`
	TotalInvoices         int     `db:"total_invoices" json:"total_invoices"`
	TotalRevenue          float64 `db:"total_revenue" json:"total_revenue"`
}

// StatusCount — количество заявок в разрезе статуса.
type StatusCount struct {
	Status string `db:"status" json:"status"`
	Count  int    `db:"count" json:"count"`
}

// SalesAnalytics — финансовая аналитика по счетам-фактурам за период.
type SalesAnalytics struct {
	From         *time.Time `json:"from"`
	To           *time.Time `json:"to"`
	InvoiceCount int        `db:"invoice_count" json:"invoice_count"`
	TotalAmount  float64    `db:"total_amount" json:"total_amount"`
	PaidAmount   float64    `db:"paid_amount" json:"paid_amount"`
	UnpaidAmount float64    `db:"unpaid_amount" json:"unpaid_amount"`
}

// TopProduct — наиболее востребованный продукт по заявкам на закупку.
type TopProduct struct {
	ProductID     int    `db:"product_id" json:"product_id"`
	Name          string `db:"name" json:"name"`
	TotalQuantity int    `db:"total_quantity" json:"total_quantity"`
	RequestCount  int    `db:"request_count" json:"request_count"`
}

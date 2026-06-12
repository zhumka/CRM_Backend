package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// SaleRepository — доступ к данным продаж.
type SaleRepository struct {
	db *sqlx.DB
}

func NewSaleRepository(db *sqlx.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) Create(ctx context.Context, s *model.Sale) error {
	const q = `
		INSERT INTO sales (invoice_id, product_name, quantity, amount, tax_rate, installation_status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, sold_at, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, q,
		s.InvoiceID, s.ProductName, s.Quantity, s.Amount, s.TaxRate, s.InstallationStatus).
		Scan(&s.ID, &s.SoldAt, &s.CreatedAt, &s.UpdatedAt)
	return mapFKError(err)
}

func (r *SaleRepository) List(ctx context.Context) ([]model.Sale, error) {
	items := []model.Sale{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM sales ORDER BY id DESC`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SaleRepository) GetByID(ctx context.Context, id int) (*model.Sale, error) {
	var s model.Sale
	if err := r.db.GetContext(ctx, &s, `SELECT * FROM sales WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *SaleRepository) Update(ctx context.Context, id int, in model.SaleInput, status string, taxRate float64) (*model.Sale, error) {
	const q = `
		UPDATE sales
		SET invoice_id = $2, product_name = $3, quantity = $4, amount = $5,
		    tax_rate = $6, installation_status = $7, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var s model.Sale
	err := r.db.GetContext(ctx, &s, q, id, in.InvoiceID, in.ProductName, in.Quantity, in.Amount, taxRate, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, mapFKError(err)
	}
	return &s, nil
}

func (r *SaleRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM sales WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

// ListByInvoice возвращает продажи, привязанные к счёту (позиции для документа).
func (r *SaleRepository) ListByInvoice(ctx context.Context, invoiceID int) ([]model.Sale, error) {
	items := []model.Sale{}
	const q = `SELECT * FROM sales WHERE invoice_id = $1 ORDER BY id`
	if err := r.db.SelectContext(ctx, &items, q, invoiceID); err != nil {
		return nil, err
	}
	return items, nil
}

// SumAmountByInvoice возвращает сумму продаж, привязанных к счёту,
// исключая продажу excludeID (0 — ничего не исключать, например при создании).
func (r *SaleRepository) SumAmountByInvoice(ctx context.Context, invoiceID, excludeID int) (float64, error) {
	var total float64
	// Сумма продаж по счёту считается с налогом (amount + налог).
	const q = `SELECT COALESCE(SUM(amount + amount * tax_rate / 100), 0) FROM sales WHERE invoice_id = $1 AND id <> $2`
	if err := r.db.GetContext(ctx, &total, q, invoiceID, excludeID); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *SaleRepository) Count(ctx context.Context) (int, error) {
	var n int
	if err := r.db.GetContext(ctx, &n, `SELECT count(*) FROM sales`); err != nil {
		return 0, err
	}
	return n, nil
}

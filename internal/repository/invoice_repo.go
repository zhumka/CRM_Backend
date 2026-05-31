package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// InvoiceRepository — доступ к данным счетов-фактур.
type InvoiceRepository struct {
	db *sqlx.DB
}

func NewInvoiceRepository(db *sqlx.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(ctx context.Context, inv *model.Invoice) error {
	const q = `
		INSERT INTO invoices (number, purchase_request_id, amount, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, issued_date, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, q, inv.Number, inv.PurchaseRequestID, inv.Amount, inv.Status).
		Scan(&inv.ID, &inv.IssuedDate, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return model.ErrAlreadyExists
		}
		return mapFKError(err)
	}
	return nil
}

func (r *InvoiceRepository) List(ctx context.Context) ([]model.Invoice, error) {
	items := []model.Invoice{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM invoices ORDER BY id DESC`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id int) (*model.Invoice, error) {
	var inv model.Invoice
	if err := r.db.GetContext(ctx, &inv, `SELECT * FROM invoices WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func (r *InvoiceRepository) Update(ctx context.Context, id int, in model.InvoiceInput) (*model.Invoice, error) {
	const q = `
		UPDATE invoices
		SET number = $2, purchase_request_id = $3, amount = $4,
		    status = COALESCE(NULLIF($5, ''), status), updated_at = now()
		WHERE id = $1
		RETURNING *`
	var inv model.Invoice
	err := r.db.GetContext(ctx, &inv, q, id, in.Number, in.PurchaseRequestID, in.Amount, in.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		if isUniqueViolation(err) {
			return nil, model.ErrAlreadyExists
		}
		return nil, mapFKError(err)
	}
	return &inv, nil
}

func (r *InvoiceRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM invoices WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

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
		INSERT INTO sales (invoice_id, product_name, quantity, amount, installation_status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sold_at, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, q,
		s.InvoiceID, s.ProductName, s.Quantity, s.Amount, s.InstallationStatus).
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

func (r *SaleRepository) Update(ctx context.Context, id int, in model.SaleInput, status string) (*model.Sale, error) {
	const q = `
		UPDATE sales
		SET invoice_id = $2, product_name = $3, quantity = $4, amount = $5,
		    installation_status = $6, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var s model.Sale
	err := r.db.GetContext(ctx, &s, q, id, in.InvoiceID, in.ProductName, in.Quantity, in.Amount, status)
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

func (r *SaleRepository) Count(ctx context.Context) (int, error) {
	var n int
	if err := r.db.GetContext(ctx, &n, `SELECT count(*) FROM sales`); err != nil {
		return 0, err
	}
	return n, nil
}

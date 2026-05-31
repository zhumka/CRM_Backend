package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// SupplierRepository — доступ к данным поставщиков.
type SupplierRepository struct {
	db *sqlx.DB
}

func NewSupplierRepository(db *sqlx.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) Create(ctx context.Context, s *model.Supplier) error {
	const q = `
		INSERT INTO suppliers (name, contact_name, phone, email, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, q, s.Name, s.ContactName, s.Phone, s.Email, s.Address).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SupplierRepository) List(ctx context.Context) ([]model.Supplier, error) {
	items := []model.Supplier{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM suppliers ORDER BY id`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SupplierRepository) GetByID(ctx context.Context, id int) (*model.Supplier, error) {
	var s model.Supplier
	if err := r.db.GetContext(ctx, &s, `SELECT * FROM suppliers WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepository) Update(ctx context.Context, id int, in model.SupplierInput) (*model.Supplier, error) {
	const q = `
		UPDATE suppliers
		SET name = $2, contact_name = $3, phone = $4, email = $5, address = $6, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var s model.Supplier
	err := r.db.GetContext(ctx, &s, q, id, in.Name, in.ContactName, in.Phone, in.Email, in.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM suppliers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

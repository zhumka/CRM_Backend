package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// TaxRepository — доступ к данным налоговых ставок.
type TaxRepository struct {
	db *sqlx.DB
}

func NewTaxRepository(db *sqlx.DB) *TaxRepository {
	return &TaxRepository{db: db}
}

func (r *TaxRepository) Create(ctx context.Context, t *model.Tax) error {
	const q = `
		INSERT INTO taxes (name, rate, active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, q, t.Name, t.Rate, t.Active).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TaxRepository) List(ctx context.Context) ([]model.Tax, error) {
	items := []model.Tax{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM taxes ORDER BY id`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TaxRepository) GetByID(ctx context.Context, id int) (*model.Tax, error) {
	var t model.Tax
	if err := r.db.GetContext(ctx, &t, `SELECT * FROM taxes WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TaxRepository) Update(ctx context.Context, id int, name string, rate float64, active bool) (*model.Tax, error) {
	const q = `
		UPDATE taxes
		SET name = $2, rate = $3, active = $4, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var t model.Tax
	if err := r.db.GetContext(ctx, &t, q, id, name, rate, active); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TaxRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM taxes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

func (r *TaxRepository) Count(ctx context.Context) (int, error) {
	var n int
	if err := r.db.GetContext(ctx, &n, `SELECT count(*) FROM taxes`); err != nil {
		return 0, err
	}
	return n, nil
}

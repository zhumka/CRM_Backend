package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// UnitRepository — доступ к данным единиц измерения.
type UnitRepository struct {
	db *sqlx.DB
}

func NewUnitRepository(db *sqlx.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) Create(ctx context.Context, u *model.Unit) error {
	const q = `
		INSERT INTO units (name, short_name, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, q, u.Name, u.ShortName, u.Description).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UnitRepository) List(ctx context.Context) ([]model.Unit, error) {
	items := []model.Unit{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM units ORDER BY id`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *UnitRepository) GetByID(ctx context.Context, id int) (*model.Unit, error) {
	var u model.Unit
	if err := r.db.GetContext(ctx, &u, `SELECT * FROM units WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UnitRepository) Update(ctx context.Context, id int, in model.UnitInput) (*model.Unit, error) {
	const q = `
		UPDATE units
		SET name = $2, short_name = $3, description = $4, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var u model.Unit
	if err := r.db.GetContext(ctx, &u, q, id, in.Name, in.ShortName, in.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UnitRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM units WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

func (r *UnitRepository) Count(ctx context.Context) (int, error) {
	var n int
	if err := r.db.GetContext(ctx, &n, `SELECT count(*) FROM units`); err != nil {
		return 0, err
	}
	return n, nil
}

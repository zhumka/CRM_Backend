package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// CategoryRepository — доступ к данным категорий.
type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) error {
	const q = `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, q, c.Name, c.Description).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	items := []model.Category{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM categories ORDER BY id`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	var c model.Category
	if err := r.db.GetContext(ctx, &c, `SELECT * FROM categories WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int, in model.CategoryInput) (*model.Category, error) {
	const q = `
		UPDATE categories
		SET name = $2, description = $3, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var c model.Category
	if err := r.db.GetContext(ctx, &c, q, id, in.Name, in.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

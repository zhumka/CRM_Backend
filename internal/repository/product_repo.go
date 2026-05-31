package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// ProductRepository — доступ к данным продуктов.
type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	const q = `
		INSERT INTO products (name, description, price, category_id, supplier_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, q, p.Name, p.Description, p.Price, p.CategoryID, p.SupplierID).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return mapFKError(err)
}

// List возвращает продукты с опциональной фильтрацией по категории.
func (r *ProductRepository) List(ctx context.Context, categoryID *int) ([]model.Product, error) {
	items := []model.Product{}
	if categoryID != nil {
		const q = `SELECT * FROM products WHERE category_id = $1 ORDER BY id`
		if err := r.db.SelectContext(ctx, &items, q, *categoryID); err != nil {
			return nil, err
		}
		return items, nil
	}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM products ORDER BY id`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
	var p model.Product
	if err := r.db.GetContext(ctx, &p, `SELECT * FROM products WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, in model.ProductInput) (*model.Product, error) {
	const q = `
		UPDATE products
		SET name = $2, description = $3, price = $4, category_id = $5, supplier_id = $6, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var p model.Product
	err := r.db.GetContext(ctx, &p, q, id, in.Name, in.Description, in.Price, in.CategoryID, in.SupplierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, mapFKError(err)
	}
	return &p, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

// mapFKError превращает нарушение внешнего ключа в ErrNotFound (несуществующая категория/поставщик).
func mapFKError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "23503") ||
		strings.Contains(strings.ToLower(err.Error()), "foreign key") {
		return model.ErrNotFound
	}
	return err
}

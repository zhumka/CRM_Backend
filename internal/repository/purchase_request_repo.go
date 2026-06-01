package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// PurchaseRequestRepository — доступ к данным заявок на закупку.
type PurchaseRequestRepository struct {
	db *sqlx.DB
}

func NewPurchaseRequestRepository(db *sqlx.DB) *PurchaseRequestRepository {
	return &PurchaseRequestRepository{db: db}
}

func (r *PurchaseRequestRepository) Create(ctx context.Context, pr *model.PurchaseRequest) error {
	// Номер заявки формируется из её id (REQ-0001): вставляем, затем проставляем
	// номер вторым запросом в одной транзакции.
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	const insertQ = `
		INSERT INTO purchase_requests (user_id, client_name, product_id, quantity, status, comment)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	if err := tx.QueryRowxContext(ctx, insertQ,
		pr.UserID, pr.ClientName, pr.ProductID, pr.Quantity, pr.Status, pr.Comment).
		Scan(&pr.ID, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
		return mapFKError(err)
	}

	pr.Number = fmt.Sprintf("REQ-%04d", pr.ID)
	if _, err := tx.ExecContext(ctx, `UPDATE purchase_requests SET number = $2 WHERE id = $1`, pr.ID, pr.Number); err != nil {
		return err
	}
	return tx.Commit()
}

// List возвращает заявки. Если ownerID != nil — только заявки этого пользователя.
func (r *PurchaseRequestRepository) List(ctx context.Context, ownerID *int) ([]model.PurchaseRequest, error) {
	items := []model.PurchaseRequest{}
	if ownerID != nil {
		const q = `SELECT * FROM purchase_requests WHERE user_id = $1 ORDER BY id DESC`
		if err := r.db.SelectContext(ctx, &items, q, *ownerID); err != nil {
			return nil, err
		}
		return items, nil
	}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM purchase_requests ORDER BY id DESC`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PurchaseRequestRepository) GetByID(ctx context.Context, id int) (*model.PurchaseRequest, error) {
	var pr model.PurchaseRequest
	if err := r.db.GetContext(ctx, &pr, `SELECT * FROM purchase_requests WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &pr, nil
}

func (r *PurchaseRequestRepository) UpdateStatus(ctx context.Context, id int, status string) (*model.PurchaseRequest, error) {
	const q = `
		UPDATE purchase_requests
		SET status = $2, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var pr model.PurchaseRequest
	if err := r.db.GetContext(ctx, &pr, q, id, status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &pr, nil
}

func (r *PurchaseRequestRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM purchase_requests WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

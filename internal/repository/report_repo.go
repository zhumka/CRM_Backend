package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// ReportRepository — доступ к данным отчётов.
type ReportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(ctx context.Context, rep *model.Report) error {
	const q = `
		INSERT INTO reports (category_id, title, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_date, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, q, rep.CategoryID, rep.Title, rep.Content).
		Scan(&rep.ID, &rep.CreatedDate, &rep.CreatedAt, &rep.UpdatedAt)
	return mapFKError(err)
}

func (r *ReportRepository) List(ctx context.Context) ([]model.Report, error) {
	items := []model.Report{}
	if err := r.db.SelectContext(ctx, &items, `SELECT * FROM reports ORDER BY id DESC`); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ReportRepository) GetByID(ctx context.Context, id int) (*model.Report, error) {
	var rep model.Report
	if err := r.db.GetContext(ctx, &rep, `SELECT * FROM reports WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &rep, nil
}

func (r *ReportRepository) Update(ctx context.Context, id int, in model.ReportInput) (*model.Report, error) {
	const q = `
		UPDATE reports
		SET title = $2, content = $3, category_id = $4, updated_at = now()
		WHERE id = $1
		RETURNING *`
	var rep model.Report
	err := r.db.GetContext(ctx, &rep, q, id, in.Title, in.Content, in.CategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, mapFKError(err)
	}
	return &rep, nil
}

func (r *ReportRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM reports WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

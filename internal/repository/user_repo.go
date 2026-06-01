package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"crm/internal/model"
)

// UserRepository — доступ к данным пользователей.
type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create сохраняет нового пользователя и возвращает его id, created_at, updated_at.
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	if u.Status == "" {
		u.Status = model.UserStatusActive
	}
	const q = `
		INSERT INTO users (username, password_hash, full_name, email, role, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRowxContext(ctx, q, u.Username, u.PasswordHash, u.FullName, u.Email, u.Role, u.Status).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return model.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	const q = `SELECT * FROM users WHERE username = $1`
	if err := r.db.GetContext(ctx, &u, q, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	var u model.User
	const q = `SELECT * FROM users WHERE id = $1`
	if err := r.db.GetContext(ctx, &u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	users := []model.User{}
	const q = `SELECT * FROM users ORDER BY id`
	if err := r.db.SelectContext(ctx, &users, q); err != nil {
		return nil, err
	}
	return users, nil
}

// Update обновляет переданные поля пользователя; nil-поля не изменяются.
func (r *UserRepository) Update(ctx context.Context, id int, f model.UserUpdate) (*model.User, error) {
	const q = `
		UPDATE users
		SET password_hash = COALESCE($2, password_hash),
		    full_name     = COALESCE($3, full_name),
		    email         = COALESCE($4, email),
		    role          = COALESCE($5, role),
		    status        = COALESCE($6, status),
		    updated_at    = now()
		WHERE id = $1
		RETURNING *`
	var u model.User
	if err := r.db.GetContext(ctx, &u, q, id, f.PasswordHash, f.FullName, f.Email, f.Role, f.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return ensureAffected(res)
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var n int
	if err := r.db.GetContext(ctx, &n, `SELECT count(*) FROM users`); err != nil {
		return 0, err
	}
	return n, nil
}

// isUniqueViolation определяет нарушение уникального ограничения PostgreSQL (код 23505).
func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "23505") ||
		strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}

// ensureAffected превращает «0 строк затронуто» в ErrNotFound.
func ensureAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

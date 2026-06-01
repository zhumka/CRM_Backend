package model

import "time"

// Роли пользователей системы.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Статусы учётной записи.
const (
	UserStatusActive  = "active"
	UserStatusBlocked = "blocked"
)

// User — пользователь системы (актёр «Пользователь»/«Администратор» из ТЗ).
type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FullName     string    `db:"full_name" json:"full_name"`
	Email        string    `db:"email" json:"email"`
	Role         string    `db:"role" json:"role"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// RegisterInput — данные для регистрации/создания пользователя.
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	FullName string `json:"full_name" binding:"omitempty,max=150"`
	Email    string `json:"email" binding:"omitempty,email,max=150"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

// LoginInput — данные для авторизации.
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserInput — обновление пользователя администратором.
type UpdateUserInput struct {
	Password *string `json:"password" binding:"omitempty,min=6,max=72"`
	FullName *string `json:"full_name" binding:"omitempty,max=150"`
	Email    *string `json:"email" binding:"omitempty,email,max=150"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin user"`
	Status   *string `json:"status" binding:"omitempty,oneof=active blocked"`
}

// UserUpdate — изменяемые поля пользователя на уровне хранилища; nil оставляет прежнее.
type UserUpdate struct {
	PasswordHash *string
	FullName     *string
	Email        *string
	Role         *string
	Status       *string
}

// AuthResponse — ответ при успешной авторизации/регистрации.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

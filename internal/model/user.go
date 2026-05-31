package model

import "time"

// Роли пользователей системы.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User — пользователь системы (актёр «Пользователь»/«Администратор» из ТЗ).
type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// RegisterInput — данные для регистрации/создания пользователя.
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
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
	Role     *string `json:"role" binding:"omitempty,oneof=admin user"`
}

// AuthResponse — ответ при успешной авторизации/регистрации.
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

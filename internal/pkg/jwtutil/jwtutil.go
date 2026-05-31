// Package jwtutil отвечает за выпуск и разбор JWT access-токенов.
package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken возвращается при невалидном или просроченном токене.
var ErrInvalidToken = errors.New("invalid token")

// Claims — полезная нагрузка токена.
type Claims struct {
	UserID int    `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Manager выпускает и валидирует токены с заданным секретом и временем жизни.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager создаёт менеджер токенов.
func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// Generate выпускает подписанный токен для пользователя.
func (m *Manager) Generate(userID int, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse проверяет подпись и срок действия, возвращая claims.
func (m *Manager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

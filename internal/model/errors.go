package model

import "errors"

// Доменные ошибки, общие для слоёв сервиса и репозитория.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrForbidden          = errors.New("access forbidden")
)

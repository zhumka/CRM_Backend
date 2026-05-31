package service

import (
	"context"

	"crm/internal/model"
	"crm/internal/pkg/hash"
	"crm/internal/pkg/jwtutil"
)

// UserStore — зависимость сервиса аутентификации от хранилища пользователей.
type UserStore interface {
	Create(ctx context.Context, u *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, id int, passwordHash, role *string) (*model.User, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

// AuthService реализует регистрацию и вход с выпуском JWT.
type AuthService struct {
	users UserStore
	jwt   *jwtutil.Manager
}

func NewAuthService(users UserStore, jwt *jwtutil.Manager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

// Register создаёт пользователя с хешированным паролем и выдаёт токен.
func (s *AuthService) Register(ctx context.Context, in model.RegisterInput) (*model.AuthResponse, error) {
	role := in.Role
	if role == "" {
		role = model.RoleUser
	}

	pwdHash, err := hash.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	u := &model.User{Username: in.Username, PasswordHash: pwdHash, Role: role}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}

	return s.tokenResponse(u)
}

// Login проверяет учётные данные и выдаёт токен.
func (s *AuthService) Login(ctx context.Context, in model.LoginInput) (*model.AuthResponse, error) {
	u, err := s.users.GetByUsername(ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.ErrInvalidCredentials
		}
		return nil, err
	}

	if !hash.Check(in.Password, u.PasswordHash) {
		return nil, model.ErrInvalidCredentials
	}

	return s.tokenResponse(u)
}

func (s *AuthService) tokenResponse(u *model.User) (*model.AuthResponse, error) {
	token, err := s.jwt.Generate(u.ID, u.Role)
	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{Token: token, User: *u}, nil
}

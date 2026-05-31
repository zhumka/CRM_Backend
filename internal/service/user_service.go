package service

import (
	"context"

	"crm/internal/model"
	"crm/internal/pkg/hash"
)

// UserService — управление пользователями (только администратор).
type UserService struct {
	users UserStore
}

func NewUserService(users UserStore) *UserService {
	return &UserService{users: users}
}

func (s *UserService) List(ctx context.Context) ([]model.User, error) {
	return s.users.List(ctx)
}

func (s *UserService) Get(ctx context.Context, id int) (*model.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, in model.RegisterInput) (*model.User, error) {
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
	return u, nil
}

// Update меняет пароль и/или роль пользователя.
func (s *UserService) Update(ctx context.Context, id int, in model.UpdateUserInput) (*model.User, error) {
	var pwdHash *string
	if in.Password != nil {
		h, err := hash.Hash(*in.Password)
		if err != nil {
			return nil, err
		}
		pwdHash = &h
	}
	return s.users.Update(ctx, id, pwdHash, in.Role)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.users.Delete(ctx, id)
}

// EnsureAdmin создаёт администратора по умолчанию, если в системе нет пользователей.
func (s *UserService) EnsureAdmin(ctx context.Context, username, password string) error {
	n, err := s.users.Count(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	pwdHash, err := hash.Hash(password)
	if err != nil {
		return err
	}
	return s.users.Create(ctx, &model.User{
		Username:     username,
		PasswordHash: pwdHash,
		Role:         model.RoleAdmin,
	})
}

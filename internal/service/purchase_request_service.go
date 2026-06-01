package service

import (
	"context"

	"crm/internal/model"
)

// PurchaseRequestStore — зависимость от хранилища заявок.
type PurchaseRequestStore interface {
	Create(ctx context.Context, pr *model.PurchaseRequest) error
	List(ctx context.Context, ownerID *int) ([]model.PurchaseRequest, error)
	GetByID(ctx context.Context, id int) (*model.PurchaseRequest, error)
	UpdateStatus(ctx context.Context, id int, status string) (*model.PurchaseRequest, error)
	Delete(ctx context.Context, id int) error
}

// PurchaseRequestService реализует подсистему обработки заявок на закупку.
type PurchaseRequestService struct {
	repo     PurchaseRequestStore
	products ProductStore
}

func NewPurchaseRequestService(repo PurchaseRequestStore, products ProductStore) *PurchaseRequestService {
	return &PurchaseRequestService{repo: repo, products: products}
}

// Create создаёт заявку от имени пользователя; проверяет существование продукта.
func (s *PurchaseRequestService) Create(ctx context.Context, userID int, in model.PurchaseRequestInput) (*model.PurchaseRequest, error) {
	if _, err := s.products.GetByID(ctx, in.ProductID); err != nil {
		return nil, err // ErrNotFound, если продукт недоступен
	}
	pr := &model.PurchaseRequest{
		UserID:     userID,
		ClientName: in.ClientName,
		ProductID:  in.ProductID,
		Quantity:   in.Quantity,
		Status:     model.PurchaseStatusPending,
		Comment:    in.Comment,
	}
	if err := s.repo.Create(ctx, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

// List возвращает все заявки для админа либо только собственные для пользователя.
func (s *PurchaseRequestService) List(ctx context.Context, requesterID int, isAdmin bool) ([]model.PurchaseRequest, error) {
	if isAdmin {
		return s.repo.List(ctx, nil)
	}
	return s.repo.List(ctx, &requesterID)
}

// Get возвращает заявку с проверкой прав доступа.
func (s *PurchaseRequestService) Get(ctx context.Context, id, requesterID int, isAdmin bool) (*model.PurchaseRequest, error) {
	pr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isAdmin && pr.UserID != requesterID {
		return nil, model.ErrForbidden
	}
	return pr, nil
}

// UpdateStatus меняет статус заявки (только администратор — проверяется в обработчике).
func (s *PurchaseRequestService) UpdateStatus(ctx context.Context, id int, status string) (*model.PurchaseRequest, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

// Delete удаляет заявку с проверкой прав доступа.
func (s *PurchaseRequestService) Delete(ctx context.Context, id, requesterID int, isAdmin bool) error {
	pr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !isAdmin && pr.UserID != requesterID {
		return model.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}

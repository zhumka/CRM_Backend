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

// SaleCreator — создание продажи (реализуется SaleService).
type SaleCreator interface {
	Create(ctx context.Context, in model.SaleInput) (*model.Sale, error)
}

// PurchaseRequestService реализует подсистему обработки заявок на закупку.
type PurchaseRequestService struct {
	repo     PurchaseRequestStore
	products ProductStore
	sales    SaleCreator
}

func NewPurchaseRequestService(repo PurchaseRequestStore, products ProductStore, sales SaleCreator) *PurchaseRequestService {
	return &PurchaseRequestService{repo: repo, products: products, sales: sales}
}

// CreateSale оформляет продажу по одобренной заявке и переводит заявку в completed.
// Данные продажи берутся из заявки и продукта; ставка налога — из продукта.
func (s *PurchaseRequestService) CreateSale(ctx context.Context, id int) (*model.Sale, error) {
	pr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Продажу можно оформить только по одобренной заявке.
	if pr.Status != model.PurchaseStatusApproved {
		return nil, model.ErrRequestNotApproved
	}
	prod, err := s.products.GetByID(ctx, pr.ProductID)
	if err != nil {
		return nil, err
	}
	// Сумма — база без налога: цена × количество. Налог и склад обрабатывает SaleService.
	in := model.SaleInput{
		ProductName: prod.Name,
		Quantity:    pr.Quantity,
		Amount:      prod.Price * float64(pr.Quantity),
	}
	sale, err := s.sales.Create(ctx, in)
	if err != nil {
		return nil, err // в т.ч. ErrInsufficientStock — заявка остаётся approved
	}
	// Автопереход статуса: заявка выполнена.
	if _, err := s.repo.UpdateStatus(ctx, id, model.PurchaseStatusCompleted); err != nil {
		return nil, err
	}
	return sale, nil
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

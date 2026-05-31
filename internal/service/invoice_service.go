package service

import (
	"context"

	"crm/internal/model"
)

// InvoiceStore — зависимость от хранилища счетов-фактур.
type InvoiceStore interface {
	Create(ctx context.Context, inv *model.Invoice) error
	List(ctx context.Context) ([]model.Invoice, error)
	GetByID(ctx context.Context, id int) (*model.Invoice, error)
	Update(ctx context.Context, id int, in model.InvoiceInput) (*model.Invoice, error)
	Delete(ctx context.Context, id int) error
}

// InvoiceService — управление счетами-фактурами.
type InvoiceService struct{ repo InvoiceStore }

func NewInvoiceService(repo InvoiceStore) *InvoiceService { return &InvoiceService{repo: repo} }

func (s *InvoiceService) Create(ctx context.Context, in model.InvoiceInput) (*model.Invoice, error) {
	status := in.Status
	if status == "" {
		status = model.InvoiceStatusUnpaid
	}
	inv := &model.Invoice{
		Number:            in.Number,
		PurchaseRequestID: in.PurchaseRequestID,
		Amount:            in.Amount,
		Status:            status,
	}
	if err := s.repo.Create(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}
func (s *InvoiceService) List(ctx context.Context) ([]model.Invoice, error) { return s.repo.List(ctx) }
func (s *InvoiceService) Get(ctx context.Context, id int) (*model.Invoice, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *InvoiceService) Update(ctx context.Context, id int, in model.InvoiceInput) (*model.Invoice, error) {
	return s.repo.Update(ctx, id, in)
}
func (s *InvoiceService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

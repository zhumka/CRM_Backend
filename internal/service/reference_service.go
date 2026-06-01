package service

import (
	"context"

	"crm/internal/model"
)

// --- Налоги ---

type TaxStore interface {
	Create(ctx context.Context, t *model.Tax) error
	List(ctx context.Context) ([]model.Tax, error)
	GetByID(ctx context.Context, id int) (*model.Tax, error)
	Update(ctx context.Context, id int, name string, rate float64, active bool) (*model.Tax, error)
	Delete(ctx context.Context, id int) error
}

type TaxService struct{ repo TaxStore }

func NewTaxService(repo TaxStore) *TaxService { return &TaxService{repo: repo} }

func (s *TaxService) Create(ctx context.Context, in model.TaxInput) (*model.Tax, error) {
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	t := &model.Tax{Name: in.Name, Rate: in.Rate, Active: active}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}
func (s *TaxService) List(ctx context.Context) ([]model.Tax, error) { return s.repo.List(ctx) }
func (s *TaxService) Get(ctx context.Context, id int) (*model.Tax, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *TaxService) Update(ctx context.Context, id int, in model.TaxInput) (*model.Tax, error) {
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	return s.repo.Update(ctx, id, in.Name, in.Rate, active)
}
func (s *TaxService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

// --- Единицы измерения ---

type UnitStore interface {
	Create(ctx context.Context, u *model.Unit) error
	List(ctx context.Context) ([]model.Unit, error)
	GetByID(ctx context.Context, id int) (*model.Unit, error)
	Update(ctx context.Context, id int, in model.UnitInput) (*model.Unit, error)
	Delete(ctx context.Context, id int) error
}

type UnitService struct{ repo UnitStore }

func NewUnitService(repo UnitStore) *UnitService { return &UnitService{repo: repo} }

func (s *UnitService) Create(ctx context.Context, in model.UnitInput) (*model.Unit, error) {
	u := &model.Unit{Name: in.Name, ShortName: in.ShortName, Description: in.Description}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
func (s *UnitService) List(ctx context.Context) ([]model.Unit, error) { return s.repo.List(ctx) }
func (s *UnitService) Get(ctx context.Context, id int) (*model.Unit, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *UnitService) Update(ctx context.Context, id int, in model.UnitInput) (*model.Unit, error) {
	return s.repo.Update(ctx, id, in)
}
func (s *UnitService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

// --- Продажи ---

type SaleStore interface {
	Create(ctx context.Context, sale *model.Sale) error
	List(ctx context.Context) ([]model.Sale, error)
	GetByID(ctx context.Context, id int) (*model.Sale, error)
	Update(ctx context.Context, id int, in model.SaleInput, status string) (*model.Sale, error)
	Delete(ctx context.Context, id int) error
}

type SaleService struct{ repo SaleStore }

func NewSaleService(repo SaleStore) *SaleService { return &SaleService{repo: repo} }

func (s *SaleService) Create(ctx context.Context, in model.SaleInput) (*model.Sale, error) {
	status := in.InstallationStatus
	if status == "" {
		status = model.InstallationNotRequired
	}
	sale := &model.Sale{
		InvoiceID:          in.InvoiceID,
		ProductName:        in.ProductName,
		Quantity:           in.Quantity,
		Amount:             in.Amount,
		InstallationStatus: status,
	}
	if err := s.repo.Create(ctx, sale); err != nil {
		return nil, err
	}
	return sale, nil
}
func (s *SaleService) List(ctx context.Context) ([]model.Sale, error) { return s.repo.List(ctx) }
func (s *SaleService) Get(ctx context.Context, id int) (*model.Sale, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *SaleService) Update(ctx context.Context, id int, in model.SaleInput) (*model.Sale, error) {
	status := in.InstallationStatus
	if status == "" {
		status = model.InstallationNotRequired
	}
	return s.repo.Update(ctx, id, in, status)
}
func (s *SaleService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

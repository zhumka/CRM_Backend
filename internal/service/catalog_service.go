package service

import (
	"context"

	"crm/internal/model"
)

// --- Категории ---

type CategoryStore interface {
	Create(ctx context.Context, c *model.Category) error
	List(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, id int) (*model.Category, error)
	Update(ctx context.Context, id int, in model.CategoryInput) (*model.Category, error)
	Delete(ctx context.Context, id int) error
}

type CategoryService struct{ repo CategoryStore }

func NewCategoryService(repo CategoryStore) *CategoryService { return &CategoryService{repo: repo} }

func (s *CategoryService) Create(ctx context.Context, in model.CategoryInput) (*model.Category, error) {
	c := &model.Category{Name: in.Name, Description: in.Description}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
func (s *CategoryService) List(ctx context.Context) ([]model.Category, error) {
	return s.repo.List(ctx)
}
func (s *CategoryService) Get(ctx context.Context, id int) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *CategoryService) Update(ctx context.Context, id int, in model.CategoryInput) (*model.Category, error) {
	return s.repo.Update(ctx, id, in)
}
func (s *CategoryService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

// --- Поставщики ---

type SupplierStore interface {
	Create(ctx context.Context, sup *model.Supplier) error
	List(ctx context.Context) ([]model.Supplier, error)
	GetByID(ctx context.Context, id int) (*model.Supplier, error)
	Update(ctx context.Context, id int, in model.SupplierInput) (*model.Supplier, error)
	Delete(ctx context.Context, id int) error
}

type SupplierService struct{ repo SupplierStore }

func NewSupplierService(repo SupplierStore) *SupplierService { return &SupplierService{repo: repo} }

func (s *SupplierService) Create(ctx context.Context, in model.SupplierInput) (*model.Supplier, error) {
	sup := &model.Supplier{
		Name:        in.Name,
		ContactName: in.ContactName,
		Phone:       in.Phone,
		Email:       in.Email,
		Address:     in.Address,
	}
	if err := s.repo.Create(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}
func (s *SupplierService) List(ctx context.Context) ([]model.Supplier, error) {
	return s.repo.List(ctx)
}
func (s *SupplierService) Get(ctx context.Context, id int) (*model.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *SupplierService) Update(ctx context.Context, id int, in model.SupplierInput) (*model.Supplier, error) {
	return s.repo.Update(ctx, id, in)
}
func (s *SupplierService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

// --- Продукты ---

type ProductStore interface {
	Create(ctx context.Context, p *model.Product) error
	List(ctx context.Context, categoryID *int) ([]model.Product, error)
	GetByID(ctx context.Context, id int) (*model.Product, error)
	Update(ctx context.Context, id int, in model.ProductInput) (*model.Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductService struct{ repo ProductStore }

func NewProductService(repo ProductStore) *ProductService { return &ProductService{repo: repo} }

func (s *ProductService) Create(ctx context.Context, in model.ProductInput) (*model.Product, error) {
	p := &model.Product{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		Unit:        in.Unit,
		TaxRate:     in.TaxRate,
		CategoryID:  in.CategoryID,
		SupplierID:  in.SupplierID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
func (s *ProductService) List(ctx context.Context, categoryID *int) ([]model.Product, error) {
	return s.repo.List(ctx, categoryID)
}
func (s *ProductService) Get(ctx context.Context, id int) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *ProductService) Update(ctx context.Context, id int, in model.ProductInput) (*model.Product, error) {
	return s.repo.Update(ctx, id, in)
}
func (s *ProductService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

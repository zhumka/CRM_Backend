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
	Update(ctx context.Context, id int, in model.SaleInput, status string, taxRate float64) (*model.Sale, error)
	Delete(ctx context.Context, id int) error
	SumAmountByInvoice(ctx context.Context, invoiceID, excludeID int) (float64, error)
}

// SaleStockStore — операции со складом и ставкой продукта для продажи.
type SaleStockStore interface {
	DecreaseStockByName(ctx context.Context, name string, qty int) error
	IncreaseStockByName(ctx context.Context, name string, qty int) error
	TaxRateByName(ctx context.Context, name string) (float64, bool, error)
}

// SaleInvoiceStore — чтение счетов для проверки привязки продажи.
type SaleInvoiceStore interface {
	GetByID(ctx context.Context, id int) (*model.Invoice, error)
}

// SaleRequestStore — чтение заявок для проверки статуса «одобрено».
type SaleRequestStore interface {
	GetByID(ctx context.Context, id int) (*model.PurchaseRequest, error)
}

// SaleTaxStore — чтение справочника налоговых ставок.
type SaleTaxStore interface {
	GetByID(ctx context.Context, id int) (*model.Tax, error)
}

type SaleService struct {
	repo     SaleStore
	products SaleStockStore
	invoices SaleInvoiceStore
	requests SaleRequestStore
	taxes    SaleTaxStore
}

func NewSaleService(repo SaleStore, products SaleStockStore, invoices SaleInvoiceStore, requests SaleRequestStore, taxes SaleTaxStore) *SaleService {
	return &SaleService{repo: repo, products: products, invoices: invoices, requests: requests, taxes: taxes}
}

// resolveTaxRate определяет ставку налога для продажи:
// явная ставка из ввода > ставка из справочника по tax_id > ставка продукта по имени.
func (s *SaleService) resolveTaxRate(ctx context.Context, in model.SaleInput) (float64, error) {
	switch {
	case in.TaxRate != nil:
		return *in.TaxRate, nil
	case in.TaxID != nil:
		t, err := s.taxes.GetByID(ctx, *in.TaxID)
		if err != nil {
			return 0, err
		}
		return t.Rate, nil
	default:
		rate, _, err := s.products.TaxRateByName(ctx, in.ProductName)
		if err != nil {
			return 0, err
		}
		return rate, nil
	}
}

// validateInvoice проверяет привязку продажи к счёту:
//   - за счётом должна стоять одобренная заявка;
//   - сумма всех продаж по счёту не должна превышать его сумму.
//
// excludeID исключает текущую продажу из подсчёта (при обновлении); 0 — при создании.
// Если счёт не указан, проверки не выполняются.
func (s *SaleService) validateInvoice(ctx context.Context, invoiceID *int, amount float64, excludeID int) error {
	if invoiceID == nil {
		return nil
	}
	inv, err := s.invoices.GetByID(ctx, *invoiceID)
	if err != nil {
		return err // ErrNotFound, если счёта нет
	}
	// За счётом должна стоять одобренная заявка.
	if inv.PurchaseRequestID == nil {
		return model.ErrRequestNotApproved
	}
	req, err := s.requests.GetByID(ctx, *inv.PurchaseRequestID)
	if err != nil {
		return err
	}
	if req.Status != model.PurchaseStatusApproved {
		return model.ErrRequestNotApproved
	}
	// Остаток по счёту: уже привязанные продажи + новая не должны превышать сумму счёта.
	used, err := s.repo.SumAmountByInvoice(ctx, *invoiceID, excludeID)
	if err != nil {
		return err
	}
	const eps = 1e-9
	if used+amount > inv.Amount+eps {
		return model.ErrInvoiceAmountExceeded
	}
	return nil
}

// Create списывает товар со склада и регистрирует продажу.
// Если на складе недостаточно товара, возвращает ErrInsufficientStock.
func (s *SaleService) Create(ctx context.Context, in model.SaleInput) (*model.Sale, error) {
	status := in.InstallationStatus
	if status == "" {
		status = model.InstallationNotRequired
	}
	rate, err := s.resolveTaxRate(ctx, in)
	if err != nil {
		return nil, err
	}
	total := in.Amount * (1 + rate/100)
	// Проверяем привязку к счёту (одобренная заявка + лимит суммы с налогом) до списания склада.
	if err := s.validateInvoice(ctx, in.InvoiceID, total, 0); err != nil {
		return nil, err
	}
	// Списываем со склада — так нехватка отсекается до создания продажи.
	if err := s.products.DecreaseStockByName(ctx, in.ProductName, in.Quantity); err != nil {
		return nil, err
	}
	sale := &model.Sale{
		InvoiceID:          in.InvoiceID,
		ProductName:        in.ProductName,
		Quantity:           in.Quantity,
		Amount:             in.Amount,
		TaxRate:            rate,
		InstallationStatus: status,
	}
	if err := s.repo.Create(ctx, sale); err != nil {
		// Продажа не записалась — возвращаем списанное на склад.
		_ = s.products.IncreaseStockByName(ctx, in.ProductName, in.Quantity)
		return nil, err
	}
	sale.ApplyTax()
	return sale, nil
}
func (s *SaleService) List(ctx context.Context) ([]model.Sale, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].ApplyTax()
	}
	return items, nil
}
func (s *SaleService) Get(ctx context.Context, id int) (*model.Sale, error) {
	sale, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	sale.ApplyTax()
	return sale, nil
}

// Update корректирует остатки склада на разницу между старой и новой продажей.
func (s *SaleService) Update(ctx context.Context, id int, in model.SaleInput) (*model.Sale, error) {
	status := in.InstallationStatus
	if status == "" {
		status = model.InstallationNotRequired
	}
	old, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	rate, err := s.resolveTaxRate(ctx, in)
	if err != nil {
		return nil, err
	}
	total := in.Amount * (1 + rate/100)

	// Проверяем привязку к счёту, исключая текущую продажу из суммы по счёту.
	if err := s.validateInvoice(ctx, in.InvoiceID, total, id); err != nil {
		return nil, err
	}

	if old.ProductName == in.ProductName {
		// Тот же товар — корректируем на дельту количества.
		switch delta := in.Quantity - old.Quantity; {
		case delta > 0:
			if err := s.products.DecreaseStockByName(ctx, in.ProductName, delta); err != nil {
				return nil, err
			}
		case delta < 0:
			if err := s.products.IncreaseStockByName(ctx, in.ProductName, -delta); err != nil {
				return nil, err
			}
		}
	} else {
		// Товар сменился: сначала списываем новый (может не хватить), затем возвращаем старый.
		if err := s.products.DecreaseStockByName(ctx, in.ProductName, in.Quantity); err != nil {
			return nil, err
		}
		_ = s.products.IncreaseStockByName(ctx, old.ProductName, old.Quantity)
	}

	updated, err := s.repo.Update(ctx, id, in, status, rate)
	if err != nil {
		// Откат складских изменений при неудачном обновлении.
		if old.ProductName == in.ProductName {
			if delta := in.Quantity - old.Quantity; delta > 0 {
				_ = s.products.IncreaseStockByName(ctx, in.ProductName, delta)
			} else if delta < 0 {
				_ = s.products.DecreaseStockByName(ctx, in.ProductName, -delta)
			}
		} else {
			_ = s.products.IncreaseStockByName(ctx, in.ProductName, in.Quantity)
			_ = s.products.DecreaseStockByName(ctx, old.ProductName, old.Quantity)
		}
		return nil, err
	}
	updated.ApplyTax()
	return updated, nil
}

// Delete отменяет продажу и возвращает товар на склад.
func (s *SaleService) Delete(ctx context.Context, id int) error {
	sale, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.products.IncreaseStockByName(ctx, sale.ProductName, sale.Quantity)
	return nil
}

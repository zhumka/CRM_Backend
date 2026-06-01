// Package seed наполняет базу демонстрационными данными при первом запуске.
// Каждый блок идемпотентен и наполняется, только если его таблица пуста,
// поэтому сидер безопасно отрабатывает и на уже частично заполненной БД.
package seed

import (
	"context"
	"fmt"
	"log"

	"crm/internal/model"
	"crm/internal/pkg/hash"
	"crm/internal/repository"
)

// Seeder наполняет справочники и операционные данные демо-значениями.
type Seeder struct {
	users      *repository.UserRepository
	categories *repository.CategoryRepository
	suppliers  *repository.SupplierRepository
	products   *repository.ProductRepository
	purchases  *repository.PurchaseRequestRepository
	invoices   *repository.InvoiceRepository
	taxes      *repository.TaxRepository
	units      *repository.UnitRepository
	sales      *repository.SaleRepository
}

func NewSeeder(
	users *repository.UserRepository,
	categories *repository.CategoryRepository,
	suppliers *repository.SupplierRepository,
	products *repository.ProductRepository,
	purchases *repository.PurchaseRequestRepository,
	invoices *repository.InvoiceRepository,
	taxes *repository.TaxRepository,
	units *repository.UnitRepository,
	sales *repository.SaleRepository,
) *Seeder {
	return &Seeder{
		users:      users,
		categories: categories,
		suppliers:  suppliers,
		products:   products,
		purchases:  purchases,
		invoices:   invoices,
		taxes:      taxes,
		units:      units,
		sales:      sales,
	}
}

// Run выполняет идемпотентное наполнение по блокам.
func (s *Seeder) Run(ctx context.Context) error {
	demoUser, err := s.ensureUser(ctx)
	if err != nil {
		return fmt.Errorf("seed user: %w", err)
	}
	if err := s.seedCategories(ctx); err != nil {
		return fmt.Errorf("seed categories: %w", err)
	}
	if err := s.seedSuppliers(ctx); err != nil {
		return fmt.Errorf("seed suppliers: %w", err)
	}
	if err := s.seedProducts(ctx); err != nil {
		return fmt.Errorf("seed products: %w", err)
	}
	if err := s.seedRequestsAndInvoices(ctx, demoUser.ID); err != nil {
		return fmt.Errorf("seed requests/invoices: %w", err)
	}
	if err := s.seedTaxes(ctx); err != nil {
		return fmt.Errorf("seed taxes: %w", err)
	}
	if err := s.seedUnits(ctx); err != nil {
		return fmt.Errorf("seed units: %w", err)
	}
	if err := s.seedSales(ctx); err != nil {
		return fmt.Errorf("seed sales: %w", err)
	}
	log.Println("seed: demo data ensured")
	return nil
}

// ensureUser создаёт демо-пользователя (роль user), если его ещё нет.
func (s *Seeder) ensureUser(ctx context.Context) (*model.User, error) {
	const username = "ivanov"
	if u, err := s.users.GetByUsername(ctx, username); err == nil {
		return u, nil
	} else if err != model.ErrNotFound {
		return nil, err
	}

	pwd, err := hash.Hash("user123")
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Username:     username,
		PasswordHash: pwd,
		FullName:     "Иванов Иван",
		Email:        "ivanov@example.com",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Seeder) seedCategories(ctx context.Context) error {
	existing, err := s.categories.List(ctx)
	if err != nil || len(existing) > 0 {
		return err
	}
	data := []model.CategoryInput{
		{Name: "Вентиляционные трубы", Description: "Трубы для вентиляционных систем"},
		{Name: "Фитинги", Description: "Соединительные элементы"},
		{Name: "Крепёж", Description: "Хомуты, кронштейны, метизы"},
	}
	for _, in := range data {
		if err := s.categories.Create(ctx, &model.Category{Name: in.Name, Description: in.Description}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) seedSuppliers(ctx context.Context) error {
	existing, err := s.suppliers.List(ctx)
	if err != nil || len(existing) > 0 {
		return err
	}
	data := []model.Supplier{
		{Name: "ООО ВентСнаб", ContactName: "Петров П.П.", Phone: "+996700112233", Email: "sales@ventsnab.kg", Address: "г. Бишкек, ул. Промышленная, 5"},
		{Name: "ТД Климат", ContactName: "Сидорова А.И.", Phone: "+996555998877", Email: "info@klimat.kg", Address: "г. Ош, ул. Заводская, 12"},
	}
	for i := range data {
		if err := s.suppliers.Create(ctx, &data[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) seedProducts(ctx context.Context) error {
	existing, err := s.products.List(ctx, nil)
	if err != nil || len(existing) > 0 {
		return err
	}
	cats, err := s.categories.List(ctx)
	if err != nil {
		return err
	}
	sups, err := s.suppliers.List(ctx)
	if err != nil {
		return err
	}
	if len(cats) < 3 || len(sups) < 2 {
		return nil // справочники не готовы — пропускаем
	}
	cat := func(i int) *int { return &cats[i].ID }
	sup := func(i int) *int { return &sups[i].ID }

	data := []model.Product{
		{Name: "Труба круглая d100", Description: "Оцинкованная сталь, L=1м", Price: 350, Stock: 80, Unit: "м", TaxRate: 12, CategoryID: cat(0), SupplierID: sup(0)},
		{Name: "Труба круглая d125", Description: "Оцинкованная сталь, L=1м", Price: 420, Stock: 15, Unit: "м", TaxRate: 12, CategoryID: cat(0), SupplierID: sup(0)},
		{Name: "Труба прямоугольная 60x120", Description: "Воздуховод, L=1м", Price: 510, Stock: 42, Unit: "м", TaxRate: 12, CategoryID: cat(0), SupplierID: sup(1)},
		{Name: "Отвод 90° d100", Description: "Соединительный отвод", Price: 180, Stock: 120, Unit: "шт", TaxRate: 12, CategoryID: cat(1), SupplierID: sup(0)},
		{Name: "Тройник d125", Description: "Тройник вентиляционный", Price: 260, Stock: 8, Unit: "шт", TaxRate: 12, CategoryID: cat(1), SupplierID: sup(1)},
		{Name: "Хомут d100", Description: "Крепёжный хомут с резиной", Price: 45, Stock: 300, Unit: "шт", TaxRate: 12, CategoryID: cat(2), SupplierID: sup(1)},
	}
	for i := range data {
		if err := s.products.Create(ctx, &data[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) seedRequestsAndInvoices(ctx context.Context, userID int) error {
	existing, err := s.purchases.List(ctx, nil)
	if err != nil || len(existing) > 0 {
		return err
	}
	products, err := s.products.List(ctx, nil)
	if err != nil || len(products) < 6 {
		return err
	}

	requests := []struct {
		product int
		qty     int
		status  string
		client  string
		comment string
	}{
		{product: 0, qty: 50, status: model.PurchaseStatusCompleted, client: "ОсОО СтройМир", comment: "Партия для объекта №1"},
		{product: 2, qty: 20, status: model.PurchaseStatusApproved, client: "ИП Асанов", comment: "Срочный заказ"},
		{product: 3, qty: 100, status: model.PurchaseStatusPending, client: "ОсОО ТеплоДом", comment: "Плановая закупка"},
		{product: 5, qty: 200, status: model.PurchaseStatusChecking, client: "Частное лицо", comment: "Крепёж на склад"},
	}

	i := 0
	for _, r := range requests {
		i++
		pr := &model.PurchaseRequest{
			UserID:     userID,
			ClientName: r.client,
			ProductID:  products[r.product].ID,
			Quantity:   r.qty,
			Status:     r.status,
			Comment:    r.comment,
		}
		if err := s.purchases.Create(ctx, pr); err != nil {
			return err
		}
		// Для завершённых/одобренных заявок выставляем счёт-фактуру.
		if r.status == model.PurchaseStatusCompleted || r.status == model.PurchaseStatusApproved {
			amount := products[r.product].Price * float64(r.qty)
			status := model.InvoiceStatusIssued
			if r.status == model.PurchaseStatusCompleted {
				status = model.InvoiceStatusPaid
			}
			inv := &model.Invoice{
				Number:            fmt.Sprintf("INV-2026-%03d", i),
				PurchaseRequestID: &pr.ID,
				Amount:            amount,
				Status:            status,
			}
			if err := s.invoices.Create(ctx, inv); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Seeder) seedTaxes(ctx context.Context) error {
	existing, err := s.taxes.List(ctx)
	if err != nil || len(existing) > 0 {
		return err
	}
	data := []model.Tax{
		{Name: "НДС", Rate: 12, Active: true},
		{Name: "Налог с продаж", Rate: 2, Active: true},
		{Name: "Без налога", Rate: 0, Active: false},
	}
	for i := range data {
		if err := s.taxes.Create(ctx, &data[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) seedUnits(ctx context.Context) error {
	existing, err := s.units.List(ctx)
	if err != nil || len(existing) > 0 {
		return err
	}
	data := []model.Unit{
		{Name: "Штука", ShortName: "шт", Description: "Поштучный учёт"},
		{Name: "Метр", ShortName: "м", Description: "Погонный метр трубы"},
		{Name: "Комплект", ShortName: "компл", Description: "Набор изделий"},
	}
	for i := range data {
		if err := s.units.Create(ctx, &data[i]); err != nil {
			return err
		}
	}
	return nil
}

// seedSales формирует продажи из оплаченных счетов-фактур.
func (s *Seeder) seedSales(ctx context.Context) error {
	existing, err := s.sales.List(ctx)
	if err != nil || len(existing) > 0 {
		return err
	}
	invoices, err := s.invoices.List(ctx)
	if err != nil {
		return err
	}
	for i := range invoices {
		inv := invoices[i]
		if inv.Status != model.InvoiceStatusPaid {
			continue
		}
		id := inv.ID
		sale := &model.Sale{
			InvoiceID:          &id,
			ProductName:        "Продажа по счёту " + inv.Number,
			Quantity:           1,
			Amount:             inv.Amount,
			InstallationStatus: model.InstallationCompleted,
		}
		if err := s.sales.Create(ctx, sale); err != nil {
			return err
		}
	}
	return nil
}

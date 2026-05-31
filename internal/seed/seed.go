// Package seed наполняет базу демонстрационными данными при первом запуске.
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
}

func NewSeeder(
	users *repository.UserRepository,
	categories *repository.CategoryRepository,
	suppliers *repository.SupplierRepository,
	products *repository.ProductRepository,
	purchases *repository.PurchaseRequestRepository,
	invoices *repository.InvoiceRepository,
) *Seeder {
	return &Seeder{
		users:      users,
		categories: categories,
		suppliers:  suppliers,
		products:   products,
		purchases:  purchases,
		invoices:   invoices,
	}
}

// Run выполняет идемпотентное наполнение: каждый блок пропускается, если данные уже есть.
func (s *Seeder) Run(ctx context.Context) error {
	demoUser, err := s.ensureUser(ctx)
	if err != nil {
		return fmt.Errorf("seed user: %w", err)
	}

	cats, err := s.seedCategories(ctx)
	if err != nil {
		return fmt.Errorf("seed categories: %w", err)
	}
	if cats == nil {
		// Категории уже были — считаем, что демо-данные загружены ранее.
		log.Println("seed: demo data already present, skipping")
		return nil
	}

	sups, err := s.seedSuppliers(ctx)
	if err != nil {
		return fmt.Errorf("seed suppliers: %w", err)
	}

	prods, err := s.seedProducts(ctx, cats, sups)
	if err != nil {
		return fmt.Errorf("seed products: %w", err)
	}

	if err := s.seedRequestsAndInvoices(ctx, demoUser.ID, prods); err != nil {
		return fmt.Errorf("seed requests/invoices: %w", err)
	}

	log.Println("seed: demo data created")
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
	u := &model.User{Username: username, PasswordHash: pwd, Role: model.RoleUser}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// seedCategories возвращает nil, если категории уже существуют (признак ранее загруженных демо-данных).
func (s *Seeder) seedCategories(ctx context.Context) ([]model.Category, error) {
	existing, err := s.categories.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, nil
	}

	names := []model.CategoryInput{
		{Name: "Вентиляционные трубы", Description: "Трубы для вентиляционных систем"},
		{Name: "Фитинги", Description: "Соединительные элементы"},
		{Name: "Крепёж", Description: "Хомуты, кронштейны, метизы"},
	}
	out := make([]model.Category, 0, len(names))
	for _, in := range names {
		c := &model.Category{Name: in.Name, Description: in.Description}
		if err := s.categories.Create(ctx, c); err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, nil
}

func (s *Seeder) seedSuppliers(ctx context.Context) ([]model.Supplier, error) {
	data := []model.Supplier{
		{Name: "ООО ВентСнаб", ContactName: "Петров П.П.", Phone: "+996700112233", Email: "sales@ventsnab.kg", Address: "г. Бишкек, ул. Промышленная, 5"},
		{Name: "ТД Климат", ContactName: "Сидорова А.И.", Phone: "+996555998877", Email: "info@klimat.kg", Address: "г. Ош, ул. Заводская, 12"},
	}
	out := make([]model.Supplier, 0, len(data))
	for i := range data {
		if err := s.suppliers.Create(ctx, &data[i]); err != nil {
			return nil, err
		}
		out = append(out, data[i])
	}
	return out, nil
}

func (s *Seeder) seedProducts(ctx context.Context, cats []model.Category, sups []model.Supplier) ([]model.Product, error) {
	cat := func(i int) *int { return &cats[i].ID }
	sup := func(i int) *int { return &sups[i].ID }

	data := []model.Product{
		{Name: "Труба круглая d100", Description: "Оцинкованная сталь, L=1м", Price: 350.00, CategoryID: cat(0), SupplierID: sup(0)},
		{Name: "Труба круглая d125", Description: "Оцинкованная сталь, L=1м", Price: 420.00, CategoryID: cat(0), SupplierID: sup(0)},
		{Name: "Труба прямоугольная 60x120", Description: "Воздуховод, L=1м", Price: 510.00, CategoryID: cat(0), SupplierID: sup(1)},
		{Name: "Отвод 90° d100", Description: "Соединительный отвод", Price: 180.00, CategoryID: cat(1), SupplierID: sup(0)},
		{Name: "Тройник d125", Description: "Тройник вентиляционный", Price: 260.00, CategoryID: cat(1), SupplierID: sup(1)},
		{Name: "Хомут d100", Description: "Крепёжный хомут с резиной", Price: 45.00, CategoryID: cat(2), SupplierID: sup(1)},
	}
	out := make([]model.Product, 0, len(data))
	for i := range data {
		if err := s.products.Create(ctx, &data[i]); err != nil {
			return nil, err
		}
		out = append(out, data[i])
	}
	return out, nil
}

func (s *Seeder) seedRequestsAndInvoices(ctx context.Context, userID int, prods []model.Product) error {
	requests := []struct {
		product int
		qty     int
		status  string
		comment string
	}{
		{product: 0, qty: 50, status: model.PurchaseStatusCompleted, comment: "Партия для объекта №1"},
		{product: 2, qty: 20, status: model.PurchaseStatusApproved, comment: "Срочный заказ"},
		{product: 3, qty: 100, status: model.PurchaseStatusNew, comment: "Плановая закупка"},
		{product: 5, qty: 200, status: model.PurchaseStatusInProgress, comment: "Крепёж на склад"},
	}

	for i, r := range requests {
		pr := &model.PurchaseRequest{
			UserID:    userID,
			ProductID: prods[r.product].ID,
			Quantity:  r.qty,
			Status:    r.status,
			Comment:   r.comment,
		}
		if err := s.purchases.Create(ctx, pr); err != nil {
			return err
		}
		// Для завершённых/одобренных заявок выставляем счёт-фактуру.
		if r.status == model.PurchaseStatusCompleted || r.status == model.PurchaseStatusApproved {
			amount := prods[r.product].Price * float64(r.qty)
			status := model.InvoiceStatusUnpaid
			if r.status == model.PurchaseStatusCompleted {
				status = model.InvoiceStatusPaid
			}
			inv := &model.Invoice{
				Number:            fmt.Sprintf("INV-2026-%03d", i+1),
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

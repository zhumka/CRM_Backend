package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

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

// InvoiceSaleStore — выборка продаж-позиций по счёту для документа.
type InvoiceSaleStore interface {
	ListByInvoice(ctx context.Context, invoiceID int) ([]model.Sale, error)
}

// InvoiceRequestStore — чтение связанной заявки (клиент, товар).
type InvoiceRequestStore interface {
	GetByID(ctx context.Context, id int) (*model.PurchaseRequest, error)
}

// InvoiceProductStore — чтение продукта связанной заявки.
type InvoiceProductStore interface {
	GetByID(ctx context.Context, id int) (*model.Product, error)
}

// InvoiceService — управление счетами-фактурами.
type InvoiceService struct {
	repo     InvoiceStore
	sales    InvoiceSaleStore
	requests InvoiceRequestStore
	products InvoiceProductStore
}

func NewInvoiceService(repo InvoiceStore, sales InvoiceSaleStore, requests InvoiceRequestStore, products InvoiceProductStore) *InvoiceService {
	return &InvoiceService{repo: repo, sales: sales, requests: requests, products: products}
}

func (s *InvoiceService) Create(ctx context.Context, in model.InvoiceInput) (*model.Invoice, error) {
	status := in.Status
	if status == "" {
		status = model.InvoiceStatusIssued
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

// --- Документ счёта-фактуры (HTML для скачивания/печати) ---

type invoiceDocLine struct {
	Name      string
	Qty       int
	Amount    string
	TaxRate   string
	TaxAmount string
	Total     string
}

type invoiceDocView struct {
	Number     string
	Status     string
	IssuedDate string
	DueDate    string
	ClientName string
	Lines      []invoiceDocLine
	SubTotal   string
	TaxTotal   string
	GrandTotal string
	Amount     string
}

var invoiceStatusLabels = map[string]string{
	model.InvoiceStatusDraft:   "Черновик",
	model.InvoiceStatusIssued:  "Выставлен",
	model.InvoiceStatusPaid:    "Оплачен",
	model.InvoiceStatusOverdue: "Просрочен",
}

// Document формирует HTML-документ счёта-фактуры и имя файла для скачивания.
func (s *InvoiceService) Document(ctx context.Context, id int) (string, []byte, error) {
	inv, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", nil, err
	}

	view := invoiceDocView{
		Number:     inv.Number,
		Status:     statusLabel(inv.Status),
		IssuedDate: docDate(&inv.IssuedDate),
		DueDate:    docDate(inv.DueDate),
		Amount:     money(inv.Amount),
	}

	// Клиент из связанной заявки, если она есть.
	if inv.PurchaseRequestID != nil {
		if pr, err := s.requests.GetByID(ctx, *inv.PurchaseRequestID); err == nil {
			view.ClientName = pr.ClientName
		}
	}

	// Позиции — продажи, привязанные к счёту.
	lines, err := s.sales.ListByInvoice(ctx, id)
	if err != nil {
		return "", nil, err
	}
	var subTotal, taxTotal, grandTotal float64
	for i := range lines {
		lines[i].ApplyTax()
		subTotal += lines[i].Amount
		taxTotal += lines[i].TaxAmount
		grandTotal += lines[i].Total
		view.Lines = append(view.Lines, invoiceDocLine{
			Name:      lines[i].ProductName,
			Qty:       lines[i].Quantity,
			Amount:    money(lines[i].Amount),
			TaxRate:   money(lines[i].TaxRate),
			TaxAmount: money(lines[i].TaxAmount),
			Total:     money(lines[i].Total),
		})
	}
	view.SubTotal = money(subTotal)
	view.TaxTotal = money(taxTotal)
	view.GrandTotal = money(grandTotal)

	var buf bytes.Buffer
	if err := invoiceTemplate.Execute(&buf, view); err != nil {
		return "", nil, err
	}
	filename := fmt.Sprintf("invoice_%s.html", inv.Number)
	return filename, buf.Bytes(), nil
}

func statusLabel(status string) string {
	if l, ok := invoiceStatusLabels[status]; ok {
		return l
	}
	return status
}

func docDate(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format("02.01.2006")
}

var invoiceTemplate = template.Must(template.New("invoice").Parse(`<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="utf-8">
<title>Счёт-фактура {{.Number}}</title>
<style>
  body { font-family: Arial, sans-serif; color: #1a1a1a; margin: 40px; }
  h1 { font-size: 22px; margin-bottom: 4px; }
  .meta { margin-bottom: 24px; color: #444; }
  .meta div { margin: 2px 0; }
  table { border-collapse: collapse; width: 100%; margin-top: 12px; }
  th, td { border: 1px solid #ccc; padding: 8px 10px; text-align: left; font-size: 14px; }
  th { background: #f2f2f2; }
  td.num, th.num { text-align: right; }
  tfoot td { font-weight: bold; }
  .total { margin-top: 18px; font-size: 16px; }
  @media print { body { margin: 0; } }
</style>
</head>
<body>
  <h1>Счёт-фактура № {{.Number}}</h1>
  <div class="meta">
    <div>Статус: {{.Status}}</div>
    <div>Дата выставления: {{.IssuedDate}}</div>
    <div>Срок оплаты: {{.DueDate}}</div>
    {{if .ClientName}}<div>Клиент: {{.ClientName}}</div>{{end}}
  </div>

  <table>
    <thead>
      <tr>
        <th>Наименование</th>
        <th class="num">Кол-во</th>
        <th class="num">Сумма без налога</th>
        <th class="num">Ставка, %</th>
        <th class="num">Налог</th>
        <th class="num">Итого</th>
      </tr>
    </thead>
    <tbody>
      {{range .Lines}}
      <tr>
        <td>{{.Name}}</td>
        <td class="num">{{.Qty}}</td>
        <td class="num">{{.Amount}}</td>
        <td class="num">{{.TaxRate}}</td>
        <td class="num">{{.TaxAmount}}</td>
        <td class="num">{{.Total}}</td>
      </tr>
      {{else}}
      <tr><td colspan="6">Позиции по счёту отсутствуют</td></tr>
      {{end}}
    </tbody>
    <tfoot>
      <tr>
        <td colspan="2">Итого</td>
        <td class="num">{{.SubTotal}}</td>
        <td class="num"></td>
        <td class="num">{{.TaxTotal}}</td>
        <td class="num">{{.GrandTotal}}</td>
      </tr>
    </tfoot>
  </table>

  <p class="total">Сумма счёта: {{.Amount}}</p>
</body>
</html>
`))

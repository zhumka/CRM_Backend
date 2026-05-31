package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"crm/internal/model"
)

// ReportStore — зависимость от хранилища отчётов.
type ReportStore interface {
	Create(ctx context.Context, rep *model.Report) error
	List(ctx context.Context) ([]model.Report, error)
	GetByID(ctx context.Context, id int) (*model.Report, error)
	Update(ctx context.Context, id int, in model.ReportInput) (*model.Report, error)
	Delete(ctx context.Context, id int) error
}

// ReportService — формирование, хранение и экспорт отчётов.
type ReportService struct {
	repo      ReportStore
	analytics *AnalyticsService
}

func NewReportService(repo ReportStore, analytics *AnalyticsService) *ReportService {
	return &ReportService{repo: repo, analytics: analytics}
}

func (s *ReportService) Create(ctx context.Context, in model.ReportInput) (*model.Report, error) {
	rep := &model.Report{Title: in.Title, Content: in.Content, CategoryID: in.CategoryID}
	if err := s.repo.Create(ctx, rep); err != nil {
		return nil, err
	}
	return rep, nil
}

func (s *ReportService) List(ctx context.Context) ([]model.Report, error) { return s.repo.List(ctx) }

func (s *ReportService) Get(ctx context.Context, id int) (*model.Report, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ReportService) Update(ctx context.Context, id int, in model.ReportInput) (*model.Report, error) {
	return s.repo.Update(ctx, id, in)
}

func (s *ReportService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }

// GenerateSales формирует отчёт о продажах за период на основе реальной аналитики
// и сохраняет его (содержимое — готовый CSV).
func (s *ReportService) GenerateSales(ctx context.Context, in model.GenerateSalesReportInput) (*model.Report, error) {
	sales, err := s.analytics.Sales(ctx, in.From, in.To)
	if err != nil {
		return nil, err
	}
	byStatus, err := s.analytics.RequestsByStatus(ctx)
	if err != nil {
		return nil, err
	}

	content, err := buildSalesCSV(sales, byStatus)
	if err != nil {
		return nil, err
	}

	title := in.Title
	if title == "" {
		title = "Отчёт о продажах " + periodLabel(in.From, in.To)
	}

	rep := &model.Report{Title: title, Content: content}
	if err := s.repo.Create(ctx, rep); err != nil {
		return nil, err
	}
	return rep, nil
}

// buildSalesCSV формирует CSV-содержимое отчёта (разделитель «;» — удобно для Excel).
func buildSalesCSV(sa *model.SalesAnalytics, byStatus []model.StatusCount) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)
	w.Comma = ';'

	rows := [][]string{
		{"Финансовая аналитика по счетам"},
		{"Период с", dateOrDash(sa.From)},
		{"Период по", dateOrDash(sa.To)},
		{"Сформирован", time.Now().Format(dateLayout)},
		{},
		{"Показатель", "Значение"},
		{"Количество счетов", strconv.Itoa(sa.InvoiceCount)},
		{"Общая сумма", money(sa.TotalAmount)},
		{"Оплачено", money(sa.PaidAmount)},
		{"Не оплачено", money(sa.UnpaidAmount)},
		{},
		{"Статус заявки", "Количество"},
	}
	for _, sc := range byStatus {
		rows = append(rows, []string{sc.Status, strconv.Itoa(sc.Count)})
	}

	if err := w.WriteAll(rows); err != nil {
		return "", err
	}
	w.Flush()
	return sb.String(), w.Error()
}

func periodLabel(from, to string) string {
	switch {
	case from == "" && to == "":
		return "за всё время"
	case from != "" && to != "":
		return fmt.Sprintf("за %s — %s", from, to)
	case from != "":
		return "с " + from
	default:
		return "по " + to
	}
}

func dateOrDash(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format(dateLayout)
}

func money(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

package service

import (
	"context"
	"fmt"
	"time"

	"crm/internal/model"
)

const dateLayout = "2006-01-02"

// AnalyticsStore — зависимость от агрегированных выборок.
type AnalyticsStore interface {
	Summary(ctx context.Context) (*model.Summary, error)
	SalesAnalytics(ctx context.Context, from, to *time.Time) (*model.SalesAnalytics, error)
	RequestsByStatus(ctx context.Context) ([]model.StatusCount, error)
	TopProducts(ctx context.Context, limit int) ([]model.TopProduct, error)
}

// AnalyticsService — подсистема аналитики и отчётности (KPI, финансы, рейтинги).
type AnalyticsService struct {
	repo AnalyticsStore
}

func NewAnalyticsService(repo AnalyticsStore) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) Summary(ctx context.Context) (*model.Summary, error) {
	return s.repo.Summary(ctx)
}

// Sales рассчитывает финансовую аналитику за период. Пустые fromStr/toStr — без границы.
func (s *AnalyticsService) Sales(ctx context.Context, fromStr, toStr string) (*model.SalesAnalytics, error) {
	from, err := parseDate(fromStr)
	if err != nil {
		return nil, err
	}
	to, err := parseDate(toStr)
	if err != nil {
		return nil, err
	}
	return s.repo.SalesAnalytics(ctx, from, to)
}

func (s *AnalyticsService) RequestsByStatus(ctx context.Context) ([]model.StatusCount, error) {
	return s.repo.RequestsByStatus(ctx)
}

func (s *AnalyticsService) TopProducts(ctx context.Context, limit int) ([]model.TopProduct, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.repo.TopProducts(ctx, limit)
}

// parseDate разбирает дату YYYY-MM-DD; пустая строка → nil (без ограничения).
func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return nil, fmt.Errorf("invalid date %q, expected YYYY-MM-DD", s)
	}
	return &t, nil
}

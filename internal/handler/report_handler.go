package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listReports godoc
// @Summary  Список отчётов
// @Tags     reports
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Report
// @Router   /reports [get]
func (h *Handler) listReports(c *gin.Context) {
	items, err := h.reports.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getReport godoc
// @Summary  Получить отчёт
// @Tags     reports
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID отчёта"
// @Success  200  {object}  model.Report
// @Failure  404  {object}  errorResponse
// @Router   /reports/{id} [get]
func (h *Handler) getReport(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.reports.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createReport godoc
// @Summary  Создать отчёт вручную
// @Tags     reports
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.ReportInput  true  "Данные отчёта"
// @Success  201    {object}  model.Report
// @Failure  400    {object}  errorResponse
// @Router   /reports [post]
func (h *Handler) createReport(c *gin.Context) {
	var in model.ReportInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.reports.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// generateSalesReport godoc
// @Summary  Сформировать отчёт о продажах
// @Description Генерирует отчёт за период на основе реальной аналитики и сохраняет его (содержимое — CSV)
// @Tags     reports
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.GenerateSalesReportInput  true  "Период (поля опциональны)"
// @Success  201    {object}  model.Report
// @Failure  400    {object}  errorResponse
// @Router   /reports/generate-sales [post]
func (h *Handler) generateSalesReport(c *gin.Context) {
	var in model.GenerateSalesReportInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.reports.GenerateSales(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// deleteReport godoc
// @Summary  Удалить отчёт
// @Tags     reports
// @Security BearerAuth
// @Param    id   path  int  true  "ID отчёта"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /reports/{id} [delete]
func (h *Handler) deleteReport(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.reports.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// exportReport godoc
// @Summary  Экспорт отчёта в CSV
// @Description Отдаёт содержимое отчёта как CSV-файл (UTF-8 BOM, для Excel)
// @Tags     reports
// @Produce  text/csv
// @Security BearerAuth
// @Param    id   path      int  true  "ID отчёта"
// @Success  200  {string}  string  "CSV-файл"
// @Failure  404  {object}  errorResponse
// @Router   /reports/{id}/export [get]
func (h *Handler) exportReport(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	rep, err := h.reports.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	filename := fmt.Sprintf("report_%d.csv", rep.ID)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	// BOM, чтобы Excel корректно распознал UTF-8 (кириллицу).
	c.Data(http.StatusOK, "text/csv; charset=utf-8", append([]byte("\xEF\xBB\xBF"), []byte(rep.Content)...))
}

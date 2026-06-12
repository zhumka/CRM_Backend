package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listInvoices godoc
// @Summary  Список счетов-фактур
// @Tags     invoices
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Invoice
// @Router   /invoices [get]
func (h *Handler) listInvoices(c *gin.Context) {
	items, err := h.invoices.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getInvoice godoc
// @Summary  Получить счёт-фактуру
// @Tags     invoices
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID счёта"
// @Success  200  {object}  model.Invoice
// @Failure  404  {object}  errorResponse
// @Router   /invoices/{id} [get]
func (h *Handler) getInvoice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.invoices.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// downloadInvoice godoc
// @Summary  Скачать счёт-фактуру
// @Description Отдаёт счёт-фактуру как HTML-документ (можно открыть и распечатать/сохранить в PDF)
// @Tags     invoices
// @Produce  text/html
// @Security BearerAuth
// @Param    id   path      int  true  "ID счёта"
// @Success  200  {string}  string  "HTML-документ"
// @Failure  404  {object}  errorResponse
// @Router   /invoices/{id}/document [get]
func (h *Handler) downloadInvoice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	filename, content, err := h.invoices.Document(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

// createInvoice godoc
// @Summary  Создать счёт-фактуру
// @Tags     invoices
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.InvoiceInput  true  "Данные счёта"
// @Success  201    {object}  model.Invoice
// @Failure  400    {object}  errorResponse
// @Failure  409    {object}  errorResponse
// @Router   /invoices [post]
func (h *Handler) createInvoice(c *gin.Context) {
	var in model.InvoiceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.invoices.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateInvoice godoc
// @Summary  Обновить счёт-фактуру
// @Tags     invoices
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int                 true  "ID счёта"
// @Param    input  body      model.InvoiceInput  true  "Данные счёта"
// @Success  200    {object}  model.Invoice
// @Failure  404    {object}  errorResponse
// @Router   /invoices/{id} [put]
func (h *Handler) updateInvoice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.InvoiceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.invoices.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteInvoice godoc
// @Summary  Удалить счёт-фактуру
// @Tags     invoices
// @Security BearerAuth
// @Param    id   path  int  true  "ID счёта"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /invoices/{id} [delete]
func (h *Handler) deleteInvoice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.invoices.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listTaxes godoc
// @Summary  Список налоговых ставок
// @Tags     taxes
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Tax
// @Router   /taxes [get]
func (h *Handler) listTaxes(c *gin.Context) {
	items, err := h.taxes.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getTax godoc
// @Summary  Получить налоговую ставку
// @Tags     taxes
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID"
// @Success  200  {object}  model.Tax
// @Failure  404  {object}  errorResponse
// @Router   /taxes/{id} [get]
func (h *Handler) getTax(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.taxes.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createTax godoc
// @Summary  Создать налоговую ставку
// @Tags     taxes
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.TaxInput  true  "Данные"
// @Success  201    {object}  model.Tax
// @Failure  400    {object}  errorResponse
// @Router   /taxes [post]
func (h *Handler) createTax(c *gin.Context) {
	var in model.TaxInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.taxes.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateTax godoc
// @Summary  Обновить налоговую ставку
// @Tags     taxes
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int             true  "ID"
// @Param    input  body      model.TaxInput  true  "Данные"
// @Success  200    {object}  model.Tax
// @Failure  404    {object}  errorResponse
// @Router   /taxes/{id} [put]
func (h *Handler) updateTax(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.TaxInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.taxes.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteTax godoc
// @Summary  Удалить налоговую ставку
// @Tags     taxes
// @Security BearerAuth
// @Param    id   path  int  true  "ID"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /taxes/{id} [delete]
func (h *Handler) deleteTax(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.taxes.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listSales godoc
// @Summary  Список продаж
// @Tags     sales
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Sale
// @Router   /sales [get]
func (h *Handler) listSales(c *gin.Context) {
	items, err := h.sales.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getSale godoc
// @Summary  Получить продажу
// @Tags     sales
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID"
// @Success  200  {object}  model.Sale
// @Failure  404  {object}  errorResponse
// @Router   /sales/{id} [get]
func (h *Handler) getSale(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.sales.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createSale godoc
// @Summary  Создать продажу
// @Tags     sales
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.SaleInput  true  "Данные"
// @Success  201    {object}  model.Sale
// @Failure  400    {object}  errorResponse
// @Router   /sales [post]
func (h *Handler) createSale(c *gin.Context) {
	var in model.SaleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.sales.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateSale godoc
// @Summary  Обновить продажу
// @Tags     sales
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int              true  "ID"
// @Param    input  body      model.SaleInput  true  "Данные"
// @Success  200    {object}  model.Sale
// @Failure  404    {object}  errorResponse
// @Router   /sales/{id} [put]
func (h *Handler) updateSale(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.SaleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.sales.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteSale godoc
// @Summary  Удалить продажу
// @Tags     sales
// @Security BearerAuth
// @Param    id   path  int  true  "ID"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /sales/{id} [delete]
func (h *Handler) deleteSale(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.sales.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

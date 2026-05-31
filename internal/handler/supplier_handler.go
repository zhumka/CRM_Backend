package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listSuppliers godoc
// @Summary  Список поставщиков
// @Tags     suppliers
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Supplier
// @Router   /suppliers [get]
func (h *Handler) listSuppliers(c *gin.Context) {
	items, err := h.suppliers.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getSupplier godoc
// @Summary  Получить поставщика
// @Tags     suppliers
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID поставщика"
// @Success  200  {object}  model.Supplier
// @Failure  404  {object}  errorResponse
// @Router   /suppliers/{id} [get]
func (h *Handler) getSupplier(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.suppliers.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createSupplier godoc
// @Summary  Создать поставщика
// @Tags     suppliers
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.SupplierInput  true  "Данные поставщика"
// @Success  201    {object}  model.Supplier
// @Failure  400    {object}  errorResponse
// @Router   /suppliers [post]
func (h *Handler) createSupplier(c *gin.Context) {
	var in model.SupplierInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.suppliers.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateSupplier godoc
// @Summary  Обновить поставщика
// @Tags     suppliers
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int                  true  "ID поставщика"
// @Param    input  body      model.SupplierInput  true  "Данные поставщика"
// @Success  200    {object}  model.Supplier
// @Failure  404    {object}  errorResponse
// @Router   /suppliers/{id} [put]
func (h *Handler) updateSupplier(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.SupplierInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.suppliers.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteSupplier godoc
// @Summary  Удалить поставщика
// @Tags     suppliers
// @Security BearerAuth
// @Param    id   path  int  true  "ID поставщика"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /suppliers/{id} [delete]
func (h *Handler) deleteSupplier(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.suppliers.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

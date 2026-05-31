package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listProducts godoc
// @Summary  Список продуктов
// @Tags     products
// @Produce  json
// @Security BearerAuth
// @Param    category_id  query     int  false  "Фильтр по категории"
// @Success  200          {array}   model.Product
// @Router   /products [get]
func (h *Handler) listProducts(c *gin.Context) {
	var categoryID *int
	if raw := c.Query("category_id"); raw != "" {
		if id, err := strconv.Atoi(raw); err == nil {
			categoryID = &id
		}
	}
	items, err := h.products.List(c.Request.Context(), categoryID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getProduct godoc
// @Summary  Получить продукт
// @Tags     products
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID продукта"
// @Success  200  {object}  model.Product
// @Failure  404  {object}  errorResponse
// @Router   /products/{id} [get]
func (h *Handler) getProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.products.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createProduct godoc
// @Summary  Создать продукт
// @Tags     products
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.ProductInput  true  "Данные продукта"
// @Success  201    {object}  model.Product
// @Failure  400    {object}  errorResponse
// @Failure  404    {object}  errorResponse
// @Router   /products [post]
func (h *Handler) createProduct(c *gin.Context) {
	var in model.ProductInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.products.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateProduct godoc
// @Summary  Обновить продукт
// @Tags     products
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int                 true  "ID продукта"
// @Param    input  body      model.ProductInput  true  "Данные продукта"
// @Success  200    {object}  model.Product
// @Failure  404    {object}  errorResponse
// @Router   /products/{id} [put]
func (h *Handler) updateProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.ProductInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.products.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteProduct godoc
// @Summary  Удалить продукт
// @Tags     products
// @Security BearerAuth
// @Param    id   path  int  true  "ID продукта"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /products/{id} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.products.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listCategories godoc
// @Summary  Список категорий
// @Tags     categories
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Category
// @Router   /categories [get]
func (h *Handler) listCategories(c *gin.Context) {
	items, err := h.categories.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getCategory godoc
// @Summary  Получить категорию
// @Tags     categories
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID категории"
// @Success  200  {object}  model.Category
// @Failure  404  {object}  errorResponse
// @Router   /categories/{id} [get]
func (h *Handler) getCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.categories.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createCategory godoc
// @Summary  Создать категорию
// @Tags     categories
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.CategoryInput  true  "Данные категории"
// @Success  201    {object}  model.Category
// @Failure  400    {object}  errorResponse
// @Router   /categories [post]
func (h *Handler) createCategory(c *gin.Context) {
	var in model.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.categories.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateCategory godoc
// @Summary  Обновить категорию
// @Tags     categories
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int                  true  "ID категории"
// @Param    input  body      model.CategoryInput  true  "Данные категории"
// @Success  200    {object}  model.Category
// @Failure  404    {object}  errorResponse
// @Router   /categories/{id} [put]
func (h *Handler) updateCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.categories.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteCategory godoc
// @Summary  Удалить категорию
// @Tags     categories
// @Security BearerAuth
// @Param    id   path  int  true  "ID категории"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /categories/{id} [delete]
func (h *Handler) deleteCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.categories.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

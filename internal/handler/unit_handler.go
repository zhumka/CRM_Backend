package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listUnits godoc
// @Summary  Список единиц измерения
// @Tags     units
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.Unit
// @Router   /units [get]
func (h *Handler) listUnits(c *gin.Context) {
	items, err := h.units.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getUnit godoc
// @Summary  Получить единицу измерения
// @Tags     units
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID"
// @Success  200  {object}  model.Unit
// @Failure  404  {object}  errorResponse
// @Router   /units/{id} [get]
func (h *Handler) getUnit(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.units.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createUnit godoc
// @Summary  Создать единицу измерения
// @Tags     units
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.UnitInput  true  "Данные"
// @Success  201    {object}  model.Unit
// @Failure  400    {object}  errorResponse
// @Router   /units [post]
func (h *Handler) createUnit(c *gin.Context) {
	var in model.UnitInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.units.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateUnit godoc
// @Summary  Обновить единицу измерения
// @Tags     units
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int              true  "ID"
// @Param    input  body      model.UnitInput  true  "Данные"
// @Success  200    {object}  model.Unit
// @Failure  404    {object}  errorResponse
// @Router   /units/{id} [put]
func (h *Handler) updateUnit(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.UnitInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.units.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deleteUnit godoc
// @Summary  Удалить единицу измерения
// @Tags     units
// @Security BearerAuth
// @Param    id   path  int  true  "ID"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /units/{id} [delete]
func (h *Handler) deleteUnit(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.units.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

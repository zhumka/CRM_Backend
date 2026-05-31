package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listPurchaseRequests godoc
// @Summary  Список заявок на закупку
// @Description Администратор видит все заявки, пользователь — только свои
// @Tags     purchase-requests
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.PurchaseRequest
// @Router   /purchase-requests [get]
func (h *Handler) listPurchaseRequests(c *gin.Context) {
	items, err := h.purchases.List(c.Request.Context(), currentUserID(c), isAdmin(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// getPurchaseRequest godoc
// @Summary  Получить заявку
// @Tags     purchase-requests
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID заявки"
// @Success  200  {object}  model.PurchaseRequest
// @Failure  403  {object}  errorResponse
// @Failure  404  {object}  errorResponse
// @Router   /purchase-requests/{id} [get]
func (h *Handler) getPurchaseRequest(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.purchases.Get(c.Request.Context(), id, currentUserID(c), isAdmin(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// createPurchaseRequest godoc
// @Summary  Создать заявку на закупку
// @Tags     purchase-requests
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.PurchaseRequestInput  true  "Данные заявки"
// @Success  201    {object}  model.PurchaseRequest
// @Failure  400    {object}  errorResponse
// @Failure  404    {object}  errorResponse
// @Router   /purchase-requests [post]
func (h *Handler) createPurchaseRequest(c *gin.Context) {
	var in model.PurchaseRequestInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.purchases.Create(c.Request.Context(), currentUserID(c), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updatePurchaseRequestStatus godoc
// @Summary  Сменить статус заявки
// @Description Доступно только администратору
// @Tags     purchase-requests
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int                       true  "ID заявки"
// @Param    input  body      model.PurchaseStatusInput true  "Новый статус"
// @Success  200    {object}  model.PurchaseRequest
// @Failure  400    {object}  errorResponse
// @Failure  404    {object}  errorResponse
// @Router   /purchase-requests/{id}/status [patch]
func (h *Handler) updatePurchaseRequestStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.PurchaseStatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	item, err := h.purchases.UpdateStatus(c.Request.Context(), id, in.Status)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// deletePurchaseRequest godoc
// @Summary  Удалить заявку
// @Tags     purchase-requests
// @Security BearerAuth
// @Param    id   path  int  true  "ID заявки"
// @Success  204  "No Content"
// @Failure  403  {object}  errorResponse
// @Failure  404  {object}  errorResponse
// @Router   /purchase-requests/{id} [delete]
func (h *Handler) deletePurchaseRequest(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.purchases.Delete(c.Request.Context(), id, currentUserID(c), isAdmin(c)); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

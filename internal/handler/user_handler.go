package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// listUsers godoc
// @Summary  Список пользователей
// @Tags     users
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}   model.User
// @Failure  401  {object}  errorResponse
// @Failure  403  {object}  errorResponse
// @Router   /users [get]
func (h *Handler) listUsers(c *gin.Context) {
	users, err := h.users.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// getUser godoc
// @Summary  Получить пользователя
// @Tags     users
// @Produce  json
// @Security BearerAuth
// @Param    id   path      int  true  "ID пользователя"
// @Success  200  {object}  model.User
// @Failure  404  {object}  errorResponse
// @Router   /users/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	u, err := h.users.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

// createUser godoc
// @Summary  Создать пользователя
// @Tags     users
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input  body      model.RegisterInput  true  "Данные пользователя (можно задать роль)"
// @Success  201    {object}  model.User
// @Failure  400    {object}  errorResponse
// @Failure  409    {object}  errorResponse
// @Router   /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var in model.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	u, err := h.users.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}

// updateUser godoc
// @Summary  Обновить пользователя
// @Description Меняет пароль и/или роль
// @Tags     users
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path      int                   true  "ID пользователя"
// @Param    input  body      model.UpdateUserInput true  "Поля для обновления"
// @Success  200    {object}  model.User
// @Failure  404    {object}  errorResponse
// @Router   /users/{id} [put]
func (h *Handler) updateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in model.UpdateUserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	u, err := h.users.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

// deleteUser godoc
// @Summary  Удалить пользователя
// @Tags     users
// @Security BearerAuth
// @Param    id   path  int  true  "ID пользователя"
// @Success  204  "No Content"
// @Failure  404  {object}  errorResponse
// @Router   /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.users.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

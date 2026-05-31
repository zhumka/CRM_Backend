package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя с ролью user и возвращает JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      model.RegisterInput  true  "Данные регистрации"
// @Success      201    {object}  model.AuthResponse
// @Failure      400    {object}  errorResponse
// @Failure      409    {object}  errorResponse
// @Router       /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var in model.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	// Через публичную регистрацию роль admin назначить нельзя.
	in.Role = model.RoleUser

	resp, err := h.auth.Register(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// login godoc
// @Summary      Вход в систему
// @Description  Проверяет учётные данные и возвращает JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      model.LoginInput  true  "Логин и пароль"
// @Success      200    {object}  model.AuthResponse
// @Failure      400    {object}  errorResponse
// @Failure      401    {object}  errorResponse
// @Router       /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var in model.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c, err)
		return
	}
	resp, err := h.auth.Login(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// me godoc
// @Summary      Текущий пользователь
// @Description  Возвращает данные авторизованного пользователя
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.User
// @Failure      401  {object}  errorResponse
// @Router       /me [get]
func (h *Handler) me(c *gin.Context) {
	u, err := h.users.Get(c.Request.Context(), currentUserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
)

// errorResponse — единый формат ошибки API.
type errorResponse struct {
	Error string `json:"error"`
}

// respondError сопоставляет доменную ошибку с HTTP-статусом.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, model.ErrAlreadyExists):
		c.JSON(http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, model.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
	case errors.Is(err, model.ErrForbidden):
		c.JSON(http.StatusForbidden, errorResponse{Error: err.Error()})
	case errors.Is(err, model.ErrInsufficientStock),
		errors.Is(err, model.ErrRequestNotApproved),
		errors.Is(err, model.ErrInvoiceAmountExceeded):
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

// badRequest отвечает 400 с текстом ошибки.
func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
}

// parseID извлекает числовой параметр пути :id.
func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		badRequest(c, errors.New("invalid id"))
		return 0, false
	}
	return id, true
}

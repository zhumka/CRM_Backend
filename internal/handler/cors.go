package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// corsMiddleware разрешает кросс-доменные запросы от фронтенда.
// Origin отражается обратно (а не "*"), чтобы при необходимости работали
// и запросы с конкретного домена. Аутентификация — через Bearer-токен в
// заголовке (не cookie), поэтому Allow-Credentials не требуется.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crm/internal/model"
	"crm/internal/pkg/jwtutil"
)

const (
	ctxUserID = "userID"
	ctxRole   = "role"
)

// authMiddleware проверяет JWT из заголовка Authorization: Bearer <token>.
func authMiddleware(jwt *jwtutil.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{Error: "missing bearer token"})
			return
		}

		claims, err := jwt.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{Error: "invalid token"})
			return
		}

		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

// adminOnly разрешает доступ только администраторам.
func adminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if currentRole(c) != model.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{Error: "admin access required"})
			return
		}
		c.Next()
	}
}

func currentUserID(c *gin.Context) int {
	if v, ok := c.Get(ctxUserID); ok {
		if id, ok := v.(int); ok {
			return id
		}
	}
	return 0
}

func currentRole(c *gin.Context) string {
	if v, ok := c.Get(ctxRole); ok {
		if r, ok := v.(string); ok {
			return r
		}
	}
	return ""
}

func isAdmin(c *gin.Context) bool {
	return currentRole(c) == model.RoleAdmin
}

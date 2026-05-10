package middleware

import (
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员权限中间件，必须在 JWTAuth 之后使用
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

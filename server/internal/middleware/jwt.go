package middleware

import (
	"net/http"
	"strings"

	"devprompt-ai/internal/response"
	"devprompt-ai/internal/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件，从 Authorization header 提取并验证 token
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(secret, tokenStr)
		if err != nil {
			response.Unauthorized(c, "token 已过期或无效")
			c.Abort()
			return
		}

		// 将用户信息注入 Gin context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// OptionalAuth 可选认证：有 token 时解析用户信息，没有 token 时继续
func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(secret, tokenStr)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
		}
		c.Next()
	}
}

// GetUserID 从 context 安全获取 user_id
func GetUserID(c *gin.Context) (uint64, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

// Cors CORS 跨域中间件
func Cors(webURL, adminURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := false

		for _, u := range []string{webURL, adminURL, "http://localhost:3000", "http://localhost:5173"} {
			if origin == u || u == "*" {
				allowed = true
				break
			}
		}
		// 开发模式宽松处理
		if !allowed && strings.HasPrefix(origin, "http://localhost") {
			allowed = true
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

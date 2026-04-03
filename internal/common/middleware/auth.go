package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
	appjwt "goal-planner/internal/infra/jwt"
)

// AuthMiddleware 校验请求头中的 JWT，并把用户信息写入上下文。
func AuthMiddleware(manager *appjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		tokenText := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenText == "" {
			response.Fail(c, http.StatusUnauthorized, "无效的认证信息")
			c.Abort()
			return
		}

		claims, err := manager.ParseToken(tokenText)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "token无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

package middleware

import (
	"doudian/internal/config"
	"doudian/internal/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired 认证中间件，检查用户是否已登录
// 优先从 Authorization: Bearer <token> header 获取 token
// 降级方案：从 cookie 中获取 session
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 方案 1：从 Authorization header 获取 Bearer token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			tokenString = strings.TrimSpace(tokenString)

			claims, err := util.ParseToken(tokenString)
			if err == nil {
				// token 验证通过，设置用户信息到 context
				c.Set("user_id", strconv.FormatUint(uint64(claims.UserID), 10))
				c.Set("username", claims.Username)
				c.Next()
				return
			}
			// token 无效，返回 401（不降级到 cookie，因为客户端明确用了 Bearer 方式）
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "token 无效或已过期: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 方案 2：从 cookie 中获取 session（降级方案）
		session, err := c.Cookie("doudian_session")
		if err != nil || session == "" {
			// 检查是否是 AJAX 请求
			if isAjaxRequest(c) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "未登录或登录已过期",
				})
			} else {
				// 重定向到登录页（包含 secret path 前缀）
				cfg := config.Get()
				secretPath := strings.Trim(cfg.SecretPath, "/")
				loginPath := "/login"
				if secretPath != "" {
					loginPath = "/" + secretPath + "/login"
				}
				c.Redirect(http.StatusFound, loginPath)
			}
			c.Abort()
			return
		}

		// 简单验证 session 格式（实际项目中可使用 JWT 或 Redis session）
		// 这里我们假设 session 中存的是 user_id:username
		if !strings.Contains(session, ":") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的会话",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		parts := strings.SplitN(session, ":", 2)
		if len(parts) == 2 {
			c.Set("user_id", parts[0])
			c.Set("username", parts[1])
		}

		c.Next()
	}
}

// SecretPath 中间件：检查路径前缀是否匹配 SecretPath
func SecretPath() gin.HandlerFunc {
	cfg := config.Get()
	secretPath := cfg.SecretPath

	return func(c *gin.Context) {
		// 如果没有配置 SecretPath，直接放行
		if secretPath == "" {
			c.Next()
			return
		}

		// 获取请求路径
		path := c.Request.URL.Path

		// 检查路径是否以 secretPath 开头
		expectedPrefix := "/" + strings.Trim(secretPath, "/")
		if !strings.HasPrefix(path, expectedPrefix) {
			// 返回 404，伪装路径不存在
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Not Found",
			})
			c.Abort()
			return
		}

		// 剥离 secret path 前缀，供后续路由处理
		// 注意：Gin 的路由是在注册时确定的，这里我们通过修改 URL Path 来实现
		// 但更简单的方式是在路由注册时就加上前缀
		c.Next()
	}
}

// isAjaxRequest 判断是否为 AJAX 请求
func isAjaxRequest(c *gin.Context) bool {
	return c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
		c.GetHeader("Content-Type") == "application/json" ||
		strings.HasPrefix(c.Request.URL.Path, "/api/")
}

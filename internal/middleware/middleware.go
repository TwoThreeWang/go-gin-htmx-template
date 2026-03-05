package middleware

import (
	"github.com/gin-gonic/gin"
)

// Security 安全头中间件
// 添加常见的安全响应头，防止XSS、点击劫持等攻击
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")           // 防止点击劫持
		c.Header("X-Content-Type-Options", "nosniff") // 防止MIME类型嗅探
		c.Header("X-XSS-Protection", "1; mode=block") // 启用XSS保护
		c.Next()
	}
}

// CORS 跨域资源共享中间件
// 允许前端跨域访问API
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源（生产环境应限制具体域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头（包含HTMX特有的头）
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, HX-Request, HX-Current-URL, Hx-Target")

		// 预检请求直接返回204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

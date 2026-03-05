package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明结构体
// 包含用户ID、用户名、邮箱等基本信息
type Claims struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Role     string `json:"role"`     // 用户角色
	jwt.RegisteredClaims
}

// RequireAuth 必须登录中间件
// 验证JWT Token，未登录用户会被重定向到登录页或返回401错误
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie或Header中提取Claims
		claims, err := extractClaims(c, jwtSecret)
		if err != nil {
			// 如果是页面请求，重定向到登录页
			if strings.Contains(c.GetHeader("Accept"), "text/html") {
				c.Redirect(http.StatusFound, "/auth/login?redirect="+c.Request.URL.Path)
				c.Abort()
				return
			}
			// API请求返回401未授权
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续处理器使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// 滑动续期逻辑：如果Token过期时间消耗超过一半，则刷新
		if shouldRefresh(claims) {
			tokenExpiry := claims.RegisteredClaims.ExpiresAt.Sub(claims.RegisteredClaims.IssuedAt.Time)
			newToken, err := GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role, jwtSecret, tokenExpiry)
			if err == nil {
				// 设置新的Cookie，HttpOnly=true防止XSS攻击
				c.SetCookie("token", newToken, int(tokenExpiry.Seconds()), "/", "", false, true)
			}
		}

		c.Next()
	}
}

// OptionalAuth 可选登录中间件
// 不强制要求登录，但如果用户已登录则提取用户信息并执行滑动续期
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractClaims(c, jwtSecret)
		if err == nil {
			// 用户已登录，存入上下文
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			// 滑动续期逻辑
			if shouldRefresh(claims) {
				tokenExpiry := claims.RegisteredClaims.ExpiresAt.Sub(claims.RegisteredClaims.IssuedAt.Time)
				newToken, err := GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role, jwtSecret, tokenExpiry)
				if err == nil {
					c.SetCookie("token", newToken, int(tokenExpiry.Seconds()), "/", "", false, true)
				}
			}
		}
		c.Next()
	}
}

// RequireAdmin 管理员权限中间件
// 验证当前用户是否为管理员角色
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// extractClaims 从Cookie或Authorization Header中提取JWT Claims
// 优先从Cookie获取，其次从Bearer Token获取
func extractClaims(c *gin.Context, jwtSecret string) (*Claims, error) {
	var tokenString string

	// 优先从Cookie获取Token
	if cookie, err := c.Cookie("token"); err == nil {
		tokenString = cookie
	} else {
		// 从Authorization Header获取Bearer Token
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		return nil, jwt.ErrTokenMalformed
	}

	// 解析并验证Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// GetUserID 从上下文获取用户ID（未登录返回0）
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从上下文获取用户名（未登录返回空字符串）
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// IsLoggedIn 判断用户是否已登录
func IsLoggedIn(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// IsAdmin 判断当前用户是否为管理员
func IsAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	return exists && role == "admin"
}

// GenerateToken 生成JWT Token
// 参数: userID-用户ID, username-用户名, email-邮箱, role-角色, jwtSecret-密钥, expiry-过期时长
func GenerateToken(userID uint, username, email, role, jwtSecret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// shouldRefresh 判断是否需要刷新Token
// 滑动续期逻辑：如果已经消耗了总有效期的50%以上，则建议刷新
func shouldRefresh(claims *Claims) bool {
	if claims.ExpiresAt == nil || claims.IssuedAt == nil {
		return false
	}

	totalDuration := claims.ExpiresAt.Sub(claims.IssuedAt.Time)
	elapsedDuration := time.Since(claims.IssuedAt.Time)

	// 消耗超过50%时刷新
	return elapsedDuration > totalDuration/2
}

package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// Cache 全局缓存实例
var Cache *cache.Cache

// InitCache 初始化缓存
// 默认缓存5分钟，每10分钟清理过期条目
func InitCache() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

// Response 统一API响应结构
// 用于规范所有API接口的返回格式
type Response struct {
	Code    int         `json:"code"`    // HTTP状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
	Success bool        `json:"success"` // 是否成功
}

// Success 返回成功响应
// 自动设置状态码200和success为true
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
		Success: true,
	})
}

// SuccessWithMessage 返回成功响应并自定义消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: message,
		Data:    data,
		Success: true,
	})
}

// Error 返回错误响应
// 参数: code-HTTP状态码, message-错误消息
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Success: false,
	})
}

// BadRequest 返回400错误（客户端请求错误）
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "请求参数错误"
	}
	Error(c, 400, message)
}

// Unauthorized 返回401错误（未授权/未登录）
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未登录"
	}
	Error(c, 401, message)
}

// Forbidden 返回403错误（禁止访问/无权限）
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "无权限访问"
	}
	Error(c, 403, message)
}

// NotFound 返回404错误（资源不存在）
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Error(c, 404, message)
}

// InternalServerError 返回500错误（服务器内部错误）
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	Error(c, 500, message)
}

package config

import (
	"os"
	"strconv"
)

// Config 应用配置结构体
// 包含服务器运行所需的所有配置项
type Config struct {
	Port        string // 服务端口
	Env         string // 运行环境（development/production）
	AppSecret   string // 应用密钥（用于JWT签名）
	TimeZone    string // 时区设置
	DatabaseURL string // 数据库连接URL
}

// Load 加载配置
// 从环境变量读取配置，未设置则使用默认值
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "5007"),
		Env:         getEnv("ENV", "development"),
		AppSecret:   getEnv("APP_SECRET", "your-secret-key-change-in-production"),
		TimeZone:    getEnv("TIME_ZONE", "Asia/Shanghai"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ginhtmx?sslmode=disable"),
	}
}

// SiteName 获取站点名称
func (c *Config) SiteName() string {
	return getEnv("SITE_NAME", "Gin HTMX Template")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整型环境变量，如果不存在或解析失败则返回默认值
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}
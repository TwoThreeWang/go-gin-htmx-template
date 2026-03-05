package router

import (
	"gin-htmx-template/internal/handler"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
// 注意：全局中间件、静态文件服务、模板加载已在main.go中注册
func RegisterRoutes(r *gin.Engine, h *handler.Handler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// ==================== 公开页面 ====================
	r.GET("/", h.Home)       // 首页
	r.GET("/about", h.About) // 关于页面

	// ==================== 用户管理 ====================
	r.GET("/users", h.ListUsers)          // 用户列表页面
	r.GET("/users/list", h.GetUsers)      // 用户列表片段（HTMX）
	r.POST("/users", h.CreateUser)        // 创建用户
	r.DELETE("/users/:id", h.DeleteUser)  // 删除用户
	r.GET("/users/count", h.GetUserCount) // 用户数量

	// ==================== 联系表单 ====================
	r.GET("/contact", h.ContactForm)    // 联系表单页面
	r.POST("/contact", h.SubmitContact) // 提交联系表单

	// ==================== API接口 ====================
	api := r.Group("/api")
	{
		api.GET("/users", h.GetAPIUsers) // 获取用户列表API
	}

	// ==================== 404 处理 ====================
	r.NoRoute(h.NotFound)
}

// LoadTemplates 使用multitemplate加载模板，解决模板继承问题
// 支持自定义模板函数，如dict、default、js、add、sub等
func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 获取布局模板文件列表
	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err)
	}

	// 获取局部模板文件列表（可复用的组件）
	partials, err := filepath.Glob(templatesDir + "/partials/*.html")
	if err != nil {
		panic(err)
	}

	// assemble 组装模板文件列表
	// 将布局、局部模板和页面模板合并成一个完整的模板
	assemble := func(view string) []string {
		files := make([]string, 0)
		files = append(files, layouts...)
		files = append(files, partials...)
		files = append(files, view)
		return files
	}

	// 自定义模板函数
	funcMap := template.FuncMap{
		// default 默认值函数，当值为空时返回默认值
		"default": func(defaultValue, value interface{}) interface{} {
			switch v := value.(type) {
			case string:
				if v == "" {
					return defaultValue
				}
			case int:
				if v == 0 {
					return defaultValue
				}
			case nil:
				return defaultValue
			}
			return value
		},
		// js 将字符串转换为template.JS类型，用于安全地嵌入JavaScript代码
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		// add 加法运算
		"add": func(a, b int) int {
			return a + b
		},
		// sub 减法运算
		"sub": func(a, b int) int {
			return a - b
		},
		// mul 乘法运算
		"mul": func(a, b int) int {
			return a * b
		},
		// div 整除运算
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		// contains 检查字符串是否包含子串
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		// lower 转小写
		"lower": strings.ToLower,
		// upper 转大写
		"upper": strings.ToUpper,
		// trim 去除首尾空白
		"trim": strings.TrimSpace,
	}

	// 注册所有页面模板
	pages, err := filepath.Glob(templatesDir + "/pages/*.html")
	if err != nil {
		panic(err)
	}
	for _, page := range pages {
		pageName := strings.TrimSuffix(filepath.Base(page), ".html")
		log.Printf("Registering template: %s", pageName)
		r.AddFromFilesFuncs(pageName, funcMap, assemble(page)...)
	}

	// 注册局部模板（用于HTMX局部更新）
	for _, partial := range partials {
		name := "partials/" + strings.TrimSuffix(filepath.Base(partial), ".html")
		// 构建文件列表：当前partial放第一位（作为主入口），其他partials随后（作为依赖）
		files := make([]string, 0, len(partials))
		files = append(files, partial)
		for _, p := range partials {
			if p != partial {
				files = append(files, p)
			}
		}
		r.AddFromFilesFuncs(name, funcMap, files...)
	}

	return r
}

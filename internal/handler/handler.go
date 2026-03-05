package handler

import (
	"fmt"
	"gin-htmx-template/internal/config"
	"gin-htmx-template/internal/middleware"
	"gin-htmx-template/internal/model"
	"gin-htmx-template/internal/repository"
	"gin-htmx-template/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler HTTP处理器结构体
// 封装了Repository和Config，提供统一的处理入口
type Handler struct {
	repos *repository.Repositories // 数据仓库
	cfg   *config.Config           // 配置信息
}

// NewHandler 创建Handler实例
func NewHandler(repos *repository.Repositories, cfg *config.Config) *Handler {
	return &Handler{
		repos: repos,
		cfg:   cfg,
	}
}

// RenderData 统一封装公共渲染数据
// 自动注入站点配置、用户信息、菜单高亮等公共数据
func (h *Handler) RenderData(c *gin.Context, data gin.H) gin.H {
	// 基础数据
	res := gin.H{
		"SiteName": h.cfg.SiteName(), // 站点名称
		"Path":     c.Request.URL.Path,
		"FullPath": c.Request.RequestURI,
		"Referer":  c.Request.Referer(),
	}

	// 注入用户信息（如果已登录）
	userID := middleware.GetUserID(c)
	if userID > 0 {
		res["UserID"] = userID
		res["Username"] = middleware.GetUsername(c)
		res["IsLoggedIn"] = true
		res["IsAdmin"] = middleware.IsAdmin(c)
	} else {
		res["IsLoggedIn"] = false
		res["IsAdmin"] = false
	}

	// 菜单高亮逻辑
	res["ActiveMenu"] = h.getActiveMenu(c)

	// 合并传入的数据
	for k, v := range data {
		res[k] = v
	}

	return res
}

// getActiveMenu 根据路径判断当前高亮菜单
func (h *Handler) getActiveMenu(c *gin.Context) string {
	path := c.Request.URL.Path
	switch {
	case path == "/":
		return "home"
	case path == "/about":
		return "about"
	case len(path) >= 6 && path[:6] == "/users":
		return "users"
	case len(path) >= 8 && path[:8] == "/contact":
		return "contact"
	default:
		return ""
	}
}

// Home 首页处理器
func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home", h.RenderData(c, gin.H{
		"Title": "首页",
	}))
}

// GetUsers 获取用户列表片段（用于HTMX局部更新）
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.repos.User.FindAll()
	if err != nil {
		c.String(http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	c.HTML(http.StatusOK, "partials/users-list", h.RenderData(c, gin.H{
		"Users": users,
	}))
}

// ListUsers 用户列表页面
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.repos.User.FindAll()
	if err != nil {
		c.String(http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	count, _ := h.repos.User.Count()

	c.HTML(http.StatusOK, "users", h.RenderData(c, gin.H{
		"Title": "用户列表",
		"Users": users,
		"Count": count,
	}))
}

// CreateUserInput 创建用户输入结构体
type CreateUserInput struct {
	Username string `form:"username" binding:"required,min=3,max=50"` // 用户名（3-50字符）
	Email    string `form:"email" binding:"required,email"`           // 邮箱（必须为有效邮箱格式）
}

// CreateUser 创建用户处理器
func (h *Handler) CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "partials/users-list", h.RenderData(c, gin.H{
			"Error": err.Error(),
		}))
		return
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Active:   true,
	}

	if err := h.repos.User.Create(user); err != nil {
		c.HTML(http.StatusInternalServerError, "partials/users-list", h.RenderData(c, gin.H{
			"Error": "创建用户失败",
		}))
		return
	}

	users, _ := h.repos.User.FindAll()

	c.HTML(http.StatusOK, "partials/users-list", h.RenderData(c, gin.H{
		"Users":   users,
		"Success": "用户创建成功！",
	}))
}

// DeleteUser 删除用户处理器
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var userID uint
	if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
		c.String(http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := h.repos.User.Delete(userID); err != nil {
		c.String(http.StatusInternalServerError, "删除用户失败")
		return
	}

	users, _ := h.repos.User.FindAll()

	c.HTML(http.StatusOK, "partials/users-list", h.RenderData(c, gin.H{
		"Users":   users,
		"Success": "用户删除成功！",
	}))
}

// GetUserCount 获取用户数量
func (h *Handler) GetUserCount(c *gin.Context) {
	count, _ := h.repos.User.Count()
	c.String(http.StatusOK, "%d", count)
}

// GetAPIUsers 获取用户列表API（返回JSON格式）
func (h *Handler) GetAPIUsers(c *gin.Context) {
	users, err := h.repos.User.FindAll()
	if err != nil {
		utils.InternalServerError(c, "获取用户列表失败")
		return
	}
	utils.Success(c, gin.H{
		"users": users,
		"count": len(users),
	})
}

// ContactForm 联系表单页面
func (h *Handler) ContactForm(c *gin.Context) {
	c.HTML(http.StatusOK, "contact", h.RenderData(c, gin.H{
		"Title": "联系我们",
	}))
}

// ContactInput 联系表单输入结构体
type ContactInput struct {
	Name    string `form:"name" binding:"required"`           // 姓名（必填）
	Email   string `form:"email" binding:"required,email"`    // 邮箱（必填，有效邮箱格式）
	Message string `form:"message" binding:"required,min=10"` // 消息内容（必填，最少10字符）
}

// SubmitContact 提交联系表单处理器
func (h *Handler) SubmitContact(c *gin.Context) {
	var input ContactInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "partials/contact-form", h.RenderData(c, gin.H{
			"Error": err.Error(),
		}))
		return
	}

	c.HTML(http.StatusOK, "partials/contact-form", h.RenderData(c, gin.H{
		"Success": "感谢您的留言，我们会尽快回复！",
		"Name":    input.Name,
	}))
}

// About 关于页面处理器
func (h *Handler) About(c *gin.Context) {
	c.HTML(http.StatusOK, "about", h.RenderData(c, gin.H{
		"Title": "关于我们",
	}))
}

// NotFound 404 页面处理器
func (h *Handler) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404", h.RenderData(c, gin.H{
		"Title": "页面未找到 - " + h.cfg.SiteName(),
	}))
}

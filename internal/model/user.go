package model

import (
	"gorm.io/gorm"
)

// User 用户模型
// 存储用户基本信息
type User struct {
	gorm.Model
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名（唯一）
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`   // 邮箱（唯一）
	Password string `gorm:"size:255" json:"-"`                            // 密码（加密存储，不返回给前端）
	Role     string `gorm:"size:20;default:user" json:"role"`             // 角色（user/admin）
	Active   bool   `gorm:"default:true" json:"active"`                   // 是否激活
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// SessionUser 会话用户信息
// 用于存储在JWT Token或Session中的用户信息
type SessionUser struct {
	ID       uint   `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Role     string `json:"role"`     // 角色
	Active   bool   `json:"active"`   // 是否激活
}

// IsAdmin 判断用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActive 判断用户是否已激活
func (u *User) IsActive() bool {
	return u.Active
}

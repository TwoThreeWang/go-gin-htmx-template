package repository

import (
	"gin-htmx-template/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repositories 数据仓库集合
// 统一管理所有数据仓库实例
type Repositories struct {
	User *UserRepository // 用户数据仓库
}

// NewRepositories 创建数据仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}

// UserRepository 用户数据仓库
// 提供用户相关的数据库操作
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindAll 查询所有活跃用户
func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Where("active = ?", true).Find(&users).Error
	return users, err
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查询用户
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 软删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// Count 统计活跃用户数量
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("active = ?", true).Count(&count).Error
	return count, err
}

// InitDB 初始化数据库连接
// 参数: databaseURL - PostgreSQL数据库连接字符串
func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
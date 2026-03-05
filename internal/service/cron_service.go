package service

import (
	"gin-htmx-template/internal/repository"
	"log"
	"time"
)

// CronTaskService 定时任务服务
// 用于执行定期任务
type CronTaskService struct {
	repos *repository.Repositories
}

// NewCronTaskService 创建定时任务服务
func NewCronTaskService(repos *repository.Repositories) *CronTaskService {
	return &CronTaskService{repos: repos}
}

// Start 启动定时任务服务
// 每24小时执行一次任务
func (s *CronTaskService) Start() {
	go func() {
		// 创建定时器，每24小时触发一次
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("正在执行定时任务...")
			// 在这里添加任务逻辑
			// 例如：清理过期的会话、日志、临时文件等
		}
	}()
}

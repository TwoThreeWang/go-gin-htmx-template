package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-htmx-template/internal/config"
	"gin-htmx-template/internal/handler"
	"gin-htmx-template/internal/middleware"
	"gin-htmx-template/internal/repository"
	"gin-htmx-template/internal/router"
	"gin-htmx-template/internal/service"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载.env环境变量文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 加载应用配置
	cfg := config.Load()

	// 设置全局时区
	if loc, err := time.LoadLocation(cfg.TimeZone); err == nil {
		time.Local = loc
	} else {
		log.Printf("加载时区 %s 失败，使用系统默认时区: %v", cfg.TimeZone, err)
	}

	// 初始化数据库连接
	db, err := repository.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 获取底层数据库连接，用于延迟关闭
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// 创建数据仓库集合
	repos := repository.NewRepositories(db)

	// 生产环境设置Gin为发布模式
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 注册全局中间件
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // Gzip压缩
	r.Use(middleware.Security())              // 安全头
	r.Use(middleware.CORS())                  // 跨域支持

	// 加载HTML模板（使用multitemplate支持模板继承）
	r.HTMLRender = router.LoadTemplates("./web/templates")

	// 静态文件服务
	r.Static("/static", "./web/static")

	// 创建HTTP处理器
	h := handler.NewHandler(repos, cfg)

	// 启动定时任务服务
	cronTaskSvc := service.NewCronTaskService(repos)
	cronTaskSvc.Start()

	// 注册路由
	router.RegisterRoutes(r, h)

	// 获取服务端口
	port := cfg.Port
	if port == "" {
		port = "5007"
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second, // 读取超时
		WriteTimeout:   10 * time.Second, // 写入超时
		MaxHeaderBytes: 1 << 20,          // 最大请求头大小（1MB）
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("服务器启动于 http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 监听系统信号（优雅关闭）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 设置5秒超时的优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已退出")
}

package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"goal-planner/internal/common/response"
	"goal-planner/internal/config"
	"goal-planner/internal/goal"
	"goal-planner/internal/infra/db"
	"goal-planner/internal/infra/logger"
)

func main() {
	// 启动 API 服务前，先读取运行配置。
	cfg := config.Load()

	// 初始化项目统一日志。
	log := logger.New()

	// 先尝试连接 MySQL，确认数据库配置可用。
	database, err := db.NewMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Error("mysql connect failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	log.Info("mysql connected")

	// 创建 Gin 路由引擎，并开启请求日志和异常恢复中间件。
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 健康检查接口，后续本地调试和部署探活都会用到。
	router.GET("/healthz", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"service": "api",
		})
	})

	// 注册目标模块路由。
	goalHandler := goal.NewHandler(database)
	goalHandler.RegisterRoutes(router)

	log.Info("api server starting", "addr", cfg.HTTPAddr)

	// 启动 HTTP 服务；如果启动失败则记录日志并退出程序。
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Error("api server stopped", "error", err)
		os.Exit(1)
	}
}

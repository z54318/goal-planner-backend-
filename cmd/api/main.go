package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "goal-planner/docs"
	"goal-planner/internal/auth"
	"goal-planner/internal/common/middleware"
	"goal-planner/internal/common/response"
	"goal-planner/internal/config"
	"goal-planner/internal/goal"
	appai "goal-planner/internal/infra/ai"
	"goal-planner/internal/infra/db"
	appjwt "goal-planner/internal/infra/jwt"
	"goal-planner/internal/infra/logger"
	"goal-planner/internal/menu"
	"goal-planner/internal/phase"
	"goal-planner/internal/plan"
	"goal-planner/internal/rbac"
	"goal-planner/internal/task"
	"goal-planner/internal/user"
)

// @title Goal Planner Backend API
// @version 1.0
// @description Goal Planner 后端接口文档
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 启动 API 服务前，先读取运行配置。
	cfg := config.Load()

	// 初始化项目统一日志。
	log := logger.New()
	slog.SetDefault(log)

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

	// 初始化 JWT 管理器，后续登录和鉴权都会复用。
	jwtManager := appjwt.NewManager(cfg.JWTSecret)
	aiClient := appai.NewClient(cfg.AIAPIKey, cfg.AIBaseURL, cfg.AIModel)

	// 健康检查接口，后续本地调试和部署探活都会用到。
	router.GET("/healthz", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"service": "api",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 注册认证模块路由。
	authHandler := auth.NewHandler(database, jwtManager)
	authHandler.RegisterRoutes(router)

	// 提供一个最小受保护接口，用来验证 JWT 中间件是否生效。
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	protected.GET("/auth/profile", authHandler.Profile)
	protected.GET("/auth/menus", authHandler.Menus)

	// 注册目标模块受保护路由。
	goalHandler := goal.NewHandler(database)
	goalHandler.RegisterProtectedRoutes(protected)

	// 注册计划模块受保护路由。
	planHandler := plan.NewHandler(database, aiClient)
	planHandler.RegisterProtectedRoutes(protected)

	// 注册阶段模块受保护路由。
	phaseHandler := phase.NewHandler(database, aiClient)
	phaseHandler.RegisterProtectedRoutes(protected)

	// 注册任务模块受保护路由。
	taskHandler := task.NewHandler(database, aiClient)
	taskHandler.RegisterProtectedRoutes(protected)

	// 注册菜单管理受保护路由。
	menuHandler := menu.NewHandler(database)
	menuHandler.RegisterProtectedRoutes(protected)

	// 注册 RBAC 管理受保护路由。
	rbacHandler := rbac.NewHandler(database)
	rbacHandler.RegisterProtectedRoutes(protected)

	// 注册用户管理受保护路由。
	userHandler := user.NewHandler(database)
	userHandler.RegisterProtectedRoutes(protected)

	log.Info("api server starting", "addr", cfg.HTTPAddr)

	// 启动 HTTP 服务；如果启动失败则记录日志并退出程序。
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Error("api server stopped", "error", err)
		os.Exit(1)
	}
}

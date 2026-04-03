package main

import (
	"time"

	"goal-planner/internal/config"
	"goal-planner/internal/infra/logger"
)

func main() {
	// 启动 worker 前，先读取运行配置。
	cfg := config.Load()

	// 初始化项目统一日志。
	log := logger.New()
	log.Info("worker started", "name", cfg.WorkerName)

	// 先用心跳日志模拟后台常驻进程，确认 worker 能持续运行。
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Info("worker heartbeat")
	}
}

package logger

import (
	"log/slog"
	"os"
)

// New 创建项目统一使用的结构化日志实例。
func New() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

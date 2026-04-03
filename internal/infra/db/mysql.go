package db

import (
	"database/sql"
	"fmt"
	"time"

	// 引入 MySQL 驱动，让 database/sql 知道如何连接 MySQL。
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQL 根据传入的 DSN 创建数据库连接。
func NewMySQL(dsn string) (*sql.DB, error) {
	// sql.Open 不会立刻真正连接数据库，它只是先创建一个连接对象。
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("创建 MySQL 连接对象失败: %w", err)
	}

	// 设置连接池参数，当前先用比较保守的默认值。
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(time.Hour)

	// 通过 Ping 主动测试数据库是否真的可用。
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	return database, nil
}

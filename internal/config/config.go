package config

import "os"

// Config 用于保存从环境变量读取到的运行配置。
type Config struct {
	AppEnv     string
	HTTPAddr   string
	WorkerName string

	MySQLDSN  string
	RedisAddr string
	JWTSecret string
	AIAPIKey  string
}

// Load 读取项目需要的环境变量，并在未设置时使用默认值。
func Load() Config {
	appEnv := getEnv("APP_ENV", "dev")
	httpAddr := getEnv("HTTP_ADDR", ":8080")
	workerName := getEnv("WORKER_NAME", "ai-worker")

	mysqlDSN := getEnv("MYSQL_DSN", "")
	redisAddr := getEnv("REDIS_ADDR", "")
	jwtSecret := getEnv("JWT_SECRET", "")
	aiAPIKey := getEnv("AI_API_KEY", "")

	return Config{
		AppEnv:     appEnv,
		HTTPAddr:   httpAddr,
		WorkerName: workerName,
		MySQLDSN:   mysqlDSN,
		RedisAddr:  redisAddr,
		JWTSecret:  jwtSecret,
		AIAPIKey:   aiAPIKey,
	}
}

// getEnv 在环境变量为空时返回备用值。
func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

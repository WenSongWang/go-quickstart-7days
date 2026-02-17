package config

import (
	"os"
	"strconv"
)

// Config Day7 应用配置：端口、可选 API Key（空则不做鉴权）
type Config struct {
	HTTPPort int
	APIKey   string
}

// Load 从环境变量读取；HTTP_PORT 默认 8080，API_KEY 空表示不校验
func Load() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	return &Config{
		HTTPPort: port,
		APIKey:   os.Getenv("API_KEY"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

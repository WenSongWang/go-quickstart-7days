package config

import (
	"os"
	"strconv"
	"time"
)

// Config Redis + JWT 项目配置（进阶项目建议用 viper；此处先用最小可用实现）
type Config struct {
	HTTPPort int

	// RedisAddr 为空则不启用 Redis（方便没有 Redis 时也能跑通基础接口）
	RedisAddr string

	// JWTSecret 用于签名；生产请使用强随机值并通过环境变量注入
	JWTSecret string

	// JWTTTL 访问 token 的过期时间
	JWTTTL time.Duration
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	ttlSeconds, _ := strconv.Atoi(getEnv("JWT_TTL_SECONDS", "3600"))
	if ttlSeconds <= 0 {
		ttlSeconds = 3600
	}

	return &Config{
		HTTPPort:  port,
		RedisAddr: os.Getenv("REDIS_ADDR"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret"),
		JWTTTL:    time.Duration(ttlSeconds) * time.Second,
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}


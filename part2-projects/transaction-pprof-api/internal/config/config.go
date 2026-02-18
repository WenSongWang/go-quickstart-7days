package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort int
	SQLiteDSN string
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8081"))
	return &Config{
		HTTPPort:  port,
		SQLiteDSN: getEnv("SQLITE_DSN", "file:tx_demo?mode=memory&cache=shared"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}


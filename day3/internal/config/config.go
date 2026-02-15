// config 包：从环境变量读取配置（支持 viper 多源、类型解析），供 main 或其它包使用
package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config 应用配置结构体：启动时加载一次，后面各处用指针或值传递
// mapstructure 标签：Unmarshal 时 viper 用「viper 内部的 key」填到对应字段，这里约定用下划线小写（app_env → AppEnv）
type Config struct {
	AppEnv   string `mapstructure:"app_env"`   // 环境名，如 development / production
	HTTPPort int    `mapstructure:"http_port"` // HTTP 监听端口
	DBDSN    string `mapstructure:"db_dsn"`   // 数据库连接串，空表示不用数据库
}

// Load 从环境变量读取配置；没有设的用默认值。先用 viper 读 env（godotenv 已在 main 里把 .env 载入 env）
func Load() *Config {
	viper.SetDefault("app_env", "development")
	viper.SetDefault("http_port", 8080)
	viper.SetDefault("db_dsn", "file:./day4/data.db") // 默认 SQLite，无需安装数据库
	viper.AutomaticEnv() // 会按 key 转大写+下划线去读环境变量；下面 BindEnv 显式绑定「viper key ↔ 环境变量名」，与 .env 里 APP_ENV 等一致
	viper.BindEnv("app_env", "APP_ENV")   // viper 内部 key "app_env" 的值 ← 环境变量 APP_ENV
	viper.BindEnv("http_port", "HTTP_PORT")
	viper.BindEnv("db_dsn", "DB_DSN")

	var cfg Config
	// Unmarshal：把 viper 里 key 为 app_env/http_port/db_dsn 的值，按 mapstructure 标签填进 Config 字段
	_ = viper.Unmarshal(&cfg)
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 8080
	}
	if cfg.AppEnv == "" {//如果AppEnv为空，则设置为development
		cfg.AppEnv = "development"
	}
	return &cfg//&是取地址符，用于取结构体的地址
}

// LoadWithOS 仅用 os.Getenv 的简易版（不依赖 viper），便于对比理解。本目录 main 用的是 Load，此处不参与运行。
func LoadWithOS() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))//getEnv是getenv的函数，用于获取环境变量,比如HTTP_PORT映射到HTTPPort
	return &Config{
		AppEnv:   getEnv("APP_ENV", "development"),//getEnv是getenv的函数，用于获取环境变量,比如APP_ENV映射到AppEnv
		HTTPPort: port,
		DBDSN:    getEnv("DB_DSN", ""),//getEnv是getenv的函数，用于获取环境变量,比如DB_DSN映射到DBDSN
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {//getenv是getenv的函数，用于获取环境变量,比如HTTP_PORT映射到HTTPPort
		return v
	}
	return defaultVal
}

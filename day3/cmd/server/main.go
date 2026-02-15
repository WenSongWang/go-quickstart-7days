// Day3 示例：带配置的 HTTP 服务（项目结构：cmd + config）
// 本文件学习：.env、godotenv、配置结构体、用配置里的端口起服务
package main

// go.mod 中 module github.com/go-quickstart-7days 表示本项目的模块路径，包路径 = 模块路径 + 相对目录
// 所以本项目里的 day3/internal/config 的 import 路径就是 github.com/go-quickstart-7days/day3/internal/config
// 开发时 Go 从【本地目录】解析，不会从 GitHub 拉代码
import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-quickstart-7days/day3/internal/config" //配置加载逻辑
	"github.com/joho/godotenv"                             //从 .env 文件加载环境变量
)

func main() {
	// 可选：从 day3/config/.env 文件加载键值对到环境变量，本地开发方便
	_ = godotenv.Load("day3/config/.env") // _ = godotenv.Load("day3/config/.env") 是忽略返回值的赋值，因为godotenv.Load("day3/config/.env")返回值为error，所以用_忽略

	// 从环境变量读进配置结构体（端口、环境名等）
	cfg := config.Load()//这里是调用的config.go里的Load函数
	log.Printf("启动服务 env=%s port=%d", cfg.AppEnv, cfg.HTTPPort)

	// 注册 /health：返回 JSON 表示服务正常（健康检查）
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 注册 /info：返回当前配置（体现「配置注入」被用到）
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]interface{}{"env": cfg.AppEnv, "port": cfg.HTTPPort})
		w.Write(body)
	})

	// 用配置里的端口拼成 ":8080" 这种地址，再监听
	addr := ":" + strconv.Itoa(cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(addr, nil)) // Fatal 遇到错误会直接退出程序
}

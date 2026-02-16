// Day5 示例：中间件链（日志 + 简单 API Key 鉴权）
// 本文件学习：ServeMux、中间件包装 Handler、StripPrefix、链式调用
package main

import (
	"net/http"
	"time"

	"github.com/go-quickstart-7days/day5/internal/middleware"
)

func main() {
	// 用 NewServeMux 替代默认的 DefaultServeMux，方便挂子路由
	mux := http.NewServeMux()// 最外层, 可以挂多个路由
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})// /health 路由的处理函数在函数体内实现

	// 只对 /api/ 下的路径做鉴权：先 StripPrefix 去掉 /api，再交给带鉴权的 apiHandler
	api := http.NewServeMux()
	api.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"hello"}`))
	})
	// 链：先超时(5s)，再鉴权，再进 api；Timeout 给请求注入带取消的 context，下游可检测
	// 还有非链式调用方式：apiHandler := middleware.APIKeyAuth("secret")(middleware.Timeout(5 * time.Second)(api))
	// 进入 api 之前先经超时、再鉴权；"secret" 为鉴权中间件参数，生产环境应从配置（如 Day 3 的 cfg）读取
	apiHandler := middleware.Timeout(5 * time.Second)(middleware.APIKeyAuth("secret")(api))
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))//挂载api路由,并去掉/api前缀

	// 最外层再包一层日志：请求先打日志，再进 mux（再进超时、鉴权、业务）
	handler := middleware.Logging(mux)
	http.ListenAndServe(":8080", handler)
}

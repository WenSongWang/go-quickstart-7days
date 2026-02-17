// Day7 综合示例：配置 + 中间件 + 内存存储 + 用户 API + 优雅关闭（面试常问）
// 本文件学习：路由、中间件链、slog、Shutdown 优雅关闭、signal 监听
package main

import (
	"context"//上下文，可以传递取消与超时信号，常用来做超时，也可以用来做多个goroutine之间的同步
	"log/slog"//slog: 结构化日志，类似于Python的logging
	"net/http"//http: HTTP请求处理
	"os"//os: 操作系统相关的函数
	"os/signal"//os/signal: 信号处理
	"strconv"//strconv: 字符串转换
	"syscall"//syscall: 系统调用
	"time"//time: 时间

	"github.com/go-quickstart-7days/day7/internal/config"//config: 配置		
	"github.com/go-quickstart-7days/day7/internal/handler"//handler: 处理函数
	"github.com/go-quickstart-7days/day7/internal/middleware"//middleware: 中间件
	"github.com/go-quickstart-7days/day7/internal/store"//store: 存储
)

func main() {
	cfg := config.Load()//Load: 加载配置
	st := store.NewMemoryStore()                    // 内存版“数据库”，存用户
	userHandler := &handler.UserHandler{Store: st}  // 用户相关接口都交给 UserHandler，&符号表示取地址，即取变量的内存地址

	mux := http.NewServeMux() //NewServeMux: 创建一个多路复用器，用于处理HTTP请求
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {//HandleFunc: 注册一个处理函数，用于处理HTTP请求
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})//func(w http.ResponseWriter, r *http.Request): 处理函数，用于处理HTTP请求
	// /api/users：GET 列表，POST 创建
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {//switch: 选择语句, r.Method: 请求方法, case: 分支语句, http.MethodGet: GET请求, http.MethodPost: POST请求, default: 默认请求
		case http.MethodGet:
			userHandler.List(w, r)//List: 列表, userHandler: 用户处理函数, w: 响应写入器, r: 请求
		case http.MethodPost:
			userHandler.Create(w, r)//Create: 创建, userHandler: 用户处理函数, w: 响应写入器, r: 请求
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)//Method Not Allowed: 方法不允许, http.StatusMethodNotAllowed: 405状态码
		}
	})//func(w http.ResponseWriter, r *http.Request): 处理函数，用于处理HTTP请求
	// /api/users/1：GET 单个用户
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {//HandleFunc: 注册一个处理函数，用于处理HTTP请求
		if r.Method == http.MethodGet {//if: 条件语句, r.Method: 请求方法, http.MethodGet: GET请求
			userHandler.GetByID(w, r)//GetByID: 获取单个用户, userHandler: 用户处理函数, w: 响应写入器, r: 请求
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)//Method Not Allowed: 方法不允许, http.StatusMethodNotAllowed: 405状态码
	})//func(w http.ResponseWriter, r *http.Request): 处理函数，用于处理HTTP请求

	// 中间件链：先日志，若配置了 APIKey 再鉴权，最后进 mux
	chain := middleware.Logging(mux)
	if cfg.APIKey != "" {
		chain = middleware.APIKey(cfg.APIKey)(chain)//链
	}

	addr := ":" + strconv.Itoa(cfg.HTTPPort)
	srv := &http.Server{Addr: addr, Handler: chain}

	// 在 goroutine 里起服务，不阻塞 main，方便后面等信号
	go func() {
		slog.Info("server listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭：收到 Ctrl+C 或 kill 后，等正在处理的请求完成再退出（面试常问）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // 阻塞直到收到信号
	slog.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}

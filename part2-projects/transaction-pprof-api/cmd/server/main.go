// Part2 项目 2：事务 + pprof API（进阶）
// 真实使用 gin + SQLite（modernc）+ 事务示例 + pprof：
// - GET  /api/accounts
// - POST /api/transfer
package main

import (
	"fmt"
	"net/http/pprof"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"github.com/go-quickstart-7days/part2-projects/transaction-pprof-api/internal/config"
	"github.com/go-quickstart-7days/part2-projects/transaction-pprof-api/internal/handler"
	"github.com/go-quickstart-7days/part2-projects/transaction-pprof-api/internal/store"
)

func main() {
	cfg := config.Load()
	st, err := store.OpenSQLite(cfg.SQLiteDSN)
	if err != nil {
		panic(err)
	}
	defer st.DB.Close()

	healthHandler := &handler.HealthHandler{}
	accountHandler := &handler.AccountHandler{Store: st}

	r := gin.Default()
	r.GET("/health", healthHandler.Health)

	// 暴露 pprof：CPU、堆、goroutine 等（gin 包装标准库 pprof）
	// 使用：go tool pprof http://localhost:8081/debug/pprof/profile?seconds=5
	debug := r.Group("/debug/pprof")
	{
		debug.GET("/", gin.WrapF(pprof.Index))
		debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.GET("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
	}

	api := r.Group("/api")
	{
		api.GET("/accounts", accountHandler.List)
		api.POST("/transfer", accountHandler.Transfer)
	}

	addr := ":" + strconv.Itoa(cfg.HTTPPort)
	fmt.Println("listening on", addr)
	_ = r.Run(addr)
}

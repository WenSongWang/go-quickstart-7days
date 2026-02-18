package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logging 用 slog 打结构化日志（Go 1.21+ 标准库，复习可重点看）
// 记录每个请求的 method、path、耗时（毫秒），便于排查和监控
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

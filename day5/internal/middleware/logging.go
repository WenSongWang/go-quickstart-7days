package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging 中间件：记录每个请求的方法、路径和耗时
// 接收“下一个 Handler”，返回一个新的 Handler：请求来时先记时间，调 next，再打日志
// 这里只有一层 return：返回的是 http.HandlerFunc(...)，即一个 Handler，不是「返回一个函数」（那种是 Timeout 的写法）
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()// 记录请求开始时间
		next.ServeHTTP(w, r) // 把请求交给下一层（鉴权或业务）
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))// time.Since(start) 计算请求处理时间,Since是time包中的函数,表示时间差
	})
}

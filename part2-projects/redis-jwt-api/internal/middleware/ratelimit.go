// RateLimit 基于令牌桶的限流中间件（单机），便于复习/自测「限流怎么实现」。
// 多实例需用 Redis 做分布式限流（INCR + EXPIRE 或 Lua）。
package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 按 IP 的令牌桶限流（每 IP 每秒 r 个令牌，桶容量 b）。
// 超过速率返回 429 Too Many Requests。
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	mu := sync.Mutex{}
	m := make(map[string]*rate.Limiter)

	getLimiter := func(key string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		if lim, ok := m[key]; ok {
			return lim
		}
		lim := rate.NewLimiter(r, b)
		m[key] = lim
		return lim
	}

	return func(c *gin.Context) {
		key := c.ClientIP()
		lim := getLimiter(key)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// GlobalRateLimiter 全局限流（不按 IP），适合「全站 QPS 上限」的简单场景。
func GlobalRateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	lim := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

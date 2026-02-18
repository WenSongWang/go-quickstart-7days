// Part2 项目 1：Redis + JWT API（进阶）
// 真实使用 gin + JWT（HMAC）+ Redis（可选）：
// - POST /api/login 获取 token
// - GET /api/me 需 Authorization: Bearer <token>
// - GET /api/users 带 Redis 缓存（REDIS_ADDR 非空时启用）
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-quickstart-7days/part2-projects/redis-jwt-api/internal/config"
	"github.com/go-quickstart-7days/part2-projects/redis-jwt-api/internal/handler"
	"github.com/go-quickstart-7days/part2-projects/redis-jwt-api/internal/middleware"
	"github.com/go-quickstart-7days/part2-projects/redis-jwt-api/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		redisClient = redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	}

	mem := store.NewMemoryStore()
	cached := &store.CachedUserStore{
		Base:    mem,
		ListTTL: 30 * time.Second,
	}
	if redisClient != nil {
		cached.Cache = &store.RedisCache{Client: redisClient}
	}

	healthHandler := &handler.HealthHandler{Redis: redisClient}
	authHandler := &handler.AuthHandler{Secret: cfg.JWTSecret, TTL: cfg.JWTTTL}
	userHandler := &handler.UserHandler{Store: cached}

	r := gin.Default()
	r.GET("/health", healthHandler.Health)

	// 限流：全站每秒 100 请求、桶容量 200（可按需改为按 IP：middleware.RateLimiter(rate.Limit(10), 20)）
	// r.Use(middleware.GlobalRateLimiter(100, 200))

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/users", userHandler.List) // 列表可公开（展示缓存）

		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		protected.GET("/me", userHandler.Me)
		protected.POST("/users", userHandler.Create)
	}

	addr := ":" + strconv.Itoa(cfg.HTTPPort)
	fmt.Println("listening on", addr)
	_ = r.Run(addr)
}

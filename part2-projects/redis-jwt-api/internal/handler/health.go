package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	Redis *redis.Client // 可为空
}

func (h *HealthHandler) Health(c *gin.Context) {
	resp := gin.H{"status": "ok", "project": "redis-jwt-api"}

	if h.Redis == nil {
		resp["redis"] = "disabled"
		c.JSON(http.StatusOK, resp)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 300*time.Millisecond)
	defer cancel()
	if err := h.Redis.Ping(ctx).Err(); err != nil {
		resp["redis"] = "unavailable"
		resp["redis_error"] = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp["redis"] = "ok"
	c.JSON(http.StatusOK, resp)
}


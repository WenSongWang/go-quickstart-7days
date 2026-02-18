package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-quickstart-7days/part2-projects/redis-jwt-api/internal/store"
)

type UserHandler struct {
	Store *store.CachedUserStore
}

// GET /api/users
// 返回 {"cached":true/false,"users":[...]}
func (h *UserHandler) List(c *gin.Context) {
	users, cached, err := h.Store.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list users failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cached": cached, "users": users})
}

// POST /api/users（需 JWTAuth 中间件）
// body: {"name":"xxx"}，name 必填、长度 1～100（validator 示例）
func (h *UserHandler) Create(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required" validate:"min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if len(body.Name) < 1 || len(body.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name length 1-100"})
		return
	}
	u, err := h.Store.Create(c.Request.Context(), body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
		return
	}
	c.JSON(http.StatusCreated, u)
}

// GET /api/me（需要 JWTAuth 中间件）
func (h *UserHandler) Me(c *gin.Context) {
	claimsAny, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"claims": claimsAny})
}


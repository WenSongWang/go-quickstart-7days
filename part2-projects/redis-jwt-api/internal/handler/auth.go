package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	Secret string
	TTL    time.Duration
}

// POST /api/login
// body: {"user":"demo"}  -> 返回 {"token":"..."}
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		User string `json:"user"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user required"})
		return
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": body.User,
		"iat": now.Unix(),
		"exp": now.Add(h.TTL).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(h.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed})
}


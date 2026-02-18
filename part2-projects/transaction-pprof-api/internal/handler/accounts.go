package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-quickstart-7days/part2-projects/transaction-pprof-api/internal/resp"
	"github.com/go-quickstart-7days/part2-projects/transaction-pprof-api/internal/store"
)

type AccountHandler struct {
	Store *store.SQLiteStore
}

// GET /api/accounts
func (h *AccountHandler) List(c *gin.Context) {
	accounts, err := h.Store.ListAccounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Fail(resp.CodeInternal, "list accounts failed"))
		return
	}
	c.JSON(http.StatusOK, resp.OK(accounts))
}

// POST /api/transfer
// body: {"from":1,"to":2,"amount":100}
func (h *AccountHandler) Transfer(c *gin.Context) {
	var body struct {
		From   int `json:"from"`
		To     int `json:"to"`
		Amount int `json:"amount"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, resp.Fail(resp.CodeBadRequest, "invalid json"))
		return
	}
	if err := h.Store.Transfer(c.Request.Context(), body.From, body.To, body.Amount); err != nil {
		c.JSON(http.StatusBadRequest, resp.Fail(resp.CodeInsufficientBalance, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.OK(gin.H{"transferred": body.Amount}))
}


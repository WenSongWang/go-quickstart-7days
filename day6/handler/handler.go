// handler 包：Day6 用于测试的 HTTP Handler 示例
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// User 示例实体，JSON 字段名为 id / name
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetUser 根据路径中的 id 返回对应用户的 JSON（如 /api/users/1）
// 用于 Day6 的 httptest 测试示例
func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(User{ID: id, Name: "用户" + strconv.Itoa(id)})
}

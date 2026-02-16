package middleware

import (
	"net/http"
)

// APIKeyAuth 返回一个“包装函数”：只接受带正确 X-API-Key 的请求，否则返回 401
// 用法：APIKeyAuth("secret")(nextHandler)，表示用 "secret" 当合法 Key
func APIKeyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" || key != apiKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized) // 401
				w.Write([]byte(`{"error":"unauthorized"}`))
				return // 如果key为空或不等于apiKey,则返回401
			}
			next.ServeHTTP(w, r) // 校验通过，交给下一个 Handler
		})
	}
}

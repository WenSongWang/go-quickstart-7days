// 请求超时中间件：为每个请求注入带超时的 context（复习常考：如何做请求超时）
package middleware

import (
	"context" // 上下文, 可以传递取消与超时信号, 常用来做超时, 也可以用来做多个goroutine之间的同步
	"net/http"
	"time"
)

// Timeout 限制单个请求的最大处理时间，超时后 ctx 被取消，下游可检测 ctx.Done()
// 用法：Timeout(5*time.Second)(nextHandler)
//
// 嵌套等价理解（按「请求来了会发生什么」看）：
//   step1 := Timeout(5 * time.Second)   // 得到「一个函数：接收 next，返回 Handler」
//   step2 := step1(业务Handler)          // 把业务传进去，得到「包装后的 Handler」
//   请求来时 step2 被调用，做三件事：
//     ① ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
//     ② defer cancel()
//     ③ next.ServeHTTP(w, r.WithContext(ctx))  // 用带超时的 ctx 调下一层
//   「下一层」拿到的 r.Context() 就是这个 ctx，超时后会被取消。
//
// 两层 return 怎么理解：
//   - 函数签名是 Timeout(d) 返回「func(http.Handler) http.Handler」，所以第一层 return 的必须是「一个函数」；
//     我们 return 的就是 func(next http.Handler) http.Handler { ... }，即「接收 next，返回 Handler」的那个函数。
//   - 这个匿名函数自己的返回值类型是 http.Handler，所以第二层 return 的必须是「一个 Handler」；
//     我们 return http.HandlerFunc(...)，即真正处理请求、带超时逻辑的那一坨。
//   - 总结：第一层 return「造中间件的函数」，第二层 return「造出来的那个 Handler」。
func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			// 把带超时的 ctx 放进请求里，下一层 Handler 用 r.Context() 就能收到取消信号
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

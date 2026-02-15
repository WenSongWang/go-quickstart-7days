// Day2 示例：标准库 HTTP 服务 + 简单路由 + JSON
// 本文件学习：HandleFunc、Request/ResponseWriter、写 JSON、读路径参数
package main
/**
**
HandleFunc：注册路由与处理函数
Request/ResponseWriter：处理请求与响应
写 JSON：json.NewEncoder(w).Encode(user)
读路径参数：strconv.Atoi(path)
监听端口：http.ListenAndServe(":8080", nil)
优雅关闭：context.WithTimeout(context.Background(), 10*time.Second)
**/
import (
	"crypto/rand"//随机数
	"encoding/hex"//hex编码
	"encoding/json"
	"net/http"//HTTP服务
	"strconv"//字符串转换为整数
	"strings"//字符串操作
)

// setTraceID 从请求头读取 X-Trace-Id，没有则生成一个，并写回响应头（便于调用方和日志串联）
func setTraceID(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-Trace-Id")
	if id == "" {
		b := make([]byte, 8)
		rand.Read(b)
		id = hex.EncodeToString(b)
	}
	w.Header().Set("X-Trace-Id", id)
}

// User 返回给前端的结构体，`json:"id"` 表示转成 JSON 时字段名叫 "id"
type User struct {
	ID   int    `json:"id"` //json tag：表示转成 JSON 时字段名叫 "id", 比如{"id":1,"name":"张三"}
	Name string `json:"name"` //json tag：表示转成 JSON 时字段名叫 "name"
}

// handleHello 处理 /hello 路径：设置响应头、状态码、写响应体
func handleHello(w http.ResponseWriter, r *http.Request) {
	setTraceID(w, r) // 响应头带上 trace_id，便于调用方和日志追踪
	if r.URL.Path != "/hello" {
		http.NotFound(w, r) // NotFound（w, r）：w 是 ResponseWriter，r 是 Request
		return //返回 404
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // 告诉浏览器是纯文本
	w.WriteHeader(http.StatusOK)                                // 200 状态码，还有其他状态码，比如404、500、502等，详见https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status
	w.Write([]byte("侬好, Go HTTP!"))
}

// handleGetUser 处理 /api/users/1 这种路径：从路径里取出 id，查“用户”，返回 JSON
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	setTraceID(w, r) // 响应头带上 trace_id
	if r.Method != http.MethodGet {//如果方法不是 GET，则报错
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)// "Method Not Allowed" 是错误信息，http.StatusMethodNotAllowed 是 405 状态码
		return
	}
	// 从 /api/users/1 里取出 "1"，再转成数字
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")//TrimPrefix 是 strings 包的函数，用于去除字符串的前缀, trim(修剪，削减，装饰，整齐地)
	id, err := strconv.Atoi(path)//Atoi 是 strconv 包的函数，用于将字符串转换为整数,atoi是ascii to integer的缩写
	if err != nil || id <= 0 {//如果转换失败或者 id 小于等于 0，则报错
		w.Header().Set("Content-Type", "application/json") // 须在 WriteHeader 前设置头
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)// "invalid id" 是错误信息，http.StatusBadRequest 是 400 状态码
		return
	}
	user := User{ID: id, Name: "用户" + strconv.Itoa(id)} // User 结构体；:= 短声明；Itoa = integer to string（整型转十进制字符串）；string + string 拼接
	w.Header().Set("Content-Type", "application/json")//设置响应头
	w.WriteHeader(http.StatusOK)//200 状态码
	json.NewEncoder(w).Encode(user) // 把结构体编码成 JSON 写到 w,json.NewEncoder(w).Encode(user) 是 json 包的函数，用于将结构体编码成 JSON 写到 w
}

func main() {
	// 注册路径与处理函数：/hello 为精确匹配；/api/users/ 末尾有斜杠，为标准库的「前缀匹配」，故 /api/users/1、/api/users/99 都会进 handleGetUser
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/api/users/", handleGetUser)
	// 监听 8080 端口，nil 表示用默认的路由（就是上面 HandleFunc 注册的）
	http.ListenAndServe(":8080", nil)
}

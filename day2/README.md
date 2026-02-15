# Day 2：HTTP 服务

- 使用标准库 `net/http` 起一个 HTTP 服务
- 路由：`http.HandleFunc` 或自己用 `mux` 匹配
- 读取请求：Method、URL、Body、Header
- 返回响应：状态码、JSON、Write

## 运行

```bash
go run ./day2/server
# 访问 http://localhost:8080/hello
# 访问 http://localhost:8080/api/users/1
```

## 本日注意点与易踩坑

- **响应头顺序**：返回 JSON 错误时要先 `w.Header().Set("Content-Type", "application/json")`，再 `http.Error()`，否则头已发出无法再改。
- **路由匹配**：注册 `/api/users/`（末尾有斜杠）是**前缀匹配**，`/api/users/1` 会进同一 handler；无末尾斜杠才是精确匹配。
- 更多见：[docs/PITFALLS_AND_SOLUTIONS.md](../docs/PITFALLS_AND_SOLUTIONS.md)（HTTP 与响应）。

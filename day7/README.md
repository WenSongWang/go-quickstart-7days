# Day 7：综合实战（入门）

整合前 6 天内容：标准项目结构、配置、HTTP 路由、中间件、内存版「数据层」，实现一个小型 REST API，并带 **slog** 结构化日志与**优雅关闭**（复习常考）。本日定位为**入门级**综合实战；进阶实战可跟进后续 part。

## 功能

| 接口 | 说明 |
| :--- | :--- |
| `GET /health` | 健康检查，返回 `{"status":"ok"}` |
| `GET /api/users` | 用户列表（内存 store） |
| `GET /api/users/:id` | 用户详情 |
| `POST /api/users` | 创建用户，Body 为 JSON `{"name":"xxx"}` |
| 中间件 | 请求日志（**slog**）、可选 API Key 鉴权（配置里 `API_KEY` 非空时开启） |
| 优雅关闭 | 收到 SIGINT/SIGTERM（如 Ctrl+C）后，等当前请求处理完再退出（`Shutdown`） |

## 目录结构

| 目录/文件 | 说明 |
| :--- | :--- |
| **cmd/server/main.go** | 入口：加载配置、挂路由、中间件链（日志 → 可选鉴权）、`http.Server` + `Shutdown`、`signal.Notify` 监听退出信号 |
| **internal/config** | 配置：HTTP_PORT、API_KEY（空则不做鉴权） |
| **internal/store** | 内存版「数据层」：NewMemoryStore，存用户；后续可换成 Day 4 的真实 DB |
| **internal/handler** | 用户接口：List、Create、GetByID（读写 JSON） |
| **internal/middleware** | 请求日志（slog）、API Key 鉴权 |

```
day7/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── store/memory.go
│   ├── handler/users.go
│   └── middleware/
│       ├── logging.go
│       └── auth.go
├── README.md
└── csdn.md
```

## 运行

**以下命令均在项目根目录执行。**

### 1. 直接运行

```bash
go run ./day7/cmd/server
```

可选：设置环境变量（或使用 Day 3 的 .env + godotenv，本目录 config 为简易版 os.Getenv）：

- `HTTP_PORT=8080`（默认 8080）
- `API_KEY=secret`：非空时开启鉴权，请求需带 `X-API-Key: secret`

### 2. 示例请求（curl）

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/users
curl http://localhost:8080/api/users/1
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"name":"李四"}'
```

若开启了 API Key，需在请求头加：`-H "X-API-Key: secret"`（或你配置的值）。

**浏览器快速测试**：

- **GET**：无鉴权时可直接在地址栏打开查看 JSON：
  - `http://localhost:8080/health`
  - `http://localhost:8080/api/users`
  - `http://localhost:8080/api/users/1`
- **若开启了 API Key**：浏览器无法带自定义头，上述地址会返回 401。可在开发者工具（F12）→ Console 执行（若提示不要粘贴可先输入 `allow pasting` 回车）：
  ```js
  fetch("http://localhost:8080/api/users", { headers: { "X-API-Key": "secret" } }).then(r=>r.json()).then(console.log)
  ```
- **POST 创建用户**（浏览器只能用 Console）：
  ```js
  fetch("http://localhost:8080/api/users", { method: "POST", headers: { "Content-Type": "application/json" }, body: '{"name":"李四"}' }).then(r=>r.json()).then(console.log)
  ```
  若已开启鉴权，在 `headers` 里加上 `"X-API-Key": "secret"`。

### 3. 优雅关闭怎么验证

启动服务后，在运行它的终端按 **Ctrl+C**。观察日志：应先输出 `shutting down...`，等待正在处理的请求结束后再输出 `server stopped` 并退出；若直接退出或报错，说明 Shutdown 未正确等待。

关闭时终端预期日志示例（slog 格式）：

```
2026/02/17 13:14:37 INFO shutting down...
2026/02/17 13:14:37 INFO server stopped
```

### 4. Docker（与 Day 6 共用 Dockerfile）

```bash
docker build -f day6/Dockerfile -t go-7days-app .
docker run -p 8080:8080 go-7days-app
```

镜像内即本日 server，访问 `http://localhost:8080` 验证。

## 本日注意点与易踩坑

- **路由顺序**：`/api/users` 与 `/api/users/` 要分开注册；带尾斜杠的会匹配 `/api/users/1` 等，标准库为前缀匹配。
- **优雅关闭**：`ListenAndServe` 在 goroutine 里跑，main 里用 `signal.Notify` 等信号，收到后调 `srv.Shutdown(ctx)`；`Shutdown` 会停接新请求并等当前请求完成，超时由 `context.WithTimeout` 控制（本示例 10 秒）。
- **API Key**：配置中 `API_KEY` 为空则不鉴权；非空则所有请求需带正确 `X-API-Key`，否则 401。
- 更多见 [docs/PITFALLS_AND_SOLUTIONS.md](../docs/PITFALLS_AND_SOLUTIONS.md)。

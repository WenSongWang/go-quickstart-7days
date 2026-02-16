# Day 5：中间件与业务分层

- **中间件**：包装 `http.Handler`，在「进业务前」统一做日志、鉴权、**请求超时（context）**。
- **本日链式顺序**：请求 → 日志 →（/api 下）超时 → 鉴权 → 业务 Handler。
- **简单鉴权**：从 Header 取 `X-API-Key`，与配置比对；未通过返回 401。示例里写死 `"secret"`，**生产环境应从配置（如 Day 3 的 cfg）读取**。
- **请求超时**：用 `context.WithTimeout` 给请求注入带超时的 context，下游用 `r.Context()` 可检测取消（面试常问）。
- **分层**：handler → service → repository（本日以 handler + 中间件为主，service/repository 在 Day 7 展开）。

## 目录结构

| 目录/文件 | 说明 |
| :--- | :--- |
| **cmd/server/** | 入口 main：挂路由（/health、/api/）、拼中间件链（Logging → Timeout → APIKeyAuth）、ListenAndServe :8080。 |
| **internal/middleware/** | 中间件包：logging（记耗时）、auth（X-API-Key 鉴权）、context_timeout（请求超时注入 context）。 |

```
day5/
├── cmd/server/main.go
├── internal/middleware/
│   ├── logging.go
│   ├── auth.go
│   └── context_timeout.go
├── README.md
└── csdn.md
```

## 运行

```bash
# 项目根目录执行
go run ./day5/cmd/server
```

- 健康检查（无鉴权）：`curl http://localhost:8080/health`，或浏览器直接打开 `http://localhost:8080/health`
- 带鉴权访问 API：`curl -H "X-API-Key: secret" http://localhost:8080/api/hello`
- 不带 Key 会返回 401：`curl http://localhost:8080/api/hello`，或浏览器直接打开 `http://localhost:8080/api/hello`（会看到 401，因浏览器不会带 X-API-Key）

**浏览器里带 Key 测试**：地址栏无法加自定义 Header，可打开开发者工具（F12）→ Console，执行：
`fetch("http://localhost:8080/api/hello", { headers: { "X-API-Key": "secret" } }).then(r=>r.json()).then(console.log)`
若控制台提示「Don't paste code...」：先输入 `allow pasting` 回车，再粘贴上述代码执行。

## 本日注意点与易踩坑

- **中间件顺序**：先包装的先执行。当前顺序为「日志 → 超时 → 鉴权 → 业务」，改顺序会改变行为（例如先鉴权再超时 vs 先超时再鉴权）。
- **写响应前别重复 WriteHeader**：在中间件里若已 `w.WriteHeader(401)` 并 `return`，就不要再调用 `next.ServeHTTP`；业务里也不要再改状态码。
- **超时与 context**：Timeout 中间件把 `r.WithContext(ctx)` 传给下一层，业务或 service 里应用 `r.Context()` 做 DB/RPC 调用，以便超时能取消下游。
- **API Key**：示例中 `"secret"` 写死在代码里，仅为演示；上线应从配置或环境变量读取。
- **延伸阅读**：[接口鉴权综述：签名、AK/SK、加解密与常见方式](https://blog.csdn.net/weixin_40959890/article/details/157911873)（CSDN，含 Python 示例）。
- 更多见 [docs/PITFALLS_AND_SOLUTIONS.md](../docs/PITFALLS_AND_SOLUTIONS.md)。

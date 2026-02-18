# 第二大部分：项目巩固（方案 B）

本目录是 **方案 B**：在完成**第一大部分（7 天 + 方案 A）** 之后，用 **1～2 个小项目** 练手，把 Redis、JWT、事务、pprof 等在真实场景里用一遍，巩固中级。

---

## 与第一大部分的关系

| 部分 | 位置 | 内容 |
| :--- | :--- | :--- |
| **第一大部分** | 根目录 `day0`～`day7` | **7 天入门**（标准库为主，Day 7 为入门级综合实战：net/http、slog、优雅关闭、内存 store） |
| **第二大部分** | 本目录 `part2-projects/` | **方案 B**：小项目实战，接 **gin**、Redis、JWT、事务、pprof 等 |

学完第一大部分后，用本目录下的项目做「第二次 7 天」或按需练手，不要求连续 7 天，重在**在项目里用熟**。

---

## 项目列表

| 项目 | 目录 | 目标 | 技术点 |
| :--- | :--- | :--- | :--- |
| **Redis + JWT API** | [redis-jwt-api/](redis-jwt-api/) | 小 API 接 Redis 与 JWT 鉴权 | **gin**、Redis 缓存/会话、JWT 签发与校验 |
| **事务 + pprof API** | [transaction-pprof-api/](transaction-pprof-api/) | 小 API 接数据库事务与性能分析 | **gin**、事务、统一错误码、pprof 排查 |

每个项目可单独运行；依赖统一使用根目录 `go.mod`，新增依赖在根目录执行 `go get`。

---

## 实际开发中的框架使用（Part2 统一约定）

Part2 项目**推荐使用 HTTP 框架**，贴近真实后端开发与复习常考的「用啥框架？」。

| 框架 | 说明 | Part2 建议 |
| :--- | :--- | :--- |
| **gin** | 国内最常用、中间件生态丰富、与 JWT/Redis 结合示例多 | **推荐**，两个 part2 项目均按 gin 写路由与中间件 |
| echo | 轻量、性能好、API 与 gin 类似 | 可选替代 |
| go-chi | 标准库风格、轻量 | 可选替代 |

- **当前**：part2 两个项目的 **`cmd/server/main.go` 已实际使用 gin**（`gin.Default()`、`r.GET`、`c.JSON`）；transaction-pprof-api 并用 `gin.WrapF` 挂载了 `/debug/pprof/`。根目录 `go.mod` 已包含 gin，`go run ./part2-projects/.../cmd/server` 即可跑通。
- 后续接入 Redis、JWT、事务时，均在现有 gin 路由与中间件基础上扩展。

---

## 使用建议

1. 先完成第一大部分（day0～day7，含方案 A 加料）。
2. 从本目录中选一个项目，按各项目 README 在**项目根目录**跑通（如 `go run ./part2-projects/redis-jwt-api/cmd/server`），再按 README 的「后续步骤建议」接入 Redis/JWT 或事务/pprof。
3. 复习时可把「第一大部分 Day 7」+「本目录某项目」一起梳理成你的练手经历。

**当前状态**：两个项目都已具备「最小可运行闭环」：

- `redis-jwt-api`：gin + JWT +（可选）Redis 缓存 + 可跑通的登录/鉴权/列表接口
- `transaction-pprof-api`：gin + SQLite（modernc）+ 事务示例 + 统一响应体 + pprof

后续按各项目 README 的「后续可以加什么」继续加深即可。

详见根目录 [README.md](../README.md#课程规划与技术栈)。

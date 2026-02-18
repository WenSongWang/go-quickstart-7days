# Part2 项目 2：事务 + pprof API（进阶）

在 Day 4 / Day 7 的基础上，用一个小型 API 项目练 **数据库事务** 与 **pprof 性能分析**，并配合统一错误码与响应体。

---

## 目标

| 项 | 说明 |
| :--- | :--- |
| **前置** | 已完成第一大部分 Day 0～Day 7，建议已学过 Day 4 数据层。 |
| **产出** | 可运行的 HTTP 服务：至少一个多步写操作用事务保证一致性；暴露 pprof 便于排查性能。 |
| **技术点** | 事务（Begin/Commit/Rollback）、统一错误码/响应体、net/http/pprof。 |
| **框架** | **gin**：项目代码已使用 gin（路由、c.JSON），并已用 `gin.WrapF` 挂载 /debug/pprof；后续在现有基础上加 DB、事务与统一响应体。 |

## 技术点与自测要点

| 技术点 | 本项目中用途 | 自测可讲 |
| :--- | :--- | :--- |
| **事务** | 多步写操作（如转账、扣库存）要么全成功要么全回滚 | 隔离级别、何时 Rollback、context 超时与事务 |
| **统一响应体** | `{ code, message, data }`，错误码表 | 前后端约定、错误码分层 |
| **pprof** | CPU/内存/goroutine 分析，找热点与泄漏 | 如何排查性能问题、benchmark、火焰图 |

## 目录结构（规划）

```
part2-projects/transaction-pprof-api/
├── cmd/server/          # 入口：路由、挂载 /debug/pprof、事务示例接口
├── internal/
│   ├── config/          # 配置：HTTP_PORT、DB_DSN
│   ├── store/           # 数据层：带事务的写操作（Tx）
│   ├── handler/         # 接口：统一响应体、调用带事务的 store
│   └── middleware/      # 日志、可选鉴权
├── README.md
└── （沿用根 go.mod）
```

当前已实现 **gin + SQLite（modernc）+ 事务示例 + 统一错误码/响应体 + pprof** 的闭环，可直接运行并联调；错误码集中定义，便于当手册查。

### 错误码表（当手册查）

| code | 常量 | 含义 |
| :--- | :--- | :--- |
| 0 | CodeOK | 成功 |
| 1 | CodeInternal | 内部错误（如 list 失败） |
| 2 | CodeBadRequest | 参数错误（如 JSON 无效） |
| 3 | CodeInsufficientBalance | 业务错误（如余额不足） |
| 4 | CodeNotFound | 未找到（预留） |

定义位置：`internal/resp/resp.go`。

### pprof 使用步骤（当手册查）

1. 启动服务：`go run ./part2-projects/transaction-pprof-api/cmd/server`
2. **CPU 采样**（默认 30 秒）：  
   `go tool pprof http://localhost:8081/debug/pprof/profile?seconds=30`  
   进入后输入 `top`、`web`（需 graphviz）看热点。
3. **堆内存**：  
   `go tool pprof http://localhost:8081/debug/pprof/heap`
4. **goroutine 数量**：  
   `go tool pprof http://localhost:8081/debug/pprof/goroutine`
5. 压测后采样更准：如 `wrk -t4 -c100 -d30s http://localhost:8081/api/accounts` 再采 CPU/heap。

### 深度说明（代码位置）

| 技术点 | 本项目中的实现 | 代码位置 |
| :--- | :--- | :--- |
| **事务** | Transfer：BeginTx → 查余额 → 扣款 → 加款 → Commit，任一步 err 则 Rollback | `internal/store/sqlite.go` Transfer |
| **统一响应体** | `{ code, message, data }`，resp.OK / resp.Fail | `internal/resp/resp.go`、handler 里统一返回 |
| **pprof** | gin.WrapF 挂 net/http/pprof | `cmd/server/main.go`、/debug/pprof/ |

## 运行

在**项目根目录**执行：

```bash
go run ./part2-projects/transaction-pprof-api/cmd/server
```

默认监听 `:8081`（与 redis-jwt-api 的 8080 错开）。

### 环境变量

- `HTTP_PORT`：默认 `8081`
- `SQLITE_DSN`：默认 `file:tx_demo?mode=memory&cache=shared`（内存库，进程退出即清空）

### 接口一览（可直接跑通）

- `GET /health`：健康检查（gin 返回 JSON）
- `GET /debug/pprof/`：pprof 索引（CPU、heap、goroutine 等）
- `go tool pprof http://localhost:8081/debug/pprof/profile?seconds=5`：采 5 秒 CPU
- `GET /api/accounts`：账户列表（默认 seed：id=1、2，balance=1000）
- `POST /api/transfer`：事务转账（示例 body：`{"from":1,"to":2,"amount":100}`）

### 统一响应体示例

本项目接口返回统一格式：`{ code, message, data }`：

- 成功：`{"code":0,"message":"ok","data":[...]}`  
- 失败：`{"code":3,"message":"insufficient balance"}`

## 后续可以加什么（进阶点）

1. 把 SQLite 换成 Postgres，并给转账加**幂等键**防重复提交。
2. 更复杂事务：扣库存 + 创建订单 + 写流水（多表一 tx）。
3. 压测 + pprof：`wrk` 打 /api/accounts 或 /api/transfer，再采 CPU/heap 看热点与分配。

详见根目录 [README.md](../../README.md#课程规划与技术栈) 与 [docs/INTERVIEW_TOPICS.md](../../docs/INTERVIEW_TOPICS.md)。

# Part2 项目 1：Redis + JWT API（进阶）

在 Day 7 的基础上，用一个小型 API 项目接入 **Redis** 与 **JWT 鉴权**，在真实请求链路里用熟缓存与会话、无状态鉴权。

---

## 目标

| 项 | 说明 |
| :--- | :--- |
| **前置** | 已完成第一大部分 Day 0～Day 7。 |
| **产出** | 可运行的 HTTP 服务：部分接口用 Redis 缓存、部分接口需 JWT 鉴权。 |
| **技术点** | Redis 连接与 Get/Set、JWT 签发与校验、与现有路由/中间件结合。 |
| **框架** | **gin**：项目代码已使用 gin 作为 HTTP 框架（路由、c.JSON），后续直接在此基础上加 JWT 中间件与 Redis。 |

## 技术点与自测要点

| 技术点 | 本项目中用途 | 自测可讲 |
| :--- | :--- | :--- |
| **Redis** | 缓存热点数据、会话或限流计数 | 连接池、缓存穿透/击穿/雪崩、分布式锁（SET NX） |
| **JWT** | 无状态鉴权，替代或补充 API Key | Access/Refresh、签名算法、过期与刷新 |

## 目录结构（规划）

```
part2-projects/redis-jwt-api/
├── cmd/server/          # 入口：路由、中间件（含 JWT）、Redis 客户端
├── internal/
│   ├── config/          # 配置：HTTP_PORT、REDIS_ADDR、JWT_SECRET 等
│   ├── store/           # 数据层（可先内存，再接 Redis 缓存层）
│   ├── handler/         # 接口：需缓存的 GET、需 JWT 的写操作
│   └── middleware/      # JWT 校验、可选 Redis 健康检查
├── README.md
└── （后续）go.mod 可选，或沿用根 go.mod
```

当前已实现 **gin + JWT + Redis（可选）缓存 + 写后删 + 限流中间件（可选）+ 请求体验证** 的闭环，可直接运行并联调；需要时可开启限流、接 Refresh Token 等。

### 深度说明（当手册查）

| 技术点 | 本项目中的实现 | 代码位置 |
| :--- | :--- | :--- |
| **Cache-Aside 读** | 先查缓存，miss 再查 store 并回填 | `internal/store/cache.go` List |
| **写后删** | Create 后删 `users:all` key，下次 List 再回填 | `internal/store/cache.go` Create、`RedisCache.Delete` |
| **JWT 签发/校验** | Login 发 token，JWTAuth 中间件校验 Bearer | `internal/handler/auth.go`、`internal/middleware/jwt.go` |
| **限流（令牌桶）** | 单机按 IP 或全局限流，`golang.org/x/time/rate` | `internal/middleware/ratelimit.go`；main 里注释了 `r.Use(...)`，可按需打开 |
| **请求体验证** | Create 时 name 必填、长度 1～100 | `internal/handler/users.go` Create body 校验 |
| **健康检查带 Redis** | `/health` 里 ping Redis，返回 redis: ok/unavailable | `internal/handler/health.go` |

## 运行

在**项目根目录**执行：

```bash
go run ./part2-projects/redis-jwt-api/cmd/server
```

默认监听 `:8080`。

### 环境变量

- `HTTP_PORT`：默认 `8080`
- `JWT_SECRET`：默认 `dev-secret`（生产务必改为强随机值）
- `JWT_TTL_SECONDS`：默认 `3600`
- `REDIS_ADDR`：为空则**不启用 Redis**；示例：`localhost:6379`

### 接口一览（可直接跑通）

- **GET `/health`**：返回服务状态；若启用 Redis，会带 `redis: ok/unavailable`
- **POST `/api/login`**：登录获取 token
- **GET `/api/users`**：用户列表；启用 Redis 时会做列表缓存，并返回 `cached: true/false`
- **GET `/api/me`**：需要 `Authorization: Bearer <token>`，返回解析出的 claims
- **POST `/api/users`**：需要 `Authorization: Bearer <token>`，创建用户（并刷新缓存）

### 快速测试（curl）

```bash
# 1) 登录拿 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"user":"demo"}' | jq -r .token)

# 2) 鉴权接口
curl -s http://localhost:8080/api/me -H "Authorization: Bearer $TOKEN"

# 3) 列表（无鉴权）
curl -s http://localhost:8080/api/users

# 4) 创建（鉴权）
curl -s -X POST http://localhost:8080/api/users -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"name":"王五"}'
```

## 后续可以加什么（进阶点）

1. 登录改成真实用户体系（密码哈希、**Refresh Token**）。
2. Redis 作为会话/黑名单（登出时让 token 失效）。
3. 缓存击穿/雪崩：互斥回填、TTL 随机、布隆过滤器（本项目已做写后删与 Cache-Aside）。
4. 分布式限流：用 Redis INCR + EXPIRE 或 Lua 替代单机 rate.Limiter。

详见根目录 [README.md](../../README.md#课程规划与技术栈) 与 [docs/INTERVIEW_TOPICS.md](../../docs/INTERVIEW_TOPICS.md)。

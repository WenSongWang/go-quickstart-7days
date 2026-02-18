# 7 天从零快速入门 Go 语言后端开发

一个按天拆分的 Go 后端入门项目，每天一个主题，从环境搭建到能写一个完整的 REST API，并**覆盖中级与复习中最常见的技术点**（并发、context、错误处理、中间件、优雅关闭、测试、Docker 等）。

---

## 两大目录说明

| 部分 | 目录 | 内容 |
| :--- | :--- | :--- |
| **第一大部分** | 根目录 `day0`～`day7` | **7 天入门**：Day 7 为**入门级**综合实战（标准库、slog、优雅关闭、内存 store、API Key）。学完能写完整小 API。 |
| **第二大部分** | [part2-projects/](part2-projects/) | **进阶实战**：已实现两个项目——**redis-jwt-api**（gin + JWT + Redis）、**transaction-pprof-api**（gin + 事务 + 统一响应体 + pprof）。 |

- **学完能变成中级工程师吗？** 7 天 = 入门桥梁；part2 用 gin、JWT、Redis、事务、pprof 巩固。见下方「课程规划与技术栈」与 [docs/INTERVIEW_TOPICS.md](docs/INTERVIEW_TOPICS.md)。
- **能覆盖复习常考的吗？** 能覆盖大部分入门～中级常考点；文档按主题列出覆盖清单与代码位置。

## 第一大部分：学习路线（7 天）

| 天数 | 主题 | 内容 |
| :--- | :--- | :--- |
| **Day 0** | 预习 | 关键词、语法、常用包一览（不写代码，先混个眼熟） |
| **Day 1** | Go 基础 | 环境、语法、包与模块、基础类型与错误处理 |
| **Day 2** | HTTP 服务 | `net/http`、路由、请求/响应、JSON |
| **Day 3** | 项目结构 & 配置 | 目录组织、环境变量、viper 配置加载 |
| **Day 4** | 数据层 | 数据库连接、查询、简单 CRUD（SQLite/PostgreSQL） |
| **Day 5** | 中间件与业务分层 | 日志、API Key 鉴权、请求超时 |
| **Day 6** | 测试 & 部署 | 单元测试、Docker 多阶段构建 |
| **Day 7** | 综合实战（入门） | 小型 REST API（标准库、slog、优雅关闭、内存 store） |

## 第二大部分：进阶实战（part2-projects）

学完 7 天后，用 **part2-projects/** 里的小项目巩固：接 **gin**、JWT、Redis、事务、pprof 等，贴近真实后端与复习常考点。详见 [part2-projects/README.md](part2-projects/README.md)。

| 项目 | 目录 | 技术点 |
| :--- | :--- | :--- |
| **Redis + JWT API** | [part2-projects/redis-jwt-api/](part2-projects/redis-jwt-api/) | gin、JWT 签发/校验、Redis 缓存/会话 |
| **事务 + pprof API** | [part2-projects/transaction-pprof-api/](part2-projects/transaction-pprof-api/) | gin、SQLite 事务、统一响应体、pprof 性能分析 |

- 每个项目可单独运行，依赖使用根目录 `go.mod`。
- 运行示例：`go run ./part2-projects/redis-jwt-api/cmd/server`、`go run ./part2-projects/transaction-pprof-api/cmd/server`。

## 课程规划与技术栈

- **第一大部分**：入门到能写完整小 API（标准库为主，viper 配置，database/sql 数据层）。**第二大部分**：part2 用 gin、JWT、Redis、事务、pprof 巩固，贴近实战。
- **技术栈简表**：Day 1～2 基础与 HTTP；Day 3 godotenv + viper；Day 4 database/sql + SQLite/PostgreSQL；Day 5 中间件与超时；Day 6 testing + Docker；Day 7 net/http + slog + 优雅关闭 + 内存 store。part2：redis-jwt-api 用 gin、JWT、go-redis；transaction-pprof-api 用 gin、SQLite、事务、统一响应体、pprof。
- **与中级的关系**：7 天 + part2 能把常见技术点过一遍，达到「能写、能讲、能扩展」；不能保证直接胜任中级（需真实项目经验）。限流、熔断、gRPC、可观测等扩展点见 [docs/INTERVIEW_TOPICS.md](docs/INTERVIEW_TOPICS.md) 第十一节。

## 环境要求

- **Go 1.21+**（代码按 Go 1.21 编写，1.21.x / 1.22.x 均可，如 `go1.21.12`）
- （可选）Docker，用于 Day 6/7
- **国内用户**：若 `go mod tidy` 连不上官方代理（报 connectex 超时），可执行 `go env -w GOPROXY=https://goproxy.cn,direct` 后重试；详见 [docs/PITFALLS_AND_SOLUTIONS.md](docs/PITFALLS_AND_SOLUTIONS.md)「环境与依赖」。
- **import 写 `github.com/go-quickstart-7days/...` 是本地包路径**，不是从 GitHub 在线拉代码；说明见 [Day 0 六、Go 模块与 import 路径](day0/README.md#六go-模块与-import-路径新手易懵)，以及根目录 [go.mod](go.mod) 顶部注释。

## 快速开始

```bash
# 克隆仓库（HTTPS）
git clone https://github.com/WenSongWang/go-quickstart-7days.git
cd go-quickstart-7days

# 拉取依赖（国内若超时见上方「环境要求」）
go mod tidy

# 按天运行示例（第一大部分）
go run ./day1/hello
go run ./day2/server
go run ./day7/cmd/server

# 进阶项目示例（第二大部分，需先 go mod tidy）
go run ./part2-projects/redis-jwt-api/cmd/server
go run ./part2-projects/transaction-pprof-api/cmd/server
```

## 目录结构

```
.
├── README.md                  # 项目说明与两大目录、学习路线
├── go.mod                     # Go 模块定义与依赖声明
├── docs/
│   ├── INTERVIEW_TOPICS.md     # 技术点梳理与自测要点
│   └── PITFALLS_AND_SOLUTIONS.md  # 实战踩坑与解法
│
├── day0/                     # 【第一大部分】预习：关键词、语法、常用包
├── day1/                     # Go 基础（含 errors、并发与 context）
├── day2/                     # HTTP 服务
├── day3/                     # 项目结构与配置（viper、godotenv）
├── day4/                     # 数据层（SQLite、PostgreSQL）
├── day5/                     # 中间件（鉴权、请求超时）
├── day6/                     # 测试与 Docker
├── day7/                     # 综合实战（入门）：标准库、slog、优雅关闭、内存 store
│
└── part2-projects/            # 【第二大部分】进阶：gin、JWT、Redis、事务、pprof（已实现）
    ├── redis-jwt-api/        # gin + JWT + Redis 缓存
    ├── transaction-pprof-api/ # gin + SQLite 事务 + 统一响应体 + pprof
    └── README.md
```


## 使用建议

1. **Day 0 可先预习**：看 [day0/README.md](day0/README.md) 把关键词、语法、常用包过一遍，再看代码不会懵。
2. **按顺序学习**：每天依赖前几天的概念，建议从 Day 1 开始。
3. **动手改代码**：每个 `dayN` 下都有可运行示例，改一改再跑，加深理解。
4. **Day 7** 是前面内容的整合，可作为你的第一个「小项目」模板。
5. 学完第一大部分后，用 **part2-projects/** 做进阶巩固：两个小项目分别练 **gin + JWT + Redis** 和 **gin + 事务 + pprof**，详见 [part2-projects/README.md](part2-projects/README.md)。

## License

MIT

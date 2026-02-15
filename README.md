# 7 天从零快速入门 Go 语言后端开发

一个按天拆分的 Go 后端入门项目，每天一个主题，从环境搭建到能写一个完整的 REST API，并**覆盖中级/面试中最常见的技术点**（并发、context、错误处理、中间件、优雅关闭、测试、Docker 等）。

---

## 两大目录说明

| 部分 | 目录 | 内容 |
| :--- | :--- | :--- |
| **第一大部分** | 根目录 `day0`～`day7` | **7 天入门到中级** + **方案 A**（事务、统一错误码、JWT、validator 等已纳入对应天）。学完即可覆盖中级面试常考的知识面。 |
| **第二大部分** | [part2-projects/](part2-projects/) | **方案 B**：7 天后的**项目实战巩固**，用 1～2 个小项目接 Redis、JWT、事务、pprof 等，在项目里用熟。 |

- **学完能变成中级工程师吗？** 第一大部分 = 入门到中级的桥梁（不是「只有初级」）；第二大部分用项目巩固。详见 [docs/CURRICULUM.md](docs/CURRICULUM.md) 与 [中级与面试技术点说明](docs/INTERVIEW_TOPICS.md)。
- **能覆盖面试常问的吗？** 能覆盖大部分入门～中级常考点；文档中按主题列出了覆盖清单与对应代码位置。

## 第一大部分：学习路线（7 天 + 方案 A）

| 天数 | 主题 | 内容 |
| :--- | :--- | :--- |
| **Day 0** | 预习 | 关键词、语法、常用包一览（不写代码，先混个眼熟） |
| **Day 1** | Go 基础 | 环境、语法、包与模块、基础类型与错误处理 |
| **Day 2** | HTTP 服务 | `net/http`、路由、请求/响应、JSON |
| **Day 3** | 项目结构 & 配置 | 目录组织、环境变量、配置加载 |
| **Day 4** | 数据层 | 数据库连接、查询、简单 CRUD |
| **Day 5** | 中间件与业务分层 | 日志、认证、错误处理中间件 |
| **Day 6** | 测试 & 部署 | 单元测试、Docker 镜像 |
| **Day 7** | 综合实战 | 小型 API 项目整合 |

## 环境要求

- **Go 1.21+**（代码按 Go 1.21 编写，1.21.x / 1.22.x 均可，如 `go1.21.12`）
- （可选）Docker，用于 Day 6/7

## 快速开始

```bash
# 克隆仓库（HTTPS）
git clone https://github.com/WenSongWang/go-quickstart-7days.git
cd go-quickstart-7days

# 拉取依赖
go mod tidy

# 按天运行示例（示例）
go run ./day1/hello
go run ./day2/server
go run ./day7/app
```

## 目录结构

```
.
├── README.md
├── go.mod
├── docs/
│   ├── CURRICULUM.md           # 课程规划、两大目录、方案 A/B、可加技术点
│   ├── INTERVIEW_TOPICS.md     # 中级/面试技术点覆盖说明（必读）
│   └── PITFALLS_AND_SOLUTIONS.md  # 实战踩坑注意点与解法（按主题汇总）
│
├── day0/                     # 【第一大部分】预习：关键词、语法、常用包
├── day1/                     # Go 基础（含 errors、并发与 context）
├── day2/                     # HTTP 服务
├── day3/                     # 项目结构与配置（含 viper）
├── day4/                     # 数据层（含 sqlx、方案 A：事务）
├── day5/                     # 中间件（鉴权、请求超时）
├── day6/                     # 测试与 Docker（含 testify）
├── day7/                     # 综合实战（gin、优雅关闭、方案 A：统一错误码/JWT/validator）
│
└── part2-projects/           # 【第二大部分】方案 B：项目巩固（Redis、JWT、事务、pprof 等）
    └── README.md             # 项目列表与使用说明
```


## 使用建议

1. **Day 0 可先预习**：看 [day0/README.md](day0/README.md) 把关键词、语法、常用包过一遍，再看代码不会懵。
2. **按顺序学习**：每天依赖前几天的概念，建议从 Day 1 开始。
3. **动手改代码**：每个 `dayN` 下都有可运行示例，改一改再跑，加深理解。
4. **Day 7** 是前面内容的整合，可作为你的第一个「小项目」模板。
5. 学完第一大部分后，用 **part2-projects/** 里的项目做方案 B 巩固（见 [part2-projects/README.md](part2-projects/README.md)）。

## License

MIT

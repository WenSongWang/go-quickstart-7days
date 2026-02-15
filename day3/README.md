# Day 3：项目结构与配置

- **常见目录**（本日涉及）：
  - **cmd/**：放可执行程序的入口，一个子目录一个 main 包（如 `cmd/server`）；`go run ./day3/cmd/server` 即运行该入口。
  - **internal/**：仅本项目内部使用的包，外部无法 import；如 `internal/config` 放配置加载逻辑，避免在 main 里堆满读 env 的代码。
  - **config/**：放配置相关文件，如 `.env.example` 模板；复制为 `.env` 后本地使用，不提交。
  - **pkg/**（本日未用）：可供其他项目 import 的公共包；若没有对外复用的代码可不建。
- 环境变量：`os.Getenv`、`.env` 文件（godotenv）、**viper**（多源、类型解析）
- 配置结构体：启动时加载，注入到 handler 或 service

**本日说明**：Day 3 只做「项目结构 + 配置」，**不连数据库、不写 SQL**。Config 里的 `DB_DSN` 是**留给 Day 4 数据层用的**：默认写成 `file:./day4/data.db`（SQLite），到 **Day 4** 才会用这份配置真正连库、做 CRUD。这样先在本日把「配置项」定好，下一天直接拿来用。

## 运行

```bash
# 首次或依赖有变更时执行：go mod tidy 根据代码里的 import 同步 go.mod（拉取缺失依赖、移除未用依赖）；本日会拉取 viper（类似 Python 的 dynaconf / pydantic-settings）等
go mod tidy

# 可选：复制示例环境并修改（勿提交 .env 到 Git）
cp day3/config/.env.example day3/config/.env

go run ./day3/cmd/server
# 访问 http://localhost:8080/health 或 http://localhost:8080/info（info 会返回当前 env、port，体现配置注入）
```

## 本日注意点与易踩坑

- **配置加载顺序**：若有 `.env`，需在 `config.Load()` 之前执行 `godotenv.Load(...)`，否则 viper 读到的环境变量里没有 .env 里的值。
- **敏感信息**：`.env` 里常放端口、DB 连接串等，务必加入 `.gitignore`，不要提交到仓库；提交 `.env.example` 作为模板即可。
- **viper 与 os.Getenv**：viper 支持默认值、类型解析（如 HTTP_PORT 直接解析成 int），生产里多环境、多源配置更常用。本目录 main 实际调用的是 `config.Load()`；`config.LoadWithOS()` 仅用 os.Getenv 实现，未在 main 中调用，保留便于对比。
- **import 路径**：代码里 `github.com/go-quickstart-7days/day3/...` 是本项目**本地包路径**，不是从 GitHub 在线拉；详见 [Day 0 六、Go 模块与 import 路径](../day0/README.md#六go-模块与-import-路径新手易懵)。
- **go mod tidy 超时**：若连不上 `proxy.golang.org`（国内常见），可设置 `go env -w GOPROXY=https://goproxy.cn,direct` 后重试；详见 [docs/PITFALLS_AND_SOLUTIONS.md](../docs/PITFALLS_AND_SOLUTIONS.md)「环境与依赖」。
- 更多见：[docs/PITFALLS_AND_SOLUTIONS.md](../docs/PITFALLS_AND_SOLUTIONS.md)（配置/部署相关可后续补充）。

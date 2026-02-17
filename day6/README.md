# Day 6：测试 & 部署

- **单元测试**：`*_test.go`、`testing` 包、**表驱动测试**（多组输入期望用 `[]struct` + 循环）。
- **HTTP 测试**：`net/http/httptest`：`NewRequest` 构造请求、`NewRecorder` 录响应，**不真正起服务**即可测 Handler。
- **Docker**：多阶段构建（先 go build 再拷二进制到小镜像），镜像体积小、部署快；本目录 Dockerfile 用于构建 **Day 7** 的 server。

## 目录结构

| 目录/文件 | 说明 |
| :--- | :--- |
| **handler/** | 被测的 Handler（GetUser）与对应测试（handler_test.go）。 |
| **handler/handler.go** | GetUser：按路径 `/api/users/{id}` 解析 id，合法返回 User JSON，非法返回 400。 |
| **handler/handler_test.go** | 表驱动：多组 path + 期望状态码，用 httptest 调 GetUser 并断言 `rec.Code`。 |
| **Dockerfile** | 多阶段：builder 阶段编译 `./day7/cmd/server`，最终阶段用 alpine 只放二进制 + 运行。 |

```
day6/
├── handler/
│   ├── handler.go        # 被测 Handler
│   └── handler_test.go    # 表驱动 + httptest
├── Dockerfile             # 多阶段构建（产物为 day7 的 server）
├── README.md
└── csdn.md
```

## 运行

**以下命令均在项目根目录执行。**

### 1. 跑测试（必做，无需 Docker）

```bash
go test ./day6/...
```

**`...` 的含义**：表示「当前目录及所有子目录下的包」。本项目中，`.go` 文件不在 `day6/` 这一层，而在子目录 `day6/handler/` 里；若只写 `go test ./day6`，Go 只看 `day6/` 这一层，发现没有 `.go` 也就没有「包」，不会去测 `handler`。写成 `./day6/...` 才会把 `day6/handler` 等子目录里的包都算上，从而跑到 `handler_test.go`。

预期：`PASS`，表示表驱动里多组用例（如 `/api/users/1` 200、`/api/users/abc` 400）均通过。

带详细输出（`-v` 打印每个测试名）：

```bash
go test -v ./day6/...
```

预期输出示例：

```
=== RUN   TestGetUser
--- PASS: TestGetUser (0.00s)
PASS
ok      github.com/go-quickstart-7days/day6/handler     0.229s
```

### 2. 构建并运行 Docker 镜像（可选，需本机已装 Docker）

```bash
docker build -f day6/Dockerfile -t go-7days-app .
docker run -p 8080:8080 go-7days-app
```

镜像里跑的是 **Day 7** 的 server，启动后可访问 `http://localhost:8080` 验证（Day 7 再做完整接口联调）。

## 本日注意点与易踩坑

- **为什么测试里能直接写 GetUser，不用 import？** 因为 `handler_test.go` 和 `handler.go` 都在同一目录且都是 `package handler`，属于**同一个包**。Go 规定：同一包内的所有 .go 文件共享命名空间，所以 `handler.go` 里导出的 `GetUser` 在 `handler_test.go` 里**直接可见**，无需 import，也不能写 `handler.GetUser`（那是「从别的包调」的写法）。和 Python 不同：Go 没有 `__init__.py`，目录下同名的 `package` 声明即构成一个包。
- **测试文件命名**：测试文件必须 `_test.go` 结尾，测试函数必须 `func TestXxx(t *testing.T)`，否则 `go test` 不会执行。
- **表驱动**：把多组「输入 + 期望」放进 `[]struct`，循环里调被测函数并断言，方便加用例、易读。
- **httptest 不占端口**：`NewRequest` + `NewRecorder` 直接调 Handler，不启动 HTTP 服务，速度快、适合 CI。
- **Dockerfile 构建的是 day7**：当前 Dockerfile 里 `go build ... ./day7/cmd/server`，所以本目录侧重「测试写法 + 镜像写法」，实际跑的是 Day 7 应用；若只想验证镜像构建，可先 `docker build`，再 `docker run` 看 8080 是否可访问。
- 更多见 [docs/PITFALLS_AND_SOLUTIONS.md](../docs/PITFALLS_AND_SOLUTIONS.md)。

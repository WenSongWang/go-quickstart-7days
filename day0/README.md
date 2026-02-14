# Day 0：预习——关键词、语法与常用包一览

**这份文档是干啥的**：不写代码、不跑程序，先把 7 天里会出现的**关键词、语法、标准库包**过一遍，看懂**每个是啥、干啥用**。这样后面看 Day 1～7 的代码时，不会一眼全是“天书”。

**适合谁**：完全没写过 Go、甚至没写过代码也没关系。看不懂的地方先跳过，看 Day 1 时再回来查；**不需要一次全懂**，有个印象即可。

---

## 一、7 天路线一览（会碰到什么）

| 天数 | 主题 | 会接触的关键词 / 包 |
| :--- | :--- | :--- |
| **Day1** | Go 基础 | **fmt**、Println/Printf、main、nil、变量、函数、结构体、error、errors.Is/As、make、defer、range、goroutine、channel、WaitGroup、select、context |
| **Day2** | HTTP 服务 | net/http、ListenAndServe、HandleFunc、Request、ResponseWriter、JSON、Marshal/Unmarshal |
| **Day3** | 项目与配置 | cmd、internal、config、环境变量、os.Getenv、.env、godotenv |
| **Day4** | 数据层 | database/sql、连接池、Query/Exec、预编译、SQLite、PostgreSQL |
| **Day5** | 中间件与分层 | 中间件、Handler 包装、鉴权、API Key、context 超时、handler→service→repository |
| **Day6** | 测试与部署 | go test、表驱动测试、httptest、Dockerfile、多阶段构建 |
| **Day7** | 综合实战 | REST API、slog、优雅关闭、Shutdown、健康检查 |

**说明**：每天的内容会依赖前面几天。Day 1 打语言基础，Day 2 用标准库写 HTTP，Day 3 学怎么组织目录和读配置，后面再叠数据库、中间件、测试和部署，Day 7 把前面串成一个完整小项目。

**Day 1 代码里会遇到的词**：打开 `day1/hello/main.go` 会看到 `fmt`、`Println`；`basics` 里还有 `Printf`、`nil`、`var`、`:=`；`errors` 里用 `fmt.Errorf`、`errors.Is`、`errors.As`；`concurrency` 里用 `make`、`defer`、`range`、`go`、`chan`、`context`。下面二、三、四节都能查到对应解释。

---

## 二、关键词速览（按主题，带解释）

### 2.1 基础

| 关键词 | 解释 |
| :--- | :--- |
| **package** | 每个 .go 文件都属于一个包。可执行程序必须是 `package main`，别的包用来被 import。 |
| **import** | 引入其他包，才能用里面的函数和类型，比如 `import "fmt"` 后用 `fmt.Println`。 |
| **fmt（包）** | Day 1 第一个示例就会用到，是 Go 自带的“打印”包。在终端里输出文字就靠它：`fmt.Println(...)` 打一行并换行；`fmt.Printf("格式", 参数)` 按“占位符”打（占位符就是先把位置占住，后面再填数，比如 `%d` 填数字、`%s` 填字符串）。 |
| **Println / Printf** | 都在 `fmt` 包里。Println 简单打一行；Printf 可以按格式打，例如 `"用户: ID=%d Name=%s"` 里 %d、%s 会被后面的参数替换。 |
| **main** | 可执行程序的入口函数必须叫 `main()`，且必须在 `package main` 里，否则 `go run` 不知道从哪执行。 |
| **nil** | 表示“空、没有”。很多类型没赋值时就是 nil。判断“有没有出错”就写 `if err != nil { ... }`（err 不是 nil 就说明出错了）。 |
| **var / :=** | 声明变量。`var x int = 1` 显式类型，`x := 1` 由右边自动推断类型，写起来更短。 |
| **func** | 定义函数。Go 里函数可以返回多个值，通常最后一个会是 `error`，调用方要判断 `if err != nil`。 |
| **struct / type** | `type 名字 struct { 字段 类型 }` 定义“一种数据类型”，把几条信息打包在一起。比如“用户”可以包成：ID + 名字，后面当做一个整体用。 |

### 2.2 错误

| 关键词 | 解释 |
| :--- | :--- |
| **error** | Go 里不用“抛异常”，而是把错误当返回值。函数返回 `(结果, error)`，出错时 error 不是 nil，调用方要自己检查。 |
| **errors.New** | 造一个错误值，比如“未找到”“参数不合法”，后面可以统一用这个名字判断。 |
| **fmt.Errorf 与 %w** | 在已有错误外面再包一层说明，如“查询用户时出错：xxx”。用 `%w` 可以把里面的错误保留住，后面还能用 `errors.Is` / `errors.As` 认出是哪种错。 |
| **errors.Is / errors.As** | `errors.Is(err, 目标)` 判断“是不是某个已知错误”；`errors.As(err, &变量)` 把错误转成更具体的类型，方便取里面的字段（比如哪个字段校验没过）。面试常问。 |

### 2.3 并发

| 关键词 | 解释 |
| :--- | :--- |
| **make** | 用来创建 channel、slice、map。Day 1 会看到 `make(chan int, 2)` 创建能存 2 个整数的 channel，`make(chan struct{})` 创建“信号”用的空 channel。 |
| **close** | 关闭 channel。关闭后不能再往 channel 发数据；接收方 `for v := range ch` 会收完已有数据后自动结束。 |
| **defer** | 把一句代码推迟到当前函数 return 之前执行。常用于 `defer cancel()`、`defer wg.Done()`，保证资源释放或计数减一。 |
| **range** | 用在 channel 上：`for v := range ch` 会一直从 ch 收数据，直到 ch 被 close；Day 1 的 concurrency 示例里会用到。 |
| **goroutine** | 用 `go 函数()` 就能让这个函数“在后台跑”，不卡住当前代码。可以理解成轻量级线程，一个程序里能跑很多个。 |
| **channel** | 用来在多个 goroutine 之间“传数据”。`ch <- 1` 表示往 ch 里发一个 1，`<-ch` 表示从 ch 里收一个数。可以带缓冲（先存几个再收），用完了可以 `close`。 |
| **sync.WaitGroup** | 要等“多个 goroutine 都干完”再往下走时用：`Add(1)`、`Done()`、`Wait()`，避免 main 提前退出。 |
| **select** | 在多个 channel 上“等谁先到”，常和 `time.After` 一起做超时：等业务结果或等超时，谁先到执行谁。 |
| **context** | 在“整条请求链路”里传“取消/超时”的信号。比如设一个 3 秒超时，下游查数据库、调别人接口都能收到“时间到了别干了”。`context.Background()` 是起点；`WithTimeout` 设超时时间；`ctx.Done()`、`ctx.Err()` 用来判断是不是已经取消或超时了。面试高频。 |

### 2.4 HTTP

| 关键词 | 解释 |
| :--- | :--- |
| **ListenAndServe** | 用标准库“起一个网站服务”：监听某个端口（如 8080），有人访问就把请求交给你注册好的处理函数（Handler）。 |
| **HandleFunc** | 给“路径”和“处理函数”做绑定。比如把 `/hello` 绑到一个函数上，以后有人访问 `/hello` 就自动调这个函数。 |
| **Request / ResponseWriter** | 处理函数会拿到“请求”和“写响应的笔”。从 Request 里能读：用什么方法访问的、访问的地址、请求头、请求体；用 ResponseWriter 可以写状态码和返回给前端的文字或 JSON。 |
| **Marshal / Unmarshal** | 把 Go 里的结构体变成 JSON 字符串发出去叫 Marshal；把收到的 JSON 字符串变回结构体叫 Unmarshal。和前端、手机端交换数据时常用。 |

### 2.5 配置

| 关键词 | 解释 |
| :--- | :--- |
| **os.Getenv** | 读环境变量，如端口、数据库连接串。部署时改环境变量即可，不用改代码。 |
| **.env / godotenv** | 本地开发把配置写在 `.env` 文件里，用 godotenv 加载到环境变量，方便且不提交到 Git。 |
| **配置结构体 / 注入** | 启动时把端口、数据库地址等读进一个“配置结构体”，再把这个结构体当作参数传给各个 handler、service，这样不用在代码里到处写 `os.Getenv`。 |

### 2.6 数据库

| 关键词 | 解释 |
| :--- | :--- |
| **sql.Open / DSN** | 用“驱动名 + 连接串（DSN）”连数据库。DSN 就是一串描述“连哪台机器、哪个库、用户名密码”的字符串，一般从配置里读，不写死在代码里。标准库会帮你维护“连接池”（多个连接复用）。 |
| **Query / QueryRow / Exec** | 查多行用 `Query`，查一行用 `QueryRow`，增删改用 `Exec`。查出来的结果用 `Scan` 扫进变量或结构体里。 |
| **预编译 / 参数化查询** | 写 SQL 时用 `?` 或 `$1` 当占位符，参数单独传进去，不要用字符串拼 SQL（既危险又难维护）。 |

### 2.7 中间件

| 关键词 | 解释 |
| :--- | :--- |
| **包装 Handler** | 中间件就是一个函数：你给它“下一个处理函数”，它返回“一个新的处理函数”。新函数里先干自己的事（打日志、检查是否登录、设超时），再调用“下一个”。 |
| **链式调用** | 一个请求会依次经过：先经过日志中间件 → 再经过鉴权中间件 → 再经过超时中间件 → 最后才到真正的业务逻辑，像一条链。 |
| **鉴权 / API Key** | 从请求头（Header）里读 Token 或 API Key，检查是不是合法的。通过才继续调业务，不通过就直接返回 401（未授权）。 |

### 2.8 测试与部署

| 关键词 | 解释 |
| :--- | :--- |
| **_test.go / testing** | 测试代码写在 `xxx_test.go` 里，函数名以 `Test` 开头，用 `go test` 命令跑。“表驱动”就是把多组“输入 + 期望输出”放一张表里，循环跑，少写重复代码。 |
| **httptest** | 测 HTTP 处理函数时，不用真的起一个网站。用 httptest 可以“假发请求、录下响应”，在内存里就测完。 |
| **Dockerfile / 多阶段构建** | Dockerfile 用来把程序打成“镜像”，到处都能跑。多阶段构建就是：先在一个阶段里把 Go 编译成二进制，再在另一个阶段里只拷贝这个二进制到一个很小的镜像里，这样最终镜像又小又容易部署。 |

---

## 三、语法速览（写代码时会反复看到）

| 语法 | 长什么样 | 解释 |
| :--- | :--- | :--- |
| 包声明 | `package main` 或 `package xxx` | 可执行程序必须是 `main` 包且要有 `main()`；被别的包 import 的用其他名字（如 `package config`）。 |
| 导入 | `import "fmt"` 或 `import ( "a" "b" )` | 引入后可用该包里的公开函数、类型。括号写法一次导入多个包。 |
| 变量 | `var x int = 1` 或 `x := 1` | 前者显式类型，后者自动推断。函数内常用 `:=`，包级变量常用 `var`。 |
| 函数 | `func 名(参数) 返回值 { }`，多返回值 `(int, error)` | Go 习惯把错误作为最后一个返回值，调用处必须检查 `err != nil`，否则容易漏错。 |
| 结构体 | `type User struct { ID int; Name string }` | 把一组字段绑成一种类型，用来表示请求体、响应、数据库一行等。 |
| 错误处理 | `if err != nil { return err }`、`fmt.Errorf("xxx: %w", err)` | 几乎每个返回 error 的函数都要处理；包装时用 `%w` 才能被 `errors.Is`/`As` 识别。 |
| 并发 | `go f()`、`ch := make(chan int)`、`select { case <-ch: }` | `go` 让函数在后台跑；channel 在多个“后台任务”之间传数据；`select` 可以“等好几个 channel，谁先到就处理谁”。 |
| context | `ctx, cancel := context.WithTimeout(...)`、`<-ctx.Done()` | 创建一个“带超时”的 context，在整条调用链里传。`ctx.Done()` 能收到“该停了”的信号（超时或取消）。 |

---

## 四、常用标准库包（本系列会用到）

| 包名 | 常见用途 | 解释 |
| :--- | :--- | :--- |
| **fmt** | 打印、格式化、`Errorf` | Day 1 的 hello 就会用到：`fmt.Println` 打一行，`Printf` 按格式打；`Sprintf` 拼字符串不输出，`Errorf` 包装错误（配合 `%w`）。 |
| **errors** | `New`、`Is`、`As` | 定义和判断错误，和 `fmt.Errorf` 的 `%w` 一起用，实现错误链。 |
| **context** | 超时、取消 | 请求超时、下游取消都靠 context 传递，避免泄漏 goroutine。 |
| **sync** | `WaitGroup`、Mutex | 等多协程结束用 WaitGroup；多协程抢同一资源时用 Mutex（本系列 Day 1 会用到 WaitGroup）。 |
| **net/http** | 服务、路由、Request/Response | 用 Go 写网站、写接口就靠这个包，不用装别的框架也能写“按 URL 区分接口”的 API。 |
| **encoding/json** | Marshal、Unmarshal | 把结构体变成 JSON 字符串，或把 JSON 字符串变回结构体。和前端、别的服务交换数据时常用。 |
| **database/sql** | 连接、Query、Exec | 通用数据库接口，配不同驱动（如 pq、sqlite3）连不同数据库。 |
| **log/slog** | 结构化日志 | Go 1.21 起标准库自带，打出来的日志便于采集、检索，替代散乱的 `Printf`。 |
| **os** | 环境变量、信号 | `Getenv` 读配置；监听 SIGINT/SIGTERM 做优雅关闭时会用到。 |
| **time** | 睡眠、超时 | `time.Sleep`、`time.After` 常和 `select` 一起做超时。 |
| **testing** | 单元测试 | `TestXxx`、`t.Run`、表驱动测试，本系列 Day 6 会写。 |

---

## 五、本系列用到的第三方包（仅作了解）

| 包 | 用途 | 解释 |
| :--- | :--- | :--- |
| **github.com/joho/godotenv** | 加载 .env | 把 `.env` 文件里的键值对加载进环境变量，本地开发不用在系统里配一堆 env。 |
| **github.com/lib/pq** | PostgreSQL 驱动 | 实现 `database/sql` 的接口，连 Postgres 时必须 `import _ "github.com/lib/pq"` 注册驱动。 |
| **github.com/mattn/go-sqlite3** | SQLite 驱动 | 同上，连 SQLite 用，本系列 Day 4 可选 SQLite 无需装数据库。 |

---

## 六、建议（怎么用这份文档）

1. **先从头到尾扫一遍**，不用背。目标是：看到某个词时，能想起来“在 Day 0 哪一节见过、大概是干啥的”。
2. **看 Day 1 代码时**：打开 `day1/hello/main.go`，遇到 `package`、`import`、`fmt`、`Println` 就回本页「二、关键词」或「三、语法」「四、常用包」查一下，一行一行对上号。
3. **后面某天卡住**：直接回本页，用 Ctrl+F 搜那个词（如 `context`、`Marshal`），找到对应解释再看一遍。

**不用一次全懂**，有个印象、知道“去哪查”就够。看完就可以开始 [Day 1：Go 基础](../day1/README.md) 了。

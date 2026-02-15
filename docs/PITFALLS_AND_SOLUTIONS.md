# 实战开发踩坑注意点与解法

本文档按主题整理本系列及 Go 后端实战里**常见坑、原因和解法**，便于开发时少踩坑、面试时能讲清「遇到过什么问题、怎么解决的」。

**融入日常**：各天 `dayN/README.md` 会有「本日注意点与易踩坑」小节，列出当日相关的 1～3 条并链到本文档对应主题，学当天内容时就能看到，不必单独翻 docs。

---

## 容易踩的坑 / 高频 bug 覆盖总表

| 类别 | 坑 / 高频 bug | 是否涉及 | 本项目中的位置或计划 |
| :--- | :--- | :--- | :--- |
| **HTTP** | 先写 Body 再设 Content-Type | ✅ 已涉及 | day2 handleGetUser、博文逐段解读 |
| **HTTP** | 路由末尾斜杠（前缀 vs 精确匹配） | ✅ 已涉及 | day2 main、博文 |
| **HTTP** | handler 返回后仍写 ResponseWriter（如 goroutine 里写） | ⚠️ 后续补 | Day 5/7 或第四节 |
| **语言** | 闭包捕获循环变量（goroutine 里打出 3,3,3） | ✅ 已涉及 | day1/concurrency、博文 |
| **语言** | Itoa 误解为 ASCII、select 与 Done 混淆、打印顺序误当 bug | ✅ 已涉及 | day1/day2、rule/skill |
| **错误** | 包装错误用 %v 导致 Is/As 失效 | ✅ 已涉及 | day1/errors、博文 |
| **错误** | 返回「无结果」用值类型导致零值歧义 | ✅ 已涉及 | day1/basics 返回 *User |
| **错误** | 未检查 err 就用返回值（nil 指针） | ✅ 已涉及 | day1 多处 if err != nil |
| **并发** | context 未 cancel 导致 goroutine 泄漏 | ⚠️ 后续补 | Day 5/7 中间件、第四节 |
| **并发** | 多 goroutine 写同一 map 未加锁 | ⚠️ 后续补 | 第四节或 part2 |
| **并发** | channel 未 close 导致 for range 死等 | ⚠️ 后续补 | day1 有 close 示例，可单列坑 |
| **数据库** | SQL 注入（拼接 SQL）、事务未 Commit/Rollback | ⚠️ 后续补 | Day 4、第四节 |
| **配置** | godotenv 在 config.Load 之后执行，viper 读不到 .env | ✅ 已涉及 | day3 main：先 godotenv 再 Load；博文易踩坑 |
| **配置/部署** | 敏感信息进代码或 Git、未优雅关闭 | ✅/.env 已提醒 | Day 3 README/csdn：.env 勿提交；第四节 |
| **环境/依赖** | go mod tidy 连不上 proxy.golang.org（国内常见） | ✅ 已涉及 | docs 本节「环境与依赖」、day3 README、根 README 环境要求 |

**结论**：**容易踩的坑**和**高频 bug** 里，与 Day 1～2 强相关的（HTTP 头顺序、路由、闭包、错误链、nil 指针）已写入正文并融入教学演练；**并发/数据库/部署/gin** 相关的高频坑已在「第四节」列出，会随 Day 4～7 和 part2 推进时补进正文表格。

---

## 一、HTTP 与响应

| 坑 | 现象 / 原因 | 解法 | 本项目位置 |
| :--- | :--- | :--- | :--- |
| **先写 Body 再设 Content-Type** | 想返回 JSON 错误却先调了 `http.Error()`，再 `w.Header().Set("Content-Type", "application/json")`，浏览器仍按纯文本解析。原因：`http.Error` 内部会调 `WriteHeader()`，**响应头必须在第一次 WriteHeader 之前全部设好**，之后不可再改。 | 先 `w.Header().Set("Content-Type", "application/json")`，再 `http.Error(...)`。 | day2/server handleGetUser 错误分支 |
| **路由末尾有没有斜杠搞混** | 注册 `/api/users/`（带斜杠）以为只能匹配 `/api/users/`，结果 `/api/users/1` 也能进。标准库规定：pattern 末尾有 `/` 为**前缀匹配**，无 `/` 为**精确匹配**。 | 明确要前缀匹配就写 `/api/users/`；要精确就写 `/api/users`。 | day2/server main、day2 csdn 解读 |

---

## 二、语言与标准库

| 坑 | 现象 / 原因 | 解法 | 本项目位置 |
| :--- | :--- | :--- | :--- |
| **Itoa 记成 "integer to ASCII"** | 注释或口述写错，误导他人。`strconv.Itoa` 是 **Integer to string**（整型转十进制字符串），返回 `string`，不是 ASCII 类型。 | 记准：Itoa = integer to string；拼接是 string + string。 | day2 注释规范、rule/skill |
| **闭包捕获循环变量** | `for i := 0; i < 3; i++ { go func() { fmt.Println(i) }() }` 可能打出 3,3,3。goroutine 启动时捕获的是**同一个 i**（循环结束后的值）。 | 传参：`go func(id int) { ... }(i)`，把当前 i 的值传进去。 | day1/concurrency WaitGroup 示例 |
| **select 里变量名与 wg.Done 混淆** | 用 `done := make(chan struct{})` 时，和前面的 `wg.Done()` 方法名易混，读代码时误解。 | 用不同名字，如 `ready := make(chan struct{})` 表示「就绪」信号。 | day1/concurrency、day1 csdn 易混淆点 |
| **goroutine 打印顺序不固定** | 多个 goroutine 并发打印 "worker 0/1/2"，顺序每次运行可能不同，误以为 bug。 | 这是正常并发行为，不是错误。 | day1/concurrency、day1 csdn |

---

## 三、错误处理

| 坑 | 现象 / 原因 | 解法 | 本项目位置 |
| :--- | :--- | :--- | :--- |
| **包装错误没用 %w** | 用 `fmt.Errorf("xxx: %v", err)` 包装后，`errors.Is(err, ErrNotFound)` 匹配不到。原因：只有 `%w` 才会保留错误链。 | 包装时用 `fmt.Errorf("xxx: %w", err)`，调用方才能用 `errors.Is`/`errors.As`。 | day1/errors |
| **返回 nil 却用值类型** | 出错时想返回「没有结果」；若返回类型是 `User`（值），无法表示「没有」，只能返回零值，调用方难以区分「查到了零值」和「没查到」。 | 返回指针 `*User`，出错时返回 `nil, err`；调用方先判 err 再使用返回值。 | day1/basics FindUser |

---

## 四、后续可补充（按 Day 4～7 与 part2 推进时添加）

以下坑可在写到对应天或 part2 项目时，按实际遇到再补进上表或本节约简版说明：

- **数据库**：连接未关闭/未用连接池、SQL 注入（未参数化）、事务里忘了 Commit/Rollback、context 超时后连接未归还。
- **并发**：channel 未 close 导致 range 死等、多 goroutine 写同一 map 未加锁、context 未 cancel 导致泄漏。
- **配置**：**加载顺序**：godotenv 必须在 config.Load（viper 读 env）之前执行，否则 .env 未灌入 env；敏感信息写进代码或提交到 Git（.env 勿提交，用 .env.example 做模板）。见 day3。
- **部署**：镜像里用 root、未做优雅关闭、健康检查路径未注册或返回非 200。
- **gin**：绑定与校验失败未统一返回格式、中间件里修改 ResponseWriter 后再次 WriteHeader。

---

## 环境与依赖：go mod tidy 超时（国内常见）

| 坑 | 现象 / 原因 | 解法 | 说明 |
| :--- | :--- | :--- | :--- |
| **go mod tidy 连不上 proxy.golang.org** | 报错 `Get "https://proxy.golang.org/...": dial tcp ... connectex: ... failed to respond`。国内或受限网络访问官方代理超时。 | 使用国内镜像并**持久生效**：`go env -w GOPROXY=https://goproxy.cn,direct`，再执行 `go mod tidy`。可选镜像：`goproxy.io`、`https://mirrors.aliyun.com/goproxy/`。 | 不要保留 `proxy.golang.org` 在 GOPROXY 里，否则会回退到官方源继续超时；用 `,direct` 即可。 |
| **若曾改用 go-sqlite3 出现 CGO 报错** | 本仓库 Day 4 SQLite 已改用**纯 Go 驱动**（modernc.org/sqlite），无需 CGo/gcc，直接 `go run ./day4/sqlite` 即可。若自己改成 go-sqlite3 后报 CGO，需安装 gcc 并设 `CGO_ENABLED=1`。 | 使用仓库默认的 modernc 驱动即可；或见上条安装 gcc 并开启 CGO。 | day4/sqlite 默认 modernc.org/sqlite，无 CGo 要求。 |
| **protocol error: received DATA after END_STREAM** | 拉取依赖时偶现，多为代理/网络瞬时问题。 | 一般无需处理，**再次执行 `go mod tidy`** 即可；若反复出现可换一个镜像（如 goproxy.io）。 | 本经验来自国内环境实践。 |

---

## 五、使用建议

1. **开发时**：写某一块前扫一眼对应主题，少犯表中列出的错。  
2. **排错时**：现象类似可先对表查「坑」和「解法」。  
3. **面试时**：可说「在练手项目里遇到过 XX 坑，用 YY 解法解决」，并指向本项目具体位置（如 day2 handleGetUser）。  
4. **协作**：新增坑与解法时按上表格式追加，并注明对应 day 或 part2 项目与文件。

本系列 day 或 part2 的博文/README 里若写了「易混淆点」「注意点」，也会在本文档中做对应条目或引用，保持一处汇总、多处可查。

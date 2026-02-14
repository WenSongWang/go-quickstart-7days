# Day 1：Go 基础

第一天打好环境、语法、错误处理与并发基础，后续写 HTTP、中间件、数据层都会用到。

## 内容概览

| 模块 | 内容 |
| :--- | :--- |
| 环境 | 安装 Go、`go run` / `go build` |
| 语法 | 变量、函数、结构体、多返回值、`error` |
| 错误处理 | `errors.Is` / `errors.As`、错误包装 `fmt.Errorf("%w", err)`（面试常问） |
| 并发 | goroutine、channel、`sync.WaitGroup`、`select`、`context` 取消与超时（面试高频） |

## 示例与知识点

| 示例目录 | 主要知识点 |
| :--- | :--- |
| `hello/` | 最简程序：`package main`、`main()`、`fmt.Println` |
| `basics/` | 变量与短声明、结构体、多返回值 + error、`if err != nil` |
| `errors/` | 错误包装 `%w`、`errors.Is` / `errors.As`、自定义错误类型 |
| `concurrency/` | goroutine、channel、`sync.WaitGroup`、`select`、`context.WithTimeout` |

## 运行

在**项目根目录**（`go-quickstart-7days`）下执行：

```bash
go run ./day1/hello
go run ./day1/basics
go run ./day1/errors
go run ./day1/concurrency
```

建议按顺序跑：hello → basics → errors → concurrency。

## 学习建议

- 先跑通再改：改变量名、返回值或错误信息，再 `go run` 看效果。
- 重点看 `errors` 的错误链与 `concurrency` 的 context，面试常考。

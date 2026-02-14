// Day1 示例：最简 Go 程序
// 认识：main 包、import、main 函数、fmt.Println

package main // 可执行程序必须用 main 包，否则无法 go run / go build

import "fmt" // 导入标准库 fmt，用来在终端里“打印”文字

func main() { // 程序入口：从 main 函数开始执行
	fmt.Println("Hello, Go! 7天入门go后端开发") // 打印一行字符串到终端，并换行
}

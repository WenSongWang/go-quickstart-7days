// Day1 示例：变量、函数、结构体、错误处理
// 本文件学习：var/:=、多返回值+error、结构体、if err != nil
package main

import (
	"errors"
	"fmt"
)

// User 结构体：把“用户”的多个字段绑在一起，后端里常用来表示请求/响应或数据库一行
type User struct {
	ID   int    // 用户 ID
	Name string // 用户名
}

// FindUser 根据 id 查用户。
//
// 【重点】返回的是 *User（指针），不是 User（值）。原因：
//   - 出错时没有“用户”可返回，只能返回 nil；只有指针类型才能表示“没有”（值类型没有 nil）。
//   - 调用方拿到 (*User, error)，先判断 err != nil，再安全地用 u.ID、u.Name（此时 u 一定非 nil）。
func FindUser(id int) (*User, error) {
	if id <= 0 {
		// errors.New 创建一个错误值，表示“参数不合法”
		return nil, errors.New("invalid user id")
	}
	// &User{...} 取结构体的地址，得到 *User 指针，因为函数签名要求返回 *User
	return &User{ID: id, Name: "张三"}, nil
}

func main() {
	// ---------- 变量与短声明 ----------
	// var 变量名 类型 = 值：显式写类型
	var a int = 1
	// 短声明 变量名 := 值：由右边自动推断类型，函数内常用
	b := 2
	fmt.Println("a+b =", a+b)

	// ---------- 调用多返回值函数并处理错误 ----------
	// u 的类型是 *User（指针），err 是 error。err 是 Go 里错误返回值的常规命名（大家都会这么写），不是关键字。
	// 出错时 u 为 nil，所以必须先判 err 再用 u。:= 短声明会根据 FindUser 的返回值自动推断 u 为 *User、err 为 error。
	u, err := FindUser(1)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}
	// 这里 u 一定非 nil，可以安全用 u.ID、u.Name（指针用 . 访问字段，和值一样）,%d 是数字占位符，%s 是字符串占位符,还有\n换行符
	fmt.Printf("用户: ID=%d Name=%s\n", u.ID, u.Name)

	// ---------- 故意触发错误分支 ----------
	// 用 _ 忽略第一个返回值，只关心 err
	_, err = FindUser(-1)
	if err != nil {
		fmt.Println("预期错误:", err)
	}
}

// Day1 示例：错误包装、errors.Is / errors.As（面试常问）
// 本文件学习：定义“已知错误”、用 %w 包装、用 Is/As 判断错误类型
package main

import (
	"errors"
	"fmt"
)

// 定义“已知错误”常量，方便在代码里统一判断（如“是不是未找到”）
var ErrNotFound = errors.New("not found")
var ErrInvalidID = errors.New("invalid id")

// ValidationError 自定义错误类型：带字段，便于 errors.As 取出具体信息

type ValidationError struct {
	Field string // 哪个字段校验失败
}

// Error 实现 error 接口，这样 *ValidationError 就是一种 error
// e 是 ValidationError 类型的变量，Error() 方法返回一个字符串，表示错误信息
func (e *ValidationError) Error() string {
	return "validation error: " + e.Field
}

// FindUser 根据 id 返回不同错误，用 fmt.Errorf("%w", err) 包装，保留错误链
func FindUser(id int) error {
	if id <= 0 {
		// %w 会把 ErrInvalidID 包在里层，外面还能用 errors.Is 认出来
		return fmt.Errorf("find user: %w", ErrInvalidID)
	}
	if id > 100 {
		return fmt.Errorf("find user: %w", ErrNotFound)
	}
	if id == 50 {
		return fmt.Errorf("find user: %w", &ValidationError{Field: "id"})
	}
	return nil
}

func main() {
	// ---------- errors.Is：判断是不是某个“已知错误”（即使被 %w 包了一层也能认） ----------
	err := FindUser(-1)
	if errors.Is(err, ErrInvalidID) {
		fmt.Println("检测到 ErrInvalidID:", err)
	}

	err = FindUser(101)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("检测到 ErrNotFound:", err)
	}

	// ---------- errors.As：把错误转成具体类型，取出里面的字段 ----------
	err = FindUser(50)
	var valErr *ValidationError//valErr 是 ValidationError 类型的变量，用来存储错误信息
	if errors.As(err, &valErr) {//errors.As(err, &valErr) 把错误转成具体类型，取出里面的字段
		fmt.Println("校验错误字段:", valErr.Field)
	}
}

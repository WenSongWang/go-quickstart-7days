// Day1 示例：goroutine、channel、sync.WaitGroup、select、context（面试高频）
package main
// 本文件学习：go	、make(chan)、defer、range、WaitGroup(wg)、select(ch)、context
// go 启动一个后台任务(go func()),可以理解为轻量级线程，一个程序里能跑很多个
// make(chan) 创建一个 channel,channel是用来在多个goroutine之间传递数据的
// defer 延迟执行(defer func()),可以理解为栈，先进后出，常用来释放资源，比如关闭文件、关闭数据库连接等
// range 遍历 channel
// WaitGroup(wg) 等待多个 goroutine 完成(wg.Add(1)、wg.Done()、wg.Wait())，可以理解为计数器，用来等待多个goroutine完成，常用来做多个goroutine之间的同步	
// select 多路复用，谁先到执行谁，常用来做超时，也可以用来做多个goroutine之间的同步(select { case <-ch: }、select { case <-time.After(100 * time.Millisecond): })
// context 取消与超时(context.WithTimeout(context.Background(), 100 * time.Millisecond))，可以理解为上下文，用来传递取消与超时信号，常用来做超时，也可以用来做多个goroutine之间的同步

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// ---------- 1. goroutine + channel：后台发数据，主流程收 ----------
	// make(chan int, 2) 创建能存 2 个 int 的 channel，带缓冲(ch := make(chan int, 2))
	ch := make(chan int, 2)
	go func() {
		ch <- 1 // 往 channel 里发 1
		ch <- 2 // 往 channel 里发 2
		close(ch) // 发完后关闭，接收方 for range 会收完已有的就结束
	}()
	// for v := range ch 会一直从 ch 收数据，直到 ch 被 close(for v := range ch { ... })
	for v := range ch {
		fmt.Println("channel:", v)
	}

	// ---------- 2. sync.WaitGroup：等 3 个 goroutine 都干完再继续 ----------
	// 多个 goroutine 是并发执行的，所以 "worker 0/1/2" 的打印顺序不固定，每次运行可能不同（如 2 0 1）
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			fmt.Println("worker", id)
		}(i) // 把 i 传进去，避免闭包捕获同一个 i
	}
	wg.Wait()
	fmt.Println("all workers done")

	// ---------- 3. select：多路等待，谁先到执行谁（常用来做超时） ----------
	// 变量名用 ready 而不是 done，避免和上面的 wg.Done() 方法混淆（Done 是方法名，这里是 channel 变量）
	ready := make(chan struct{}) // 空结构体 channel 只当“信号”用，不传数据
	go func() {
		time.Sleep(150 * time.Millisecond)//50 和150 对比下面的超时时间100
		close(ready) // 关闭 channel，接收方 <-ready 会收到“已关闭”的信号
	}()
	select {
	case <-ready://从channel里接收数据，相当于收到通知，我已经干完了, <-ready是关键，相当于从channel里接收数据
		fmt.Println("done")
	case <-time.After(100 * time.Millisecond):
		// 100ms 后还没收到 ready 信号，就走这里（超时）
		fmt.Println("timeout")
	}

	// ---------- 4. context：带超时的上下文，下游可检测“该停了” ----------
	// 流程简述：创建一个“20ms 后自动取消”的 context，传给 doWithContext；里面用 select 等“完成”或“取消”，谁先到执行谁。
	// context.Background() 是根 context（无超时）；WithTimeout 包一层，加上 20ms 的“死线”。
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel() // 用完必须调 cancel() 释放内部计时器，defer 保证函数退出前一定会执行

	result := doWithContext(ctx)
	fmt.Println("context result:", result)
}

// doWithContext 模拟一个“依赖 context”的操作：select 等两路——要么 10ms 内“干完”，要么 ctx 先超时/取消。
// 因为 ctx 是 20ms 超时，而这里 10ms 就“完成”，所以通常走 time.After 分支返回 "ok"；
// 若把 10 改成 30（模拟干活更久），就会先触发 ctx 超时，走 ctx.Done() 返回 "cancelled"。
func doWithContext(ctx context.Context) string {
	select {
	case <-ctx.Done():
		// ctx 超时或被取消时，Done() 会收到信号；Err() 里是原因（如 context.DeadlineExceeded）
		return "cancelled: " + ctx.Err().Error()
	case <-time.After(30 * time.Millisecond):
		// 模拟 10ms 干完活；实际项目里这里可能是查 DB、调 RPC，它们会监听 ctx 决定是否提前结束
		return "ok"
	}
}

package safe

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
)

// WaitGroup 提供带错误收集和panic恢复的并发等待组
// 包含互斥锁用于保护错误切片的并发写入
type WaitGroup struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	errs []error
}

// NewWaitGroup 创建并返回新的安全等待组实例
func NewWaitGroup() *WaitGroup {
	return &WaitGroup{}
}

// Go 在独立goroutine中执行任务函数，提供上下文传递和错误捕获功能
// 参数:
//
//	ctx: 传递给任务函数的上下文对象，用于跨goroutine的取消信号传递
//	f: 需要并发执行的任务函数，接收context.Context参数并返回error
//
// 机制:
//  1. 使用双defer结构确保panic恢复处理在等待组计数器递减前执行
//  2. 自动捕获任务函数的panic和返回错误，记录带堆栈信息的PanicError
//  3. 线程安全地聚合多个goroutine的错误信息
func (i *WaitGroup) Go(ctx context.Context, f func(ctx context.Context) error) {
	i.wg.Add(1)

	go func() {
		var err error
		defer i.wg.Done()

		// 双defer执行顺序：
		// 1. 当前defer处理错误收集
		// 2. 外层defer执行Done操作
		defer func() {
			// panic恢复处理：将panic转换为带有堆栈跟踪的错误对象
			if r := recover(); r != nil {
				err = &PanicError{Stack: string(debug.Stack()), Value: r}
			}

			// 并发安全地追加错误到共享切片
			if err != nil {
				i.mu.Lock()
				i.errs = append(i.errs, err)
				i.mu.Unlock()
			}
		}()

		// 执行用户定义的任务函数并捕获常规错误
		err = f(ctx)
	}()
}

// Wait 阻塞直到所有goroutine完成，返回合并的错误结果
// 返回:
//
//	error: 使用errors.Join组合的错误对象，包含所有任务错误和panic信息
//
// 特性:
//  1. 等待期间会持续阻塞直到所有关联的goroutine执行完毕
//  2. 返回的错误包含原始错误顺序和嵌套结构
func (i *WaitGroup) Wait() error {
	i.wg.Wait()
	return errors.Join(i.errs...)
}

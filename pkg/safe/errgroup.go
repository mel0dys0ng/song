package safe

import (
	"context"
	"runtime/debug"

	"golang.org/x/sync/errgroup"
)

// ErrorGroup 包装了errgroup.Group和context，提供安全的并发控制机制
// 在标准errgroup基础上增加panic捕获功能，防止goroutine崩溃
type ErrorGroup struct {
	eg  *errgroup.Group
	ctx context.Context
}

// NewErrorGroup 创建新的并发错误控制组
// ctx: 父级上下文，用于传播取消信号和超时控制
// 返回: 初始化后的ErrorGroup指针，包含派生的上下文和错误组
func NewErrorGroup(ctx context.Context) *ErrorGroup {
	eg, ctx := errgroup.WithContext(ctx)
	return &ErrorGroup{
		eg:  eg,
		ctx: ctx,
	}
}

// Go 启动一个受保护的goroutine执行任务
// f: 需要执行的上下文感知函数，接收派生上下文并返回可能发生的错误
// 注：自动捕获任务执行过程中的panic，转换为可处理的错误对象
func (i *ErrorGroup) Go(f func(ctx context.Context) (err error)) {
	i.eg.Go(func() (err error) {
		// 防御性恢复机制，确保goroutine崩溃不会导致整个程序退出
		// 捕获panic后生成包含堆栈信息的复合错误，便于问题追踪
		defer func() {
			if r := recover(); r != nil {
				err = &PanicError{Stack: string(debug.Stack()), Value: r}
			}
		}()
		return f(i.ctx)
	})
}

// Wait 等待所有goroutine完成并返回首个错误
// 返回: 遇到的第一个错误（包含可能的panic转换错误），正常完成时返回nil
func (i *ErrorGroup) Wait() error {
	return i.eg.Wait()
}

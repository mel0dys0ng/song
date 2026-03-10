package safe

import (
	"context"
	"runtime/debug"
)

// F 是一个通用函数执行包装器，用于捕获和处理panic，并返回带有堆栈信息的错误
// 参数:
//   - ctx context.Context: 上下文对象，用于传递请求作用域的值、取消信号和截止时间
//   - f  func(context.Context) (T, error): 需要执行的目标函数，接收上下文并返回泛型结果和错误
//
// 返回值:
//   - res T: 目标函数执行成功后的返回结果
//   - err error: 执行过程中发生的错误，包含panic转换的错误及堆栈信息
func F[T any](ctx context.Context, f func(ctx context.Context) (T, error)) (res T, err error) {
	// 使用defer实现panic捕获机制，确保程序不会因未处理的panic而崩溃
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{Stack: string(debug.Stack()), Value: r}
		}
	}()
	// 执行目标函数并正常返回结果
	return f(ctx)
}

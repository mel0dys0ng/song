package result

type (
	// Result 是一个泛型结构体，包含 data、err 字段
	Result[T any] struct {
		data T
		err  error
	}

	Interface[T any] interface {
		Data() T
		SetData(value T)
		Err() error
		SetErr(err error)
	}
)

func New[T any](data T, err error) Interface[T] {
	return &Result[T]{data: data, err: err}
}

func Success[T any](data T) Interface[T] {
	return &Result[T]{data: data, err: nil}
}

func Error[T any](err error) Interface[T] {
	var zero T
	return &Result[T]{data: zero, err: err}
}

// Data 返回 data 字段的值，如果 receiver 为 nil，则返回零值
func (r *Result[T]) Data() T {
	if r == nil {
		var zero T
		return zero
	}
	return r.data
}

// SetData 设置 data 字段的值，如果 receiver 为 nil，则不设置
func (r *Result[T]) SetData(value T) {
	if r == nil {
		return
	}
	r.data = value
}

// Err 返回 err 字段的值，如果 receiver 为 nil，则返回 nil
func (r *Result[T]) Err() error {
	if r == nil {
		return nil
	}
	return r.err
}

// SetErr 设置 err 字段的值，如果 receiver 为 nil，则不设置
func (r *Result[T]) SetErr(value error) {
	if r == nil {
		return
	}
	r.err = value
}

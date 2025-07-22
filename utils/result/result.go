package result

type (
	// Result 是一个泛型结构体，包含 data、ok 和 err 字段
	Result[T any] struct {
		data T
		err  error
	}

	Interface[T any] interface {
		GetData() T
		SetData(value T)
		GetErr() error
		SetErr(err error)
		// Ok if err == nil return true else return false
		Ok() bool
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

// GetData 返回 data 字段的值，如果 receiver 为 nil，则返回零值
func (r *Result[T]) GetData() T {
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

// GetErr 返回 err 字段的值，如果 receiver 为 nil，则返回 nil
func (r *Result[T]) GetErr() error {
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

// Ok if err == nil return true else return false
func (r *Result[T]) Ok() bool {
	return r.GetErr() == nil
}

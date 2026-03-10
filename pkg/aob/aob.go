package aob

// VarOrVar if c == true return a else return b
func VarOrVar[T any](c bool, a, b T) T {
	if c {
		return a
	}
	return b
}

// FuncOrFunc if c == true return a() else return b()
func FuncOrFunc[T any](c bool, a, b func() T) T {
	if c {
		return a()
	}
	return b()
}

// VarOrFunc if c == true return v else return f()
func VarOrFunc[T any](c bool, v T, f func() T) T {
	if c {
		return v
	}
	return f()
}

// FuncOrVar if c == true return f() else return v
func FuncOrVar[T any](c bool, f func() T, v T) T {
	if c {
		return f()
	}
	return v
}

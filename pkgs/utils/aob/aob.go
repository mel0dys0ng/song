package aob

// Aorb if c == true return a else return b
func Aorb[T any](c bool, a, b T) T {
	if c {
		return a
	}
	return b
}

// AorbFunc if c == true return a() else return b()
func AorbFunc[T any](c bool, a, b func() T) T {
	if c {
		return a()
	}
	return b()
}

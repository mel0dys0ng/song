package safe

import "fmt"

type PanicError struct {
	Value any
	Stack string
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic occurred: %+v\nstack trace:\n%s", e.Value, e.Stack)
}

func (e *PanicError) Unwrap() error {
	err, _ := e.Value.(error)
	return err
}

func (e *PanicError) StackTrace() string {
	return e.Stack
}

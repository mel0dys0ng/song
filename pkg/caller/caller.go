package caller

import (
	"fmt"
	"runtime"
)

type Caller struct {
	pc   uintptr
	file string
	line int
	ok   bool
}

// New return the caller that reports file and line number information about
// function invocations on the calling goroutine's stack
func New(skip int) *Caller {
	pc, file, line, ok := runtime.Caller(skip)
	return &Caller{pc: pc, file: file, line: line, ok: ok}
}

// Location return the caller location
func Location(skip int, isContainFile bool) string {
	c := New(skip)
	if isContainFile {
		return fmt.Sprintf("%s:\n  %s:%d", c.File(), c.Func(), c.Line())
	} else {
		return fmt.Sprintf("%s:%d", c.Func(), c.Line())
	}
}

// File return the file path of caller
func (c *Caller) File() string {
	return c.file
}

// Line return the line number of caller
func (c *Caller) Line() int {
	return c.line
}

// Func return the func name of caller
func (c *Caller) Func() string {
	if c.ok {
		if fn := runtime.FuncForPC(c.pc); fn != nil {
			return fn.Name()
		}
	}
	return ""
}

// Stack return the stack trace of caller
func (c *Caller) Stack() string {
	var pcs [32]uintptr             // 预分配一个固定大小的数组来存储程序计数器
	n := runtime.Callers(1, pcs[:]) // 跳过当前函数，从上一层开始记录

	if n == 0 {
		return "no callers"
	}

	var result string
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		result += fmt.Sprintf("%s:%d in %s\n", frame.File, frame.Line, frame.Function)

		if !more {
			break
		}
	}

	return result
}

func (c *Caller) String() string {
	return fmt.Sprintf("%s:%d in %s", c.File(), c.Line(), c.Func())
}

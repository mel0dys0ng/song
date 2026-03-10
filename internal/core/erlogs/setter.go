package erlogs

import "go.uber.org/zap"

// setLevel 设置日志级别，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setLevel(level Level) {
	if e == nil {
		return
	}
	e.level = level
}

// setErr 设置原始错误，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setErr(err error) {
	if e == nil {
		return
	}
	e.err = err
}

// setKind 设置日志类型，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setKind(kind Kind) {
	if e == nil {
		return
	}
	e.kind = kind
}

// setBizID 设置业务/模块 ID，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setBizID(bizID int32) {
	if e == nil {
		return
	}
	e.bizID = bizID
}

// setBizName 设置业务/模块名称，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setBizName(bizName string) {
	if e == nil {
		return
	}
	e.bizName = bizName
}

// setCode 设置错误编码，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setCode(code int64) {
	if e == nil {
		return
	}
	e.code = code
}

// setMsg 设置概述信息，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setMsg(msg string) {
	if e == nil {
		return
	}
	e.msg = msg
}

// setContent 设置详细内容，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setContent(content string) {
	if e == nil {
		return
	}
	e.content = content
}

// setAt 设置时间戳（纳秒），如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setAt(at int64) {
	if e == nil {
		return
	}
	e.at = at
}

// setSkip 设置跳过的调用栈层数，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setSkip(skip int) {
	if e == nil {
		return
	}
	e.skip = skip
}

// setFields 设置 zap 字段列表，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setFields(fields []zap.Field) {
	if e == nil {
		return
	}
	e.fields = fields
}

// appendFields 追加一个或多个 zap 字段到字段列表，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) appendFields(fields ...zap.Field) {
	if e == nil {
		return
	}
	e.fields = append(e.fields, fields...)
}

// setPCs 设置 program counters 列表，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) setPCs(pcs []uintptr) {
	if e == nil {
		return
	}
	e.pcs = pcs
}

// appendPC 追加单个 program counter 到列表，如果 ErLog 为 nil 则不执行任何操作
func (e *ErLog) appendPC(pc uintptr) {
	if e == nil {
		return
	}
	e.pcs = append(e.pcs, pc)
}

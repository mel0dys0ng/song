package erlogs

import "go.uber.org/zap"

// GetLevel 获取日志级别，如果 ErLog 为 nil 则返回 LevelUnknown
func (e *ErLog) GetLevel() Level {
	if e == nil {
		return LevelUnknown
	}
	return e.level
}

// GetErr 获取原始错误，如果 ErLog 为 nil 则返回 nil
func (e *ErLog) GetErr() error {
	if e == nil {
		return nil
	}
	return e.err
}

// GetKind 获取日志类型，如果 ErLog 为 nil 则返回空字符串
func (e *ErLog) GetKind() Kind {
	if e == nil {
		return ""
	}
	return e.kind
}

// GetBizID 获取业务/模块 ID，如果 ErLog 为 nil 则返回 0
func (e *ErLog) GetBizID() int32 {
	if e == nil {
		return 0
	}
	return e.bizID
}

// GetBizName 获取业务/模块名称，如果 ErLog 为 nil 则返回空字符串
func (e *ErLog) GetBizName() string {
	if e == nil {
		return ""
	}
	return e.bizName
}

// GetCode 获取错误编码，如果 ErLog 为 nil 则返回 0
func (e *ErLog) GetCode() int64 {
	if e == nil {
		return 0
	}
	return e.code
}

// GetMsg 获取概述信息，如果 ErLog 为 nil 则返回空字符串
func (e *ErLog) GetMsg() string {
	if e == nil {
		return ""
	}
	return e.msg
}

// GetContent 获取详细内容，如果 ErLog 为 nil 则返回空字符串
func (e *ErLog) GetContent() string {
	if e == nil {
		return ""
	}
	return e.content
}

// GetAt 获取时间戳（纳秒），如果 ErLog 为 nil 则返回 0
func (e *ErLog) GetAt() int64 {
	if e == nil {
		return 0
	}
	return e.at
}

// GetFields 获取 zap 字段列表，如果 ErLog 为 nil 则返回 nil
func (e *ErLog) GetFields() []zap.Field {
	if e == nil {
		return nil
	}
	return e.fields
}

// GetSkip 获取跳过的调用栈层数，如果 ErLog 为 nil 则返回 0
func (e *ErLog) GetSkip() int {
	if e == nil {
		return 0
	}
	return e.skip
}

// GetPCs 获取 program counters 列表，如果 ErLog 为 nil 则返回 nil
func (e *ErLog) GetPCs() []uintptr {
	if e == nil {
		return nil
	}
	return e.pcs
}

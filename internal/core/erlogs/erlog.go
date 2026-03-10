package erlogs

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/mel0dys0ng/song/internal/core/metas"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	CodeSuccess int64  = 0              // CodeSuccess 成功编码
	CodeUnknown int64  = -1             // CodeUnknown 未知错误编码
	MsgUnknown  string = "系统网络错误，请稍后重试" // MsgUnknown 未知错误概述
	MsgSuccess  string = "success"      // MsgSuccess 成功概述

	KindBiz    Kind = "biz"    // KindBiz 业务日志
	KindSystem Kind = "system" // KindSystem 系统（框架）日志
	KindTrace  Kind = "trace"  // KindTrace Trace日志

	SkipDefault = 3 // SkipDefault 默认跳过调用栈层数
)

type Kind = string

type ErLog struct {
	err     error
	level   Level
	kind    Kind
	bizID   int32
	bizName string
	code    int64
	msg     string
	content string
	at      int64
	fields  []zap.Field
	pcs     []uintptr
	skip    int
}

// Constructor 创建一个新的 ErLog 实例，使用默认配置并应用可选参数
func Constructor(opts ...Option) *ErLog {
	e := &ErLog{}
	e.setErr(errors.New(MsgUnknown))
	e.setLevel(LevelError)
	e.setKind(KindBiz)
	e.setCode(CodeUnknown)
	e.setMsg(MsgUnknown)
	e.setContent(MsgUnknown)
	e.setSkip(SkipDefault)
	e.buildOptions(opts...)
	return e
}

// buildOptions 应用可选参数配置
func (e *ErLog) buildOptions(opts ...Option) {
	for _, opt := range opts {
		opt(e)
	}
}

// Status 设置状态码和消息，返回一个新的 ErLog 副本，日志级别设为 Warn
func (e *ErLog) Status(code int64, msg string) *ErLog {
	if e == nil {
		return nil
	}

	c := e.Clone()
	c.setCode(code)
	c.setMsg(msg)
	c.setLevel(LevelWarn)
	c.setContent(msg)

	return c
}

// Statusf 格式化设置状态码和消息，返回一个新的 ErLog 副本，日志级别设为 Warn
func (e *ErLog) Statusf(code int64, format string, args ...any) *ErLog {
	if e == nil {
		return nil
	}

	text := fmt.Sprintf(format, args...)
	c := e.Clone()
	c.setCode(code)
	c.setMsg(text)
	c.setLevel(LevelWarn)
	c.setContent(text)

	return c
}

// UseStatusIfNot 如果当前ErLog的code和msg为未知值，则使用status的code和msg
func (e *ErLog) UseStatusIfNot(err *ErLog) *ErLog {
	if e == nil || err == nil {
		return e
	}

	if e.GetCode() == CodeUnknown {
		e.setCode(err.GetCode())
		e.setMsg(err.GetMsg())
	}

	return e
}

// Options 应用可选参数配置，返回自身以支持链式调用
func (e *ErLog) Options(opts []Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	return e
}

// AppendFields 添加自定义字段，返回自身以支持链式调用
func (e *ErLog) AppendFields(fields ...zap.Field) *ErLog {
	if e == nil {
		return nil
	}
	e.fields = append(e.fields, fields...)
	return e
}

// Wrap 包装错误，添加额外的上下文信息
func (e *ErLog) Wrap(text string, fields ...zap.Field) *ErLog {
	if e == nil {
		return nil
	}
	e.setContent(fmt.Sprintf("%s: %s", text, e.content))
	e = e.AppendFields(fields...)
	return e
}

// Wrapf 格式化包装错误，添加额外的上下文信息
func (e *ErLog) Wrapf(format string, args ...any) *ErLog {
	if e == nil {
		return nil
	}
	text := fmt.Sprintf(format, args...)
	e.setContent(fmt.Sprintf("%s: %s", text, e.content))
	return e
}

// WrapE 包装错误，添加额外的上下文信息和自定义字段
func (e *ErLog) WrapE(err error, fields ...zap.Field) *ErLog {
	if e == nil {
		return nil
	}
	e.setErr(err)
	e.setContent(fmt.Sprintf("%s: %s", err.Error(), e.content))
	e = e.AppendFields(fields...)
	e.capturePC()
	return e
}

// Unwrap 实现 errors.Unwrap 接口，返回包装的底层错误
func (e *ErLog) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.GetErr()
}

// Is 实现 errors.Is 接口，判断错误是否为指定类型
func (e *ErLog) Is(target error) bool {
	if e == nil {
		return false
	}
	return errors.Is(e.GetErr(), target)
}

// As 实现 errors.As 接口，尝试将错误转换为指定类型
func (e *ErLog) As(target any) bool {
	if e == nil {
		return false
	}
	return errors.As(e.GetErr(), target)
}

// Debug 设置日志级别为 Debug，返回自身以支持链式调用
func (e *ErLog) Debug(opts ...Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	e.level = LevelDebug
	return e
}

// Info 设置日志级别为 Info，返回自身以支持链式调用
func (e *ErLog) Info(opts ...Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	e.level = LevelInfo
	return e
}

// Warn 设置日志级别为 Warn，返回自身以支持链式调用
func (e *ErLog) Warn(opts ...Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	e.level = LevelWarn
	return e
}

// Erorr 设置日志级别为 Error，返回自身以支持链式调用
func (e *ErLog) Erorr(opts ...Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	e.level = LevelError
	return e
}

// Panic 设置日志级别为 Panic，返回自身以支持链式调用
func (e *ErLog) Panic(opts ...Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	e.level = LevelPanic
	return e
}

// Fatal 设置日志级别为 Fatal，返回自身以支持链式调用
func (e *ErLog) Fatal(opts ...Option) *ErLog {
	if e == nil {
		return nil
	}
	e.buildOptions(opts...)
	e.level = LevelFatal
	return e
}

// DebugLog 记录 Debug 级别的日志
func (e *ErLog) DebugLog(ctx context.Context, opts ...Option) {
	if e == nil {
		return
	}
	e.buildOptions(opts...)
	e.level = LevelDebug
	e.RecordLog(ctx)
}

// InfoLog 记录 Info 级别的日志
func (e *ErLog) InfoLog(ctx context.Context, opts ...Option) {
	if e == nil {
		return
	}
	e.buildOptions(opts...)
	e.level = LevelInfo
	e.RecordLog(ctx)
}

// WarnLog 记录 Warn 级别的日志
func (e *ErLog) WarnLog(ctx context.Context, opts ...Option) {
	if e == nil {
		return
	}
	e.buildOptions(opts...)
	e.level = LevelWarn
	e.RecordLog(ctx)
}

// ErrorLog 记录 Error 级别的日志
func (e *ErLog) ErrorLog(ctx context.Context, opts ...Option) {
	if e == nil {
		return
	}
	e.buildOptions(opts...)
	e.level = LevelError
	e.RecordLog(ctx)
}

// PanicLog 记录 Panic 级别的日志
func (e *ErLog) PanicLog(ctx context.Context, opts ...Option) {
	if e == nil {
		return
	}
	e.buildOptions(opts...)
	e.level = LevelPanic
	e.RecordLog(ctx)
}

// FatalLog 记录 Fatal 级别的日志
func (e *ErLog) FatalLog(ctx context.Context, opts ...Option) {
	if e == nil {
		return
	}
	e.buildOptions(opts...)
	e.level = LevelFatal
	e.RecordLog(ctx)
}

// capturePC 捕获当前调用栈的 program counter，用于记录 ErLog 传递路径
func (e *ErLog) capturePC() {
	if e == nil {
		return
	}
	pc, _, _, ok := runtime.Caller(e.GetSkip())
	if ok {
		e.pcs = append(e.pcs, pc)
	}
}

// RecordLog 记录日志，根据日志级别输出到对应的日志级别
func (e *ErLog) RecordLog(ctx context.Context, opts ...Option) {
	_ = e.Log(ctx, opts...)
}

// Log 记录日志，根据日志级别输出到对应的日志级别
func (e *ErLog) Log(ctx context.Context, opts ...Option) *ErLog {
	if e == nil {
		return e
	}

	e.buildOptions(opts...)

	e.setAt(time.Now().UnixNano())
	if e.GetCode() == CodeUnknown && (e.GetLevel() == LevelDebug || e.GetLevel() == LevelInfo || e.GetLevel() == LevelWarn) {
		e.setCode(CodeSuccess)

		if e.GetMsg() == "" || e.GetMsg() == MsgUnknown {
			e.setMsg(MsgSuccess)
		}

		if e.GetContent() == "" || e.GetContent() == MsgUnknown {
			e.setContent(MsgSuccess)
		}
	}

	var callers []map[string]any
	for _, pc := range e.GetPCs() {
		f := runtime.FuncForPC(pc)
		if f == nil {
			continue
		}
		file, line := f.FileLine(pc)
		callers = append(callers, map[string]any{
			"func": f.Name(),
			"file": file,
			"line": line,
		})
	}

	logFields := []zap.Field{
		zap.String("type", e.GetKind()),
		zap.Int64("code", e.GetCode()),
		zap.String("content", e.GetContent()),
		zap.Int32("biz_id", e.GetBizID()),
		zap.String("biz_name", e.GetBizName()),
		zap.Int64("at", e.GetAt()),
		zap.Any("callers", callers),
	}

	fields := e.baseFields()
	fields = append(fields, logFields...)
	fields = append(fields, e.fields...)
	fields = dedupFields(fields)

	logger := e.getLogger()
	if logger == nil {
		addLogToBuffer(e.GetLevel(), e.GetMsg(), fields)
		return e
	}

	logByLevel(logger, e.GetLevel(), e.GetMsg(), fields...)

	return e
}

// getLogger 获取 zap.Logger 单例实例
func (e *ErLog) getLogger() *zap.Logger {
	config := GetConfig()
	if e == nil || config == nil {
		return nil
	}
	return zap.New(newZapCore(config))
}

// baseFields 获取基础日志字段，包含应用元数据信息
func (e *ErLog) baseFields() []zap.Field {
	if e == nil {
		return nil
	}

	mt := metas.Metadata()
	return []zap.Field{
		zap.String("app", mt.App()),
		zap.String("kind", mt.Kind().String()),
		zap.String("mode", mt.Mode().String()),
		zap.String("ip", mt.Ip()),
		zap.String("node", mt.Node()),
		zap.String("region", mt.Region()),
		zap.String("zone", mt.Zone()),
		zap.String("provider", mt.Provider()),
	}
}

func (e *ErLog) Clone() *ErLog {
	if e == nil {
		return nil
	}

	fields := make([]zap.Field, len(e.fields))
	copy(fields, e.fields)

	pcs := make([]uintptr, len(e.pcs))
	copy(pcs, e.pcs)

	return &ErLog{
		level:   e.level,
		kind:    e.kind,
		bizID:   e.bizID,
		bizName: e.bizName,
		code:    e.code,
		msg:     e.msg,
		content: e.content,
		at:      e.at,
		fields:  fields,
		pcs:     pcs,
	}
}

// Error 实现 error 接口，返回错误字符串表示
func (e *ErLog) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s %s", e.level.String(), e.msg, e.String())
}

// String 实现 fmt.Stringer 接口，返回 ErLog 的字符串表示
func (e *ErLog) String() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf(
		`{"level":"%s","kind":"%s","bizID":%d,"bizName":"%s","code":%d,"msg":"%s","content":"%s","at":%d}`,
		e.level.String(),
		e.kind,
		e.bizID,
		e.bizName,
		e.code,
		e.msg,
		e.content,
		e.at,
	)
}

// dedupFields 对字段进行去重处理，重复的字段名添加序号后缀
func dedupFields(fields []zap.Field) []zap.Field {
	if len(fields) <= 1 {
		return fields
	}

	seen := make(map[string]int, len(fields))
	result := make([]zap.Field, 0, len(fields))

	for _, f := range fields {
		key := f.Key
		count := seen[key]
		seen[key] = count + 1
		if count > 0 {
			f.Key = key + "_" + cast.ToString(count)
		}
		result = append(result, f)
	}

	return result
}

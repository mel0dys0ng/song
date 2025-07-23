package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mel0dys0ng/song/metas"
	"github.com/mel0dys0ng/song/utils/aob"
	"github.com/mel0dys0ng/song/utils/caller"
	"go.uber.org/zap"
)

const (
	TypeBiz     = "biz"    // 业务日志
	TypeSystem  = "system" // 系统（框架）日志
	TypeTrace   = "trace"  // Trace日志
	SkipDefault = 4
)

var (
	existingFieldNames = []string{
		"name", "mode", "ip", "node", "region", "zone", "provider",
		"type", "code", "content", "at", "caller",
	}
)

type ErLog struct {
	// Logger
	logger *zap.Logger
	// 等级
	level Level
	// 类型
	typet string
	// 业务/模块 ID
	bizId int32
	// 业务/模块 name
	bizName string
	// 编码（对外展示）
	code int64
	// 概述（对外展示）
	msg string
	// 内容
	content string
	// 时间（纳秒）
	at int64
	// caller
	caller string
	// 调用链
	chain []*ErLog
	// fields zap.Field
	fields []zap.Field
	// 描述格式
	format string
	// skip caller skip
	skip int
	// 是否记录日志
	log bool
	// built
	built bool
}

func New(opts []Option) ErLogInterface {
	e := &ErLog{code: -1, level: LevelError, typet: TypeBiz}
	e.buildOptions(opts)
	return e
}

func (e *ErLog) Level() string {
	if e != nil {
		return e.level.String()
	}
	return ""
}

func (e *ErLog) Type() string {
	if e != nil {
		return e.typet
	}
	return ""
}

func (e *ErLog) BizId() int32 {
	if e != nil {
		return e.bizId
	}
	return 0
}

func (e *ErLog) BizName() string {
	if e != nil {
		return e.bizName
	}
	return ""
}

func (e *ErLog) Code() int64 {
	if e != nil {
		return e.code
	}
	return 0
}

func (e *ErLog) Msg() string {
	if e != nil {
		return e.msg
	}
	return ""
}

func (e *ErLog) Content() string {
	if e != nil {
		return e.content
	}
	return ""
}

func (e *ErLog) At() int64 {
	if e != nil {
		return e.at
	}
	return 0
}

func (e *ErLog) Caller() string {
	if e != nil {
		return e.caller
	}
	return ""
}

func (e *ErLog) Format() string {
	if e != nil {
		return e.format
	}
	return ""
}

func (e *ErLog) Chain() []*ErLog {
	if e != nil {
		return e.chain
	}
	return nil
}

func (e *ErLog) SetLogger(config *Config) {
	if e != nil && config != nil {
		mt := metas.Mt()
		fields := []zap.Field{
			zap.String("app", mt.App()),
			zap.String("product", mt.Product()),
			zap.String("mode", mt.Mode().String()),
			zap.String("ip", mt.Ip()),
			zap.String("node", mt.Node()),
			zap.String("region", mt.Region()),
			zap.String("zone", mt.Zone()),
			zap.String("provider", mt.Provider()),
		}

		config.Dir = metas.LogDir()
		if err := config.Check(); err != nil {
			config = DefaultConfig()
			fields = append(fields, zap.String("setLoggerError", err.Error()))
		}

		e.logger = zap.New(newZapCore(config)).With(fields...)
	}
}

func (e *ErLog) SetLevel(level Level) {
	if e != nil {
		e.level = level
	}
}

func (e *ErLog) SetType(types string) {
	if e != nil {
		e.typet = types
	}
}

func (e *ErLog) SetLog(b bool) {
	if e != nil {
		e.log = b
	}
}

func (e *ErLog) SetBiz(id int32, name string) {
	if e != nil {
		e.bizId = id
		e.bizName = name
	}
}

func (e *ErLog) SetCode(code int64) {
	if e != nil {
		e.code = code
	}
}

func (e *ErLog) SetMsg(msg string) {
	if e != nil {
		e.msg = msg
	}
}

func (e *ErLog) SetContent(content string) {
	if e != nil {
		e.content = content
	}
}

func (e *ErLog) SetFields(fields ...zap.Field) {
	if e != nil {
		e.fields = fields
	}
}

func (e *ErLog) AddFields(fields ...zap.Field) {
	if e != nil && len(fields) > 0 {
		e.fields = append(e.fields, fields...)
	}
}

func (e *ErLog) SetFormat(format string) {
	if e != nil && len(format) > 0 {
		e.format = format
	}
}

func (e *ErLog) SetSkip(skip int) {
	if e != nil {
		e.skip = skip
	}
}

// OK if err == nil || err.Level <= LevelWarn return true, else return false
func (e *ErLog) OK() bool {
	return e == nil || e.level <= LevelWarn
}

// WithOptions returns a new *ErLog with new options. the new option replaces the old one.
func (e *ErLog) WithOptions(opts ...Option) *ErLog {
	c := e.new()
	c.buildOptions(opts)
	return c
}

// WithStatus returns a new *ErLog with new code and msg options.  the new option replaces the old one.
func (e *ErLog) WithStatus(code int64, msg string) *ErLog {
	c := e.new()
	c.SetCode(code)
	c.SetMsg(msg)
	return c
}

// WithStatusf returns a new *ErLog with new code and msg format options.  the new option replaces the old one.
func (e *ErLog) WithStatusf(code int64, format string, values ...any) *ErLog {
	c := e.new()
	c.SetCode(code)
	if len(format) > 0 {
		c.SetMsg(fmt.Sprintf(format, values...))
		c.SetFormat(format)
	}
	return c
}

func (e *ErLog) Debug(ctx context.Context, opts ...Option) *ErLog {
	return e.new().build(ctx, LevelDebug, opts)
}

func (e *ErLog) DebugL(ctx context.Context, opts ...Option) {
	c := e.new()
	opts = append(opts, Log(true))
	_ = c.build(ctx, LevelDebug, opts)
}

func (e *ErLog) DebugE(ctx context.Context, opts ...Option) *ErLog {
	c := e.new()
	opts = append(opts, Log(false))
	return c.build(ctx, LevelDebug, opts)
}

func (e *ErLog) Info(ctx context.Context, opts ...Option) *ErLog {
	return e.new().build(ctx, LevelInfo, opts)
}

func (e *ErLog) InfoL(ctx context.Context, opts ...Option) {
	c := e.new()
	opts = append(opts, Log(true))
	_ = c.build(ctx, LevelInfo, opts)
}

func (e *ErLog) InfoE(ctx context.Context, opts ...Option) *ErLog {
	c := e.new()
	opts = append(opts, Log(false))
	return c.build(ctx, LevelInfo, opts)
}

func (e *ErLog) Warn(ctx context.Context, opts ...Option) *ErLog {
	return e.new().build(ctx, LevelWarn, opts)
}

func (e *ErLog) WarnL(ctx context.Context, opts ...Option) {
	c := e.new()
	opts = append(opts, Log(true))
	_ = c.build(ctx, LevelWarn, opts)
}

func (e *ErLog) WarnE(ctx context.Context, opts ...Option) *ErLog {
	c := e.new()
	opts = append(opts, Log(false))
	return c.build(ctx, LevelWarn, opts)
}

func (e *ErLog) Erorr(ctx context.Context, opts ...Option) *ErLog {
	return e.new().build(ctx, LevelError, opts)
}

func (e *ErLog) ErorrL(ctx context.Context, opts ...Option) {
	c := e.new()
	opts = append(opts, Log(true))
	_ = c.build(ctx, LevelError, opts)
}

func (e *ErLog) ErorrE(ctx context.Context, opts ...Option) *ErLog {
	c := e.new()
	opts = append(opts, Log(false))
	return c.build(ctx, LevelError, opts)
}

func (e *ErLog) Panic(ctx context.Context, opts ...Option) *ErLog {
	return e.new().build(ctx, LevelPanic, opts)
}

func (e *ErLog) PanicL(ctx context.Context, opts ...Option) {
	c := e.new()
	opts = append(opts, Log(true))
	_ = c.build(ctx, LevelPanic, opts)
}

func (e *ErLog) PanicE(ctx context.Context, opts ...Option) *ErLog {
	c := e.new()
	opts = append(opts, Log(false))
	return c.build(ctx, LevelPanic, opts)
}

func (e *ErLog) Fatal(ctx context.Context, opts ...Option) *ErLog {
	return e.new().build(ctx, LevelFatal, opts)
}

func (e *ErLog) FatalL(ctx context.Context, opts ...Option) {
	c := e.new()
	opts = append(opts, Log(true))
	_ = c.build(ctx, LevelFatal, opts)
}

func (e *ErLog) FatalE(ctx context.Context, opts ...Option) *ErLog {
	c := e.new()
	opts = append(opts, Log(false))
	return c.build(ctx, LevelFatal, opts)
}

func (e *ErLog) Error() string {
	return e.String()
}

func (e *ErLog) String() string {
	data := map[string]any{
		"lv":      e.Level(),
		"type":    e.Type(),
		"bizId":   e.BizId(),
		"bizName": e.BizName(),
		"code":    e.Code(),
		"msg":     e.Msg(),
		"content": e.Content(),
		"caller":  e.Caller(),
		"at":      e.At(),
	}
	bytes, _ := json.Marshal(&data)
	return string(bytes)
}

func (e *ErLog) new() *ErLog {
	if e == nil {
		return e
	}

	// 调用链
	c := *e
	chain := e.chain
	e.chain = nil
	c.chain = nil
	c.chain = append(c.chain, e)
	if len(chain) > 0 {
		c.chain = append(c.chain, chain...)
	}

	// 设置未未被build过，重新来过
	c.built = false

	return &c
}

func (e *ErLog) build(ctx context.Context, level Level, opts []Option) *ErLog {
	e.buildOptions(opts)
	e.SetSkip(aob.Aorb(e.skip > 0, e.skip, SkipDefault))
	e.SetLevel(level)

	e.caller = caller.Location(e.skip+1, false)
	e.at = time.Now().UnixNano()

	switch e.Type() {
	case TypeBiz, TypeSystem, TypeTrace:
	default:
		e.SetType(TypeBiz)
	}

	if e.code < 0 && e.level < LevelWarn {
		e.code = 0
	}

	if len(e.msg) == 0 && len(e.format) > 0 {
		e.msg = e.format
	}

	if e.log {
		e.RecordLog(ctx)
	}

	e.built = true

	return e
}

func (e *ErLog) buildOptions(opts []Option) {
	if len(opts) > 0 {
		for _, v := range opts {
			v.Apply(e)
		}
	}
}

func (e *ErLog) RecordLog(ctx context.Context) {
	if e.logger == nil {
		e.SetLogger(DefaultConfig())
	}

	fields := e.buildFields(ctx)

	switch e.level {
	case LevelDebug:
		e.logger.Debug(e.msg, fields...)
	case LevelInfo:
		e.logger.Info(e.msg, fields...)
	case LevelWarn:
		e.logger.Warn(e.msg, fields...)
	case LevelError:
		e.logger.Error(e.msg, fields...)
	case LevelPanic:
		e.logger.Panic(e.msg, fields...)
	case LevelFatal:
		e.logger.Fatal(e.msg, fields...)
	default:
		e.logger.Error(e.msg, fields...)
	}
}

func (e *ErLog) buildFields(ctx context.Context) (fields []zap.Field) {
	span := TraceSpanFromContext(ctx)

	fields = []zap.Field{
		zap.String("type", e.Type()),
		zap.Int32("bizId", e.BizId()),
		zap.String("bizName", e.BizName()),
		zap.Int64("code", e.Code()),
		zap.String("content", e.Content()),
		zap.Int64("at", e.At()),
		zap.String("caller", e.Caller()),
		zap.String("traceId", span.GetTraceID()),
	}

	keys := map[string]int{"chain": 0}
	for _, v := range existingFieldNames {
		keys[v] = 0
	}

	if len(e.fields) > 0 {
		for _, v := range e.fields {
			if _, ok := keys[v.Key]; ok {
				v.Key = fmt.Sprintf("%s%d", v.Key, keys[v.Key])
				keys[v.Key] += 1
			} else {
				keys[v.Key] = 0
			}
			fields = append(fields, v)
		}
	}

	fields = append(fields, zap.Any("chain", e.chain))

	return
}

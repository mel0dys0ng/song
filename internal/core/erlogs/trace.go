package erlogs

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mel0dys0ng/song/pkg/aob"
	"go.uber.org/zap"
)

type traceSpanContextValueKey struct{}

// TraceSpan 表示一个追踪跨度，用于记录请求链路中的调用信息
type TraceSpan struct {
	name         string
	traceID      string
	spanID       string
	parentSpanID string
	startAt      time.Time
	endAt        time.Time
	cost         time.Duration
}

// GetName 获取追踪跨度名称，如果 TraceSpan 为 nil 则返回空字符串
func (ts *TraceSpan) GetName() string {
	if ts == nil {
		return ""
	}
	return ts.name
}

// GetTraceID 获取追踪 ID，如果 TraceSpan 为 nil 则返回空字符串
func (ts *TraceSpan) GetTraceID() string {
	if ts == nil {
		return ""
	}
	return ts.traceID
}

// GetSpanID 获取跨度 ID，如果 TraceSpan 为 nil 则返回空字符串
func (ts *TraceSpan) GetSpanID() string {
	if ts == nil {
		return ""
	}
	return ts.spanID
}

// GetParentSpanID 获取父跨度 ID，如果 TraceSpan 为 nil 则返回空字符串
func (ts *TraceSpan) GetParentSpanID() string {
	if ts == nil {
		return ""
	}
	return ts.parentSpanID
}

// GetStartAt 获取开始时间，如果 TraceSpan 为 nil 则返回零值
func (ts *TraceSpan) GetStartAt() time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.startAt
}

// GetEndAt 获取结束时间，如果 TraceSpan 为 nil 则返回零值
func (ts *TraceSpan) GetEndAt() time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.endAt
}

// GetCost 获取耗时，如果 TraceSpan 为 nil 则返回 0
func (ts *TraceSpan) GetCost() time.Duration {
	if ts == nil {
		return 0
	}
	return ts.cost
}

// ZapFields 将 TraceSpan 转换为 zap.Field 列表，用于日志记录
func (ts *TraceSpan) ZapFields() []zap.Field {
	if ts == nil {
		return nil
	}
	return []zap.Field{
		zap.String("trace", ts.name),
		zap.String("trace_id", ts.traceID),
		zap.String("span_id", ts.spanID),
		zap.String("parent_span_id", ts.parentSpanID),
		zap.String("start_at", ts.startAt.Format(time.RFC3339Nano)),
		zap.String("end_at", ts.endAt.Format(time.RFC3339Nano)),
		zap.Int64("cost", ts.cost.Milliseconds()),
	}
}

// StartTrace 开始一个新的追踪跨度，返回包含 TraceSpan 的上下文
// 如果上下文中已存在 TraceSpan，则新跨度将作为其子跨度
func StartTrace(ctx context.Context, name string) context.Context {
	spanParent := TraceSpanFromContext(ctx)

	var traceID string
	if spanParent != nil && len(spanParent.GetTraceID()) > 0 {
		traceID = spanParent.GetTraceID()
	} else {
		traceID = genTraceID()
	}

	span := &TraceSpan{
		name:         name,
		startAt:      time.Now(),
		traceID:      traceID,
		spanID:       genTraceID(),
		parentSpanID: aob.VarOrVar(len(spanParent.GetSpanID()) > 0, spanParent.GetSpanID(), ""),
	}

	return context.WithValue(ctx, traceSpanContextValueKey{}, span)
}

// EndTrace 结束当前追踪跨度，计算耗时并返回 TraceSpan
func EndTrace(ctx context.Context, err error) {
	span := TraceSpanFromContext(ctx)
	if span == nil {
		return
	}

	span.endAt = time.Now()
	span.cost = span.endAt.Sub(span.startAt)

	// 记录追踪跨度信息
	options := []Option{OptionSkip(4), OptionKind(KindTrace), ErLogsBiz}
	if err != nil {
		Convert(err).Options(options).AppendFields(span.ZapFields()...).RecordLog(ctx)
	} else {
		New("").Options(options).AppendFields(span.ZapFields()...).InfoLog(ctx)
	}
}

// TraceSpanFromContext 从上下文中获取 TraceSpan
func TraceSpanFromContext(ctx context.Context) *TraceSpan {
	if ctx == nil {
		return nil
	}
	span, _ := ctx.Value(traceSpanContextValueKey{}).(*TraceSpan)
	return span
}

// genTraceID 生成唯一的追踪 ID，使用 UUID v4
func genTraceID() string {
	return uuid.New().String()
}

package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/mel0dys0ng/song/core/utils/aob"
	"github.com/mel0dys0ng/song/core/utils/caller"
	"github.com/mel0dys0ng/song/core/utils/crypto"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type traceSpanContextValueKey struct{}

var TraceSpanContextValueKey = traceSpanContextValueKey{}

const CallerSkip = 4

type TraceSpan struct {
	Name         string
	TraceID      string
	SpanID       string
	ParentSpanID string
	Caller       string
	StartTime    time.Time
	EndTime      time.Time
	Cost         time.Duration
}

func (ts *TraceSpan) GetName() string {
	if ts != nil {
		return ts.Name
	}
	return ""
}

func (ts *TraceSpan) GetTraceID() string {
	if ts != nil {
		return ts.TraceID
	}
	return ""
}

func (ts *TraceSpan) GetSpanID() string {
	if ts != nil {
		return ts.SpanID
	}
	return ""
}

func (ts *TraceSpan) GetParentSpanID() string {
	if ts != nil {
		return ts.ParentSpanID
	}
	return ""
}

func (ts *TraceSpan) GetCaller() string {
	if ts != nil {
		return ts.Caller
	}
	return ""
}

func (ts *TraceSpan) GetStartTime() time.Time {
	if ts != nil {
		return ts.StartTime
	}
	return time.Time{}
}

func (ts *TraceSpan) GetEndTime() time.Time {
	if ts != nil {
		return ts.EndTime
	}
	return time.Time{}
}

func (ts *TraceSpan) GetCost() time.Duration {
	if ts != nil {
		return ts.Cost
	}
	return 0
}

func (ts *TraceSpan) ZapFields() []zap.Field {
	return []zap.Field{
		zap.String("traceName", ts.GetName()),
		zap.String("traceSpanID", ts.GetSpanID()),
		zap.String("traceParentSpanID", ts.GetParentSpanID()),
		zap.String("traceCaller", ts.GetCaller()),
		zap.String("traceStartTime", ts.GetStartTime().String()),
		zap.String("traceEndTime", ts.GetEndTime().String()),
		zap.Int64("traceCost", ts.GetCost().Milliseconds()),
	}
}

func StartTrace(ctx context.Context, name string) context.Context {
	c := caller.New(CallerSkip)
	spanParent := TraceSpanFromContext(ctx)
	span := &TraceSpan{
		Name:         name,
		StartTime:    time.Now(),
		TraceID:      aob.Aorb(len(spanParent.GetTraceID()) > 0, spanParent.GetTraceID(), genTraceNameID("traceID")),
		SpanID:       genTraceNameID("spanID"),
		ParentSpanID: aob.Aorb(len(spanParent.GetSpanID()) > 0, spanParent.GetSpanID(), ""),
		Caller:       fmt.Sprintf("%s:%d", c.Func(), c.Line()),
	}
	return context.WithValue(ctx, TraceSpanContextValueKey, span)
}

func EndTrace(ctx context.Context) (span *TraceSpan) {
	if span = TraceSpanFromContext(ctx); span != nil {
		span.EndTime = time.Now()
		span.Cost = span.EndTime.Sub(span.StartTime)
	}
	return
}

func TraceSpanFromContext(ctx context.Context) *TraceSpan {
	span, _ := ctx.Value(TraceSpanContextValueKey).(*TraceSpan)
	return span
}

func genTraceNameID(name string) string {
	return crypto.MD5([]any{name, lo.RandomString(32, lo.AllCharset), time.Now().UnixNano()})
}

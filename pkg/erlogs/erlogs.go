package erlogs

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mel0dys0ng/song/internal/core/erlogs"
	"go.uber.org/zap"
)

type (
	Kind           = erlogs.Kind
	ErLogInterface = erlogs.ErLogInterface
	Config         = erlogs.Config
	TraceSpan      = erlogs.TraceSpan
	Option         = erlogs.Option
)

const (
	KindBiz    = erlogs.KindBiz
	KindSystem = erlogs.KindSystem
	KindTrace  = erlogs.KindTrace
)

// Initialize 初始化错误日志记录器，配置全局日志记录器
func Initialize(config *Config) {
	erlogs.Initialize(config)
}

// New 创建一个新的 ErLog 实例，使用给定的文本作为错误消息
func New(text string, opts ...Option) ErLogInterface {
	return erlogs.New(text, opts...)
}

// Newf 创建一个新的 ErLog 实例，使用格式化字符串作为错误消息
func Newf(format string, args ...any) ErLogInterface {
	return erlogs.New(fmt.Sprintf(format, args...))
}

// Convert 将标准 error 转换为 ErLog，如果 err 已经是 ErLog 类型则直接返回
func Convert(err error, opts ...Option) ErLogInterface {
	return erlogs.Convert(err, opts...)
}

// WithOptions 创建一个新的 ErLog 实例，使用给定的选项
func WithOptions(opts ...Option) ErLogInterface {
	return erlogs.WithOptions(opts...)
}

// StartTrace 开始一个新的跟踪 span，返回包含 span 信息的上下文
func StartTrace(ctx context.Context, name string) context.Context {
	return erlogs.StartTrace(ctx, name)
}

// StartTracef 开始一个新的跟踪 span，返回包含 span 信息的上下文
func StartTracef(ctx context.Context, format string, args ...any) context.Context {
	return erlogs.StartTrace(ctx, fmt.Sprintf(format, args...))
}

// EndTrace 结束当前跟踪 span，记录 span 信息
func EndTrace(ctx context.Context, err error) {
	erlogs.EndTrace(ctx, err)
}

// TraceSpanFromContext 从上下文获取当前跟踪 span 信息
func TraceSpanFromContext(ctx context.Context) *TraceSpan {
	return erlogs.TraceSpanFromContext(ctx)
}

// MaskFields 对 zap.Field 数组进行敏感信息脱敏
func MaskFields(fields ...zap.Field) []zap.Field {
	return erlogs.MaskFields(fields...)
}

// MaskField 对 zap.Field 进行敏感信息脱敏
func MaskField(field zap.Field) zap.Field {
	return erlogs.MaskField(field)
}

// ValidatorError 从验证错误中创建一个 ErLog 实例，包含请求参数和验证错误信息
func ValidatorError(request any, validateErr error, opts ...Option) (err error) {
	if validateErr == nil || request == nil {
		return
	}

	err = InvalidArguments.Options(opts).Info(
		erlogs.OptionContent(validateErr.Error()),
		erlogs.OptionFields(
			zap.Any("request", request),
		),
	)

	errs, ok := validateErr.(validator.ValidationErrors)
	if !ok || len(errs) == 0 {
		return
	}

	tp := reflect.TypeOf(request)
	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
	}

	if tp.Kind() != reflect.Struct {
		return
	}

	field, found := tp.FieldByName(errs[0].StructField())
	if !found {
		return
	}

	if msg := field.Tag.Get("msg"); msg != "" {
		err = Convert(err, erlogs.OptionMsg(msg))
	}

	return
}

package erlogs

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mel0dys0ng/song/erlogs/internal"
	"go.uber.org/zap"
)

const (
	BadReuestParams = "bad request params"
)

// Logger Set the logger of the erlog
func Logger(config *Config) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetLogger(config)
		},
	}
}

// Log Set the option of whether to record logs
func Log(b bool) internal.Option {
	return internal.Log(b)
}

// TypeBiz Set the type option to biz
func TypeBiz() internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetType(internal.TypeBiz)
		},
	}
}

// TypeSystem Set the type option to system
func TypeSystem() internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetType(internal.TypeSystem)
		},
	}
}

// TypeTrace Set the type option to trace
func TypeTrace() internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetType(internal.TypeTrace)
		},
	}
}

func Event(event string) zap.Field {
	return zap.String("event", event)
}

func Step(step string) zap.Field {
	return zap.String("step", step)
}

func Pkg(pkg string) zap.Field {
	return zap.String("pkg", pkg)
}

func Biz(id int32, name string) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetBiz(id, name)
		},
	}
}

func Code(code int64) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetCode(code)
		},
	}
}

func Status(code int64, msg string) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetCode(code)
			e.SetMsg(msg)
		},
	}
}

func Statusf(code int64, format string, values ...any) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetCode(code)
			if len(format) > 0 && len(values) > 0 {
				e.SetMsg(fmt.Sprintf(format, values...))
			} else {
				e.SetFormat(format)
			}
		},
	}
}

func Msg(msg string) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetMsg(msg)
		},
	}
}

func Msgf(format string, values ...any) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			if len(format) > 0 && len(values) > 0 {
				e.SetMsg(fmt.Sprintf(format, values...))
			} else {
				e.SetFormat(format)
			}
		},
	}
}

func Msgv(values ...any) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			if format := e.Format(); len(format) > 0 {
				e.SetMsg(fmt.Sprintf(format, values...))
				e.SetFormat("")
			}
		},
	}
}

// ValidatorError 根据validator的error设置Msg和Content
func ValidatorError(data any, err error) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			if err == nil {
				return
			}

			e.SetMsg(BadReuestParams)
			e.SetContent(err.Error())
			e.SetFields(zap.Any("req", data))

			_, ok := err.(*validator.InvalidValidationError)
			if ok || data == nil {
				return
			}

			errors := err.(validator.ValidationErrors)
			if len(errors) > 0 {
				typ := reflect.TypeOf(data)
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
				}

				field, _ := typ.FieldByName(errors[0].StructField())
				if msg := field.Tag.Get("msg"); len(msg) > 0 {
					e.SetMsg(msg)
				}

				return
			}
		},
	}
}

func Content(content string) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetContent(content)
		},
	}
}

func ContentError(err error) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			if err != nil {
				e.SetContent(err.Error())
			}
		},
	}
}

func Contentf(format string, values ...any) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			if len(format) > 0 {
				e.SetContent(fmt.Sprintf(format, values...))
			}
		},
	}
}

func Fields(fields ...zap.Field) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetFields(fields...)
		},
	}
}

func AddFields(fields ...zap.Field) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.AddFields(fields...)
		},
	}
}

func Skip(skip int) internal.Option {
	return internal.Option{
		Apply: func(e *internal.ErLog) {
			e.SetSkip(skip)
		},
	}
}

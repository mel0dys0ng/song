package erlogs

import (
	"reflect"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	maskTag = "mask"
)

// MaskField 对 zap.Field 进行敏感信息脱敏
func MaskField(field zap.Field) zap.Field {
	if field.Type == zapcore.StringType {
		return zap.String(field.Key, smartMask(field.String))
	}

	if field.Type == zapcore.ReflectType {
		maskReflectValue(field.Interface)
		return field
	}

	return field
}

// MaskFields 对 zap.Field 数组进行敏感信息脱敏
func MaskFields(fields ...zap.Field) []zap.Field {
	result := make([]zap.Field, len(fields))
	for i, field := range fields {
		result[i] = MaskField(field)
	}
	return result
}

// maskReflectValue 递归脱敏结构体中的字符串字段
func maskReflectValue(v any) {
	if v == nil {
		return
	}

	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return
	}

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		tag := structField.Tag.Get(maskTag)
		if tag == "" {
			continue
		}

		if field.Kind() == reflect.String {
			field.SetString(smartMask(field.String()))
		} else if field.Kind() == reflect.Struct {
			maskReflectValue(field.Interface())
		}
	}
}

// smartMask 对字符串进行智能脱敏，根据长度返回不同的掩码
func smartMask(value string) string {
	n := len(value)
	if n == 0 {
		return value
	}

	switch {
	case n <= 1:
		return "*"
	case n == 2:
		return value[:1] + "*"
	case n >= 3 && n <= 5:
		return value[:1] + strings.Repeat("*", n-2) + value[n-1:]
	case n >= 6 && n <= 10:
		return value[:2] + strings.Repeat("*", n-4) + value[n-2:]
	default:
		return value[:3] + strings.Repeat("*", n-6) + value[n-3:]
	}
}

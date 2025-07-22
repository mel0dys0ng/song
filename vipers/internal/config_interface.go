// Code by Song Struct2Interface. DO NOT EDIT.

package internal

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

type ConfigInterface interface {
	Provider() ProviderInterface

	// OnConfigChange 配置变更回调
	OnConfigChange(fn func(event fsnotify.Event, options *Options))

	IsSet(key string) bool

	AllKeys() []string

	AllSettings() map[string]any

	Get(key string, defaultValue any) any

	GetBool(key string, defaultValue bool) bool

	GetDuration(key string, defaultValue time.Duration) time.Duration

	GetFloat64(key string, defaultValue float64) float64

	GetInt(key string, defaultValue int) int

	GetInt32(key string, defaultValue int32) int32

	GetInt64(key string, defaultValue int64) int64

	GetIntSlice(key string, defaultValue []int) []int

	GetSizeInBytes(key string, defaultValue uint) uint

	GetString(key string, defaultValue string) string

	GetStringMap(key string, defaultValue map[string]any) map[string]any

	GetStringMapString(key string, defaultValue map[string]string) map[string]string

	GetStringMapStringSlice(key string, defaultValue map[string][]string) map[string][]string

	GetStringSlice(key string, defaultValue []string) []string

	GetTime(key string, defaultValue time.Time) time.Time

	GetUint(key string, defaultValue uint) uint

	GetUint32(key string, defaultValue uint32) uint32

	GetUint64(key string, defaultValue uint64) uint64

	UnmarshalKey(key string, value, defaultValue any) (err error)

	Unmarshal(value any) error

	UnmarshalExact(value any) error
}

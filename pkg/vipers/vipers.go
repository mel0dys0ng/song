package vipers

import (
	"context"
	"strings"
	"time"

	"github.com/mel0dys0ng/song/internal/core/vipers"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/metas"
	"github.com/mel0dys0ng/song/pkg/singleton"
)

func Config() vipers.ConfigInterface {
	config := singleton.Once(singleton.Key(), initialize)
	if config == nil {
		ctx := context.Background()
		erlogs.New("failed to init vipers config: config is nil").PanicLog(ctx)
	}
	return config
}

func initialize() (res vipers.ConfigInterface) {
	var err error

	ctx := erlogs.StartTrace(context.Background(), "initVipers")
	defer func() {
		erlogs.EndTrace(ctx, nil)
		if err != nil {
			erlogs.Convert(err).RecordLog(ctx)
		}
	}()

	mt := metas.Metadata()
	opts := []vipers.Option{
		vipers.Provider(mt.ConfigType()),
		vipers.Type(mt.ConfigType()),
		vipers.Endpoint(mt.ConfigAddr()),
		vipers.Path(mt.ConfigPath()),
	}

	config, err := vipers.New(opts...)
	if err != nil {
		erlogs.Convert(err).RecordLog(ctx)
		return
	}

	if err := config.Provider().Load(); err != nil {
		erlogs.Convert(err).Wrap("failed to load config").PanicLog(ctx)
		return
	}

	return config
}

func Key(names ...string) string {
	return strings.Join(names, ".")
}

func IsSet(key string) bool {
	return Config().IsSet(key)
}

func AllKeys() []string {
	return Config().AllKeys()
}

func AllSettings() map[string]any {
	return Config().AllSettings()
}

func Get(key string, defaultValue any) any {
	return Config().Get(key, defaultValue)
}

func GetBool(key string, defaultValue bool) bool {
	return Config().GetBool(key, defaultValue)
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	return Config().GetDuration(key, defaultValue)
}

func GetFloat64(key string, defaultValue float64) float64 {
	return Config().GetFloat64(key, defaultValue)
}

func GetInt(key string, defaultValue int) int {
	return Config().GetInt(key, defaultValue)
}

func GetInt32(key string, defaultValue int32) int32 {
	return Config().GetInt32(key, defaultValue)
}

func GetInt64(key string, defaultValue int64) int64 {
	return Config().GetInt64(key, defaultValue)
}

func GetIntSlice(key string, defaultValue []int) []int {
	return Config().GetIntSlice(key, defaultValue)
}

func GetSizeInBytes(key string, defaultValue uint) uint {
	return Config().GetSizeInBytes(key, defaultValue)
}

func GetString(key string, defaultValue string) string {
	return Config().GetString(key, defaultValue)
}

func GetStringMap(key string, defaultValue map[string]any) map[string]any {
	return Config().GetStringMap(key, defaultValue)
}

func GetStringMapString(key string, defaultValue map[string]string) map[string]string {
	return Config().GetStringMapString(key, defaultValue)
}

func GetStringMapStringSlice(key string, defaultValue map[string][]string) map[string][]string {
	return Config().GetStringMapStringSlice(key, defaultValue)
}

func GetStringSlice(key string, defaultValue []string) []string {
	return Config().GetStringSlice(key, defaultValue)
}

func GetTime(key string, defaultValue time.Time) time.Time {
	return Config().GetTime(key, defaultValue)
}

func GetUint(key string, defaultValue uint) uint {
	return Config().GetUint(key, defaultValue)
}

func GetUint32(key string, defaultValue uint32) uint32 {
	return Config().GetUint32(key, defaultValue)
}

func GetUint64(key string, defaultValue uint64) uint64 {
	return Config().GetUint64(key, defaultValue)
}

func UnmarshalKey(key string, value any) error {
	return Config().UnmarshalKey(key, value)
}

func Unmarshal(value any) error {
	return Config().Unmarshal(value)
}

func UnmarshalExact(value any) error {
	return Config().UnmarshalExact(value)
}

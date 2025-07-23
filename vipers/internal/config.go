package internal

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/erlogs"
	"github.com/spf13/viper"
)

const (
	ConfigProviderJson   = "json"
	ConfigProviderYaml   = "yaml"
	ConfigProviderToml   = "toml"
	ConfigProviderEtcd   = "etcd"
	ConfigProviderConsul = "consul"
)

type Config struct {
	*Options
	*viper.Viper
	provider ProviderInterface
}

func New(elg erlogs.ErLogInterface, opts ...Option) (res ConfigInterface, err error) {
	c := &Config{
		Viper:   viper.New(),
		Options: DefaultOptions(elg),
	}

	c.buildOptions(opts)
	c.provider, err = NewProvider(c.Viper, elg, c.Options)

	return c, err
}

func (c *Config) Provider() ProviderInterface {
	return c.provider
}

// OnConfigChange 配置变更回调
func (c *Config) OnConfigChange(fn func(event fsnotify.Event, options *Options)) {
	c.buildOptions([]Option{OnChangeConfig(fn)})
}

func (c *Config) buildOptions(opts []Option) {
	if len(opts) > 0 {
		for _, v := range opts {
			v.fn(c.Options)
		}
	}
}

func (c *Config) IsSet(key string) bool {
	return c.Viper.IsSet(key)
}

func (c *Config) AllKeys() []string {
	return c.Viper.AllKeys()
}

func (c *Config) AllSettings() map[string]any {
	return c.Viper.AllSettings()
}

func (c *Config) Get(key string, defaultValue any) any {
	if c.Viper.IsSet(key) {
		return c.Viper.Get(key)
	}
	return defaultValue
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	if c.Viper.IsSet(key) {
		return c.Viper.GetBool(key)
	}
	return defaultValue
}

func (c *Config) GetDuration(key string, defaultValue time.Duration) time.Duration {
	if c.Viper.IsSet(key) {
		return c.Viper.GetDuration(key)
	}
	return defaultValue
}

func (c *Config) GetFloat64(key string, defaultValue float64) float64 {
	if c.Viper.IsSet(key) {
		return c.Viper.GetFloat64(key)
	}
	return defaultValue
}

func (c *Config) GetInt(key string, defaultValue int) int {
	if c.Viper.IsSet(key) {
		return c.Viper.GetInt(key)
	}
	return defaultValue
}

func (c *Config) GetInt32(key string, defaultValue int32) int32 {
	if c.Viper.IsSet(key) {
		return c.Viper.GetInt32(key)
	}
	return defaultValue
}

func (c *Config) GetInt64(key string, defaultValue int64) int64 {
	if c.Viper.IsSet(key) {
		return c.Viper.GetInt64(key)
	}
	return defaultValue
}

func (c *Config) GetIntSlice(key string, defaultValue []int) []int {
	if c.Viper.IsSet(key) {
		return c.Viper.GetIntSlice(key)
	}
	return defaultValue
}

func (c *Config) GetSizeInBytes(key string, defaultValue uint) uint {
	if c.Viper.IsSet(key) {
		return c.Viper.GetSizeInBytes(key)
	}
	return defaultValue
}

func (c *Config) GetString(key string, defaultValue string) string {
	if c.Viper.IsSet(key) {
		return c.Viper.GetString(key)
	}
	return defaultValue
}

func (c *Config) GetStringMap(key string, defaultValue map[string]any) map[string]any {
	if c.Viper.IsSet(key) {
		return c.Viper.GetStringMap(key)
	}
	return defaultValue
}

func (c *Config) GetStringMapString(key string, defaultValue map[string]string) map[string]string {
	if c.Viper.IsSet(key) {
		return c.Viper.GetStringMapString(key)
	}
	return defaultValue
}

func (c *Config) GetStringMapStringSlice(key string, defaultValue map[string][]string) map[string][]string {
	if c.Viper.IsSet(key) {
		return c.Viper.GetStringMapStringSlice(key)
	}
	return defaultValue
}

func (c *Config) GetStringSlice(key string, defaultValue []string) []string {
	if c.Viper.IsSet(key) {
		return c.Viper.GetStringSlice(key)
	}
	return defaultValue
}

func (c *Config) GetTime(key string, defaultValue time.Time) time.Time {
	if c.Viper.IsSet(key) {
		return c.Viper.GetTime(key)
	}
	return defaultValue
}

func (c *Config) GetUint(key string, defaultValue uint) uint {
	if c.Viper.IsSet(key) {
		return c.Viper.GetUint(key)
	}
	return defaultValue
}

func (c *Config) GetUint32(key string, defaultValue uint32) uint32 {
	if c.Viper.IsSet(key) {
		return c.Viper.GetUint32(key)
	}
	return defaultValue
}

func (c *Config) GetUint64(key string, defaultValue uint64) uint64 {
	if c.Viper.IsSet(key) {
		return c.Viper.GetUint64(key)
	}
	return defaultValue
}

func (c *Config) UnmarshalKey(key string, value any) (err error) {
	if c.Viper.IsSet(key) {
		return c.Viper.UnmarshalKey(key, value)
	}
	return
}

func (c *Config) Unmarshal(value any) error {
	return c.Viper.Unmarshal(value)
}

func (c *Config) UnmarshalExact(value any) error {
	return c.Viper.UnmarshalExact(value)
}

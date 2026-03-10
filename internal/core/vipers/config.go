package vipers

import (
	"context"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/pkg/erlogs"
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

func New(opts ...Option) (res ConfigInterface, err error) {
	c := &Config{
		Viper:   viper.New(),
		Options: DefaultOptions(),
	}

	c.buildOptions(opts)
	c.provider, err = NewProvider(c.Viper, c.Options)
	if err != nil {
		_ = erlogs.Convert(err).Options(BaseELOptions()).Log(context.Background())
		return nil, err
	}

	return c, err
}

func (c *Config) Provider() ProviderInterface {
	if c == nil {
		return nil
	}
	return c.provider
}

func (c *Config) OnConfigChange(fn func(event fsnotify.Event, options *Options)) {
	if c != nil {
		c.buildOptions([]Option{OnChangeConfig(fn)})
	}
}

func (c *Config) buildOptions(opts []Option) {
	if c != nil && len(opts) > 0 {
		for _, v := range opts {
			v(c.Options)
		}
	}
}

func (c *Config) IsSet(key string) bool {
	if c == nil || c.Viper == nil {
		return false
	}
	return c.Viper.IsSet(key)
}

func (c *Config) AllKeys() []string {
	if c == nil || c.Viper == nil {
		return []string{}
	}
	return c.Viper.AllKeys()
}

func (c *Config) AllSettings() map[string]any {
	if c == nil || c.Viper == nil {
		return map[string]any{}
	}
	return c.Viper.AllSettings()
}

func (c *Config) Get(key string, defaultValue any) any {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.Get(key)
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetBool(key)
}

func (c *Config) GetDuration(key string, defaultValue time.Duration) time.Duration {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetDuration(key)
}

func (c *Config) GetFloat64(key string, defaultValue float64) float64 {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetFloat64(key)
}

func (c *Config) GetInt(key string, defaultValue int) int {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetInt(key)
}

func (c *Config) GetInt32(key string, defaultValue int32) int32 {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetInt32(key)
}

func (c *Config) GetInt64(key string, defaultValue int64) int64 {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetInt64(key)
}

func (c *Config) GetIntSlice(key string, defaultValue []int) []int {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetIntSlice(key)
}

func (c *Config) GetSizeInBytes(key string, defaultValue uint) uint {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetSizeInBytes(key)
}

func (c *Config) GetString(key string, defaultValue string) string {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetString(key)
}

func (c *Config) GetStringMap(key string, defaultValue map[string]any) map[string]any {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetStringMap(key)
}

func (c *Config) GetStringMapString(key string, defaultValue map[string]string) map[string]string {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetStringMapString(key)
}

func (c *Config) GetStringMapStringSlice(key string, defaultValue map[string][]string) map[string][]string {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetStringMapStringSlice(key)
}

func (c *Config) GetStringSlice(key string, defaultValue []string) []string {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetStringSlice(key)
}

func (c *Config) GetTime(key string, defaultValue time.Time) time.Time {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetTime(key)
}

func (c *Config) GetUint(key string, defaultValue uint) uint {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetUint(key)
}

func (c *Config) GetUint32(key string, defaultValue uint32) uint32 {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetUint32(key)
}

func (c *Config) GetUint64(key string, defaultValue uint64) uint64 {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.GetUint64(key)
}

func (c *Config) UnmarshalKey(key string, value any) (err error) {
	if c == nil || c.Viper == nil || !c.Viper.IsSet(key) {
		return nil
	}
	return c.Viper.UnmarshalKey(key, value)
}

func (c *Config) Unmarshal(value any) error {
	if c == nil || c.Viper == nil {
		return nil
	}
	return c.Viper.Unmarshal(value)
}

func (c *Config) UnmarshalExact(value any) error {
	if c == nil || c.Viper == nil {
		return nil
	}
	return c.Viper.UnmarshalExact(value)
}

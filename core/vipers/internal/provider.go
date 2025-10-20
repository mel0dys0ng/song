package internal

import (
	"fmt"

	"github.com/mel0dys0ng/song/core/erlogs"
	"github.com/spf13/viper"
)

type ProviderInterface interface {
	Load() error
}

type NewProviderFunc func(*viper.Viper, erlogs.ErLogInterface, *Options) ProviderInterface

var (
	newProviderFuncs = map[string]NewProviderFunc{
		ConfigProviderJson:   NewLocalProvider,
		ConfigProviderToml:   NewLocalProvider,
		ConfigProviderYaml:   NewLocalProvider,
		ConfigProviderEtcd:   NewRemoteProvider,
		ConfigProviderConsul: NewRemoteProvider,
	}
)

// NewProvider 返回配置类型对应的配置加载器
func NewProvider(v *viper.Viper, el erlogs.ErLogInterface, o *Options) (provider ProviderInterface, err error) {
	if newProvider, ok := newProviderFuncs[o.Provider]; ok {
		provider = newProvider(v, el, o)
		return
	}
	err = fmt.Errorf("provider %s is not supported", o.Provider)
	return
}

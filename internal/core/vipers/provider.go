package vipers

import (
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ProviderInterface interface {
	Load() error
}

type NewProviderFunc func(*viper.Viper, *Options) ProviderInterface

var (
	newProviderFuncs = map[string]NewProviderFunc{
		ConfigProviderJson:   NewLocalProvider,
		ConfigProviderToml:   NewLocalProvider,
		ConfigProviderYaml:   NewLocalProvider,
		ConfigProviderEtcd:   NewRemoteProvider,
		ConfigProviderConsul: NewRemoteProvider,
	}
)

func NewProvider(v *viper.Viper, o *Options) (p ProviderInterface, err error) {
	if newProvider, ok := newProviderFuncs[o.Provider]; ok {
		p = newProvider(v, o)
		return
	}

	err = erlogs.New("not supported provider").Panic(
		erlogs.OptionFields(
			zap.String("provider", o.Provider),
		),
	)

	return
}

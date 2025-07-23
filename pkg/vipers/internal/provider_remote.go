package internal

import (
	"context"
	"errors"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/spf13/viper"
)

/*
 *========================================**
 * 加载远程配置
 *========================================**
 */

type RemoteProvider struct {
	*Options
	*viper.Viper
	elg erlogs.ErLogInterface
}

func NewRemoteProvider(v *viper.Viper, elg erlogs.ErLogInterface, o *Options) ProviderInterface {
	return &RemoteProvider{
		Viper:   v,
		Options: o,
		elg:     elg,
	}
}

func (p *RemoteProvider) Load() (err error) {
	err = p.checkOptions()
	if err != nil {
		err = p.elg.Panic(
			context.Background(),
			erlogs.Msgv("check options failed"),
			erlogs.ContentError(err),
		)
		return
	}

	// setting
	err = p.Viper.AddRemoteProvider(p.Provider, p.Endpoint, p.Path)
	if err != nil {
		err = p.elg.Panic(
			context.Background(),
			erlogs.Msgv("add remote provider failed"),
			erlogs.ContentError(err),
		)
		return
	}

	// load
	p.Viper.SetConfigType(p.Type)
	err = p.Viper.ReadRemoteConfig()
	if err != nil {
		err = p.elg.Panic(
			context.Background(),
			erlogs.Msgv("read remote config failed"),
			erlogs.ContentError(err),
		)
		return
	}

	// watch
	p.Viper.WatchConfig()
	p.Viper.OnConfigChange(func(in fsnotify.Event) {
		if p.OnChangeConfig != nil {
			p.OnChangeConfig(in, p.Options)
		}
	})

	return
}

func (p *RemoteProvider) checkOptions() (err error) {
	if len(p.Endpoint) == 0 {
		err = errors.New("option invalid: config endpoint empty")
		return
	}

	if len(p.Path) == 0 {
		err = errors.New("option invalid: config path empty")
		return
	}

	if len(p.Type) == 0 {
		err = errors.New("option invalid: config type empty")
	}

	return
}

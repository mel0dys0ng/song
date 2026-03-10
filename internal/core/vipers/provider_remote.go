package vipers

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/metas"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type RemoteProvider struct {
	*Options
	*viper.Viper
}

func NewRemoteProvider(v *viper.Viper, o *Options) ProviderInterface {
	return &RemoteProvider{
		Viper:   v,
		Options: o,
	}
}

func (p *RemoteProvider) Load() (err error) {
	if p == nil {
		return erlogs.New("RemoteProvider is nil").Panic()
	}

	err = p.checkOptions()
	if err != nil {
		return
	}

	if p.Provider == ConfigProviderEtcd || p.Provider == ConfigProviderConsul {
		err = p.loadRemoteConfigs()
	} else {
		err = p.loadSingleRemoteConfig()
	}

	if err != nil {
		return
	}

	configMode := metas.ModeType(p.GetString("metadata.mode"))
	if !configMode.Validate() || configMode != metas.Metadata().Mode() {
		err = erlogs.New("The config file metadata mode and the command parameter metadata mode do not match.").Panic(
			erlogs.OptionFields(
				zap.String("config_mode", configMode.String()),
				zap.String("metedata_mode", metas.Metadata().Mode().String()),
			),
		)
		return
	}

	p.Viper.WatchConfig()
	p.Viper.OnConfigChange(func(in fsnotify.Event) {
		if p != nil && p.OnChangeConfig != nil {
			p.OnChangeConfig(in, p.Options)
		}
	})

	return
}

func (p *RemoteProvider) loadSingleRemoteConfig() error {
	path := p.buildConfigPath()
	err := p.Viper.AddRemoteProvider(p.Provider, p.Endpoint, path)
	if err != nil {
		err = erlogs.Convert(err).Wrap("failed to AddRemoteProvider").Panic(
			erlogs.OptionFields(
				zap.String("provider", p.Provider),
				zap.String("endpoint", p.Endpoint),
				zap.String("path", path),
			),
		)
		return err
	}

	p.Viper.SetConfigType(p.Type)
	err = p.Viper.ReadRemoteConfig()
	if err != nil {
		err = erlogs.Convert(err).Wrap("failed to ReadRemoteConfig").Panic(
			erlogs.OptionFields(
				zap.String("provider", p.Provider),
				zap.String("endpoint", p.Endpoint),
				zap.String("path", path),
			),
		)
		return err
	}

	return nil
}

func (p *RemoteProvider) loadRemoteConfigs() error {
	path := p.buildConfigPath()
	if err := p.loadRemoteConfigAtPath(path); err != nil {
		return err
	}
	return nil
}

func (p *RemoteProvider) buildConfigPath() string {
	return fmt.Sprintf("%s/%s", p.Path, metas.Metadata().Mode().String())
}

func (p *RemoteProvider) loadRemoteConfigAtPath(path string) error {
	ctx := context.Background()
	err := p.Viper.AddRemoteProvider(p.Provider, p.Endpoint, path)
	if err != nil {
		err = erlogs.Convert(err).Wrap("failed to AddRemoteProvider").Panic(
			erlogs.OptionFields(
				zap.String("provider", p.Provider),
				zap.String("endpoint", p.Endpoint),
				zap.String("path", path),
			),
		)
		return err
	}

	p.Viper.SetConfigType(p.Type)
	err = p.Viper.ReadRemoteConfig()
	if err != nil {
		err = erlogs.Convert(err).Wrap("failed to ReadRemoteConfig").Panic(
			erlogs.OptionFields(
				zap.String("provider", p.Provider),
				zap.String("endpoint", p.Endpoint),
				zap.String("path", path),
			),
		)
		return err
	}

	erlogs.New("read remote config success").Info(
		erlogs.OptionFields(
			zap.String("provider", p.Provider),
			zap.String("endpoint", p.Endpoint),
			zap.String("path", path),
		),
	).Options(BaseELOptions()).RecordLog(ctx)

	return nil
}

func (p *RemoteProvider) checkOptions() (err error) {
	if p == nil {
		return erlogs.New("RemoteProvider is nil").Panic()
	}

	if len(p.Endpoint) == 0 {
		err = erlogs.New("option invalid: config endpoint empty").Panic()
		return
	}

	if len(p.Path) == 0 {
		err = erlogs.New("option invalid: config path empty").Panic()
		return
	}

	if len(p.Type) == 0 {
		err = erlogs.New("option invalid: config type empty").Panic()
	}

	return
}

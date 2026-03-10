package vipers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	fs2 "github.com/mel0dys0ng/song/pkg/fs"
	"github.com/mel0dys0ng/song/pkg/metas"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type LocalProvider struct {
	*Options
	*viper.Viper
	configFileExt string
}

func NewLocalProvider(v *viper.Viper, o *Options) ProviderInterface {
	return &LocalProvider{
		Viper:   v,
		Options: o,
	}
}

func (p *LocalProvider) Load() (err error) {
	if p == nil {
		return erlogs.New("LocalProvider is nil").Panic()
	}

	err = p.checkOptions()
	if err != nil {
		return
	}

	files, err := p.findConfigFiles()
	if err != nil {
		return
	}

	p.Viper.AddConfigPath(p.Path)
	p.Viper.SetConfigType(p.Type)

	for i, file := range files {
		if err = p.loadConfigFile(i, file); err != nil {
			return
		}
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

func (p *LocalProvider) checkOptions() (err error) {
	if p == nil {
		return errors.New("LocalProvider is nil")
	}

	if len(p.Path) == 0 {
		err = errors.New("option invalid: config path empty")
		return
	}

	if len(p.Type) == 0 {
		err = errors.New("option invalid: config type empty")
		return
	}

	p.configFileExt = fmt.Sprintf(".%s", p.Type)
	return
}

func (p *LocalProvider) findConfigFiles() (files []string, err error) {
	if p == nil {
		err = erlogs.New("LocalProvider is nil").Panic()
		return
	}

	fileName := fmt.Sprintf("%s%s", metas.Metadata().Mode().String(), p.configFileExt)
	if strings.HasSuffix(p.Path, fileName) {
		files = []string{p.Path}
	} else {
		files = []string{fmt.Sprintf("%s/%s", p.Path, fileName)}
	}

	for _, file := range files {
		exists, errCheck := fs2.Exists(file)
		if errCheck != nil {
			err = erlogs.Convert(errCheck).Panic(
				erlogs.OptionFields(
					zap.String("file", file),
				),
			)
			return
		}
		if !exists {
			err = erlogs.New("config file not exist").Panic(
				erlogs.OptionFields(
					zap.String("file", file),
				),
			)
			return
		}
	}

	return
}

func (p *LocalProvider) loadConfigFile(index int, file string) (err error) {
	if p == nil {
		return erlogs.New("LocalProvider is nil").Panic()
	}

	ctx := context.Background()
	p.Viper.SetConfigName(strings.TrimSuffix(filepath.Base(file), p.configFileExt))
	if index == 0 {
		err = p.Viper.ReadInConfig()
	} else {
		err = p.Viper.MergeInConfig()
	}

	if err != nil {
		err = erlogs.Convert(err).Panic(
			erlogs.OptionFields(
				zap.String("file", file),
			),
		)
		return
	}

	erlogs.New("load config file success").Info(
		erlogs.OptionFields(
			zap.String("file", file),
		),
	).Options(BaseELOptions()).RecordLog(ctx)

	return
}

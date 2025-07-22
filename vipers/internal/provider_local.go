package internal

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/erlogs"
	fs2 "github.com/mel0dys0ng/song/utils/fs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
 *========================================**
 * 加载本地配置
 *========================================**
 */

type LocalProvider struct {
	*Options
	*viper.Viper
	configFileExt string // 配置文件后缀
	elg           erlogs.ErLogInterface
}

func NewLocalProvider(v *viper.Viper, elg erlogs.ErLogInterface, o *Options) ProviderInterface {
	return &LocalProvider{
		Viper:   v,
		Options: o,
		elg:     elg,
	}
}

func (p *LocalProvider) Load() (err error) {
	err = p.checkOptions()
	if err != nil {
		err = p.elg.Panic(
			context.Background(),
			erlogs.Msgv("check options failed"),
			erlogs.Content(err.Error()),
		)
		return
	}

	files, err := p.findConfigFiles()
	if err != nil {
		err = p.elg.Panic(
			context.Background(),
			erlogs.Msgv("find config files failed"),
			erlogs.Content(err.Error()),
		)
		return
	}

	// setting
	p.Viper.AddConfigPath(p.Path)
	p.Viper.SetConfigType(p.Type)

	// load
	for k, v := range files {
		if err = p.loadConfigFile(k, v); err != nil {
			return
		}
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

func (p *LocalProvider) checkOptions() (err error) {
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

// 查找配置文件
func (p *LocalProvider) findConfigFiles() (files []string, err error) {
	return fs2.WalkDir(
		p.Path,
		func(path string, info fs.FileInfo) bool {
			return !info.IsDir() && strings.HasSuffix(info.Name(), p.configFileExt)
		},
	)
}

// 加载配置文件
func (p *LocalProvider) loadConfigFile(index int, file string) (err error) {
	p.Viper.SetConfigName(strings.TrimSuffix(filepath.Base(file), p.configFileExt))
	if index == 0 {
		err = p.Viper.ReadInConfig()
	} else {
		err = p.Viper.MergeInConfig()
	}

	if err != nil {
		err = p.elg.Panic(
			context.Background(),
			erlogs.Msgv("load config file failed"),
			erlogs.Content(err.Error()),
			erlogs.Fields(zap.String("configFile", file)),
		)
		return
	}

	p.elg.InfoL(
		context.Background(),
		erlogs.Msgv("load config file success"),
		erlogs.Fields(zap.String("configFile", file)),
	)

	return
}

package internal

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/pkgs/erlogs"
)

type Options struct {
	Provider       string                         `json:"provider"` // 配置类型，仅支持yaml|json|toml|etcd|consul
	Type           string                         `json:"type"`     // 配置类型，仅支持JSON|TOML|YAML
	Endpoint       string                         `json:"endpoint"` // 远程配置endpoint
	Path           string                         `json:"path"`     // 配置路径
	OnChangeConfig func(fsnotify.Event, *Options) `json:"-"`
}

type Option struct {
	fn func(options *Options)
}

func Provider(s string) Option {
	return Option{
		func(options *Options) {
			options.Provider = s
		},
	}
}

func Type(s string) Option {
	return Option{
		func(options *Options) {
			options.Type = s
		},
	}
}

func Endpoint(s string) Option {
	return Option{
		func(options *Options) {
			options.Endpoint = s
		},
	}
}

func Path(s string) Option {
	return Option{
		func(options *Options) {
			options.Path = s
		},
	}
}

func OnChangeConfig(f func(fsnotify.Event, *Options)) Option {
	return Option{
		func(options *Options) {
			options.OnChangeConfig = f
		},
	}
}

func DefaultOptions(elg erlogs.ErLogInterface) *Options {
	return &Options{
		Provider: ConfigProviderYaml,
		Type:     ConfigProviderYaml,
		OnChangeConfig: func(event fsnotify.Event, options *Options) {
			elg.InfoL(
				context.Background(),
				erlogs.Msgv("config has changed"),
				erlogs.Content(event.String()),
			)
		},
	}
}

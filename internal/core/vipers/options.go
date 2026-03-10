package vipers

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

type Options struct {
	Provider       string                         `json:"provider"`
	Type           string                         `json:"type"`
	Endpoint       string                         `json:"endpoint"`
	Path           string                         `json:"path"`
	OnChangeConfig func(fsnotify.Event, *Options) `json:"-"`
}

type Option func(options *Options)

func Provider(s string) Option {
	return func(options *Options) {
		options.Provider = s
	}
}

func Type(s string) Option {
	return func(options *Options) {
		options.Type = s
	}
}

func Endpoint(s string) Option {
	return func(options *Options) {
		options.Endpoint = s
	}
}

func Path(s string) Option {
	return func(options *Options) {
		options.Path = s
	}
}

func OnChangeConfig(f func(fsnotify.Event, *Options)) Option {
	return func(options *Options) {
		options.OnChangeConfig = f
	}
}

func DefaultOptions() *Options {
	return &Options{
		Provider: ConfigProviderYaml,
		Type:     ConfigProviderYaml,
		OnChangeConfig: func(event fsnotify.Event, options *Options) {
			ctx := context.Background()
			erlogs.New("config has changed").Options(BaseELOptions()).InfoLog(ctx)
		},
	}
}

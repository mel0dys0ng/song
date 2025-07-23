package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/utils/sys"
	"github.com/mel0dys0ng/song/pkg/vipers"
)

const (
	ConfigKey                = "https"
	DefaultPort              = 8080
	DefaultHost              = "0.0.0.0"
	DefaultHttpsOpen         = true
	DefaultHttpsKeyFile      = ""
	DefaultHttpsKeyCert      = ""
	DefaultKeepAlive         = true
	DefaultReadTimeout       = 1 * time.Second
	DefaultReadHeaderTimeout = 1 * time.Second
	DefaultWriteTimeout      = 1 * time.Second
	DefaultIdleTimeout       = 1 * time.Second
	DefaultHammerTime        = 10 * time.Second
	DefaultMaxHeaderBytes    = 1 << 16
	DefaultTmpDir            = "./tmp"
	DefaultCorsMaxAge        = 12 * time.Hour
)

type (
	Config struct {
		Port              int            `json:"port" yaml:"port"`
		TLSOpen           bool           `json:"TLSOpen" yaml:"TLSOpen"`
		TLSKeyFile        string         `json:"TLSKeyFile" yaml:"TLSKeyFile"`
		TLSCertFile       string         `json:"TLSCertFile" yaml:"TLSCertFile"`
		KeepAlive         bool           `json:"keepAlive" yaml:"keepAlive"`
		ReadTimeout       string         `json:"readTimeout" yaml:"readTimeout"`
		ReadHeaderTimeout string         `json:"readHeaderTimeout" yaml:"readHeaderTimeout"`
		WriteTimeout      string         `json:"writeTimeout" yaml:"writeTimeout"`
		IdleTimeout       string         `json:"idleTimeout" yaml:"idleTimeout"`
		HammerTime        string         `json:"hammerTime" yaml:"hammerTime"`
		MaxHeaderBytes    int            `json:"maxHeaderBytes" yaml:"maxHeaderBytes"`
		TmpDir            string         `json:"tmpDir" yaml:"tmpDir"`
		LoggerHeaderKeys  []string       `json:"loggerHeaderKeys" yaml:"loggerHeaderKeys"`
		ErLog             *erlogs.Config `json:"erlog" yaml:"erlog"`
		Cors              *Cors          `json:"cors" yaml:"cors"`
	}

	Cors struct {
		Enable           bool     `json:"enable" yaml:"enable"`
		AllowOrigins     []string `json:"allowOrigins" yaml:"allowOrigins"`
		AllowHeaders     []string `json:"allowHeaders" yaml:"allowHeaders"`
		AllowMethods     []string `json:"allowMethods" yaml:"allowMethods"`
		ExposeHeaders    []string `json:"exposeHeaders" yaml:"exposeHeaders"`
		AllowCredentials bool     `json:"allowCredentials" yaml:"allowCredentials"`
		AllowWildcard    bool     `json:"allowWildcard" yaml:"allowWildcard"`
		MaxAge           string   `json:"maxAge" yaml:"maxAge"`
	}

	Options struct {
		// HTTP Server Config
		Config

		Addr        string
		OnStart     func()
		OnStartFail func(error)
		OnShutdown  func()
		OnExit      func()
		OnResponded func(context.Context, *RequestResponseData)
		Inits       []Init
		Defers      []Defer
		Routes      []Route
		Middlewares []Middleware
	}

	Init                 func() error
	Defer                func()
	Route                func(eng *gin.Engine)
	MiddlewareHandleFunc func(eng *gin.Engine) gin.HandlerFunc
	Middleware           struct {
		Priority int                  // 优先级
		Handle   MiddlewareHandleFunc // 中间件
	}

	Option struct {
		Apply func(*Options)
	}
)

func newOptions(opts []Option) *Options {
	options := &Options{
		Config: Config{
			Port:           DefaultPort,
			TLSOpen:        DefaultHttpsOpen,
			TLSKeyFile:     DefaultHttpsKeyFile,
			TLSCertFile:    DefaultHttpsKeyCert,
			KeepAlive:      DefaultKeepAlive,
			MaxHeaderBytes: DefaultMaxHeaderBytes,
			TmpDir:         DefaultTmpDir,
		},
	}

	err := vipers.UnmarshalKey(ConfigKey, &options.Config)
	if err != nil {
		sys.Panicf("failed to load and set the config of http server: %s", err.Error())
		return nil
	}

	for _, v := range opts {
		if v.Apply != nil {
			v.Apply(options)
		}
	}

	options.Addr = fmt.Sprintf("%s:%d", DefaultHost, options.Port)

	return options
}

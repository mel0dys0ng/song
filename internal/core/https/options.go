package https

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/sys"
	"github.com/mel0dys0ng/song/pkg/vipers"
)

const (
	ConfigKey                = "https"
	DefaultPort              = 8080
	DefaultHost              = "0.0.0.0"
	DefaultHttpsOpen         = false // 修改默认值为false，更安全
	DefaultHttpsKeyFile      = ""
	DefaultHttpsKeyCert      = ""
	DefaultKeepAlive         = true
	DefaultReadTimeout       = 30 * time.Second // 增加默认值，避免请求频繁超时
	DefaultReadHeaderTimeout = 30 * time.Second
	DefaultWriteTimeout      = 60 * time.Second
	DefaultIdleTimeout       = 60 * time.Second
	DefaultHammerTime        = 30 * time.Second // 增加关闭等待时间，确保请求完成
	DefaultMaxHeaderBytes    = 1 << 20          // 1MB, 增加默认值
	DefaultTmpDir            = "./tmp"
	DefaultCorsMaxAge        = 12 * time.Hour
)

type (
	Config struct {
		Port              int            `json:"port" yaml:"port" mapstructure:"port"`
		TLSOpen           bool           `json:"TLSOpen" yaml:"TLSOpen" mapstructure:"TLSOpen"`
		TLSKeyFile        string         `json:"TLSKeyFile" yaml:"TLSKeyFile" mapstructure:"TLSKeyFile"`
		TLSCertFile       string         `json:"TLSCertFile" yaml:"TLSCertFile" mapstructure:"TLSCertFile"`
		KeepAlive         bool           `json:"keepAlive" yaml:"keepAlive" mapstructure:"keepAlive"`
		ReadTimeout       string         `json:"readTimeout" yaml:"readTimeout" mapstructure:"readTimeout"`
		ReadHeaderTimeout string         `json:"readHeaderTimeout" yaml:"readHeaderTimeout" mapstructure:"readHeaderTimeout"`
		WriteTimeout      string         `json:"writeTimeout" yaml:"writeTimeout" mapstructure:"writeTimeout"`
		IdleTimeout       string         `json:"idleTimeout" yaml:"idleTimeout" mapstructure:"idleTimeout"`
		HammerTime        string         `json:"hammerTime" yaml:"hammerTime" mapstructure:"hammerTime"`
		MaxHeaderBytes    int            `json:"maxHeaderBytes" yaml:"maxHeaderBytes" mapstructure:"maxHeaderBytes"`
		TmpDir            string         `json:"tmpDir" yaml:"tmpDir" mapstructure:"tmpDir"`
		LoggerHeaderKeys  []string       `json:"loggerHeaderKeys" yaml:"loggerHeaderKeys" mapstructure:"loggerHeaderKeys"`
		ErLog             *erlogs.Config `json:"erlog" yaml:"erlog" mapstructure:"erlog"`
		Cors              *Cors          `json:"cors" yaml:"cors" mapstructure:"cors"`
		Csrf              *CSRF          `json:"csrf" yaml:"csrf" mapstructure:"csrf"`
		Sign              *Sign          `json:"sign" yaml:"sign" mapstructure:"sign"`
	}

	Cors struct {
		Enable                    bool     `json:"enable" yaml:"enable" mapstructure:"enable"`
		AllowOrigins              []string `json:"allowOrigins" yaml:"allowOrigins" mapstructure:"allowOrigins"`
		AllowHeaders              []string `json:"allowHeaders" yaml:"allowHeaders" mapstructure:"allowHeaders"`
		AllowMethods              []string `json:"allowMethods" yaml:"allowMethods" mapstructure:"allowMethods"`
		ExposeHeaders             []string `json:"exposeHeaders" yaml:"exposeHeaders" mapstructure:"exposeHeaders"`
		AllowCredentials          bool     `json:"allowCredentials" yaml:"allowCredentials" mapstructure:"allowCredentials"`
		AllowWildcard             bool     `json:"allowWildcard" yaml:"allowWildcard" mapstructure:"allowWildcard"`
		MaxAge                    string   `json:"maxAge" yaml:"maxAge" mapstructure:"maxAge"`
		OptionsResponseStatusCode int      `json:"optionsResponseStatusCode" yaml:"optionsResponseStatusCode" mapstructure:"optionsResponseStatusCode"`
	}

	CSRF struct {
		Enable         bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
		LookupType     string `json:"lookupType" yaml:"lookupType" mapstructure:"lookupType"`
		LookupName     string `json:"lookupName" yaml:"lookupName" mapstructure:"lookupName"`
		CookieName     string `json:"cookieName" yaml:"cookieName" mapstructure:"cookieName"`
		CookieDomain   string `json:"cookieDomain" yaml:"cookieDomain" mapstructure:"cookieDomain"`
		CookiePath     string `json:"cookiePath" yaml:"cookiePath" mapstructure:"cookiePath"`
		CookieMaxAge   int    `json:"cookieMaxAge" yaml:"cookieMaxAge" mapstructure:"cookieMaxAge"`
		CookieSecure   bool   `json:"cookieSecure" yaml:"cookieSecure" mapstructure:"cookieSecure"`
		CookieHttpOnly bool   `json:"cookieHttpOnly" yaml:"cookieHttpOnly" mapstructure:"cookieHttpOnly"`
	}

	Sign struct {
		Enable   bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
		Secret   string `json:"secret" yaml:"secret" mapstructure:"secret"`
		TTL      int    `json:"ttl" yaml:"ttl" mapstructure:"ttl"` // TTL in seconds
		Query    bool   `json:"query" yaml:"query" mapstructure:"query"`
		FormData bool   `json:"formData" yaml:"formData" mapstructure:"formData"`
		Header   bool   `json:"header" yaml:"header" mapstructure:"header"`
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
		OnRecovered func(context.Context, *RecoveredData)
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

	Option func(*Options)
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

	for _, opt := range opts {
		if opt != nil {
			opt(options)
		}
	}

	if options.Addr == "" {
		options.Addr = fmt.Sprintf("%s:%d", DefaultHost, options.Port)
	}

	return options
}

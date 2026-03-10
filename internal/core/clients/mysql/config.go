package mysql

import (
	"context"
	"time"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/vipers"
)

// Config 数据库配置
type Config struct {
	Debug                   *bool         `mapstructure:"debug" json:"debug" yaml:"debug"`                                                       // debug模式
	Master                  string        `mapstructure:"master" json:"master" yaml:"master"`                                                    // 主数据库连接DSN
	Slaves                  []string      `mapstructure:"slaves" json:"slaves" yaml:"slaves"`                                                    // 从数据库连接DSN
	TablePrefix             string        `mapstructure:"tablePrefix" json:"tablePrefix" yaml:"tablePrefix"`                                     // 表前缀
	SingularTable           *bool         `mapstructure:"singularTable" json:"singularTable" yaml:"singularTable"`                               // 是否单数表名
	MaxIdle                 int           `mapstructure:"maxIdle" json:"maxIdle" yaml:"maxIdle"`                                                 // 最大的可空闲连接数
	MaxActive               int           `mapstructure:"maxActive" json:"maxActive" yaml:"maxActive"`                                           // 最大的连接数，必须大于MaxIdle
	MaxConnLifeTime         time.Duration `mapstructure:"maxConnLifeTime" json:"maxConnLifeTime" yaml:"maxConnLifeTime"`                         // 连接最大有效期，单位:毫秒
	IdleTimeout             time.Duration `mapstructure:"idleTimeout" json:"idleTimeout" yaml:"idleTimeout"`                                     // 连接最长空闲期(空闲超时)，单位:毫秒
	LogSlow                 time.Duration `mapstructure:"logSlow" json:"logSlow" yaml:"logSlow"`                                                 // 慢查询时间上限，单位:毫秒
	LogDir                  string        `mapstructure:"logDir" json:"logDir" yaml:"logDir"`                                                    // 日志文件目录
	LogLevel                string        `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`                                              // 日志级别
	LogFileSuffix           string        `mapstructure:"logFileSuffix" json:"logFileSuffix" yaml:"logFileSuffix"`                               // 日志文件后缀
	LogIgnoreRecordNotFound *bool         `mapstructure:"logIgnoreRecordNotFound" json:"logIgnoreRecordNotFound" yaml:"logIgnoreRecordNotFound"` // 忽略查询记录不存在的错误
	LogColorful             *bool         `mapstructure:"logColorful" json:"logColorful" yaml:"logColorful"`                                     // 日志颜色
}

func ptrBool(b bool) *bool {
	return &b
}

func newConfig(ctx context.Context, key string, opts []Option) (config *Config, err error) {

	config = &Config{
		Debug:                   ptrBool(DefaultDebug),
		Master:                  DefaultMaster,
		Slaves:                  []string{DefaultSlave},
		TablePrefix:             DefaultTablePrefix,
		SingularTable:           ptrBool(DefaultSingularTable),
		MaxIdle:                 DefaultMaxIdle,
		MaxActive:               DefaultMaxActive,
		MaxConnLifeTime:         DefaultMaxConnLifeTime,
		IdleTimeout:             DefaultIdleTimeout,
		LogSlow:                 DefaultLogSlow,
		LogDir:                  DefaultLogDir,
		LogLevel:                DefaultLogLevel,
		LogFileSuffix:           DefaultLogFileSuffix,
		LogIgnoreRecordNotFound: ptrBool(DefaultLogIgnoreRecordNotFound),
		LogColorful:             ptrBool(DefaultLogColorful),
	}

	options, err := config.LoadOptions(ctx, key)
	if err != nil {
		return
	}

	opts = append(options, opts...)
	for _, v := range opts {
		v(config)
	}

	err = config.check()
	if err != nil {
		return
	}

	return
}

// LoadOptions 加载配置并返回 Options
func (c *Config) LoadOptions(ctx context.Context, key string) (opts []Option, err error) {

	values := &Config{}
	er := vipers.UnmarshalKey(key, values)
	if er != nil {
		err = erlogs.Convert(er).Wrap("failed to load options: UnmarshalKey error").Panic()
		return
	}

	if values.Debug != nil {
		opts = append(opts, Debug(*values.Debug))
	}

	if values.TablePrefix != DefaultTablePrefix {
		opts = append(opts, TablePrefix(values.TablePrefix))
	}

	if values.Master != DefaultMaster {
		opts = append(opts, Master(values.Master))
	}

	if len(values.Slaves) != 1 || (len(values.Slaves) == 1 && values.Slaves[0] != DefaultSlave) {
		opts = append(opts, Slaves(values.Slaves))
	}

	if values.IdleTimeout != DefaultIdleTimeout {
		opts = append(opts, IdleTimeout(values.IdleTimeout))
	}

	if values.MaxIdle != DefaultMaxIdle {
		opts = append(opts, MaxIdle(values.MaxIdle))
	}

	if values.MaxActive != DefaultMaxActive {
		opts = append(opts, MaxActive(values.MaxActive))
	}

	if values.SingularTable != nil {
		opts = append(opts, SingularTable(*values.SingularTable))
	}

	if values.MaxConnLifeTime != DefaultMaxConnLifeTime {
		opts = append(opts, MaxConnLifeTime(values.MaxConnLifeTime))
	}

	if values.LogColorful != nil {
		opts = append(opts, LogColorful(*values.LogColorful))
	}

	if values.LogIgnoreRecordNotFound != nil {
		opts = append(opts, LogIgnoreRecordNotFound(*values.LogIgnoreRecordNotFound))
	}

	if values.LogFileSuffix != DefaultLogFileSuffix {
		opts = append(opts, LogFileSuffix(values.LogFileSuffix))
	}

	if values.LogSlow != DefaultLogSlow {
		opts = append(opts, LogSlow(values.LogSlow))
	}

	if values.LogLevel != DefaultLogLevel {
		opts = append(opts, LogLevel(values.LogLevel))
	}

	if values.LogDir != DefaultLogDir {
		opts = append(opts, LogDir(values.LogDir))
	}

	return
}

// check validates the configuration
func (c *Config) check() (err error) {
	if len(c.Master) == 0 {
		err = erlogs.New("master is empty").Panic()
		return
	}

	if len(c.Slaves) == 0 {
		err = erlogs.New("slaves is empty").Panic()
		return
	}

	if c.MaxActive <= c.MaxIdle {
		err = erlogs.New("maxActive must > maxIdle").Panic()
		return
	}

	if c.MaxIdle <= 0 {
		err = erlogs.New("config maxIdel must > 0").Panic()
		return
	}

	if c.MaxConnLifeTime <= 0 {
		err = erlogs.New("config maxConnLifeTime must > 0").Panic()
		return
	}

	if c.IdleTimeout <= 0 {
		err = erlogs.New("config idleTimeout must > 0").Panic()
		return
	}

	if c.MaxActive <= c.MaxIdle {
		err = erlogs.New("config maxActive must > maxIdle").Panic()
		return
	}

	if c.MaxConnLifeTime <= 0 {
		err = erlogs.New("config maxConnLifeTime must > 0").Panic()
		return
	}

	if c.IdleTimeout <= 0 {
		err = erlogs.New("config idleTimeout must > 0").Panic()
		return
	}

	if c.MaxConnLifeTime <= 0 {
		err = erlogs.New("config maxConnLifeTime must > 0").Panic()
		return
	}

	if c.IdleTimeout <= 0 {
		err = erlogs.New("config idleTimeout must > 0").Panic()
		return
	}

	// 验证 Master DSN 格式
	if !isValidDSN(c.Master) {
		err = erlogs.New("master DSN format is invalid").Panic()
		return
	}

	// 验证 Slaves DSN 格式
	for _, slave := range c.Slaves {
		if !isValidDSN(slave) {
			err = erlogs.New("slave DSN format is invalid").Panic()
			return
		}
	}

	return
}

// isValidDSN checks if the given DSN is in a valid format
func isValidDSN(dsn string) bool {
	return len(dsn) > 0
}

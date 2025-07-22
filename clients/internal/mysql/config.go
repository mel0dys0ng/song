package mysql

import (
	"context"
	"time"

	"github.com/mel0dys0ng/song/erlogs"
	"github.com/mel0dys0ng/song/vipers"
)

// Config 数据库配置
type Config struct {
	Debug                   bool          `mapstructure:"debug" json:"debug"`                                     // debug模式
	Master                  string        `mapstructure:"master" json:"master"`                                   // 主数据库连接DSN
	Slaves                  []string      `mapstructure:"slaves" json:"slaves"`                                   // 从数据库连接DSN
	TablePrefix             string        `mapstructure:"tablePrefix" json:"tablePrefix"`                         // 表前缀
	SingularTable           bool          `mapstructure:"singularTable" json:"singularTable"`                     // 是否单数表名
	MaxIdle                 int           `mapstructure:"maxIdle" json:"maxIdle"`                                 // 最大的可空闲连接数
	MaxActive               int           `mapstructure:"maxActive" json:"maxActive"`                             // 最大的连接数，必须大于MaxIdle
	MaxConnLifeTime         time.Duration `mapstructure:"maxConnLifeTime" json:"maxConnLifeTime"`                 // 连接最大有效期，单位:毫秒
	IdleTimeout             time.Duration `mapstructure:"idleTimeout" json:"idleTimeout"`                         // 连接最长空闲期(空闲超时)，单位:毫秒
	LogSlow                 time.Duration `mapstructure:"logSlow" json:"logSlow"`                                 // 慢查询时间上限，单位:毫秒
	LogDir                  string        `mapstructure:"logDir" json:"logDir"`                                   // 日志文件目录
	LogLevel                string        `mapstructure:"logLevel" json:"logLevel"`                               // 日志级别
	LogFileSuffix           string        `mapstructure:"logFileSuffix" json:"logFileSuffix"`                     // 日志文件后缀
	LogIgnoreRecordNotFound bool          `mapstructure:"logIgnoreRecordNotFound" json:"logIgnoreRecordNotFound"` // 忽略查询记录不存在的错误
	LogColorful             bool          `mapstructure:"logColorful" json:"logColorful"`                         // 日志颜色
}

func newConfig(ctx context.Context, key string, elg erlogs.ErLogInterface, opts []Option) (
	config *Config, err erlogs.ErLogInterface) {

	config = &Config{
		Debug:                   DefaultDebug,
		Master:                  DefaultMaster,
		Slaves:                  []string{DefaultSlave},
		TablePrefix:             DefaultTablePrefix,
		SingularTable:           DefaultSingularTable,
		MaxIdle:                 DefaultMaxIdle,
		MaxActive:               DefaultMaxActive,
		MaxConnLifeTime:         DefaultMaxConnLifeTime,
		IdleTimeout:             DefaultIdleTimeout,
		LogSlow:                 DefaultLogSlow,
		LogDir:                  DefaultLogDir,
		LogLevel:                DefaultLogLevel,
		LogFileSuffix:           DefaultLogFileSuffix,
		LogIgnoreRecordNotFound: DefaultLogIgnoreRecordNotFound,
		LogColorful:             DefaultLogColorful,
	}

	options, err := config.LoadOptions(ctx, key, elg)
	if err != nil {
		return
	}

	opts = append(options, opts...)
	for _, v := range opts {
		v.Func(config)
	}

	err = config.check(ctx, elg)

	return
}

// LoadOptions 加载配置并返回Options
func (c *Config) LoadOptions(ctx context.Context, key string, elg erlogs.ErLogInterface) (
	opts []Option, err erlogs.ErLogInterface) {

	values := &Config{}
	er := vipers.UnmarshalKey(key, values, nil)
	if er != nil {
		err = elg.PanicE(ctx,
			erlogs.Msgv("failed to load options: UnmarshalKey error"),
			erlogs.Content(er.Error()),
		)
		return
	}

	if values == nil {
		return
	}

	if values.Debug != DefaultDebug {
		opts = append(opts, Debug(values.Debug))
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

	if values.SingularTable != DefaultSingularTable {
		opts = append(opts, SingularTable(values.SingularTable))
	}

	if values.MaxConnLifeTime != DefaultMaxConnLifeTime {
		opts = append(opts, MaxConnLifeTime(values.MaxConnLifeTime))
	}

	if values.LogColorful != DefaultLogColorful {
		opts = append(opts, LogColorful(values.LogColorful))
	}

	if values.LogIgnoreRecordNotFound != DefaultLogIgnoreRecordNotFound {
		opts = append(opts, LogIgnoreRecordNotFound(values.LogIgnoreRecordNotFound))
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

// 校验配置
func (c *Config) check(ctx context.Context, elg erlogs.ErLogInterface) (err erlogs.ErLogInterface) {
	if len(c.Master) == 0 {
		err = elg.ErorrE(ctx,
			erlogs.Msgv("mysql config invalid"),
			erlogs.Content("config master empty"),
		)
		return
	}

	if len(c.Slaves) == 0 {
		err = elg.ErorrE(ctx,
			erlogs.Msgv("mysql config invalid"),
			erlogs.Content("config slaves empty"),
		)
		return
	}

	if c.MaxIdle <= 0 {
		err = elg.ErorrE(ctx,
			erlogs.Msgv("mysql config invalid"),
			erlogs.Content("config maxIdle invalid"),
		)
		return
	}

	if c.MaxActive <= c.MaxIdle {
		err = elg.ErorrE(ctx,
			erlogs.Msgv("mysql config invalid"),
			erlogs.Content("config maxActive must > maxIdle"),
		)
		return
	}

	if c.MaxConnLifeTime <= 0 {
		err = elg.ErorrE(ctx,
			erlogs.Msgv("mysql config invalid"),
			erlogs.Content("config maxConnLifeTime invalid"),
		)
		return
	}

	if c.IdleTimeout <= 0 {
		err = elg.ErorrE(ctx,
			erlogs.Msgv("mysql config invalid"),
			erlogs.Content("config idleTimeout invalid"),
		)
		return
	}

	return
}

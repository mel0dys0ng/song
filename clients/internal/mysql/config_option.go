package mysql

import "time"

const (
	DefaultDebug                   = true
	DefaultTablePrefix             = ""
	DefaultMaster                  = ""
	DefaultSlave                   = ""
	DefaultSingularTable           = false
	DefaultMaxIdle                 = 100
	DefaultMaxActive               = 200
	DefaultMaxConnLifeTime         = 3600000
	DefaultIdleTimeout             = 3600000
	DefaultLogSlow                 = 1000
	DefaultLogDir                  = "/data/logs/mysql"
	DefaultLogLevel                = "warn"
	DefaultLogFileSuffix           = "log"
	DefaultLogIgnoreRecordNotFound = true
	DefaultLogColorful             = true
)

type Option struct {
	Func func(*Config)
}

func Debug(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.Debug = b
		},
	}
}

func TablePrefix(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.TablePrefix = s
		},
	}
}

func SingularTable(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.SingularTable = b
		},
	}
}

func Master(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.Master = s
		},
	}
}

func Slaves(ss []string) Option {
	return Option{
		Func: func(config *Config) {
			config.Slaves = ss
		},
	}
}

func MaxIdle(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxIdle = n
		},
	}
}

func MaxActive(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxActive = n
		},
	}
}

func IdleTimeout(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.IdleTimeout = t
		},
	}
}

func MaxConnLifeTime(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxConnLifeTime = t
		},
	}
}

func LogDir(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.LogDir = s
		},
	}
}

func LogSlow(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.LogSlow = t
		},
	}
}

func LogLevel(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.LogLevel = s
		},
	}
}

func LogFileSuffix(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.LogFileSuffix = s
		},
	}
}

func LogColorful(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.LogColorful = b
		},
	}
}

func LogIgnoreRecordNotFound(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.LogIgnoreRecordNotFound = b
		},
	}
}

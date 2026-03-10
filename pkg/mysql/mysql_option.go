package mysql

import (
	"time"

	"github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

func OptionDebug(b bool) mysql.Option {
	return mysql.Debug(b)
}

func OptionTablePrefix(s string) mysql.Option {
	return mysql.TablePrefix(s)
}

func OptionSingularTable(b bool) mysql.Option {
	return mysql.SingularTable(b)
}

func OptionMaster(s string) mysql.Option {
	return mysql.Master(s)
}

func OptionSlaves(ss []string) mysql.Option {
	return mysql.Slaves(ss)
}

func OptionMaxIdle(n int) mysql.Option {
	return mysql.MaxIdle(n)
}

func OptionMaxActive(n int) mysql.Option {
	return mysql.MaxActive(n)
}

func OptionIdleTimeout(t time.Duration) mysql.Option {
	return mysql.IdleTimeout(t)
}

func OptionMaxConnLifeTime(t time.Duration) mysql.Option {
	return mysql.MaxConnLifeTime(t)
}

func OptionLogDir(s string) mysql.Option {
	return mysql.LogDir(s)
}

func OptionLogSlow(t time.Duration) mysql.Option {
	return mysql.LogSlow(t)
}

func OptionLogLevel(s string) mysql.Option {
	return mysql.LogLevel(s)
}

func OptionLogFileSuffix(s string) mysql.Option {
	return mysql.LogFileSuffix(s)
}

func OptionLogColorful(b bool) mysql.Option {
	return mysql.LogColorful(b)
}

func OptionLogIgnoreRecordNotFound(b bool) mysql.Option {
	return mysql.LogIgnoreRecordNotFound(b)
}

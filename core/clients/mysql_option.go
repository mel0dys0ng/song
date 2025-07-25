package clients

import (
	"time"

	"github.com/mel0dys0ng/song/core/clients/internal/mysql"
)

func MySQLOptionDebug(b bool) mysql.Option {
	return mysql.Debug(b)
}

func MySQLOptionTablePrefix(s string) mysql.Option {
	return mysql.TablePrefix(s)
}

func MySQLOptionSingularTable(b bool) mysql.Option {
	return mysql.SingularTable(b)
}

func MySQLOptionMaster(s string) mysql.Option {
	return mysql.Master(s)
}

func MySQLOptionSlaves(ss []string) mysql.Option {
	return mysql.Slaves(ss)
}

func MySQLOptionMaxIdle(n int) mysql.Option {
	return mysql.MaxIdle(n)
}

func MySQLOptionMaxActive(n int) mysql.Option {
	return mysql.MaxActive(n)
}

func MySQLOptionIdleTimeout(t time.Duration) mysql.Option {
	return mysql.IdleTimeout(t)
}

func MySQLOptionMaxConnLifeTime(t time.Duration) mysql.Option {
	return mysql.MaxConnLifeTime(t)
}

func MySQLOptionLogDir(s string) mysql.Option {
	return mysql.LogDir(s)
}

func MySQLOptionLogSlow(t time.Duration) mysql.Option {
	return mysql.LogSlow(t)
}

func MySQLOptionLogLevel(s string) mysql.Option {
	return mysql.LogLevel(s)
}

func MySQLOptionLogFileSuffix(s string) mysql.Option {
	return mysql.LogFileSuffix(s)
}

func MySQLOptionLogColorful(b bool) mysql.Option {
	return mysql.LogColorful(b)
}

func MySQLOptionLogIgnoreRecordNotFound(b bool) mysql.Option {
	return mysql.LogIgnoreRecordNotFound(b)
}

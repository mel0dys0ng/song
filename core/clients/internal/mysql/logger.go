package mysql

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mel0dys0ng/song/core/erlogs"
	"github.com/mel0dys0ng/song/core/utils/fs"
	"gorm.io/gorm/logger"
)

// newLogger return new gorm logger interface object
func newLogger(ctx context.Context, el erlogs.ErLogInterface, config *Config) (
	lgr logger.Interface, err erlogs.ErLogInterface) {

	loggerWriter, err := newLoggerWriter(ctx, el, config)
	if err != nil {
		return
	}

	lgr = logger.New(
		log.New(loggerWriter, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             config.LogSlow * time.Millisecond, // Slow SQL threshold
			LogLevel:                  getLogLevel(config.LogLevel),      // Log level
			IgnoreRecordNotFoundError: config.LogIgnoreRecordNotFound,    // Ignore ErrRecordNotFound error for logger
			Colorful:                  config.LogColorful,                // Disable color
		},
	)

	return
}

// newLoggerWriter return new logger writer
func newLoggerWriter(ctx context.Context, el erlogs.ErLogInterface, config *Config) (
	writer *os.File, err erlogs.ErLogInterface) {

	writer = os.Stdout
	if config.Debug {
		return
	}

	dir := ""
	if len(config.LogDir) > 0 {
		dir = strings.TrimRight(config.LogDir, "/")
	}

	suffix := ""
	if len(config.LogFileSuffix) > 0 {
		suffix = strings.Trim(config.LogFileSuffix, ".")
	}

	n := time.Now().Local().Format("20060102")
	p := filepath.Join(dir, strings.Join([]string{n, suffix}, "."))
	writer, er := fs.FileWriter(p, true, 0755)
	if er != nil {
		err = el.PanicE(ctx, erlogs.Msgv("new logger writer fail"), erlogs.Content(er.Error()))
	}

	return
}

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "slient":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Error
	}
}

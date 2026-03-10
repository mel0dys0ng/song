package mysql

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/fs"
	"gorm.io/gorm/logger"
)

// newLogger return new gorm logger interface object
func newLogger(config *Config) (lgr logger.Interface, err error) {

	loggerWriter, err := newLoggerWriter(config)
	if err != nil {
		return
	}

	lgr = logger.New(
		log.New(loggerWriter, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             config.LogSlow * time.Millisecond, // Slow SQL threshold
			LogLevel:                  getLogLevel(config.LogLevel),      // Log level
			IgnoreRecordNotFoundError: *config.LogIgnoreRecordNotFound,   // Ignore ErrRecordNotFound error for logger
			Colorful:                  *config.LogColorful,               // Disable color
		},
	)

	return
}

// newLoggerWriter return new logger writer
func newLoggerWriter(config *Config) (writer io.Writer, err error) {

	// 如果是调试模式，使用标准输出
	if *config.Debug {
		writer = os.Stdout
		return
	}

	// 构建日志文件路径
	dir := ""
	if len(config.LogDir) > 0 {
		dir = strings.TrimRight(config.LogDir, "/")
	}

	suffix := ""
	if len(config.LogFileSuffix) > 0 {
		suffix = strings.Trim(config.LogFileSuffix, ".")
	}

	// 生成日志文件名（按日期）
	n := time.Now().Local().Format("20060102")
	p := filepath.Join(dir, strings.Join([]string{n, suffix}, "."))

	// 创建日志文件
	writer, err = fs.FileWriter(p, true, 0644)
	if err != nil {
		err = erlogs.Convert(err).Wrap("new logger writer fail").Panic()
		return
	}

	return
}

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
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

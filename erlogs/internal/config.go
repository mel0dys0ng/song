package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mel0dys0ng/song/metas"
	"github.com/mel0dys0ng/song/utils/systems"
)

const (
	DirDefault        = "./logs"
	FileNameDefault   = "default.log"
	LevelDefault      = LevelStringInfo
	MaxSizeDefault    = 500 // 500M
	MaxBackupsDefault = 7   // 最多备份7份
	MaxAgeDefault     = 30  // 30days
	CompressDefault   = true
)

type (
	Config struct {
		// 日志存储目录
		Dir string `yaml:"dir" json:"dir"`
		// 日志存储文件名
		FileName string `yaml:"fileName" json:"fileName"`
		// 日志记录最低级别
		Level string `yaml:"level" json:"level"`
		// 每个日志文件保存的最大尺寸 单位：M
		MaxSize int `yaml:"maxSize" json:"maxSize"`
		// 日志文件最多保存多少个备份
		MaxBackups int `yaml:"maxBackups" json:"maxBackups"`
		// 文件最多保存多少天
		MaxAge int `yaml:"maxAge" json:"maxAge"`
		// 是否压缩
		Compress bool `yaml:"compress" json:"compress"`
	}
)

func DefaultConfig() *Config {
	return &Config{
		Dir:        metas.LogDir(),
		FileName:   FileNameDefault,
		Level:      LevelDefault,
		MaxSize:    MaxSizeDefault,
		MaxBackups: MaxBackupsDefault,
		MaxAge:     MaxAgeDefault,
		Compress:   CompressDefault,
	}
}

func (c *Config) GetDir() string {
	if c != nil && len(c.Dir) > 0 {
		return c.Dir
	}
	return DirDefault
}

func (c *Config) GetFileName() string {
	if c != nil && len(c.FileName) > 0 {
		return c.FileName
	}
	return FileNameDefault
}

func (c *Config) GetFilePath() string {
	file := fmt.Sprintf("%s/%s.log", c.GetDir(), c.GetFileName())
	path, err := filepath.Abs(file)
	if err != nil {
		systems.Panic(err.Error())
	}
	return path
}

func (c *Config) GetLevel() string {
	if c != nil && len(c.Level) > 0 {
		return c.Level
	}
	return LevelDefault
}

func (c *Config) GetMaxSize() int {
	if c != nil && c.MaxSize > 0 {
		return c.MaxSize
	}
	return MaxSizeDefault
}

func (c *Config) GetMaxAge() int {
	if c != nil && c.MaxAge > 0 {
		return c.MaxAge
	}
	return MaxAgeDefault
}

func (c *Config) GetMaxBuckups() int {
	if c != nil && c.MaxBackups > 0 {
		return c.MaxBackups
	}
	return MaxBackupsDefault
}

func (c *Config) GetCompose() bool {
	if c == nil {
		return CompressDefault
	}
	return c.Compress
}

func (c *Config) Check() (err error) {
	if c == nil {
		return errors.New("config nil")
	}

	var errMsgs []string
	data := map[string]string{
		c.Level: "config.level empty",
	}

	for k, v := range data {
		if len(k) == 0 {
			errMsgs = append(errMsgs, v)
		}
	}

	maxs := map[int]string{
		c.MaxSize:    "config.maxSize <= 0",
		c.MaxBackups: "config.maxBackups <= 0",
		c.MaxAge:     "config.maxAge <= 0",
	}

	for k, v := range maxs {
		if k <= 0 {
			errMsgs = append(errMsgs, v)
		}
	}

	f, e := os.Stat(c.Dir)
	if e != nil {
		errMsgs = append(errMsgs, fmt.Sprintf("config.dir invalid: %s", e.Error()))
	}

	if f != nil && !f.IsDir() {
		errMsgs = append(errMsgs, "config.dir is not exist")
	}

	switch c.Level {
	case LevelStringDebug, LevelStringInfo, LevelStringWarn, LevelStringError, LevelStringPanic, LevelStringFatal:
	default:
		errMsgs = append(errMsgs, "config.level invalid: level should be debug|info|warn|error|panic|fatal")
	}

	if len(errMsgs) > 0 {
		err = errors.New(strings.Join(errMsgs, " and "))
	}

	return
}

package erlogs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mel0dys0ng/song/pkg/singleton"
	"github.com/mel0dys0ng/song/pkg/sys"
)

var globalConfigKey = singleton.Key()

const (
	DirDefault        = "./logs"
	FileExtDefault    = "log"
	FileNameDefault   = "default"
	LevelDefault      = LevelStringInfo
	MaxSizeDefault    = 500 // 500M
	MaxBackupsDefault = 7   // 最多备份7份
	MaxAgeDefault     = 30  // 30days
	CompressDefault   = true
)

type Config struct {
	// 日志存储目录
	Dir string `yaml:"dir" json:"dir" mapstructure:"dir"`
	// 日志存储文件扩展名
	FileExt string `yaml:"fileExt" json:"fileExt" mapstructure:"fileExt"`
	// 日志存储文件名
	FileName string `yaml:"fileName" json:"fileName" mapstructure:"fileName"`
	// 日志存储文件路径
	FilePath string `yaml:"filePath" json:"filePath" mapstructure:"filePath"`
	// 日志记录最低级别
	Level string `yaml:"level" json:"level" mapstructure:"level"`
	// 每个日志文件保存的最大尺寸 单位：M
	MaxSize int `yaml:"maxSize" json:"maxSize" mapstructure:"maxSize"`
	// 日志文件最多保存多少个备份
	MaxBackups int `yaml:"maxBackups" json:"maxBackups" mapstructure:"maxBackups"`
	// 文件最多保存多少天
	MaxAge int `yaml:"maxAge" json:"maxAge" mapstructure:"maxAge"`
	// 是否压缩
	Compress bool `yaml:"compress" json:"compress" mapstructure:"compress"`

	// 缓存字段
	dirChecked bool // 缓存目录检查结果
	dirValid   bool // 缓存目录是否有效
}

// GetConfig 获取当前配置
// 如果外部已初始化，返回初始化后的配置；否则返回 nil
func GetConfig() *Config {
	cfg, ok := singleton.Get[*Config](globalConfigKey)
	if ok && cfg != nil {
		if cfg.FileExt == "" {
			cfg.FileExt = FileExtDefault
		}

		if cfg.Dir == "" {
			absDir, err := filepath.Abs(DirDefault)
			if err != nil {
				absDir = DirDefault
			}
			cfg.Dir = absDir
		}

		if cfg.FileName == "" {
			cfg.FileName = FileNameDefault
		}

		if cfg.Level == "" {
			cfg.Level = LevelDefault
		}

		if cfg.MaxSize <= 0 {
			cfg.MaxSize = MaxSizeDefault
		}

		if cfg.MaxBackups <= 0 {
			cfg.MaxBackups = MaxBackupsDefault
		}

		if cfg.MaxAge <= 0 {
			cfg.MaxAge = MaxAgeDefault
		}
	}

	return cfg
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
	file := fmt.Sprintf("%s.%s", c.GetFileName(), c.GetFileExt())

	// 检查缓存是否有效
	if c.FilePath != "" && strings.HasSuffix(c.FilePath, file) {
		return c.FilePath
	}

	// 计算文件路径
	file = fmt.Sprintf("%s/%s", c.GetDir(), file)
	path, err := filepath.Abs(file)
	if err != nil {
		sys.Panic(err.Error())
	}

	// 缓存结果
	c.FilePath = path
	return path
}

func (c *Config) GetFileExt() string {
	if c != nil && len(c.FileExt) > 0 {
		return c.FileExt
	}
	return "log"
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

	// 检查必填字段
	if len(c.Level) == 0 {
		errMsgs = append(errMsgs, "config.level empty")
	}

	// 检查数值字段
	if c.MaxSize <= 0 {
		errMsgs = append(errMsgs, "config.maxSize <= 0")
	}

	if c.MaxBackups <= 0 {
		errMsgs = append(errMsgs, "config.maxBackups <= 0")
	}

	if c.MaxAge <= 0 {
		errMsgs = append(errMsgs, "config.maxAge <= 0")
	}

	// 检查目录是否存在（使用缓存）
	if !c.dirChecked {
		if f, e := os.Stat(c.GetDir()); e != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("config.dir invalid: %s", e.Error()))
			c.dirValid = false
		} else if !f.IsDir() {
			errMsgs = append(errMsgs, fmt.Sprintf("config.dir: [%s]. is not exist", c.GetDir()))
			c.dirValid = false
		} else {
			c.dirValid = true
		}
		c.dirChecked = true
	} else if !c.dirValid {
		errMsgs = append(errMsgs, fmt.Sprintf("config.dir: [%s]. invalid: previously checked and found invalid", c.GetDir()))
	}

	// 检查日志级别是否有效
	switch ToLevel(c.Level) {
	case LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal:
		// 有效的日志级别
	default:
		errMsgs = append(errMsgs,
			fmt.Sprintf(
				"config.level: [%s] invalid, level should be debug|info|warn|error|panic|fatal",
				c.GetLevel(),
			),
		)
	}

	if len(errMsgs) > 0 {
		err = errors.New(strings.Join(errMsgs, " and "))
	}

	return
}

package metas

import (
	"sync"

	"github.com/mel0dys0ng/song/pkgs/metas/internal"
	"github.com/mel0dys0ng/song/pkgs/utils/sys"
)

type (
	MetadataInterface = internal.MetadataInterface
	ModeType          = internal.ModeType
	Options           = internal.Options
)

var (
	mt   MetadataInterface
	once sync.Once
)

/*
Init 初始化元数据

参数：
- options
  - @param app string 应用名称（全局唯一）
  - @param product string 产品名称
  - @param config string 配置DSN，支持以下2种：
  - yaml://@./config/test or ./config/test
  - etcd://127.0.0.1:9091@config/test
*/
func Init(opts *Options) {
	once.Do(func() { mt = internal.New(opts) })
}

// Mt 返回全部metadata
func Mt() MetadataInterface {
	if mt == nil {
		sys.Panic("metadata is not initialized")
	}
	return mt
}

func Mode() ModeType {
	return Mt().Mode()
}

func Product() string {
	return Mt().Product()
}

func App() string {
	return Mt().App()
}

func Ip() string {
	return Mt().Ip()
}

func Node() string {
	return Mt().Node()
}

func Region() string {
	return Mt().Region()
}

func Zone() string {
	return Mt().Zone()
}

func Provider() string {
	return Mt().Provider()
}

// Envkey return the env key with prefix
func Envkey(name string) string {
	return Mt().Envkey(name)
}

// Getenv read and return the env value that name is name.
func Getenv(name, defaultV string) string {
	return Mt().Getenv(name, defaultV)
}

// Setenv set env
func Setenv(name, value string) error {
	return Mt().Setenv(name, value)
}

// Unsetenv unset env
func Unsetenv(name string) error {
	return Mt().Unsetenv(name)
}

// ConfigType return the type of the config
func ConfigType() string {
	return Mt().ConfigType()
}

// ConfigAddr return the addr of the remote config
func ConfigAddr() string {
	return Mt().ConfigAddr()
}

// ConfigPath return the dir path of the config
func ConfigPath() string {
	return Mt().ConfigPath()
}

// LogDir return the dir path of the log
func LogDir() string {
	return Mt().LogDir()
}

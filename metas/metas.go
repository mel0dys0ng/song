package metas

import (
	"sync"
	"time"

	"git.dreamsky.cn/song/metas/internal"
	"git.dreamsky.cn/song/utils/sljces"
	"git.dreamsky.cn/song/utils/systems"
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

// Init 初始化元数据
// @param app string 应用名称（全局唯一）
// @param product string 产品名称
// @param config string 配置DSN，支持以下2种：
// @example yaml://@./config/test or ./config/test
// @example etcd://127.0.0.1:9091@config/test
func Init(opts *Options) {
	once.Do(func() { mt = internal.New(opts) })
}

// Data 返回全部metadata
// 若metadata未初始化时ds不为空，则需要循环等待。
// 第一个时间是每次等待事件，默认为10ms；第二个时间总共等待最长时间，默认为3m。
// 最后还未初始化，则会Panic。
func Data(ds ...time.Duration) MetadataInterface {
	if len(ds) > 0 {
		sl := sljces.New(ds)
		sleep := sl.First(10 * time.Millisecond)
		maxSleep := sl.IndexOf(1, 3*time.Minute)

		var mu sync.Mutex
		startTime := time.Now()
		totalSleep := time.Duration(0)

		for {
			mu.Lock()
			currentMT := mt
			mu.Unlock()

			if currentMT != nil {
				return currentMT
			}

			if time.Since(startTime) > maxSleep {
				break
			}

			time.Sleep(sleep)
			totalSleep += sleep
		}
	}

	if mt == nil {
		systems.Panic("metadata is not initialized")
	}

	return mt
}

func Mode() ModeType {
	return Data().Mode()
}

func Product() string {
	return Data().Product()
}

func App() string {
	return Data().App()
}

func Ip() string {
	return Data().Ip()
}

func Node() string {
	return Data().Node()
}

func Region() string {
	return Data().Region()
}

func Zone() string {
	return Data().Zone()
}

func Provider() string {
	return Data().Provider()
}

// Envkey return the env key with prefix
func Envkey(name string) string {
	return Data().Envkey(name)
}

// Getenv read and return the env value that name is name.
func Getenv(name, defaultV string) string {
	return Data().Getenv(name, defaultV)
}

// Setenv set env
func Setenv(name, value string) error {
	return Data().Setenv(name, value)
}

// Unsetenv unset env
func Unsetenv(name string) error {
	return Data().Unsetenv(name)
}

// ConfigType return the type of the config
func ConfigType() string {
	return Data().ConfigType()
}

// ConfigAddr return the addr of the remote config
func ConfigAddr() string {
	return Data().ConfigAddr()
}

// ConfigPath return the dir path of the config
func ConfigPath() string {
	return Data().ConfigPath()
}

// LogDir return the dir path of the log
func LogDir() string {
	return Data().LogDir()
}

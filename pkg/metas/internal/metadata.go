package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mel0dys0ng/song/pkg/utils/aob"
	"github.com/mel0dys0ng/song/pkg/utils/fs"
	"github.com/mel0dys0ng/song/pkg/utils/ip"
	"github.com/mel0dys0ng/song/pkg/utils/sys"
)

const (
	appOrProductRegexpPattern = `^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
	ConfigDSNDeault           = "yaml://@./config/debug"
	ConfigDSNRegexpPattern    = `^(yaml|etcd)://([a-zA-Z.:0-9]*)@([a-zA-Z0-9/._-]+)$`
	ConfigPathRegexpPattern   = "^[a-zA-Z0-9.-_]+/(debug|test|pre|gray|production)/?$"
	ConfigDirDefault          = "./config/debug"
	ConfigTypeYaml            = "yaml"
	ConfigTypeEtcd            = "etcd"
	LogDirDefault             = "./logs"
)

type (
	// metadata 应用元数据
	metadata struct {
		// 产品名称
		product string
		// 应用名称（全局名称）
		app string
		// 运行模式
		mode ModeType
		// 节点ID
		node string
		// 节点Region
		region string
		// 节点Zone
		zone string
		// 服务提供方
		provider string
		// 节点IP
		ip string
		// 环境变量前缀
		envKeyPrefix string
		// 配置类型
		configType string
		// 配置addr
		configAddr string
		// 配置模式
		configMode ModeType
		// 配置目录
		configPath string
		// 日志目录
		logDir string
	}

	Options struct {
		App     string `json:"app"`
		Product string `json:"product"`
		Config  string `json:"config"`
	}
)

// New 创建并返回元数据对象
func New(opts *Options) (res MetadataInterface) {
	if opts == nil {
		sys.Panicf("init metadata opts should not be nil")
		return
	}

	for k, v := range map[string]string{"app": opts.App, "product": opts.Product} {
		if len(v) == 0 || !regexp.MustCompile(appOrProductRegexpPattern).MatchString(v) {
			f := "failed to new metadata: the %s `%s` is invalid. the %s must match the regexp pattern: `%s`"
			sys.Panicf(f, k, v, k, appOrProductRegexpPattern)
			return
		}
	}

	IP, _ := ip.GetLocalIP()

	mt := &metadata{
		app:          opts.App,
		product:      opts.Product,
		envKeyPrefix: fmt.Sprintf("%s_", strings.ToUpper(opts.App)),
		configMode:   ModeDebug,
		configType:   ConfigTypeYaml,
		configPath:   fs.Abs(ConfigDirDefault),
		ip:           IP,
	}

	mt.parseConfig(opts.Config)
	mt.mode = ModeType(mt.Getenv("SONG_MODE", ModeDebug.String()))
	if !mt.mode.Validate() {
		sys.Panicf("failed to new metadata: mode invalid. mode = %s", mt.mode.String())
		return
	}

	// 校验配置是否可以运行
	if mt.mode.IsModeTestOrPreOrDebug() && !mt.configMode.IsModeTestOrPreOrDebug() {
		sys.Panicf("线下环境不能运行线上配置: EnvMode: %s, ConfigMode: %s", mt.mode, mt.configMode)
		return
	} else if mt.mode.IsModeGray() && !mt.configMode.IsModeGray() {
		sys.Panicf("灰度环境仅能运行灰度配置: EnvMode: %s, ConfigMode: %s", mt.mode, mt.configMode)
		return
	} else if mt.mode.IsModeProduction() && !mt.configMode.IsModeProduction() {
		sys.Panicf("生产环境仅能运行生产配置: EnvMode: %s, ConfigMode: %s", mt.mode, mt.configMode)
		return
	}

	logDir, err := filepath.Abs(mt.Getenv("SONG_LOG_DIR", LogDirDefault))
	if err != nil {
		sys.Panicf("failed to new metadata: getenv log_dir error(%s)", err.Error())
		return
	}

	mt.logDir = filepath.Join(logDir, mt.app)
	mt.node = mt.Getenv("SONG_NODE", "")
	mt.region = mt.Getenv("SONG_REGION", "")
	mt.zone = mt.Getenv("SONG_ZONE", "")
	mt.provider = mt.Getenv("SONG_PROVIDER", "")

	return mt
}

func (m *metadata) parseConfig(config string) {
	if len(config) == 0 {
		config = ConfigDSNDeault
	}

	subs := regexp.MustCompile(ConfigDSNRegexpPattern).FindStringSubmatch(config)
	if len(subs) == 4 {
		m.configType = subs[1]
		m.configAddr = subs[2]
		m.configPath = fs.Abs(subs[3])
	} else {
		m.configType = ConfigTypeYaml
		m.configPath = fs.Abs(config)
	}

	if info, err := os.Stat(m.configPath); err != nil || !info.IsDir() {
		sys.Panicf("failed to new metadata: config path = %s invalid. "+
			"the config path is not a dir path or not exist",
			m.configPath,
		)
	}

	matches := regexp.MustCompile(ConfigPathRegexpPattern).FindStringSubmatch(m.configPath)
	if len(matches) == 2 {
		m.configMode = ModeType(matches[1])
	}

	if !m.configMode.Validate() {
		sys.Panicf("failed to new metadata: config mode = %s invalid. "+
			"the value range of the config mode and env mode should be the same. ",
			m.configMode,
		)
	}
}

func (m *metadata) Envkey(name string) string {
	return strings.ToUpper(fmt.Sprintf("%s%s", m.envKeyPrefix, name))
}

func (m *metadata) Getenv(name, defaultV string) string {
	value, ok := os.LookupEnv(name)
	return aob.Aorb(ok, value, defaultV)
}

func (m *metadata) Setenv(name, value string) error {
	return os.Setenv(name, value)
}

func (m *metadata) Unsetenv(name string) error {
	return os.Unsetenv(name)
}

func (m *metadata) Mode() ModeType {
	return m.mode
}

func (m *metadata) Product() string {
	return m.product
}

func (m *metadata) App() string {
	return m.app
}

func (m *metadata) Ip() string {
	return m.ip
}

func (m *metadata) Node() string {
	return m.node
}

func (m *metadata) Region() string {
	return m.region
}

func (m *metadata) Zone() string {
	return m.zone
}

func (m *metadata) Provider() string {
	return m.provider
}

func (m *metadata) ConfigType() string {
	return m.configType
}

func (m *metadata) ConfigAddr() string {
	return m.configAddr
}

func (m *metadata) ConfigPath() string {
	return m.configPath
}

func (m *metadata) LogDir() string {
	return m.logDir
}

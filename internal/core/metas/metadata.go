package metas

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mel0dys0ng/song/pkg/aob"
	"github.com/mel0dys0ng/song/pkg/fs"
	"github.com/mel0dys0ng/song/pkg/ip"
	"github.com/mel0dys0ng/song/pkg/sys"
)

const (
	FlagAppRegexpPattern    = `^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
	ConfigDSNDefault        = "yaml://@./configs/local"
	ConfigDSNRegexpPattern  = `^(yaml|etcd)://([a-zA-Z.:0-9]*)@([a-zA-Z0-9/._-]+)$`
	ConfigPathRegexpPattern = `^[a-zA-Z0-9.-_]+/(local|test|staging|prod)[/]?$`
	ConfigDirDefault        = "./configs/local"
	ConfigTypeYaml          = "yaml"
	ConfigTypeEtcd          = "etcd"
	LogDirDefault           = "./logs"

	EnvNameMode     = "SONG_MODE"
	EnvNameNode     = "SONG_NODE"
	EnvNameRegion   = "SONG_REGION"
	EnvNameZone     = "SONG_ZONE"
	EnvNameProvider = "SONG_PROVIDER"
	EnvNameLogDir   = "SONG_LOG_DIR"
)

type (
	// metadata 应用元数据
	metadata struct {
		// 应用名称
		app string
		// 应用类型
		kind KindType
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
		// 配置路径
		configPath string
		// 日志目录
		logDir string
	}

	Options struct {
		App    string   `json:"app"`    // 应用名称
		Kind   KindType `json:"kind"`   // 应用类型，如：api、job、tool、messaging
		Mode   ModeType `json:"mode"`   // 运行模式，默认local, 如：local、test、staging、prod
		Config string   `json:"config"` // 配置地址
	}
)

// New 创建并返回元数据对象
func New(opts *Options) (res MetadataInterface) {
	if opts == nil {
		sys.Panicf("init metadata opts should not be nil")
		return
	}

	if len(opts.App) == 0 || !regexp.MustCompile(FlagAppRegexpPattern).MatchString(opts.App) {
		sys.Panicf("failed to new metadata: app invalid. app = %s", opts.App)
		return
	}

	if !opts.Kind.Validate() {
		sys.Panicf("failed to new metadata: kind invalid. kind = %s", opts.Kind)
		return
	}

	if !opts.Mode.Validate() {
		sys.Panicf("failed to new metadata: mode invalid. mode = %s", opts.Mode)
		return
	}

	// 获取本地ip
	ip, _ := ip.GetLocalIp()

	configPath, err := fs.Abs(ConfigDirDefault)
	if err != nil {
		sys.Panicf("failed to new metadata: get config path error(%s)", err.Error())
		return
	}

	mt := &metadata{
		app:          opts.App,
		kind:         opts.Kind,
		mode:         opts.Mode,
		envKeyPrefix: fmt.Sprintf("%s_", strings.ToUpper(opts.App)),
		configType:   ConfigTypeYaml,
		configPath:   configPath,
		ip:           ip,
	}

	mt.parseConfigPath(opts.Config)

	logDir, err := filepath.Abs(mt.Getenv(EnvNameLogDir, LogDirDefault))
	if err != nil {
		sys.Panicf("failed to new metadata: getenv log_dir error(%s)", err.Error())
		return
	}

	mt.logDir = filepath.Join(logDir, mt.app)
	mt.node = mt.Getenv(EnvNameNode, "")
	mt.region = mt.Getenv(EnvNameRegion, "")
	mt.zone = mt.Getenv(EnvNameZone, "")
	mt.provider = mt.Getenv(EnvNameProvider, "")

	return mt
}

func (m *metadata) parseConfigPath(config string) {
	if len(config) == 0 {
		config = ConfigDSNDefault
	}

	subs := regexp.MustCompile(ConfigDSNRegexpPattern).FindStringSubmatch(config)
	if len(subs) == 4 {
		m.configType = subs[1]
		m.configAddr = subs[2]
		configPath, err := fs.Abs(subs[3])
		if err != nil {
			sys.Panicf("failed to new metadata: get config path error(%s)", err.Error())
			return
		}
		m.configPath = configPath
	} else {
		m.configType = ConfigTypeYaml
		configPath, err := fs.Abs(config)
		if err != nil {
			sys.Panicf("failed to new metadata: get config path error(%s)", err.Error())
			return
		}
		m.configPath = configPath
	}

	if info, err := os.Stat(m.configPath); err != nil || !info.IsDir() {
		sys.Panicf("failed to new metadata: config path = %s invalid. "+
			"the config path is not a dir path or not exist",
			m.configPath,
		)
		return
	}
}

func (m *metadata) Envkey(name string) string {
	if m == nil {
		return ""
	}
	return strings.ToUpper(fmt.Sprintf("%s%s", m.envKeyPrefix, name))
}

func (m *metadata) Getenv(name, defaultV string) string {
	if m == nil {
		return defaultV
	}
	value, ok := os.LookupEnv(name)
	return aob.VarOrVar(ok, value, defaultV)
}

func (m *metadata) Setenv(name, value string) error {
	if m == nil {
		return fmt.Errorf("metadata is nil")
	}
	return os.Setenv(name, value)
}

func (m *metadata) Unsetenv(name string) error {
	if m == nil {
		return fmt.Errorf("metadata is nil")
	}
	return os.Unsetenv(name)
}

func (m *metadata) Mode() ModeType {
	if m == nil {
		return ""
	}
	return m.mode
}

func (m *metadata) Kind() KindType {
	if m == nil {
		return ""
	}
	return m.kind
}

func (m *metadata) App() string {
	if m == nil {
		return ""
	}
	return m.app
}

func (m *metadata) Ip() string {
	if m == nil {
		return ""
	}
	return m.ip
}

func (m *metadata) Node() string {
	if m == nil {
		return ""
	}
	return m.node
}

func (m *metadata) Region() string {
	if m == nil {
		return ""
	}
	return m.region
}

func (m *metadata) Zone() string {
	if m == nil {
		return ""
	}
	return m.zone
}

func (m *metadata) Provider() string {
	if m == nil {
		return ""
	}
	return m.provider
}

func (m *metadata) ConfigType() string {
	if m == nil {
		return ""
	}
	return m.configType
}

func (m *metadata) ConfigAddr() string {
	if m == nil {
		return ""
	}
	return m.configAddr
}

func (m *metadata) ConfigPath() string {
	if m == nil {
		return ""
	}
	return m.configPath
}

func (m *metadata) LogDir() string {
	if m == nil {
		return ""
	}
	return m.logDir
}

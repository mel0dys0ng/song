package internal

type MetadataInterface interface {
	// Envkey return the env key with prefix
	Envkey(name string) string

	// Getenv read and return the env value that name is name
	Getenv(name, defaultV string) string

	// Setenv set env
	Setenv(name, value string) error

	// Unsetenv remove env
	Unsetenv(name string) error

	// Mode mode
	Mode() ModeType

	// Product 产品名称
	Product() string

	// App 应用名称
	App() string

	// Ip server ip
	Ip() string

	// Node server node
	Node() string

	// Region server region
	Region() string

	// Zone server zone
	Zone() string

	// Provider server provider
	Provider() string

	// ConfigType 配置类型
	ConfigType() string

	// ConfigAddr 配置Addr
	ConfigAddr() string

	// ConfigPath 配置路径
	ConfigPath() string

	// LogDir 日志路径
	LogDir() string
}

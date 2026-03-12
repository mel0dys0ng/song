# Vipers 配置管理

一个强大的配置管理系统，基于 Viper 构建，支持多种配置源（YAML、JSON、TOML、etcd、Consul）、实时配置更新和带有默认值的安全配置访问。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建配置](#创建配置)
  - [读取配置值](#读取配置值)
  - [配置提供者](#配置提供者)
  - [配置更新](#配置更新)
- [配置选项](#配置-options)
- [安全访问方法](#安全访问方法)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **多种配置源**：支持 YAML、JSON、TOML、etcd 和 Consul
- **实时更新**：配置变更时自动重新加载
- **安全访问**：带有默认值的安全访问方法
- **配置验证**：内置验证和错误处理
- **层级配置**：支持嵌套配置结构
- **环境变量集成**：自动绑定环境变量
- **远程配置**：支持分布式配置管理
- **配置变更钩子**：配置变更回调函数
- **默认值**：优雅地回退到默认值
- **统一接口**：跨所有配置源的一致 API

## 安装

vipers 包需要以下依赖：

```bash
go get github.com/spf13/viper
go get github.com/fsnotify/fsnotify
go get go.etcd.io/etcd/client/v3
go get github.com/hashicorp/consul/api
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/vipers"
)

func main() {
    // 从 YAML 文件创建配置
    config, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/app.yaml"),
    )
    if err != nil {
        panic(err)
    }
    
    // 读取配置值
    port := config.GetInt("server.port", 8080)
    host := config.GetString("server.host", "localhost")
    debug := config.GetBool("server.debug", false)
    
    fmt.Printf("服务器: %s:%d (调试: %v)\n", host, port, debug)
}
```

## 架构概览

vipers 包提供了一个统一的配置管理系统：

```
┌─────────────────┐
│    Config       │ - 主配置结构
└────────┬────────┘
         │
    ┌────┴────┐
    │ Viper   │ - 底层 Viper 实例
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Provider       │ - 配置源提供者
    └─────────────────┘
         │
    ┌────┴────┬────────┬────────┬──────────┐
    │  YAML   │  JSON  │  TOML  │  Etcd    │ Consul │
    └─────────┴────────┴────────┴──────────┴────────┘
```

**核心组件：**
- **Config**：带有 Viper 集成的 主配置结构
- **Provider**：配置源的抽象接口
- **Options**：配置选项和设置
- **ProviderInterface**：不同配置源的接口

## 使用指南

### 创建配置

创建一个新的配置实例：

```go
import (
    "github.com/mel0dys0ng/song/internal/core/vipers"
)

// YAML 配置
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderYaml),
    vipers.OnPath("./configs/app.yaml"),
)

// JSON 配置
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderJson),
    vipers.OnPath("./configs/app.json"),
)

// TOML 配置
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderToml),
    vipers.OnPath("./configs/app.toml"),
)

// Etcd 配置
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderEtcd),
    vipers.OnEndpoints("localhost:2379"),
    vipers.OnPath("config/myapp/prod"),
)

// Consul 配置
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderConsul),
    vipers.OnEndpoints("localhost:8500"),
    vipers.OnPath("config/myapp/prod"),
)
```

### 读取配置值

使用安全访问方法读取配置值：

```go
// 字符串值
host := config.GetString("server.host", "localhost")
name := config.GetString("app.name", "MyApp")

// 整数值
port := config.GetInt("server.port", 8080)
maxConn := config.GetInt("database.max_connections", 100)

// 布尔值
debug := config.GetBool("server.debug", false)
enabled := config.GetBool("features.new_ui", true)

// 浮点值
timeout := config.GetFloat64("server.timeout", 30.0)
ratio := config.GetFloat64("sampling.rate", 0.1)

// 持续时间值
readTimeout := config.GetDuration("server.read_timeout", 30*time.Second)
retryDelay := config.GetDuration("retry.delay", 5*time.Second)

// 切片值
hosts := config.GetStringSlice("server.hosts", []string{"localhost"})
ports := config.GetIntSlice("server.ports", []int{8080, 8081})

// Map 值
labels := config.GetStringMapString("labels", map[string]string{})
metadata := config.GetStringMap("metadata", map[string]any{})
```

**所有 Getter 方法：**
- `Get(key, defaultValue)` - 通用获取器
- `GetString(key, defaultValue)` - 字符串值
- `GetInt(key, defaultValue)` - 整数值
- `GetInt32(key, defaultValue)` - 32 位整数
- `GetInt64(key, defaultValue)` - 64 位整数
- `GetUint(key, defaultValue)` - 无符号整数
- `GetUint32(key, defaultValue)` - 32 位无符号整数
- `GetUint64(key, defaultValue)` - 64 位无符号整数
- `GetBool(key, defaultValue)` - 布尔值
- `GetFloat64(key, defaultValue)` - Float64 值
- `GetDuration(key, defaultValue)` - 持续时间值
- `GetTime(key, defaultValue)` - 时间值
- `GetStringSlice(key, defaultValue)` - 字符串切片
- `GetIntSlice(key, defaultValue)` - 整数切片
- `GetStringMap(key, defaultValue)` - 通用 Map
- `GetStringMapString(key, defaultValue)` - 字符串 Map
- `GetStringMapStringSlice(key, defaultValue)` - 字符串切片 Map
- `GetSizeInBytes(key, defaultValue)` - 字节大小

### 配置提供者

#### YAML 提供者

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderYaml),
    vipers.OnPath("./configs/app.yaml"),
)
```

**示例 YAML 文件：**
```yaml
server:
  host: localhost
  port: 8080
  debug: true
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: myapp
  max_connections: 100
  
features:
  new_ui: true
  beta_features:
    - feature1
    - feature2
```

#### JSON 提供者

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderJson),
    vipers.OnPath("./configs/app.json"),
)
```

**示例 JSON 文件：**
```json
{
  "server": {
    "host": "localhost",
    "port": 8080,
    "debug": true,
    "timeout": "30s"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "myapp",
    "max_connections": 100
  }
}
```

#### TOML 提供者

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderToml),
    vipers.OnPath("./configs/app.toml"),
)
```

**示例 TOML 文件：**
```toml
[server]
host = "localhost"
port = 8080
debug = true
timeout = "30s"

[database]
host = "localhost"
port = 5432
name = "myapp"
max_connections = 100
```

#### Etcd 提供者

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderEtcd),
    vipers.OnEndpoints("localhost:2379"),
    vipers.OnPath("config/myapp/prod"),
    vipers.OnUsername("etcd_user"),
    vipers.OnPassword("etcd_password"),
)
```

#### Consul 提供者

```go
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderConsul),
    vipers.OnEndpoints("localhost:8500"),
    vipers.OnPath("config/myapp/prod"),
    vipers.OnToken("consul_token"),
)
```

### 配置更新

监听配置变更：

```go
// 设置配置变更处理器
config.OnConfigChange(func(event fsnotify.Event, options *vipers.Options) {
    fmt.Println("配置已更改!")
    fmt.Printf("事件: %v\n", event)
    
    // 重新加载配置或采取行动
    newPort := config.GetInt("server.port", 8080)
    log.Printf("新端口: %d", newPort)
})
```

## 配置 Options

### Option 函数

```go
// 提供者配置
vipers.OnProvider(vipers.ConfigProviderYaml)  // 设置配置提供者
vipers.OnPath("./configs/app.yaml")           // 设置配置路径
vipers.OnEndpoints("localhost:2379")          // 设置远程端点

// 认证
vipers.OnUsername("user")                     // 设置用户名
vipers.OnPassword("password")                 // 设置密码
vipers.OnToken("token")                       // 设置认证令牌

// 变更监控
vipers.OnChangeConfig(func(event, options) {  // 设置变更回调
    // 处理配置变更
})
```

### 提供者常量

```go
const (
    ConfigProviderJson   = "json"   // JSON 提供者
    ConfigProviderYaml   = "yaml"   // YAML 提供者
    ConfigProviderToml   = "toml"   // TOML 提供者
    ConfigProviderEtcd   = "etcd"   // Etcd 提供者
    ConfigProviderConsul = "consul" // Consul 提供者
)
```

### 默认选项

```go
opts := vipers.DefaultOptions()
// 返回:
// - Provider: ConfigProviderYaml
// - Path: "./configs/local"
// - Endpoints: ""
// - Username: ""
// - Password: ""
// - Token: ""
// - OnChangeConfig: nil
```

## 安全访问方法

所有 getter 方法遵循相同的模式：

```go
func (c *Config) Get<Type>(key string, defaultValue <Type>) <Type>
```

**行为：**
- 如果键存在且可以转换为目标类型，返回值
- 如果键不存在或转换失败，返回默认值
- 从不 panic - 始终返回有效值

**示例：**
```go
// 如果键不存在则返回 8080
port := config.GetInt("server.port", 8080)

// 如果键不存在则返回 "localhost"
host := config.GetString("server.host", "localhost")

// 如果键不存在则返回 false
debug := config.GetBool("server.debug", false)
```

## 示例代码

### 完整示例：应用程序配置

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/mel0dys0ng/song/internal/core/vipers"
)

type ServerConfig struct {
    Host         string        `mapstructure:"host"`
    Port         int           `mapstructure:"port"`
    Debug        bool          `mapstructure:"debug"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host           string `mapstructure:"host"`
    Port           int    `mapstructure:"port"`
    Name           string `mapstructure:"name"`
    User           string `mapstructure:"user"`
    Password       string `mapstructure:"password"`
    MaxConnections int    `mapstructure:"max_connections"`
}

type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

func main() {
    // 加载配置
    config, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/app.yaml"),
    )
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 读取单个值
    host := config.GetString("server.host", "localhost")
    port := config.GetInt("server.port", 8080)
    debug := config.GetBool("server.debug", false)
    
    log.Printf("服务器: %s:%d (调试: %v)", host, port, debug)
    
    // 解码到结构体
    var appConfig AppConfig
    if err := config.Unmarshal(&appConfig); err != nil {
        log.Fatalf("解码配置失败: %v", err)
    }
    
    log.Printf("服务器配置: %+v", appConfig.Server)
    log.Printf("数据库配置: %+v", appConfig.Database)
    
    // 设置配置变更监控
    config.OnConfigChange(func(event fsnotify.Event, options *vipers.Options) {
        log.Println("配置文件已更改，重新加载...")
        
        // 重新读取配置
        newPort := config.GetInt("server.port", 8080)
        log.Printf("新端口: %d", newPort)
    })
    
    // 保持应用程序运行以接收配置变更
    select {}
}
```

### 示例：环境特定配置

```go
func loadEnvironmentConfig(env string) (*vipers.Config, error) {
    var path string
    var provider string
    
    switch env {
    case "local":
        path = "./configs/local.yaml"
        provider = vipers.ConfigProviderYaml
        
    case "test":
        path = "./configs/test.yaml"
        provider = vipers.ConfigProviderYaml
        
    case "staging":
        path = "config/app/staging"
        provider = vipers.ConfigProviderEtcd
        
    case "prod":
        path = "config/app/prod"
        provider = vipers.ConfigProviderEtcd
        
    default:
        return nil, fmt.Errorf("未知环境: %s", env)
    }
    
    // 构建配置选项
    opts := []vipers.Option{
        vipers.OnProvider(provider),
        vipers.OnPath(path),
    }
    
    // 为远程环境添加 etcd 端点
    if provider == vipers.ConfigProviderEtcd {
        opts = append(opts, 
            vipers.OnEndpoints("etcd.cluster:2379"),
            vipers.OnUsername("etcd_user"),
            vipers.OnPassword("etcd_password"),
        )
    }
    
    return vipers.New(opts...)
}
```

### 示例：功能开关

```go
type FeatureFlags struct {
    NewUI         bool     `mapstructure:"new_ui"`
    BetaFeatures  []string `mapstructure:"beta_features"`
    EnabledRegions []string `mapstructure:"enabled_regions"`
    RolloutRate   float64  `mapstructure:"rollout_rate"`
}

func loadFeatureFlags(config *vipers.Config) (*FeatureFlags, error) {
    var flags FeatureFlags
    
    // 使用安全访问器和默认值
    flags.NewUI = config.GetBool("features.new_ui", false)
    flags.BetaFeatures = config.GetStringSlice("features.beta_features", []string{})
    flags.EnabledRegions = config.GetStringSlice("features.enabled_regions", []string{})
    flags.RolloutRate = config.GetFloat64("features.rollout_rate", 0.0)
    
    // 或解码整个部分
    if err := config.UnmarshalKey("features", &flags); err != nil {
        return nil, err
    }
    
    return &flags, nil
}

func isEnabled(feature string, flags *FeatureFlags) bool {
    switch feature {
    case "new_ui":
        return flags.NewUI
    case "beta_feature1", "beta_feature2":
        return contains(flags.BetaFeatures, feature)
    default:
        return false
    }
}

func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
```

### 示例：数据库配置

```go
func loadDatabaseConfig(config *vipers.Config) (*DatabaseConfig, error) {
    var dbConfig DatabaseConfig
    
    // 使用默认值读取配置
    dbConfig.Host = config.GetString("database.host", "localhost")
    dbConfig.Port = config.GetInt("database.port", 5432)
    dbConfig.Name = config.GetString("database.name", "myapp")
    dbConfig.User = config.GetString("database.user", "postgres")
    dbConfig.Password = config.GetString("database.password", "")
    dbConfig.MaxConnections = config.GetInt("database.max_connections", 100)
    
    // 验证必填字段
    if dbConfig.Host == "" {
        return nil, fmt.Errorf("数据库主机是必填项")
    }
    if dbConfig.Name == "" {
        return nil, fmt.Errorf("数据库名称是必填项")
    }
    
    return &dbConfig, nil
}

func createDatabaseConnection(dbConfig *DatabaseConfig) (*sql.DB, error) {
    // 构建连接字符串
    connStr := fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
        dbConfig.Host,
        dbConfig.Port,
        dbConfig.Name,
        dbConfig.User,
        dbConfig.Password,
    )
    
    // 创建连接池
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // 配置连接池
    db.SetMaxOpenConns(dbConfig.MaxConnections)
    db.SetMaxIdleConns(dbConfig.MaxConnections / 2)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return db, nil
}
```

### 示例：多源配置

```go
func loadMultiSourceConfig() (*vipers.Config, error) {
    // 从文件加载基础配置
    baseConfig, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/base.yaml"),
    )
    if err != nil {
        return nil, err
    }
    
    // 使用环境特定配置覆盖
    env := os.Getenv("APP_ENV")
    if env != "" {
        envConfig, err := vipers.New(
            vipers.OnProvider(vipers.ConfigProviderYaml),
            vipers.OnPath(fmt.Sprintf("./configs/%s.yaml", env)),
        )
        if err != nil {
            return nil, err
        }
        
        // 合并配置（envConfig 优先）
        mergeConfigs(baseConfig, envConfig)
    }
    
    // 使用环境变量覆盖
    bindEnvironmentVariables(baseConfig)
    
    return baseConfig, nil
}

func mergeConfigs(base, override *vipers.Config) {
    // 从 override 配置获取所有键
    for _, key := range override.AllKeys() {
        value := override.Get(key, nil)
        if value != nil {
            // 在基础配置中设置（override 优先）
            base.Viper.Set(key, value)
        }
    }
}

func bindEnvironmentVariables(config *vipers.Config) {
    // 自动绑定环境变量
    config.Viper.AutomaticEnv()
    
    // 示例: SONG_SERVER_PORT 覆盖 server.port
    config.Viper.SetEnvPrefix("SONG")
    config.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
```

### 示例：配置验证

```go
type ValidatableConfig interface {
    Validate() error
}

func loadAndValidateConfig(config *vipers.Config) (*AppConfig, error) {
    var appConfig AppConfig
    
    // 解码配置
    if err := config.Unmarshal(&appConfig); err != nil {
        return nil, err
    }
    
    // 验证配置
    if err := validateAppConfig(&appConfig); err != nil {
        return nil, err
    }
    
    return &appConfig, nil
}

func validateAppConfig(config *AppConfig) error {
    // 验证服务器配置
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("无效的服务器端口: %d", config.Server.Port)
    }
    
    // 验证数据库配置
    if config.Database.MaxConnections < 1 {
        return fmt.Errorf("数据库最大连接数必须为正数")
    }
    
    if config.Database.MaxConnections > 1000 {
        return fmt.Errorf("数据库最大连接数不能超过 1000")
    }
    
    return nil
}
```

## 最佳实践

1. **使用安全访问器**：始终使用带有默认值的安全访问器方法以避免 panic。

2. **验证配置**：在使用配置之前验证关键配置值。

3. **使用结构体标签**：在解码到结构体时使用 `mapstructure` 标签以获得更好的类型安全。

4. **分离环境**：为不同环境使用不同的配置文件/源。

5. **监控变更**：设置配置变更处理器以进行动态配置更新。

6. **保护敏感数据**：永远不要在纯文本配置文件中存储敏感数据（密码、令牌）。

7. **使用环境变量**：使用环境变量进行环境特定的覆盖。

8. **记录配置**：记录所有配置选项、类型和默认值。

9. **提供合理的默认值**：始终为可选配置提供合理的默认值。

10. **测试配置**：在测试套件中测试配置加载和验证。

11. **快速失败**：如果缺少或无效的配置则提前失败。

12. **使用层级配置**：组织配置为逻辑层级（例如，`server.port`、`database.host`）。

## 相关文档

- [Song 框架文档](../../README.md)
- [Metas 元数据模块](../metas/README.md)
- [Viper 文档](https://github.com/spf13/viper)

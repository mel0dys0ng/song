# vipers 使用指南

`vipers` 是一个基于 `spf13/viper` 的配置管理包，提供了统一的配置读取接口，并支持多种配置源（本地文件、etcd、consul 等）。它简化了配置加载过程并增强了错误处理与日志记录功能。

## 📁 项目结构概览

```
.
├── internal/
│   ├── config.go                  // 主要的配置结构体及其实现
│   ├── config_interface.go        // 配置接口定义
│   ├── localprovider_interface.go // 本地配置提供者接口
│   ├── options.go                 // 配置选项
│   ├── provider.go                // 提供者的抽象工厂
│   ├── provider_local.go          // 本地配置实现
│   ├── provider_remote.go         // 远程配置实现
│   └── remoteprovider_interface.go// 远程配置提供者接口
├── README.md                      // 包说明文档
└── vipers.go                      // 对外暴露的便捷API入口
```

## 🔧 工作原理概述

### 1. 初始化流程

- 在应用启动时调用 vipers.Config() 获取全局唯一配置实例。
- 内部通过单例模式 (singleton.GetInstance) 实现懒加载和线程安全。
- 自动根据元数据（metas）决定使用的配置类型、地址和路径。

### 2. 支持的配置来源

- **本地文件**：支持 JSON/YAML/TOML 文件格式。
- **远程服务**：支持 etcd 和 Consul。

### 3. 功能特性

- 支持多个配置文件合并加载
- 支持热更新监听（watch）
- 统一的默认值获取机制
- 完善的日志跟踪与错误处理

## 📘 使用方法详解

### 1. 基本使用方式

```go
import "github.com/mel0dys0ng/song/core/vipers"

// 判断某个键是否存在
if vipers.IsSet("server.port") {
    port := vipers.GetInt("server.port", 8080)
}

// 获取字符串类型的配置项，默认为"default_value"
host := vipers.GetString("server.host", "localhost")

// 解析嵌套结构体
type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
var server ServerConfig
_ = vipers.UnmarshalKey("server", &server)
```

### 2. 监听配置变化

```go
vipers.Config().OnConfigChange(func(event fsnotify.Event, options *internal.Options) {
    log.Println("配置发生变化:", event.Name)
})
```

> 注意：此回调仅在启用 watch 的情况下生效（目前所有 provider 都已内置 watch）

### 3. 多环境配置支持

假设你有以下目录结构：

```
configs/
├── app.yaml
├── dev.yaml
└── prod.yaml
```

可以通过设置不同的 metas.ConfigPath() 来切换不同环境的配置文件。

### 4. 结构化反序列化

```go
type AppConfig struct {
    Database struct {
        DSN string `mapstructure:"dsn"`
    } `mapstructure:"database"`

    Server struct {
        Addr string `mapstructure:"addr"`
    } `mapstructure:"server"`
}

var cfg AppConfig
err := vipers.Unmarshal(&cfg)
if err != nil {
    panic(err)
}
fmt.Printf("%+v\n", cfg)
```

## 📌 最佳实践建议

1. 将配置文件放在固定位置，并使用 `metas` 模块控制运行环境。
2. 推荐使用结构体映射而不是手动一个个读取字段。
3. 合理利用 `default value` 参数防止程序因缺失配置而崩溃。
4. 开启 watch 并注册回调函数以便动态调整行为。
5. 错误应尽早捕获并打印详细日志便于调试定位问题。

## ✅ 示例代码总结

```go
// 获取配置值
name := vipers.GetString("app.name", "MyApp")
debug := vipers.GetBool("debug", false)

// 映射到结构体
type DBConfig struct {
    URL string `mapstructure:"url"`
}
var dbConf DBConfig
_ = vipers.UnmarshalKey("database", &dbConf)

// 注册变更监听器
vipers.Config().OnConfigChange(func(e fsnotify.Event, opt *internal.Options) {
    log.Printf("Config file changed: %s", e.Name)
})
```

希望这份使用指南能帮助您快速掌握 `vipers` 库的核心功能！如果您有任何疑问或者发现 bug，请及时联系维护人员。

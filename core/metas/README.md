根据提供的代码内容，我将为您梳理 `metas` 项目的逻辑并编写一个 README 文档。

## metas 项目概述

`metas` 是一个用于管理应用元数据的 Go 语言库。它提供了统一的方式来处理应用程序的基本信息、运行环境、配置和日志设置等元数据。

### 主要功能

1. **元数据管理**: 统一管理应用的基本信息，包括应用名称、产品名称、节点信息等
2. **环境配置**: 支持多种运行环境（debug、test、pre、gray、production）
3. **配置管理**: 支持本地 YAML 文件和远程 ETCD 配置
4. **环境变量处理**: 提供带前缀的环境变量读取和设置功能

### 核心组件

#### ModeType (运行模式)

支持五种运行模式:

- `debug`: 调试环境
- `test`: 测试环境
- `pre`: 预发布环境
- `gray`: 灰度环境
- `production`: 生产环境

#### MetadataInterface (元数据接口)

定义了访问元数据的标准接口，包括:

- 环境变量操作 ([Envkey], [Getenv], [Setenv], [Unsetenv])
- 基本信息获取 ([Mode], [Product], [App], [Ip], [Node]
- 配置信息获取 [ConfigType], [ConfigAddr], [ConfigPath])
- 日志路径获取 ([LogDir])

### 使用方法

1. **初始化**:

```go
opts := &Options{
    App:     "myapp",
    Product: "myproduct",
    Config:  "yaml://@./config/debug",
}
Init(opts)
```

2. **访问元数据**:

```go
// 获取应用名称
appName := App()

// 获取运行模式
mode := Mode()

// 获取环境变量
value := Getenv("MY_VAR", "default_value")
```

### 配置支持

支持两种配置方式:

1. **本地 YAML 配置**: `yaml://@./config/debug` 或 `./config/debug`
2. **ETCD 远程配置**: `etcd://127.0.0.1:9091@config/test`

### 环境验证规则

为保证环境一致性，项目实施了严格的配置与环境匹配验证:

- 线下环境(debug/test/pre)只能使用线下配置
- 灰度环境只能使用灰度配置
- 生产环境只能使用生产配置

### 约定规则

- 应用名称和产品名称必须符合正则表达式: `^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
- 环境变量会自动添加应用名称大写形式作为前缀
- 所有路径都会转换为绝对路径存储

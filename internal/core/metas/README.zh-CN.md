# Metas 元数据管理

一个全面的元数据管理系统，提供应用程序配置、环境检测和运行时信息。该包作为应用程序身份、部署环境和配置设置的中心来源。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建元数据](#创建元数据)
  - [访问元数据](#访问元数据)
  - [环境变量](#环境变量)
  - [配置管理](#配置管理)
- [元数据组件](#元数据组件)
- [配置选项](#配置选项)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **应用程序标识**：集中管理应用程序名称和类型
- **环境检测**：自动检测部署环境（本地、测试、预发布、生产）
- **运行时信息**：访问节点 ID、区域、可用区和提供商信息
- **配置管理**：支持 YAML 和 etcd 配置源
- **环境变量**：自动从环境变量加载配置
- **单例模式**：在整个应用程序中全局访问元数据
- **验证**：内置应用程序名称和配置路径验证
- **IP 检测**：自动检测本地 IP 地址
- **路径管理**：自动解析日志和配置路径

## 安装

metas 包是核心框架的一部分，没有外部依赖，仅使用标准库。

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func main() {
    // 创建元数据
    mt := metas.New(&metas.Options{
        App:  "myapp",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })
    
    // 访问元数据
    fmt.Printf("应用: %s\n", mt.App())
    fmt.Printf("类型: %s\n", mt.Kind())
    fmt.Printf("模式: %s\n", mt.Mode())
    fmt.Printf("节点: %s\n", mt.Node())
    fmt.Printf("IP: %s\n", mt.IP())
}
```

## 架构概览

metas 包提供了一个集中的元数据管理系统：

```
┌─────────────────┐
│   Metadata      │ - 全局元数据实例
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options │ - 配置选项
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Environment    │ - 环境变量
    └─────────────────┘
```

**核心组件：**
- **Metadata**：具有应用程序信息的主元数据结构
- **Options**：元数据初始化的配置选项
- **KindType**：应用程序类型（API、Job、Tool、Messaging）
- **ModeType**：部署模式（Local、Test、Staging、Prod）

## 使用指南

### 创建元数据

创建一个新的元数据实例：

```go
import (
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func main() {
    // 使用选项创建元数据
    mt := metas.New(&metas.Options{
        App:    "myapp",
        Kind:   metas.KindAPI,
        Mode:   metas.ModeLocal,
        Config: "yaml://@./configs/local",
    })
    
    // 访问全局元数据实例
    mt := metas.Metadata()
}
```

**验证规则：**
- 应用程序名称必须匹配模式：`^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
- 类型必须是以下之一：API、Job、Tool、Messaging
- 模式必须是以下之一：Local、Test、Staging、Prod
- 配置路径必须匹配模式：`^[a-zA-Z0-9.-_]+/(local|test|staging|prod)[/]?$`

### 访问元数据

访问各种元数据字段：

```go
mt := metas.Metadata()

// 应用程序信息
appName := mt.App()           // 应用程序名称
appKind := mt.Kind()          // 应用程序类型
appMode := mt.Mode()          // 部署模式

// 部署信息
node := mt.Node()             // 节点 ID
region := mt.Region()         // 区域
zone := mt.Zone()             // 可用区
provider := mt.Provider()     // 服务提供商
ip := mt.IP()                 // 本地 IP 地址

// 配置信息
configType := mt.ConfigType() // 配置类型（yaml/etcd）
configAddr := mt.ConfigAddr() // 配置地址
configPath := mt.ConfigPath() // 配置路径

// 日志信息
logDir := mt.LogDir()         // 日志目录
```

### 环境变量

包自动从环境变量读取：

```bash
# 设置环境变量
export SONG_MODE=prod
export SONG_NODE=node-1
export SONG_REGION=us-east-1
export SONG_ZONE=us-east-1a
export SONG_PROVIDER=aws
export SONG_LOG_DIR=/var/log/myapp
```

这些环境变量覆盖默认值：
- `SONG_MODE`：部署模式
- `SONG_NODE`：节点标识符
- `SONG_REGION`：部署区域
- `SONG_ZONE`：可用区
- `SONG_PROVIDER`：云提供商
- `SONG_LOG_DIR`：日志目录

### 配置管理

使用不同的配置源配置元数据：

```go
// YAML 配置
mt := metas.New(&metas.Options{
    App:  "myapp",
    Kind: metas.KindAPI,
    Mode: metas.ModeLocal,
    Config: "yaml://@./configs/local",
})

// Etcd 配置
mt := metas.New(&metas.Options{
    App:  "myapp",
    Kind: metas.KindAPI,
    Mode: metas.ModeProd,
    Config: "etcd://localhost:2379@config/myapp/prod",
})
```

**配置 DSN 格式：**
```
<类型>://[地址]@<路径>
```

其中：
- `类型`：配置类型（yaml 或 etcd）
- `地址`：配置服务器地址（yaml 可选）
- `路径`：配置文件或密钥路径

## 元数据组件

### 应用程序类型（KindType）

包支持不同的应用程序类型：

```go
// API 应用程序
metas.KindAPI      // Web API 服务

// Job 应用程序
metas.KindJob      // 后台作业处理器

// Tool 应用程序
metas.KindTool     // 命令行工具

// Messaging 应用程序
metas.KindMessaging // 消息处理器
```

### 部署模式（ModeType）

包支持不同的部署模式：

```go
// 本地开发
metas.ModeLocal    // 本地开发环境

// 测试
metas.ModeTest     // 测试环境

// 预发布
metas.ModeStaging  // 预发布/预生产环境

// 生产
metas.ModeProd     // 生产环境
```

### 模式验证

检查当前部署模式：

```go
mt := metas.Metadata()

if mt.Mode().IsModeLocal() {
    // 启用调试功能
}

if mt.Mode().IsModeProd() {
    // 使用生产设置
}

if mt.Mode().IsModeTest() || mt.Mode().IsModeStaging() {
    // 使用预发布/测试设置
}
```

## 配置 Options

### Options 结构体

```go
type Options struct {
    App    string   // 应用程序名称（必需）
    Kind   KindType // 应用程序类型（必需）
    Mode   ModeType // 部署模式（可选，默认：local）
    Config string   // 配置地址（可选）
}
```

### 配置常量

```go
const (
    // 配置 DSN
    ConfigDSNDefault       = "yaml://@./configs/local"
    ConfigDSNRegexpPattern = `^(yaml|etcd)://([a-zA-Z.:0-9]*)@([a-zA-Z0-9/._-]+)$`
    
    // 配置路径
    ConfigPathRegexpPattern = `^[a-zA-Z0-9.-_]+/(local|test|staging|prod)[/]?$`
    ConfigDirDefault        = "./configs/local"
    ConfigTypeYaml          = "yaml"
    ConfigTypeEtcd          = "etcd"
    
    // 日志目录
    LogDirDefault = "./logs"
    
    // 应用名称模式
    FlagAppRegexpPattern = `^[a-zA-Z]+[a-zA-Z0-9_-]+[a-zA-Z0-9]+$`
)
```

### 环境变量名称

```go
const (
    EnvNameMode     = "SONG_MODE"
    EnvNameNode     = "SONG_NODE"
    EnvNameRegion   = "SONG_REGION"
    EnvNameZone     = "SONG_ZONE"
    EnvNameProvider = "SONG_PROVIDER"
    EnvNameLogDir   = "SONG_LOG_DIR"
)
```

## 示例代码

### 完整示例：API 应用程序

```go
package main

import (
    "fmt"
    "log"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化元数据
    mt := metas.New(&metas.Options{
        App:  "user-api",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
        Config: "yaml://@./configs/local",
    })
    
    // 记录应用程序信息
    log.Printf("启动 %s (%s) 在 %s 模式", 
        mt.App(), mt.Kind(), mt.Mode())
    log.Printf("节点: %s, 区域: %s, 可用区: %s", 
        mt.Node(), mt.Region(), mt.Zone())
    log.Printf("IP: %s, 日志目录: %s", 
        mt.IP(), mt.LogDir())
    
    // 创建 HTTP 服务器
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "app":    mt.App(),
                    "kind":   mt.Kind().String(),
                    "mode":   mt.Mode().String(),
                    "node":   mt.Node(),
                    "status": "ok",
                })
            })
            
            eng.GET("/info", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "application": mt.App(),
                    "type":        mt.Kind().String(),
                    "environment": mt.Mode().String(),
                    "node":        mt.Node(),
                    "region":      mt.Region(),
                    "zone":        mt.Zone(),
                    "provider":    mt.Provider(),
                    "ip":          mt.IP(),
                    "config_type": mt.ConfigType(),
                    "config_path": mt.ConfigPath(),
                    "log_dir":     mt.LogDir(),
                })
            })
        }),
    })
    
    // 启动服务器
    server.Serve()
}
```

### 示例：环境特定配置

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func initializeApplication() {
    mt := metas.Metadata()
    
    // 根据模式配置
    var configPath string
    switch {
    case mt.Mode().IsModeLocal():
        configPath = "yaml://@./configs/local"
        enableDebugFeatures()
        
    case mt.Mode().IsModeTest():
        configPath = "yaml://@./configs/test"
        enableTestFixtures()
        
    case mt.Mode().IsModeStaging():
        configPath = "etcd://staging-etcd:2379@config/app/staging"
        enableStagingFeatures()
        
    case mt.Mode().IsModeProd():
        configPath = "etcd://prod-etcd:2379@config/app/prod"
        enableProductionOptimizations()
    }
    
    // 加载配置
    loadConfig(configPath)
}

func enableDebugFeatures() {
    // 启用详细日志、调试端点等
}

func enableProductionOptimizations() {
    // 启用连接池、缓存等
}
```

### 示例：多区域部署

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/metas"
)

func main() {
    // 初始化元数据
    mt := metas.New(&metas.Options{
        App:  "global-service",
        Kind: metas.KindAPI,
        Mode: metas.ModeProd,
        Config: "etcd://etcd.global:2379@config/global/prod",
    })
    
    // 检查部署位置
    region := mt.Region()
    zone := mt.Zone()
    
    fmt.Printf("部署在区域: %s, 可用区: %s\n", region, zone)
    
    // 配置区域特定设置
    switch region {
    case "us-east-1":
        configureUSEast()
    case "eu-west-1":
        configureEUWest()
    case "ap-northeast-1":
        configureAPNortheast()
    }
    
    // 配置可用区特定设置
    configureForZone(zone)
}

func configureUSEast() {
    // 美国东部配置
}

func configureEUWest() {
    // 欧盟西部配置
}

func configureAPNortheast() {
    // 亚太地区配置
}

func configureForZone(zone string) {
    // 可用区特定配置
}
```

### 示例：应用程序类型检测

```go
func initializeBasedOnKind() {
    mt := metas.Metadata()
    
    switch mt.Kind() {
    case metas.KindAPI:
        initializeAPIServer()
    case metas.KindJob:
        initializeJobProcessor()
    case metas.KindTool:
        initializeCLI()
    case metas.KindMessaging:
        initializeMessageHandler()
    }
}

func initializeAPIServer() {
    fmt.Println("启动 API 服务器...")
    // 初始化 HTTP 服务器、路由、处理器
}

func initializeJobProcessor() {
    fmt.Println("启动作业处理器...")
    // 初始化作业队列、工作线程
}

func initializeCLI() {
    fmt.Println("启动 CLI 工具...")
    // 初始化命令行界面
}

func initializeMessageHandler() {
    fmt.Println("启动消息处理器...")
    // 初始化消息订阅者、处理器
}
```

### 示例：带元数据的健康检查

```go
func HealthCheckHandler(c *gin.Context) {
    mt := metas.Metadata()
    
    response := gin.H{
        "status":     "healthy",
        "timestamp":  time.Now().UTC(),
        "app":        mt.App(),
        "version":    getVersion(),
        "environment": mt.Mode().String(),
        "node":       mt.Node(),
        "region":     mt.Region(),
        "uptime":     getUptime(),
    }
    
    // 添加提供商特定信息
    if mt.Provider() != "" {
        response["provider"] = mt.Provider()
    }
    
    c.JSON(200, response)
}

func getVersion() string {
    // 从构建信息返回应用程序版本
    return "1.0.0"
}

func getUptime() string {
    // 计算并返回运行时间
    return time.Since(startTime).String()
}
```

## 最佳实践

1. **尽早初始化**：在应用程序启动的最开始初始化元数据。

2. **使用单例模式**：使用 `metas.Metadata()` 访问全局元数据实例，而不是创建多个实例。

3. **验证配置**：始终验证配置路径和应用程序名称。

4. **使用环境变量**：利用环境变量进行特定于部署的设置。

5. **分离环境**：使用不同的模式（local、test、staging、prod）来分离环境。

6. **在日志中包含元数据**：在日志条目中包含应用程序元数据，以便更好地可观察性。

7. **用于健康检查**：在健康检查端点中公开元数据以进行监控。

8. **尊重模式设置**：根据部署模式调整行为（例如，在 local 中启用调试，在 prod 中优化）。

9. **记录配置**：记录所有配置选项和环境变量。

10. **使用一致的命名**：在整个组织中遵循应用程序的命名模式。

11. **保护敏感信息**：不要在日志或响应中公开敏感信息（如凭据）。

12. **监控配置更改**：跟踪生产环境中的配置更改。

## 相关文档

- [Song 框架文档](../../README.md)
- [Viper 配置模块](../vipers/README.md)

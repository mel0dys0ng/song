# Song 框架 - 轻量级模块化 Golang 开发框架

`Song` 是一个轻量级、模块化的 Golang 开发框架，用于构建现代云原生应用，旨在为构建现代云原生应用提供基础设施组件和最佳实践支持。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [核心组件](#核心组件)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)
- [贡献指南](#贡献指南)

## 特性

### 核心特性

- **模块化架构**：松耦合的组件，可独立使用
- **类型安全配置**：强类型配置管理，支持多数据源
- **结构化日志**：高性能结构化日志，支持上下文信息
- **数据库支持**：MySQL 支持读写分离和连接池
- **缓存支持**：Redis 客户端支持多种连接模式（单机、集群、哨兵）
- **HTTP 服务器**：基于 Gin 的 RESTful API 服务器，支持中间件
- **CLI 框架**：基于 Cobra 的命令行界面开发
- **消息订阅**：基于 Redis Streams 的发布/订阅系统
- **HTTP 客户端**：基于 Resty 的 HTTP 客户端，支持请求签名
- **元数据管理**：集中式应用元数据和环境检测

### 技术特性

- **连接池管理**：高效的数据库和 Redis 连接管理
- **读写分离**：自动路由数据库读写操作
- **错误处理**：全面的错误包装和上下文保留
- **配置热重载**：无需重启即可动态更新配置
- **环境检测**：自动检测部署环境
- **安全特性**：CORS、CSRF 保护、请求签名支持
- **可观测性**：结构化日志、指标和追踪支持

## 安装

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0+（用于数据库功能）
- Redis 6.0+（用于缓存和消息订阅）
- etcd 3.5+（可选，用于分布式配置）

### 安装依赖

```bash
go mod download
```

### 安装可选依赖

```bash
# 用于 etcd 配置支持
go get go.etcd.io/etcd/client/v3

# 用于 Consul 配置支持
go get github.com/hashicorp/consul/api
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化元数据
    metas.New(&metas.Options{
        App:  "myapp",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // 创建 HTTP 服务器
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "status": "ok",
                    "app":    metas.Metadata().App(),
                })
            })
        }),
    })

    // 启动服务器
    server.Serve()
}
```

## 项目结构

```
song/
├── cmd/                      # 命令行应用程序
│   └── main.go              # 主入口
├── docs/                     # 文档
├── examples/                 # 示例应用
│   └── demo/                # 示例应用
│       ├── cmd/             # 应用命令
│       │   ├── api/         # API 服务器命令
│       │   └── jobs/        # 作业处理器命令
│       ├── configs/         # 配置文件
│       │   ├── local/       # 本地环境
│       │   ├── test/        # 测试环境
│       │   ├── staging/     # 预发布环境
│       │   └── prod/        # 生产环境
│       ├── internal/        # 应用代码
│       │   ├── api/         # API 层
│       │   ├── cache/       # 缓存层
│       │   ├── client/      # 客户端层
│       │   ├── jobs/        # 作业处理器
│       │   ├── messaging/   # 消息处理
│       │   ├── repository/  # 数据访问层
│       │   ├── service/     # 业务逻辑层
│       │   └── tools/       # 工具类
│       └── migrations/      # 数据库迁移
├── internal/                # 内部包
│   └── core/               # 核心框架组件
│       ├── clients/        # 外部服务客户端
│       │   ├── mysql/      # MySQL 客户端
│       │   ├── pubsub/     # 消息订阅
│       │   ├── redis/      # Redis 客户端
│       │   └── resty/      # HTTP 客户端
│       ├── cobras/         # CLI 框架
│       ├── erlogs/         # 日志框架
│       ├── https/          # HTTP 服务器
│       ├── metas/          # 元数据管理
│       └── vipers/         # 配置管理
├── pkg/                     # 公共包
├── scripts/                 # 工具脚本
└── tests/                   # 测试文件
```

## 核心组件

### 配置管理 (Vipers)

`vipers` 包提供多数据源配置管理支持：

```go
import "github.com/mel0dys0ng/song/internal/core/vipers"

// 加载配置
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderYaml),
    vipers.OnPath("./configs/app.yaml"),
)

// 读取值
port := config.GetInt("server.port", 8080)
host := config.GetString("server.host", "localhost")
```

**特性：**

- 多配置源支持（YAML、JSON、TOML、etcd、Consul）
-
- 带 实时配置更新默认值的类型安全访问
- 环境变量集成

[了解更多](internal/core/vipers/README.zh-CN.md)

### 元数据管理 (Metas)

`metas` 包管理应用元数据和环境检测：

```go
import "github.com/mel0dys0ng/song/internal/core/metas"

// 初始化元数据
metas.New(&metas.Options{
    App:  "myapp",
    Kind: metas.KindAPI,
    Mode: metas.ModeLocal,
})

// 访问元数据
mt := metas.Metadata()
fmt.Printf("应用: %s, 模式: %s\n", mt.App(), mt.Mode())
```

**特性：**

- 应用身份管理
- 环境检测（本地、测试、预发布、生产）
- 运行时信息（节点、区域、可用区、提供商）
- 配置路径管理

[了解更多](internal/core/metas/README.zh-CN.md)

### 日志 (Erlogs)

`erlogs` 包提供结构化日志和错误处理：

```go
import "github.com/mel0dys0ng/song/internal/core/erlogs"

// 简单日志
erlogs.Info(ctx, "用户登录",
    erlogs.OptionFields(
        zap.String("user_id", "123"),
        zap.Int("attempt", 1),
    ),
)

// 带包装的错误日志
err := someOperation()
if err != nil {
    erlogs.Convert(err).
        Wrap("处理用户请求失败").
        ErrorLog(ctx,
            erlogs.OptionFields(
                zap.String("user_id", "123"),
            ),
        )
}
```

**特性：**

- 基于 Zap 的结构化日志
- 错误包装和上下文保留
- 多日志级别（Debug、Info、Warn、Error、Fatal、Panic）
- 上下文字段和堆栈跟踪

[了解更多](internal/core/erlogs/README.zh-CN.md)

### HTTP 服务器 (HTTPS)

`https` 包提供基于 Gin 的 HTTP 服务器：

```go
import (
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/gin-gonic/gin"
)

// 创建服务器
server := https.New([]https.Option{
    https.Port(8080),
    https.Route(func(eng *gin.Engine) {
        eng.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })
    }),
})

// 启动服务器
server.Serve()
```

**特性：**

- Gin 框架集成
- 中间件支持（CORS、CSRF、恢复）
- 请求签名和验证
- 生命周期钩子（启动/停止前后）
- TLS/HTTPS 支持

[了解更多](internal/core/https/README.zh-CN.md)

### MySQL 客户端

`mysql` 包提供读写分离的数据库访问：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/mysql"

// 创建客户端
client, err := mysql.New(&mysql.Options{
    DSN:          "user:pass@tcp(localhost:3306)/db",
    ReadDSNs:     []string{"read1:3306", "read2:3306"},
    MaxOpenConns: 100,
    MaxIdleConns: 10,
})

// 执行查询
rows, err := client.ReadDB().Query("SELECT * FROM users WHERE id = ?", userID)
```

**特性：**

- 读写分离
- 连接池
- 自动故障转移
- 查询日志和指标

[了解更多](internal/core/clients/mysql/README.zh-CN.md)

### Redis 客户端

`redis` 包提供多种连接模式的 Redis 客户端：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/redis"

// 创建客户端
client, err := redis.New(&redis.Options{
    Addr:     "localhost:6379",
    Password: "secret",
    DB:       0,
})

// 使用 Redis
err := client.Set(ctx, "key", "value", 0).Err()
val, err := client.Get(ctx, "key").Result()
```

**特性：**

- 单机、集群和哨兵模式
- 连接池
- Pub/Sub 支持
- Redis Streams 支持

[了解更多](internal/core/clients/redis/README.zh-CN.md)

### 消息订阅 (Pub/Sub)

`pubsub` 包提供基于 Redis Streams 的消息系统：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/pubsub"

// 创建消息器
messager := pubsub.NewMessager(redisClient)

// 发布消息
err := messager.Publish(ctx, "stream", &pubsub.Message{
    ID:      "msg-001",
    Payload: []byte(`{"event": "user.created"}`),
})

// 订阅消息
err := messager.Subscribe(ctx, "stream", handler)
```

**特性：**

- Redis Streams 后端
- 消息确认
- 消费者组
- 重试机制

[了解更多](internal/core/clients/pubsub/README.zh-CN.md)

### HTTP 客户端 (Resty)

`resty` 包提供带签名支持的 HTTP 客户端：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/resty"

// 创建客户端
client := resty.New(&resty.Options{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
})

// 发送请求
resp, err := client.R().
    SetBody(map[string]interface{}{"key": "value"}).
    Post("/endpoint")
```

**特性：**

- 请求签名
- 连接池
- 自动重试
- 响应缓存

[了解更多](internal/core/clients/resty/README.zh-CN.md)

### CLI 框架 (Cobras)

`cobras` 包提供基于 Cobra 的 CLI 开发：

```go
import "github.com/mel0dys0ng/song/internal/core/cobras"

// 创建命令
cmd := cobras.NewCommand("myapp", "1.0.0", "我的应用")

// 添加子命令
cmd.AddCommand(&cobra.Command{
    Use:   "start",
    Short: "启动应用",
    Run: func(cmd *cobra.Command, args []string) {
        // 启动应用
    },
})

// 执行
cmd.Execute()
```

**特性：**

- Cobra 集成
- 命令层级
- 参数解析
- 帮助生成

[了解更多](internal/core/cobras/README.zh-CN.md)

## 架构概览

```
┌─────────────────────────────────────────────────────────┐
│                    应用层                                 │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐           │
│  │  API 服务器 │  │   作业    │  │   工具    │           │
│  └───────────┘  └───────────┘  └───────────┘           │
└─────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────┐
│                   服务层                                 │
│  ┌─────────────────────────────────────────────────┐    │
│  │           业务逻辑与编排                           │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────┐
│                  仓储层                                  │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐           │
│  │   MySQL   │  │   Redis   │  │  外部 API  │           │
│  │  仓储     │  │   缓存    │  │           │           │
│  └───────────┘  └───────────┘  └───────────┘           │
└─────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────┐
│                 基础设施层                               │
│  ┌────────┐ ┌───────┐ ┌────────┐ ┌────────┐ ┌────────┐ │
│  │ 配置   │ │ 日志   │ │ 指标   │ │ 追踪   │ │ 消息   │ │
│  │ Vipers│ │Erlogs │ │        │ │        │ │ PubSub │ │
│  └───────┘ └───────┘ └────────┘ └────────┘ └────────┘ │
└─────────────────────────────────────────────────────────┘
```

## 使用指南

### 项目初始化

1. **创建项目结构**

```bash
mkdir myapp
cd myapp
go mod init github.com/myorg/myapp
```

2. **添加 Song 框架**

```bash
go get github.com/mel0dys0ng/song
```

3. **创建目录结构**

```bash
mkdir -p cmd/internal configs/local internal/core
```

### 配置设置

创建 `configs/local/app.yaml`：

```yaml
server:
  port: 8080
  host: localhost

database:
  dsn: "user:password@tcp(localhost:3306)/myapp"
  max_open_conns: 100
  max_idle_conns: 10

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
```

### 应用启动

创建 `cmd/main.go`：

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/vipers"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/gin-gonic/gin"
)

func main() {
    // 初始化配置
    config, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/local/app.yaml"),
    )
    if err != nil {
        panic(err)
    }

    // 初始化元数据
    metas.New(&metas.Options{
        App:  "myapp",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // 初始化日志
    erlogs.Initialize()

    // 创建 HTTP 服务器
    server := https.New([]https.Option{
        https.Port(config.GetInt("server.port", 8080)),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "status": "ok",
                    "app":    metas.Metadata().App(),
                })
            })
        }),
    })

    // 启动服务器
    server.Serve()
}
```

### 构建和运行

```bash
# 构建
go build -o myapp ./cmd/main.go

# 运行
./myapp
```

## 示例代码

### RESTful API 服务

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "github.com/gin-gonic/gin"
    "net/http"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // 初始化
    metas.New(&metas.Options{
        App:  "user-api",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // 创建服务器
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            // 用户处理器
            eng.GET("/users/:id", getUser)
            eng.POST("/users", createUser)
            eng.PUT("/users/:id", updateUser)
            eng.DELETE("/users/:id", deleteUser)
        }),
    })

    server.Serve()
}

func getUser(c *gin.Context) {
    id := c.Param("id")

    // 从数据库获取用户
    user := &User{
        ID:    1,
        Name:  "张三",
        Email: "zhangsan@example.com",
    }

    c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // 保存用户到数据库
    user.ID = 1

    c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
    id := c.Param("id")
    var user User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    // 更新数据库中的用户

    c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
    id := c.Param("id")

    // 从数据库删除用户

    c.JSON(http.StatusNoContent, nil)
}
```

### 后台作业处理器

```go
package main

import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "time"
)

func main() {
    // 初始化元数据
    metas.New(&metas.Options{
        App:  "job-processor",
        Kind: metas.KindJob,
        Mode: metas.ModeLocal,
    })

    // 初始化 Redis
    redisClient, err := redis.New(&redis.Options{
        Addr: "localhost:6379",
    })
    if err != nil {
        panic(err)
    }

    // 创建消息器
    messager := pubsub.NewMessager(redisClient)

    // 订阅作业队列
    ctx := context.Background()
    err = messager.Subscribe(ctx, "jobs", func(ctx context.Context, msg *pubsub.Message) error {
        // 处理作业
        return processJob(ctx, msg)
    })

    if err != nil {
        panic(err)
    }

    // 保持运行
    select {}
}

func processJob(ctx context.Context, msg *pubsub.Message) error {
    erlogs.Info(ctx, "正在处理作业",
        erlogs.OptionFields(
            zap.String("job_id", msg.ID),
            zap.String("payload", string(msg.Payload)),
        ),
    )

    // 模拟作业处理
    time.Sleep(1 * time.Second)

    erlogs.Info(ctx, "作业完成",
        erlogs.OptionFields(
            zap.String("job_id", msg.ID),
        ),
    )

    return nil
}
```

### 命令行工具

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
    "github.com/spf13/cobra"
)

func main() {
    // 创建根命令
    cmd := cobras.NewCommand("mytool", "1.0.0", "我的 CLI 工具")

    // 添加命令
    cmd.AddCommand(createGreetCommand())
    cmd.AddCommand(createVersionCommand())

    // 执行
    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
    }
}

func createGreetCommand() *cobra.Command {
    var name string

    cmd := &cobra.Command{
        Use:   "greet",
        Short: "打招呼",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("你好，%s！\n", name)
        },
    }

    cmd.Flags().StringVarP(&name, "name", "n", "世界", "要打招呼的名字")

    return cmd
}

func createVersionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "打印版本",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("版本: 1.0.0")
        },
    }
}
```

## 最佳实践

### 项目组织

1. **分离关注点**：保持 API、业务逻辑和数据访问层分离
2. **使用依赖注入**：显式传递依赖项而不是使用全局变量
3. **遵循约定**：在整个项目中保持一致的命名和结构
4. **记录公共 API**：记录所有导出的函数和类型

### 配置管理

1. **环境特定配置**：为本地、测试、预发布和生产使用不同的配置
2. **验证配置**：在启动时验证必需的配置
3. **使用默认值**：为可选配置提供合理的默认值
4. **保护敏感数据**：使用环境变量或密钥管理服务存储凭据

### 日志

1. **使用结构化日志**：始终使用带上下文字段的结构化日志
2. **适当的日志级别**：开发时使用 DEBUG，正常操作使用 INFO，错误使用 ERROR
3. **包含上下文**：在日志条目中包含请求 ID 和用户 ID
4. **避免记录敏感数据**：永远不要记录密码、令牌或个人信息

### 错误处理

1. **包装错误**：使用错误包装保留上下文
2. **显式处理错误**：不要忽略错误
3. **提前返回**：提前返回错误以减少嵌套
4. **使用自定义错误类型**：为特定领域的错误创建自定义错误类型

### 测试

1. **编写单元测试**：单独测试各个组件
2. **使用集成测试**：测试组件交互
3. **模拟外部依赖**：为数据库和外部服务使用模拟
4. **测试边界情况**：测试错误条件和边界值

### 性能

1. **使用连接池**：重用数据库和 Redis 连接
2. **实现缓存**：缓存频繁访问的数据
3. **优化查询**：使用索引并优化数据库查询
4. **监控性能**：使用指标和追踪识别瓶颈

## 贡献指南

我们欢迎贡献！请按照以下步骤操作：

1. Fork 仓库
2. 创建功能分支（`git checkout -b feature/amazing-feature`）
3. 提交更改（`git commit -m 'Add amazing feature'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 打开 Pull Request

### 代码风格

- 遵循 Go 最佳实践
- 使用 `gofmt` 格式化代码
- 编写有意义的提交消息
- 为新功能添加测试
- 更新文档

### 报告问题

- 使用 GitHub Issues 报告 bug
- 包含复现步骤
- 提供环境详情
- 包含日志和错误信息

## 许可证

本项目基于 MIT 许可证 - 参见 [LICENSE](LICENSE) 文件了解更多详情。

## 更多资源

- [核心组件文档](internal/core/README.md)
- [MySQL 客户端](internal/core/clients/mysql/README.zh-CN.md)
- [Redis 客户端](internal/core/clients/redis/README.zh-CN.md)
- [消息订阅](internal/core/clients/pubsub/README.zh-CN.md)
- [HTTP 客户端](internal/core/clients/resty/README.zh-CN.md)
- [CLI 框架](internal/core/cobras/README.zh-CN.md)
- [日志框架](internal/core/erlogs/README.zh-CN.md)
- [HTTP 服务器](internal/core/https/README.zh-CN.md)
- [元数据管理](internal/core/metas/README.zh-CN.md)
- [配置管理](internal/core/vipers/README.zh-CN.md)

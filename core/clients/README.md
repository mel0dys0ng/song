# Clients Package Documentation

`clients` 提供了多种常用客户端的统一封装，包括 MySQL、Redis、HTTP(Resty) 和 Pub/Sub 功能。这些封装简化了客户端的创建和使用，并支持灵活的配置选项，大大简化了在项目中使用这些客户端的复杂度，同时保证了良好的扩展性和维护性。

## 目录结构

```
clients/
├── internal/
│   ├── mysql/          # MySQL 客户端实现
│   ├── redis/          # Redis 客户端实现
│   ├── resty/          # HTTP 客户端实现（基于 Resty）
│   └── pubsub/         # 消息发布/订阅实现（基于 Watermill）
├── mysql.go            # MySQL 客户端对外接口
├── mysql_option.go     # MySQL 配置选项
├── redis.go            # Redis 客户端对外接口
├── redis_option.go     # Redis 配置选项
├── resty.go            # Resty 客户端对外接口
├── resty_option.go     # Resty 配置选项
├── pubub.go            # Pub/Sub 对外接口
└── pubsub_option.go    # Pub/Sub 配置选项
```

## 功能模块

### 1. MySQL 客户端

提供对 MySQL 数据库的访问支持，具备以下特性：

- 读写分离支持
- 负载均衡
- 连接池管理
- 慢查询日志记录

#### 使用方法

```go
// 创建默认 MySQL 客户端
client := clients.NewMySQLClient(ctx, "mysql.config.key")

// 创建自定义 MySQL 客户端
client := clients.CustomMySQLClient(ctx, "custom-name", "mysql.config.key",
    clients.MySQLOptionMaxIdle(50),
    clients.MySQLOptionMaxActive(100))
```

#### 主要方法

- `NewMySQLClient(ctx context.Context, key string)` - 创建基于统一配置的 MySQL 客户端
- `CustomMySQLClient(ctx context.Context, name, key string, opts ...mysql.Option)` - 创建自定义 MySQL 客户端
- `Master()` - 获取主数据库连接（写操作）
- `Slave()` - 获取从数据库连接（读操作）

### 2. Redis 客户端

提供对 Redis 的访问支持，支持多种部署模式：

- 单节点模式
- 集群模式
- 哨兵模式

#### 使用方法

```go
// 创建普通 Redis 客户端
client := clients.NewRedisClient(ctx, "redis.config.key")

// 创建通用 Redis 客户端（支持集群和哨兵）
universalClient := clients.NewRedisUniversalClient(ctx, "redis.config.key")

// 创建自定义 Redis 客户端
client := clients.CustomRedisClient(ctx, "custom-name", "redis.config.key",
    clients.RedisOptionPoolSize(20),
    clients.RedisOptionReadTimeout(5*time.Second))
```

#### 主要方法

- `NewRedisClient(ctx context.Context, key string)` - 创建普通 Redis 客户端
- `NewRedisUniversalClient(ctx context.Context, key string)` - 创建通用 Redis 客户端
- `CustomRedisClient(ctx context.Context, name, key string, options ...redis.Option)` - 创建自定义 Redis 客户端
- `CustomRedisUniversalClient(ctx context.Context, name, key string, options ...redis.Option)` - 创建自定义通用 Redis 客户端

### 3. Resty HTTP 客户端

基于 `go-resty/resty` 封装的 HTTP 客户端，具有以下特性：

- 自动签名验证
- 请求重试机制
- 链路追踪支持
- 内外网请求区分

#### 使用方法

```go
// 创建 HTTP 客户端
client := clients.New(ctx, "resty.config.key")

// 创建自定义 HTTP 客户端
client := clients.Custom(ctx, "custom-name", "resty.config.key",
    clients.RestyOptionTimeout(2*time.Second),
    clients.RestyOptionRetryCount(5))

// 发起请求
resp, err := client.R(ctx).Get("/api/endpoint")
```

#### 主要方法

- `New(ctx context.Context, key string)` - 创建基于统一配置的 HTTP 客户端
- `Custom(ctx context.Context, name, key string, options ...resty.Option)` - 创建自定义 HTTP 客户端
- [R(ctx context.Context)] - 创建带有上下文的请求对象

### 4. Pub/Sub 消息系统

基于 `Watermill` 实现的消息发布/订阅系统，当前主要支持 Redis Streams：

#### 使用方法

```go
// 创建 Redis Stream 发布者
publisher := clients.NewPubSubRedisStreamPublisher(ctx, redisClient)

// 创建 Redis Stream 订阅者
subscriber := clients.NewPubSubRedisStreamSubscriber(ctx, redisClient)

// 创建消息处理器
handler := clients.NewPubSubHandler(ctx,
    clients.PubSubHandlerName("handler-name"),
    clients.PubSubHandlerSubscriber(subscriber),
    clients.PubSubHandlerSubscriberTopic("topic-name"),
    clients.PubSubHandlerHasPublisherFunc(handlerFunc))

// 创建消息路由器并运行
messager := clients.NewPubSubMessager(ctx,
    clients.PubSubMessagerHandlers(handler))
messager.Run()
```

#### 主要组件

- [PubSubRedisStreamPublisher] - Redis Stream 发布者
- [PubSubRedisStreamSubscriber] - Redis Stream 订阅者
- [PubSubHandler] - 消息处理器
- [PubSubMessager] - 消息路由管理器

## 配置选项

各个客户端都支持丰富的配置选项，可以通过对应的方法进行自定义配置。详细配置项可参考各模块的 `*_option.go` 文件。

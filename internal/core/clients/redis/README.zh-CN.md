# Redis 客户端

一个功能强大的 Redis 客户端，支持多种连接模式（单机、集群、哨兵）、连接池管理、Pub/Sub 消息订阅和 Redis Streams 操作。该客户端基于 go-redis 库构建，提供了简洁易用的 API。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建客户端](#创建客户端)
  - [连接模式](#连接模式)
  - [基本操作](#基本操作)
  - [连接池配置](#连接池配置)
  - [Pub/Sub 订阅](#pub/sub-订阅)
- [配置选项](#配置选项)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **多种连接模式**：支持单机、集群和哨兵模式
- **连接池管理**：高效管理连接，支持连接复用
- **事务支持**：支持 Redis 事务（MULTI/EXEC）
- **管道操作**：支持管道操作，提高批量操作性能
- **Lua 脚本**：支持 Lua 脚本执行
- **Pub/Sub**：支持发布/订阅模式
- **Redis Streams**：支持 Redis Streams 操作
- **连接池监控**：提供连接池统计信息
- **自动重试**：支持连接失败自动重试

## 安装

确保已安装 Redis 服务器，然后安装客户端依赖：

```bash
go get github.com/redis/go-redis/v9
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func main() {
    // 创建 Redis 客户端
    client, err := redis.New(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 设置值
    err = client.Set("key", "value", 0).Err()
    if err != nil {
        panic(err)
    }
    
    // 获取值
    val, err := client.Get("key").Result()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("key: %s\n", val)
}
```

## 架构概览

Redis 客户端提供了统一的 Redis 访问接口：

```
┌─────────────────┐
│   Redis Client   │ - 主客户端接口
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options │ - 配置选项
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Pool Manager    │ - 连接池管理
    └─────────────────┘
         │
    ┌────┴────┬─────────┬─────────┐
    │ Standalone│ Cluster │ Sentinel│
    └──────────┴─────────┴─────────┘
```

**核心组件：**
- **Client**：主客户端结构体
- **Options**：客户端配置选项
- **Pool**：连接池管理

## 使用指南

### 创建客户端

创建一个新的 Redis 客户端实例：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/redis"

// 单机模式
client, err := redis.New(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// 集群模式
client, err := redis.NewCluster(&redis.ClusterOptions{
    Addrs: []string{
        "localhost:7000",
        "localhost:7001",
        "localhost:7002",
    },
})

// 哨兵模式
client, err := redis.NewFailover(&redis.FailoverOptions{
    MasterName:    "mymaster",
    SentinelAddrs: []string{"localhost:26379"},
})
```

### 连接模式

#### 单机模式

```go
// 最简单的单机连接
client, err := redis.New(&redis.Options{
    Addr:     "localhost:6379",
    Password: "mypassword",
    DB:       0,
})
```

#### 集群模式

```go
// Redis 集群连接
clusterClient, err := redis.NewCluster(&redis.ClusterOptions{
    Addrs: []string{
        "localhost:7000",
        "localhost:7001",
        "localhost:7002",
        "localhost:7003",
        "localhost:7004",
        "localhost:7005",
    },
    Password: "",
    
    // 集群配置
    MaxRedirects:   3,
    ReadOnly:       true,
    RouteByLatency: true,
})
```

#### 哨兵模式

```go
// 哨兵模式连接
client, err := redis.NewFailover(&redis.FailoverOptions{
    MasterName:    "mymaster",
    SentinelAddrs: []string{"localhost:26379", "localhost:26380"},
    SentinelPassword: "sentinel-password",
    Password:        "master-password",
    DB:              0,
})
```

### 基本操作

执行各种 Redis 操作：

```go
// 字符串操作
err = client.Set("key", "value", 0).Err()
val, err = client.Get("key").Result()

// 哈希操作
err = client.HSet("user:1", "name", "John").Err()
err = client.HSet("user:1", "age", "30").Err()
name, err := client.HGet("user:1", "name").Result()
all, err := client.HGetAll("user:1").Result()

// 列表操作
err = client.LPush("tasks", "task1").Err()
val, err = client.LPop("tasks").Result()

// 集合操作
err = client.SAdd("tags", "tag1", "tag2", "tag3").Err()
members, err := client.SMembers("tags").Result()

// 有序集合操作
err = client.ZAdd("leaderboard", redis.Z{Score: 100, Member: "player1"}).Err()
ranks, err := client.ZRevRangeWithScores("leaderboard", 0, 10).Result()

// 过期操作
err = client.Expire("key", time.Hour).Err()
ttl, err := client.TTL("key").Result()

// 删除操作
err = client.Del("key", "key2", "key3").Err()
```

### 连接池配置

合理配置连接池以获得最佳性能：

```go
client, err := redis.New(&redis.Options{
    Addr:         "localhost:6379",
    Password:     "",
    DB:           0,
    
    // 连接池配置
    PoolSize:     100,          // 连接池大小
    MinIdleConns: 10,           // 最小空闲连接数
    MaxConnAge:   time.Hour,    // 连接最大年龄
    PoolTimeout:  30 * time.Second, // 连接池超时
    IdleTimeout:  10 * time.Minute, // 空闲超时
})
```

### Pub/Sub 订阅

使用 Pub/Sub 订阅消息：

```go
// 订阅频道
pubsub := client.Subscribe("news")
defer pubsub.Close()

// 接收消息
for msg := range pubsub.Channel() {
    fmt.Printf("收到消息: %s from %s\n", msg.Payload, msg.Channel)
}

// 模式订阅
pubsub := client.PSubscribe("news:*")
defer pubsub.Close()

for msg := range pubsub.ChannelMode() {
    fmt.Printf("收到模式消息: %s\n", msg.Payload)
}
```

## 配置 Options

### Options 结构体

```go
type Options struct {
    // 连接配置
    Addr     string // Redis 地址
    Password string // 密码
    DB       int    // 数据库编号
    
    // 连接池配置
    PoolSize     int           // 连接池大小
    MinIdleConns int           // 最小空闲连接
    MaxConnAge   time.Duration // 连接最大年龄
    PoolTimeout  time.Duration // 连接池超时
    IdleTimeout  time.Duration // 空闲超时
    
    // 读写超时
    ReadTimeout  time.Duration // 读超时
    WriteTimeout time.Duration // 写超时
    
    // 其他配置
    TLSConfig *tls.Config // TLS 配置
}
```

### ClusterOptions 结构体

```go
type ClusterOptions struct {
    Addrs []string // 集群节点地址
    
    // 集群配置
    MaxRedirects   int  // 最大重定向次数
    ReadOnly      bool // 只读模式
    RouteByLatency bool // 按延迟路由
    
    // 连接池配置（同 Options）
    PoolSize     int
    MinIdleConns int
    // ...
}
```

## 示例代码

### 完整示例：缓存服务

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CacheService struct {
    client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
    return &CacheService{client: client}
}

func (s *CacheService) GetUser(ctx context.Context, id int) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 从缓存获取
    data, err := s.client.Get(ctx, key).Bytes()
    if err == nil {
        var user User
        json.Unmarshal(data, &user)
        return &user, nil
    }
    
    // 缓存未命中，从数据库获取
    user := &User{
        ID:    id,
        Name:  "John",
        Email: "john@example.com",
    }
    
    // 存入缓存
    jsonData, _ := json.Marshal(user)
    s.client.Set(ctx, key, jsonData, 30*time.Minute)
    
    return user, nil
}

func (s *CacheService) InvalidateUser(ctx context.Context, id int) error {
    key := fmt.Sprintf("user:%d", id)
    return s.client.Del(ctx, key).Err()
}

func main() {
    // 创建客户端
    client, err := redis.New(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 创建缓存服务
    cache := NewCacheService(client)
    ctx := context.Background()
    
    // 获取用户（首次从数据库）
    user1, err := cache.GetUser(ctx, 1)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户 1: %+v\n", user1)
    
    // 再次获取（从缓存）
    user2, err := cache.GetUser(ctx, 1)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户 1 (缓存): %+v\n", user2)
    
    // 使缓存失效
    cache.InvalidateUser(ctx, 1)
}
```

### 示例：分布式锁

```go
func acquireLock(ctx context.Context, client *redis.Client, lockKey string, 
    expiration time.Duration) (bool, error) {
    
    // 使用 SET NX 实现分布式锁
    result, err := client.SetNX(ctx, lockKey, "locked", expiration).Result()
    return result, err
}

func releaseLock(ctx context.Context, client *redis.Client, lockKey string) error {
    // 删除锁
    return client.Del(ctx, lockKey).Err()
}

// 使用分布式锁
func processWithLock(ctx context.Context, client *redis.Client) error {
    lockKey := "process:lock"
    
    // 尝试获取锁
    acquired, err := acquireLock(ctx, client, lockKey, 10*time.Second)
    if err != nil {
        return err
    }
    
    if !acquired {
        return fmt.Errorf("无法获取锁")
    }
    
    // 释放锁（使用 defer）
    defer releaseLock(ctx, client, lockKey)
    
    // 执行任务
    fmt.Println("执行任务...")
    time.Sleep(1 * time.Second)
    
    return nil
}
```

### 示例：计数器

```go
func incrementCounter(ctx context.Context, client *redis.Client, key string) (int64, error) {
    return client.Incr(ctx, key).Result()
}

func getCounter(ctx context.Context, client *redis.Client, key string) (int64, error) {
    return client.Get(ctx, key).Int64()
}

func resetCounter(ctx context.Context, client *redis.Client, key string) error {
    return client.Set(ctx, key, 0, 0).Err()
}
```

### 示例：连接池监控

```go
import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func printPoolStats(client *redis.Client) {
    stats := client.PoolStats()
    fmt.Printf("连接池统计:\n")
    fmt.Printf("  总连接数: %d\n", stats.TotalConns)
    fmt.Printf("  空闲连接数: %d\n", stats.IdleConns)
    fmt.Printf("  使用中的连接数: %d\n", stats.StaleConns)
    fmt.Printf("  点击次数: %d\n", stats.Hits)
    fmt.Printf("  未命中次数: %d\n", stats.Misses)
}
```

## 最佳实践

1. **使用连接池**：始终使用连接池以提高性能

2. **合理配置池大小**：根据应用负载设置合适的 PoolSize

3. **使用上下文**：使用 context.Context 进行超时和取消控制

4. **错误处理**：正确处理 Redis 错误，包括连接超时和命令错误

5. **关闭客户端**：使用 defer client.Close() 确保连接关闭

6. **使用批量操作**：对于多个操作，使用管道或 MGET/MSET

7. **设置过期时间**：为缓存数据设置合理的过期时间

8. **监控连接池**：定期检查连接池统计信息

9. **分布式锁**：使用 SET NX 实现分布式锁

10. **Lua 脚本**：对于复杂操作，使用 Lua 脚本保证原子性

## 相关文档

- [Song 框架文档](../../README.md)
- [MySQL 客户端](../mysql/README.md)
- [Pub/Sub 消息系统](../pubsub/README.md)

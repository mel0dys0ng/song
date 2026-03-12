# Pub/Sub 消息系统

一个基于 Redis Streams 的发布/订阅消息系统，提供消息发布、订阅、消费组支持、消息确认和重试机制。该系统构建在 Redis 之上，利用 Redis Streams 的强大功能实现可靠的消息传递。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建消息器](#创建消息器)
  - [发布消息](#发布消息)
  - [订阅消息](#订阅消息)
  - [消费组](#消费组)
  - [消息确认](#消息确认)
- [配置选项](#配置选项)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **Redis Streams 后端**：利用 Redis Streams 的高性能和可靠性
- **消费组支持**：支持多个消费者组成消费组，实现负载均衡
- **消息确认**：支持消息确认机制，确保消息被正确处理
- **重试机制**：自动重试失败的消息
- **持久化**：消息持久化存储，支持消息回溯
- **消息 ID**：自动生成唯一消息 ID
- **上下文支持**：完全支持 Go 的 context.Context

## 安装

确保已安装 Redis，然后安装必要的依赖：

```bash
go get github.com/redis/go-redis/v9
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "context"
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func main() {
    // 创建 Redis 客户端
    redisClient, err := redis.New(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    if err != nil {
        panic(err)
    }
    
    // 创建消息器
    messager := pubsub.NewMessager(redisClient)
    
    ctx := context.Background()
    
    // 发布消息
    err = messager.Publish(ctx, "orders", &pubsub.Message{
        ID:      "msg-001",
        Payload: []byte(`{"event": "order.created", "order_id": 12345}`),
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Println("消息发布成功")
}
```

## 架构概览

Pub/Sub 消息系统提供了可靠的消息传递机制：

```
┌─────────────────┐
│   Publisher     │ - 消息发布者
└────────┬────────┘
         │
    ┌────┴────┐
    │  Redis   │
    │  Streams │
    └────┬────┘
         │
┌────────┴────────┐
│ Consumer Group  │ - 消费组
└────┬────────┬───┘
     │        │
  ┌──┴──┐  ┌──┴──┐
  │ C1  │  │ C2  │ - 消费者
  └─────┘  └─────┘
```

**核心组件：**
- **Messager**：消息器主结构
- **Message**：消息结构
- **Handler**：消息处理函数

## 使用指南

### 创建消息器

创建一个新的消息器实例：

```go
import (
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

// 创建 Redis 客户端
redisClient, err := redis.New(&redis.Options{
    Addr:     "localhost:6379",
    Password: "secret",
    DB:       0,
})

// 创建消息器
messager := pubsub.NewMessager(redisClient)
```

### 发布消息

向流中发布消息：

```go
ctx := context.Background()

// 简单消息
err := messager.Publish(ctx, "orders", &pubsub.Message{
    Payload: []byte(`{"event": "order.created", "order_id": 12345}`),
})

// 带自定义 ID 的消息
err := messager.Publish(ctx, "orders", &pubsub.Message{
    ID:      "custom-id-001",
    Payload: []byte(`{"event": "order.updated", "order_id": 12345}`),
})

// 批量发布
messages := []*pubsub.Message{
    {Payload: []byte(`{"event": "order.created", "order_id": 1}`)},
    {Payload: []byte(`{"event": "order.created", "order_id": 2}`)},
    {Payload: []byte(`{"event": "order.created", "order_id": 3}`)},
}
err = messager.PublishBatch(ctx, "orders", messages)
```

### 订阅消息

订阅流中的消息：

```go
ctx := context.Background()

// 简单订阅
err := messager.Subscribe(ctx, "orders", func(ctx context.Context, msg *pubsub.Message) error {
    fmt.Printf("收到消息: %s\n", string(msg.Payload))
    return nil
})

// 使用处理器订阅
handler := &MyHandler{}
err := messager.SubscribeWithHandler(ctx, "orders", handler)
```

### 消费组

使用消费组实现负载均衡：

```go
ctx := context.Background()

// 创建消费组
err := messager.CreateGroup(ctx, "orders", "my-group")
if err != nil {
    // 消费组可能已存在
    fmt.Printf("消费组创建: %v\n", err)
}

// 使用消费组订阅
err = messager.SubscribeGroup(ctx, "orders", "my-group", "consumer-1", 
    func(ctx context.Context, msg *pubsub.Message) error {
        fmt.Printf("消费者 1 收到消息: %s\n", string(msg.Payload))
        return nil
    })
```

### 消息确认

确认已处理的消息：

```go
ctx := context.Background()

// 订阅时自动确认
err := messager.Subscribe(ctx, "orders", func(ctx context.Context, msg *pubsub.Message) error {
    // 处理消息
    processMessage(msg)
    
    // 确认消息
    return messager.Ack(ctx, "orders", msg.ID)
})

// 手动确认
err := messager.Ack(ctx, "orders", "msg-id-001")
```

## 配置选项

### Message 结构体

```go
type Message struct {
    ID        string    // 消息 ID（可选，自动生成）
    Payload   []byte    // 消息内容
    Stream    string    // 流名称
    Consumer  string    // 消费者名称
    Timestamp time.Time // 时间戳
}
```

### Handler 接口

```go
type Handler interface {
    Handle(ctx context.Context, msg *Message) error
}
```

## 示例代码

### 完整示例：订单处理系统

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

type OrderHandler struct {
    messager *pubsub.Messager
}

func (h *OrderHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
    fmt.Printf("处理订单消息: %s\n", string(msg.Payload))
    
    // 模拟处理
    time.Sleep(100 * time.Millisecond)
    
    // 确认消息
    return h.messager.Ack(ctx, "orders", msg.ID)
}

func main() {
    // 创建 Redis 客户端
    redisClient, err := redis.New(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    if err != nil {
        panic(err)
    }
    
    // 创建消息器
    messager := pubsub.NewMessager(redisClient)
    ctx := context.Background()
    
    // 创建消费组
    messager.CreateGroup(ctx, "orders", "order-processors")
    
    // 启动消费者
    handler := &OrderHandler{messager: messager}
    
    fmt.Println("订单处理器已启动，等待消息...")
    
    err = messager.SubscribeGroup(ctx, "orders", "order-processors", 
        "consumer-1", handler.Handle)
    if err != nil {
        panic(err)
    }
    
    // 保持运行
    select {}
}

// 模拟订单处理
func processMessage(msg *pubsub.Message) {
    fmt.Printf("处理中: %s\n", string(msg.Payload))
}
```

### 示例：发布订单事件

```go
func publishOrderEvent(redisClient *redis.Client, eventType string, orderID int) error {
    messager := pubsub.NewMessager(redisClient)
    ctx := context.Background()
    
    payload := fmt.Sprintf(`{"event": "%s", "order_id": %d, "timestamp": "%s"}`,
        eventType, orderID, time.Now().Format(time.RFC3339))
    
    return messager.Publish(ctx, "orders", &pubsub.Message{
        Payload: []byte(payload),
    })
}

// 使用
err := publishOrderEvent(redisClient, "order.created", 12345)
if err != nil {
    return err
}
```

### 示例：多个消费者负载均衡

```go
func startConsumers(redisClient *redis.Client, groupName string, consumerCount int) {
    messager := pubsub.NewMessager(redisClient)
    ctx := context.Background()
    
    // 创建消费组
    messager.CreateGroup(ctx, "orders", groupName)
    
    // 启动多个消费者
    for i := 0; i < consumerCount; i++ {
        consumerName := fmt.Sprintf("consumer-%d", i)
        
        go func() {
            err := messager.SubscribeGroup(ctx, "orders", groupName, 
                consumerName, handleOrder)
            if err != nil {
                fmt.Printf("消费者 %s 错误: %v\n", consumerName, err)
            }
        }()
    }
}

func handleOrder(ctx context.Context, msg *pubsub.Message) error {
    fmt.Printf("消费者处理: %s\n", string(msg.Payload))
    // 确认消息
    return nil
}
```

## 最佳实践

1. **使用消费组**：对于生产环境，使用消费组实现消息的可靠处理和负载均衡

2. **消息确认**：始终确认已处理的消息，避免消息丢失

3. **幂等处理**：消息处理应该是幂等的，以便支持重试

4. **错误处理**：正确处理消息处理中的错误，返回错误以触发重试

5. **监控队列积压**：监控流长度，避免消息积压过多

6. **合理设置消费者数量**：根据处理能力和消息量设置合适的消费者数量

7. **使用上下文**：使用 context.Context 进行超时和取消控制

8. **消息持久化**：Redis Streams 会持久化消息，无需额外配置

9. **消息大小**：注意消息大小限制，大的消息考虑拆分

10. **错误重试**：实现适当的重试机制，避免消息永久失败

## 相关文档

- [Song 框架文档](../../README.md)
- [Redis 客户端](../redis/README.md)
- [MySQL 客户端](../mysql/README.md)

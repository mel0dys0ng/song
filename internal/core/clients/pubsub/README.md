# PubSub Client

A robust publish-subscribe messaging system built on Watermill, providing reliable message streaming with Redis Streams backend. This package enables asynchronous communication between services with built-in retry mechanisms and graceful shutdown.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating a Messager](#creating-a-messager)
  - [Creating Handlers](#creating-handlers)
  - [Publishing Messages](#publishing-messages)
  - [Subscribing to Messages](#subscribing-to-messages)
  - [Running the Message Router](#running-the-message-router)
- [Redis Streams Integration](#redis-streams-integration)
- [Message Structure](#message-structure)
- [Middleware Support](#middleware-support)
- [Configuration Options](#configuration-options)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Redis Streams Backend**: Leverages Redis Streams for reliable message persistence
- **Publisher-Subscriber Pattern**: Decoupled communication between services
- **Message Retry Logic**: Built-in retry middleware for failed message processing
- **Graceful Shutdown**: Handles interrupt signals for clean application shutdown
- **Watermill Integration**: Built on the proven Watermill messaging library
- **Middleware Support**: Add custom middleware for logging, metrics, etc.
- **Multiple Handlers**: Support for multiple message handlers with different topics
- **Consumer Groups**: Redis Streams consumer group support for load balancing
- **Comprehensive Logging**: Detailed logging for message publishing and consumption

## Installation

Ensure you have the required dependencies:

```bash
go get github.com/ThreeDotsLabs/watermill
go get github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream
go get github.com/redis/go-redis/v9
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "context"
    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()
    
    // Create Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Create publisher and subscriber
    publisher := pubsub.NewRedisStreamPublisher(ctx, redisClient)
    subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient)
    
    // Create handler
    handler := pubsub.NewHandler(ctx,
        pubsub.HandlerName("order-processor"),
        pubsub.HandlerSubscriberTopic("orders.new"),
        pubsub.HandlerSubscriber(subscriber),
        pubsub.HandlerFunc(func(msg *message.Message) (events message.Events, err error) {
            // Process the message
            println("Received order:", string(msg.Payload))
            return nil, nil
        }),
    )
    
    // Create messager
    messager := pubsub.NewMessager(ctx,
        pubsub.MessagerHandlers(handler),
    )
    
    // Run the messager (blocks until shutdown)
    messager.Run()
}
```

## Architecture Overview

The PubSub package follows a publisher-subscriber architecture:

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│  Publisher  │────▶│ Redis Stream │────▶│  Subscriber  │
└─────────────┘     └──────────────┘     └──────────────┘
                           │
                    ┌──────────────┐
                    │   Message    │
                    │   Handler    │
                    └──────────────┘
```

**Components:**
- **Publisher**: Sends messages to Redis Streams
- **Redis Stream**: Persistent message storage
- **Subscriber**: Receives messages from Redis Streams
- **Handler**: Processes incoming messages
- **Messager**: Orchestrates handlers and manages the message router

## Usage Guide

### Creating a Messager

The `Messager` is the central component that manages message routing:

```go
import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
)

func main() {
    ctx := context.Background()
    
    // Create a basic messager
    messager := pubsub.NewMessager(ctx,
        pubsub.MessagerLogger(
            pubsub.LoggerWatermillLog(true),
            pubsub.LoggerWatermillMeta(true),
        ),
        pubsub.MessagerConfigCloseTimeOut(30 * time.Second),
    )
}
```

### Creating Handlers

Handlers define how messages are processed:

```go
// Basic handler
handler := pubsub.NewHandler(ctx,
    pubsub.HandlerName("email-sender"),
    pubsub.HandlerSubscriberTopic("emails.send"),
    pubsub.HandlerSubscriber(subscriber),
    pubsub.HandlerFunc(func(msg *message.Message) (message.Events, error) {
        // Extract payload
        payload := string(msg.Payload)
        
        // Process the message
        err := sendEmail(payload)
        
        if err != nil {
            return nil, err
        }
        
        return nil, nil
    }),
)

// Handler with retry middleware
handlerWithRetry := pubsub.NewHandler(ctx,
    pubsub.HandlerName("payment-processor"),
    pubsub.HandlerSubscriberTopic("payments.process"),
    pubsub.HandlerSubscriber(subscriber),
    pubsub.HandlerFunc(processPayment),
    pubsub.HandlerMiddlewareRetry(
        3,                              // max retries
        100*time.Millisecond,           // initial interval
        pubsub.MiddlewareRetryMultiplier(2.0),
        pubsub.MiddlewareRetryMaxInterval(10*time.Second),
    ),
)
```

### Publishing Messages

Create and publish messages:

```go
// Create a simple message
msg := pubsub.NewMessage(`{"order_id": 123, "user_id": 456}`)

// Create message with metadata
msg := pubsub.NewMessage(
    `{"order_id": 123}`,
    pubsub.MessageMetadata(map[string]string{
        "event_type": "order_created",
        "version":    "1.0",
    }),
)

// Publish to a topic
err := publisher.Publish(ctx, "orders.new", msg)
if err != nil {
    log.Printf("Failed to publish message: %v", err)
}

// Publish multiple messages
msgs := []*message.Message{
    pubsub.NewMessage(`{"order_id": 123}`),
    pubsub.NewMessage(`{"order_id": 124}`),
}
err := publisher.Publish(ctx, "orders.new", msgs...)
```

### Subscribing to Messages

Configure subscribers with consumer groups:

```go
// Basic subscriber
subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient)

// Subscriber with consumer group
subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient,
    pubsub.RedisStreamSubscriberConsumer("consumer-1"),
    pubsub.RedisStreamSubscriberConsumerGroup("order-processors"),
    pubsub.RedisStreamSubscriberOldestIdFirst(), // Process from beginning
)

// Or process only new messages
subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient,
    pubsub.RedisStreamSubscriberOldestIdLatest(), // Process from now
)
```

### Running the Message Router

Start the message processing:

```go
// Create handlers
handlers := []*pubsub.Handler{
    createOrderHandler(ctx, subscriber),
    createEmailHandler(ctx, subscriber),
    createPaymentHandler(ctx, subscriber),
}

// Create and run messager
messager := pubsub.NewMessager(ctx,
    pubsub.MessagerHandlers(handlers...),
)

// Run blocks until interrupt signal (Ctrl+C or SIGTERM)
messager.Run()
```

## Redis Streams Integration

### Publisher Configuration

```go
publisher := pubsub.NewRedisStreamPublisher(ctx, redisClient,
    pubsub.RedisStreamPublisherConfigDefaultMaxLen(1000), // Max messages per stream
    pubsub.RedisStreamPublisherConfigMaxLens(map[string]int64{
        "orders.new": 5000,    // Custom limit for specific topic
        "emails.send": 2000,
    }),
)
```

### Subscriber Configuration

```go
subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient,
    pubsub.RedisStreamSubscriberConsumer("worker-1"),
    pubsub.RedisStreamSubscriberConsumerGroup("order-workers"),
    pubsub.RedisStreamSubscriberOldestIdFirst(),
)
```

## Message Structure

Messages follow the Watermill message format:

```go
type Message struct {
    UUID     string            // Unique message identifier
    Metadata map[string]string // Message metadata
    Payload  []byte            // Message data
}
```

### Creating Messages with Metadata

```go
msg := pubsub.NewMessage(
    `{"user_id": 123, "action": "signup"}`,
    pubsub.MessageMetadata(map[string]string{
        "source":      "web",
        "timestamp":   "2024-01-01T12:00:00Z",
        "correlation": "abc-123-xyz",
    }),
)

// Access metadata in handler
handler := pubsub.NewHandler(ctx,
    pubsub.HandlerFunc(func(msg *message.Message) (message.Events, error) {
        source := msg.Metadata["source"]
        correlationID := msg.Metadata["correlation"]
        
        // Process message
        return nil, nil
    }),
)
```

## Middleware Support

Add custom middleware to handlers:

```go
// Logging middleware
loggingMiddleware := func(h message.HandlerFunc) message.HandlerFunc {
    return func(msg *message.Message) (message.Events, error) {
        log.Printf("Processing message: %s", msg.UUID)
        events, err := h(msg)
        log.Printf("Finished processing: %s", msg.UUID)
        return events, err
    }
}

handler := pubsub.NewHandler(ctx,
    pubsub.HandlerFunc(processMessage),
    pubsub.HandlerMiddleware(loggingMiddleware),
)

// Multiple middleware
handler := pubsub.NewHandler(ctx,
    pubsub.HandlerFunc(processMessage),
    pubsub.HandlerMiddleware(loggingMiddleware),
    pubsub.HandlerMiddleware(metricsMiddleware),
    pubsub.HandlerMiddlewareRetry(3, 100*time.Millisecond),
)
```

## Configuration Options

### Messager Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `MessagerLogger` | `...LoggerOption` | Configure logger for the messager |
| `MessagerConfigCloseTimeOut` | `time.Duration` | Timeout for closing the router |
| `MessagerHandlers` | `...*Handler` | Add handlers to the messager |

### Handler Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `HandlerName` | `string` | Set handler name (auto-generates UUID) |
| `HandlerPublisherTopic` | `string` | Topic for publishing responses |
| `HandlerPublisher` | `message.Publisher` | Publisher instance |
| `HandlerSubscriberTopic` | `string` | Topic to subscribe to |
| `HandlerSubscriber` | `message.Subscriber` | Subscriber instance |
| `HandlerFunc` | `message.HandlerFunc` | Message handler function |
| `HandlerNoPublisherFunc` | `message.NoPublishHandlerFunc` | Handler without publisher |
| `HandlerMiddleware` | `message.HandlerMiddleware` | Add middleware |
| `HandlerMiddlewareRetry` | `int, time.Duration, ...MiddlewareRetryOption` | Add retry middleware |

### Publisher Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `RedisStreamPublisherLogger` | `...LoggerOption` | Configure publisher logger |
| `RedisStreamPublisherConfigDefaultMaxLen` | `int64` | Default max messages per stream |
| `RedisStreamPublisherConfigMaxLens` | `map[string]int64` | Per-topic max message limits |

### Subscriber Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `RedisStreamSubscriberLogger` | `...LoggerOption` | Configure subscriber logger |
| `RedisStreamSubscriberConsumer` | `string` | Consumer name |
| `RedisStreamSubscriberConsumerGroup` | `string` | Consumer group name |
| `RedisStreamSubscriberOldestIdFirst` | - | Process from beginning of stream |
| `RedisStreamSubscriberOldestIdLatest` | - | Process from current position |

## Examples

### Complete Example: Order Processing System

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"
    
    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "github.com/redis/go-redis/v9"
)

type Order struct {
    OrderID int    `json:"order_id"`
    UserID  int    `json:"user_id"`
    Total   int    `json:"total"`
}

func main() {
    ctx := context.Background()
    
    // Setup Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Create publisher and subscriber
    publisher := pubsub.NewRedisStreamPublisher(ctx, redisClient)
    subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient,
        pubsub.RedisStreamSubscriberConsumerGroup("order-processors"),
        pubsub.RedisStreamSubscriberConsumer("worker-1"),
    )
    
    // Order processing handler
    orderHandler := pubsub.NewHandler(ctx,
        pubsub.HandlerName("order-processor"),
        pubsub.HandlerSubscriberTopic("orders.new"),
        pubsub.HandlerSubscriber(subscriber),
        pubsub.HandlerFunc(func(msg *message.Message) (message.Events, error) {
            var order Order
            if err := json.Unmarshal(msg.Payload, &order); err != nil {
                return nil, err
            }
            
            log.Printf("Processing order %d for user %d", order.OrderID, order.UserID)
            
            // Publish to next stage
            confirmationMsg := pubsub.NewMessage(
                `{"order_id": ` + string(rune(order.OrderID)) + `, "status": "confirmed"}`,
            )
            publisher.Publish(ctx, "orders.confirmed", confirmationMsg)
            
            return nil, nil
        }),
        pubsub.HandlerMiddlewareRetry(3, 100*time.Millisecond),
    )
    
    // Email notification handler
    emailHandler := pubsub.NewHandler(ctx,
        pubsub.HandlerName("email-notifier"),
        pubsub.HandlerSubscriberTopic("orders.confirmed"),
        pubsub.HandlerSubscriber(subscriber),
        pubsub.HandlerFunc(func(msg *message.Message) (message.Events, error) {
            log.Printf("Sending confirmation email: %s", string(msg.Payload))
            return nil, nil
        }),
    )
    
    // Create and run messager
    messager := pubsub.NewMessager(ctx,
        pubsub.MessagerHandlers(orderHandler, emailHandler),
        pubsub.MessagerLogger(pubsub.LoggerWatermillLog(true)),
    )
    
    log.Println("Starting order processing system...")
    messager.Run()
}
```

### Publishing Messages from HTTP Handler

```go
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Create order
    order := Order{OrderID: 123, UserID: 456, Total: 9999}
    orderJSON, _ := json.Marshal(order)
    
    // Publish message
    msg := pubsub.NewMessage(
        string(orderJSON),
        pubsub.MessageMetadata(map[string]string{
            "event_type": "order_created",
        }),
    )
    
    err := publisher.Publish(ctx, "orders.new", msg)
    if err != nil {
        http.Error(w, "Failed to process order", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(order)
}
```

## Best Practices

1. **Use Consumer Groups**: Always use consumer groups for production workloads to enable load balancing and fault tolerance.

2. **Implement Idempotency**: Message processing may be retried, so ensure your handlers are idempotent.

3. **Handle Errors Gracefully**: Use retry middleware for transient errors and implement dead-letter queues for permanent failures.

4. **Monitor Message Lag**: Track the difference between published and processed messages to detect bottlenecks.

5. **Set Appropriate MaxLen**: Configure stream max length based on your retention requirements and storage capacity.

6. **Use Meaningful Metadata**: Include correlation IDs, timestamps, and event types in message metadata for better observability.

7. **Graceful Shutdown**: The messager handles SIGTERM/SIGINT automatically. Ensure your handlers can complete processing within the close timeout.

8. **Partition by Topic**: Use different topics for different event types to enable independent scaling.

9. **Test with Real Traffic**: Test your message processing under realistic load conditions before deploying to production.

10. **Log Important Events**: Use the built-in logging and add custom logging for business-critical events.

## Additional Resources

- [Watermill Documentation](https://watermill.io/)
- [Redis Streams Documentation](https://redis.io/docs/data-types/streams/)
- [Song Framework Documentation](../../README.md)

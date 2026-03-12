# Redis Client

A flexible Redis client supporting standalone, cluster, and sentinel modes with comprehensive connection pooling and configuration management. Built on top of go-redis, this client provides a robust foundation for caching, session management, and distributed data operations.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Connection Modes](#connection-modes)
- [Usage Guide](#usage-guide)
  - [Creating a Client](#creating-a-client)
  - [Basic Operations](#basic-operations)
  - [Advanced Operations](#advanced-operations)
  - [Using with PubSub](#using-with-pubsub)
- [Configuration](#configuration)
- [Configuration Options](#configuration-options)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Multiple Connection Modes**: Support for standalone, cluster, and sentinel deployments
- **Connection Pooling**: Efficient connection management with configurable pool sizes
- **Automatic Retry**: Built-in retry logic for transient failures
- **Timeout Management**: Configurable dial, read, and write timeouts
- **Configuration Management**: Load configuration from YAML, JSON, or other sources
- **Singleton Pattern**: Automatic client reuse and lifecycle management
- **go-redis Integration**: Full support for go-redis/v9 features
- **Comprehensive Logging**: Detailed connection and operation logging

## Installation

Ensure you have the required dependencies:

```bash
go get github.com/redis/go-redis/v9
go get github.com/samber/lo
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "context"
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func main() {
    ctx := context.Background()
    
    // Create a Redis client
    client := redis.CreateClient(ctx, "", "database.redis")
    
    // Basic operations
    err := client.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }
    
    val, err := client.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    
    fmt.Println("key:", val)
}
```

## Connection Modes

The Redis client supports three deployment modes:

### Standalone Mode

```yaml
database:
  redis:
    addrs:
      - "localhost:6379"
    password: "your-password"
    database: 0
    poolSize: 10
```

### Cluster Mode

```yaml
database:
  redis:
    addrs:
      - "node1:6379"
      - "node2:6379"
      - "node3:6379"
    password: "your-password"
    poolSize: 20
```

### Sentinel Mode

```yaml
database:
  redis:
    masterName: "mymaster"
    addrs:
      - "sentinel1:26379"
      - "sentinel2:26379"
      - "sentinel3:26379"
    sentinelUsername: "sentinel-user"
    sentinelPassword: "sentinel-pass"
    username: "app-user"
    password: "app-pass"
    database: 0
```

## Usage Guide

### Creating a Client

Create a Redis client instance:

```go
import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func main() {
    ctx := context.Background()
    
    // Basic client creation (loads config from "database.redis")
    client := redis.CreateClient(ctx, "", "database.redis")
    
    // Client with custom name (for multiple instances)
    cacheClient := redis.CreateClient(ctx, "cache", "database.redis")
    sessionClient := redis.CreateClient(ctx, "session", "database.redis")
}
```

**Parameters:**
- `ctx`: Context for the operation
- `name`: Custom name for the client (optional, use empty string for default)
- `key`: Configuration key to load from config source
- `opts`: Optional configuration options

### Basic Operations

The client provides direct access to go-redis commands:

```go
// String operations
err := client.Set(ctx, "user:1", "John", 0).Err()
value, err := client.Get(ctx, "user:1").Result()

// Hash operations
err := client.HSet(ctx, "user:100", "name", "Alice", "age", 30).Err()
name, _ := client.HGet(ctx, "user:100", "name").Result()

// List operations
err := client.LPush(ctx, "tasks", "task1", "task2").Err()
task, _ := client.RPop(ctx, "tasks").Result()

// Set operations
err := client.SAdd(ctx, "tags", "go", "redis", "cache").Err()
members, _ := client.SMembers(ctx, "tags").Result()

// Sorted Set operations
err := client.ZAdd(ctx, "leaderboard", 
    redis.Z{Score: 100, Member: "player1"},
    redis.Z{Score: 200, Member: "player2"},
).Err()

// Key operations
exists, _ := client.Exists(ctx, "key").Result()
client.Del(ctx, "key")
ttl, _ := client.TTL(ctx, "key").Result()
```

### Advanced Operations

#### Pipelining

```go
pipe := client.Pipeline()
pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)
pipe.Get(ctx, "key1")
cmds, err := pipe.Exec(ctx)
```

#### Transactions

```go
err := client.Watch(ctx, func(tx *redis.Tx) error {
    // Get current value
    val, _ := tx.Get(ctx, "counter").Result()
    
    // Increment
    _, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
        pipe.Set(ctx, "counter", val+1, 0)
        return nil
    })
    
    return err
}, "counter")
```

#### Pub/Sub

```go
// Subscribe to a channel
pubsub := client.Subscribe(ctx, "channel1")
channel := pubsub.Channel()

// Publish to a channel
err := client.Publish(ctx, "channel1", "message").Err()

// Receive messages
go func() {
    for msg := range channel {
        fmt.Println("Received:", msg.Payload)
    }
}()
```

#### Lua Scripts

```go
// Define and execute a Lua script
script := redis.NewScript(`
    local current = redis.call('GET', KEYS[1])
    if current then
        return tonumber(current) + tonumber(ARGV[1])
    end
    return tonumber(ARGV[1])
`)

result, err := script.Run(ctx, client, []string{"counter"}, 5).Result()
```

### Using with PubSub

The Redis client integrates seamlessly with the PubSub package:

```go
import (
    "github.com/mel0dys0ng/song/internal/core/clients/pubsub"
    "github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func main() {
    ctx := context.Background()
    
    // Create Redis client
    redisClient := redis.CreateClient(ctx, "", "database.redis")
    
    // Create PubSub publisher and subscriber
    publisher := pubsub.NewRedisStreamPublisher(ctx, redisClient)
    subscriber := pubsub.NewRedisStreamSubscriber(ctx, redisClient,
        pubsub.RedisStreamSubscriberConsumerGroup("processors"),
        pubsub.RedisStreamSubscriberConsumer("worker-1"),
    )
}
```

## Configuration

### Configuration File Example (YAML)

```yaml
database:
  redis:
    # Basic settings
    debug: false
    addrs:
      - "localhost:6379"
    username: "default"
    password: "your-password"
    database: 0
    
    # Connection settings
    network: "tcp"
    dialTimeout: 5000      # 5 seconds
    readTimeout: 3000      # 3 seconds
    writeTimeout: 3000     # 3 seconds
    
    # Pool settings
    poolSize: 10
    poolTimeout: 4000      # 4 seconds
    minIdleConns: 2
    maxIdleConns: 10
    connMaxIdleTime: 300000   # 5 minutes
    connMaxLifetime: 300000   # 5 minutes
    
    # Retry settings
    maxRetries: 3
    minRetryBackoff: 8     # 8ms
    maxRetryBackoff: 512   # 512ms
    
    # Advanced settings
    poolFIFO: false
    contextTimeoutEnabled: false
```

## Configuration Options

### Programmatic Configuration

You can configure the client using functional options:

```go
client := redis.CreateClient(ctx, "", "database.redis",
    redis.Debug(true),
    redis.Addrs([]string{"localhost:6379", "localhost:6380"}),
    redis.Password("your-password"),
    redis.Database(0),
    redis.PoolSize(20),
    redis.MaxRetries(5),
    redis.DialTimeout(10000),
    redis.ReadTimeout(5000),
)
```

### Available Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `Debug` | `bool` | Enable debug mode | `false` |
| `IdentitySuffix` | `string` | Client identity suffix | `""` |
| `MasterName` | `string` | Sentinel master name | `""` |
| `ClientName` | `string` | Client name (CLIENT SETNAME) | `""` |
| `Network` | `string` | Network type (tcp/unix) | `"tcp"` |
| `Addrs` | `[]string` | Server addresses | `[]` |
| `Username` | `string` | Connection username | `""` |
| `Password` | `string` | Connection password | `""` |
| `SentinelUsername` | `string` | Sentinel username | `""` |
| `SentinelPassword` | `string` | Sentinel password | `""` |
| `Database` | `int` | Database number | `0` |
| `MaxRetries` | `int` | Maximum retry attempts | `3` |
| `MinRetryBackoff` | `time.Duration` | Minimum retry backoff | `8ms` |
| `MaxRetryBackoff` | `time.Duration` | Maximum retry backoff | `512ms` |
| `DialTimeout` | `time.Duration` | Connection timeout | `5s` |
| `ReadTimeout` | `time.Duration` | Read timeout | `3s` |
| `WriteTimeout` | `time.Duration` | Write timeout | `3s` |
| `PoolTimeout` | `time.Duration` | Pool wait timeout | `4s` |
| `PoolFIFO` | `bool` | FIFO pool mode | `false` |
| `PoolSize` | `int` | Maximum pool size | `10` |
| `MinIdleConns` | `int` | Minimum idle connections | `0` |
| `MaxIdleConns` | `int` | Maximum idle connections | `0` |
| `ConnMaxIdleTime` | `time.Duration` | Max connection idle time | `5m` |
| `ConnMaxLifetime` | `time.Duration` | Max connection lifetime | `5m` |
| `MaxRedirects` | `int` | Maximum cluster redirects | `3` |
| `ReadOnly` | `bool` | Read-only mode (cluster) | `false` |
| `RouteByLatency` | `bool` | Route by latency (cluster) | `false` |
| `RouteRandomly` | `bool` | Route randomly (cluster) | `false` |
| `ContextTimeoutEnabled` | `bool` | Enable context timeout | `false` |

## Examples

### Complete Example: Caching Layer

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

func main() {
    ctx := context.Background()
    
    // Create Redis client
    client := redis.CreateClient(ctx, "", "database.redis")
    
    // Create user
    user := User{
        ID:    1,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Serialize user
    userJSON, _ := json.Marshal(user)
    
    // Cache user with TTL (1 hour)
    key := fmt.Sprintf("user:%d", user.ID)
    err := client.Set(ctx, key, userJSON, time.Hour).Err()
    if err != nil {
        panic(err)
    }
    
    // Retrieve cached user
    val, err := client.Get(ctx, key).Result()
    if err != nil {
        panic(err)
    }
    
    // Deserialize user
    var cachedUser User
    json.Unmarshal([]byte(val), &cachedUser)
    
    fmt.Printf("Cached user: %+v\n", cachedUser)
    
    // Check if key exists
    exists, _ := client.Exists(ctx, key).Result()
    fmt.Printf("Key exists: %v\n", exists)
    
    // Get TTL
    ttl, _ := client.TTL(ctx, key).Result()
    fmt.Printf("Key TTL: %v\n", ttl)
}
```

### Example: Distributed Lock

```go
func acquireLock(ctx context.Context, client *redis.Client, lockKey string, ttl time.Duration) (bool, error) {
    // Try to set key with NX (only if not exists)
    ok, err := client.SetNX(ctx, lockKey, "locked", ttl).Result()
    if err != nil {
        return false, err
    }
    return ok, nil
}

func releaseLock(ctx context.Context, client *redis.Client, lockKey string) error {
    return client.Del(ctx, lockKey).Err()
}

// Usage
func processWithLock(ctx context.Context, client *redis.Client) {
    lockKey := "lock:resource:123"
    
    acquired, _ := acquireLock(ctx, client, lockKey, 30*time.Second)
    if !acquired {
        fmt.Println("Could not acquire lock")
        return
    }
    
    defer releaseLock(ctx, client, lockKey)
    
    // Critical section
    fmt.Println("Processing with lock...")
}
```

### Example: Rate Limiting

```go
func isRateLimited(ctx context.Context, client *redis.Client, key string, limit int, window time.Duration) (bool, error) {
    now := time.Now().UnixNano()
    windowStart := now - window.Nanoseconds()
    
    // Remove old entries
    client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
    
    // Add current request
    client.ZAdd(ctx, key, redis.Z{
        Score:  float64(now),
        Member: fmt.Sprintf("%d", now),
    })
    
    // Count requests in window
    count, err := client.ZCard(ctx, key).Result()
    if err != nil {
        return false, err
    }
    
    // Set expiry on the key
    client.Expire(ctx, key, window)
    
    return count > int64(limit), nil
}

// Usage
func handleRequest(ctx context.Context, client *redis.Client, userID string) {
    key := fmt.Sprintf("ratelimit:user:%s", userID)
    
    limited, _ := isRateLimited(ctx, client, key, 100, time.Minute)
    if limited {
        fmt.Println("Rate limit exceeded")
        return
    }
    
    fmt.Println("Processing request...")
}
```

### Example: Session Management

```go
func createSession(ctx context.Context, client *redis.Client, userID string) (string, error) {
    sessionID := uuid.New().String()
    key := fmt.Sprintf("session:%s", sessionID)
    
    sessionData := map[string]interface{}{
        "user_id":    userID,
        "created_at": time.Now().Unix(),
    }
    
    data, _ := json.Marshal(sessionData)
    err := client.Set(ctx, key, data, 24*time.Hour).Err()
    
    return sessionID, err
}

func getSession(ctx context.Context, client *redis.Client, sessionID string) (map[string]interface{}, error) {
    key := fmt.Sprintf("session:%s", sessionID)
    
    val, err := client.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    
    var session map[string]interface{}
    json.Unmarshal([]byte(val), &session)
    
    return session, nil
}
```

## Best Practices

1. **Connection Pooling**: Tune `PoolSize`, `MinIdleConns`, and `MaxIdleConns` based on your application's concurrency patterns.

2. **Set Appropriate Timeouts**: Configure `DialTimeout`, `ReadTimeout`, and `WriteTimeout` to prevent hanging connections.

3. **Use Connection Names**: Set `ClientName` to identify connections in Redis monitoring tools.

4. **Enable Retry Logic**: Use `MaxRetries` with appropriate backoff settings for transient failures.

5. **Monitor Connection Pool**: Track pool utilization and adjust settings based on metrics.

6. **Use TTLs**: Always set expiration times on cached data to prevent memory leaks.

7. **Batch Operations**: Use pipelining for multiple operations to reduce network round trips.

8. **Handle Errors Gracefully**: Implement proper error handling for Redis operations, especially in production.

9. **Use Appropriate Data Structures**: Choose the right Redis data structure (string, hash, list, set, sorted set) for your use case.

10. **Secure Connections**: Use authentication and consider TLS for production deployments.

11. **Implement Circuit Breakers**: Add circuit breaker patterns to handle Redis outages gracefully.

12. **Monitor Memory Usage**: Keep track of Redis memory usage and implement eviction policies.

## Additional Resources

- [go-redis Documentation](https://redis.uptrace.dev/)
- [Redis Documentation](https://redis.io/docs/)
- [Redis Commands Reference](https://redis.io/commands/)
- [Song Framework Documentation](../../README.md)

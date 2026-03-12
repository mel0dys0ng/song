# Song Framework - Lightweight Modular Golang Development Framework

`song` is a lightweight, modular Golang development framework designed for building modern cloud-native applications, providing infrastructure components and best practices support.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Core Components](#core-components)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Contributing](#contributing)

## Features

### Core Features

- **Modular Architecture**: Loosely coupled components that can be used independently
- **Type-Safe Configuration**: Strongly-typed configuration management with multiple sources
- **Structured Logging**: High-performance logging with contextual information
- **Database Support**: MySQL with read/write separation and connection pooling
- **Caching**: Redis client with multiple connection modes (standalone, cluster, sentinel)
- **HTTP Server**: Gin-based RESTful API server with middleware support
- **CLI Framework**: Cobra-based command-line interface development
- **Pub/Sub Messaging**: Redis Streams-based publish/subscribe system
- **HTTP Client**: Resty-based HTTP client with request signing
- **Metadata Management**: Centralized application metadata and environment detection

### Technical Features

- **Connection Pooling**: Efficient database and Redis connection management
- **Read/Write Separation**: Automatic routing for database read and write operations
- **Error Handling**: Comprehensive error wrapping and context preservation
- **Configuration Hot-Reload**: Dynamic configuration updates without restart
- **Environment Detection**: Automatic detection of deployment environment
- **Security**: CORS, CSRF protection, and request signing support
- **Observability**: Structured logging, metrics, and tracing support

## Installation

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0+ (for database features)
- Redis 6.0+ (for caching and pub/sub)
- etcd 3.5+ (optional, for distributed configuration)

### Install Dependencies

```bash
go mod download
```

### Install Optional Dependencies

```bash
# For etcd configuration support
go get go.etcd.io/etcd/client/v3

# For Consul configuration support
go get github.com/hashicorp/consul/api
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize metadata
    metas.New(&metas.Options{
        App:  "myapp",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // Create HTTP server
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

    // Start server
    server.Serve()
}
```

## Project Structure

```
song/
├── cmd/                      # Command-line applications
│   └── main.go              # Main entry point
├── docs/                     # Documentation
├── examples/                 # Example applications
│   └── demo/                # Demo application
│       ├── cmd/             # Application commands
│       │   ├── api/         # API server commands
│       │   └── jobs/        # Job processor commands
│       ├── configs/         # Configuration files
│       │   ├── local/       # Local environment
│       │   ├── test/        # Test environment
│       │   ├── staging/     # Staging environment
│       │   └── prod/        # Production environment
│       ├── internal/        # Application code
│       │   ├── api/         # API layer
│       │   ├── cache/       # Cache layer
│       │   ├── client/      # Client layer
│       │   ├── jobs/        # Job processors
│       │   ├── messaging/   # Message handlers
│       │   ├── repository/  # Data access layer
│       │   ├── service/     # Business logic layer
│       │   └── tools/       # Utility tools
│       └── migrations/      # Database migrations
├── internal/                # Internal packages
│   └── core/               # Core framework components
│       ├── clients/        # External service clients
│       │   ├── mysql/      # MySQL client
│       │   ├── pubsub/     # Pub/Sub messaging
│       │   ├── redis/      # Redis client
│       │   └── resty/      # HTTP client
│       ├── cobras/         # CLI framework
│       ├── erlogs/         # Logging framework
│       ├── https/          # HTTP server
│       ├── metas/          # Metadata management
│       └── vipers/         # Configuration management
├── pkg/                     # Public packages
├── scripts/                 # Utility scripts
└── tests/                   # Test files
```

## Core Components

### Configuration Management (Vipers)

The `vipers` package provides configuration management with support for multiple sources:

```go
import "github.com/mel0dys0ng/song/internal/core/vipers"

// Load configuration
config, err := vipers.New(
    vipers.OnProvider(vipers.ConfigProviderYaml),
    vipers.OnPath("./configs/app.yaml"),
)

// Read values
port := config.GetInt("server.port", 8080)
host := config.GetString("server.host", "localhost")
```

**Features:**

- Multiple configuration sources (YAML, JSON, TOML, etcd, Consul)
- Real-time configuration updates
- Type-safe access with default values
- Environment variable integration

[Learn more](internal/core/vipers/README.md)

### Metadata Management (Metas)

The `metas` package manages application metadata and environment detection:

```go
import "github.com/mel0dys0ng/song/internal/core/metas"

// Initialize metadata
metas.New(&metas.Options{
    App:  "myapp",
    Kind: metas.KindAPI,
    Mode: metas.ModeLocal,
})

// Access metadata
mt := metas.Metadata()
fmt.Printf("App: %s, Mode: %s\n", mt.App(), mt.Mode())
```

**Features:**

- Application identity management
- Environment detection (local, test, staging, prod)
- Runtime information (node, region, zone, provider)
- Configuration path management

[Learn more](internal/core/metas/README.md)

### Logging (Erlogs)

The `erlogs` package provides structured logging with error handling:

```go
import "github.com/mel0dys0ng/song/internal/core/erlogs"

// Simple logging
erlogs.Info(ctx, "User logged in",
    erlogs.OptionFields(
        zap.String("user_id", "123"),
        zap.Int("attempt", 1),
    ),
)

// Error logging with wrapping
err := someOperation()
if err != nil {
    erlogs.Convert(err).
        Wrap("failed to process request").
        ErrorLog(ctx,
            erlogs.OptionFields(
                zap.String("user_id", "123"),
            ),
        )
}
```

**Features:**

- Structured logging with Zap
- Error wrapping and context preservation
- Multiple log levels (Debug, Info, Warn, Error, Fatal, Panic)
- Contextual fields and stack traces

[Learn more](internal/core/erlogs/README.md)

### HTTP Server (HTTPS)

The `https` package provides a Gin-based HTTP server:

```go
import (
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/gin-gonic/gin"
)

// Create server
server := https.New([]https.Option{
    https.Port(8080),
    https.Route(func(eng *gin.Engine) {
        eng.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })
    }),
})

// Start server
server.Serve()
```

**Features:**

- Gin framework integration
- Middleware support (CORS, CSRF, recovery)
- Request signing and validation
- Lifecycle hooks (before/after start/stop)
- TLS/HTTPS support

[Learn more](internal/core/https/README.md)

### MySQL Client

The `mysql` package provides database access with read/write separation:

```go
import "github.com/mel0dys0ng/song/internal/core/clients/mysql"

// Create client
client, err := mysql.New(&mysql.Options{
    DSN:          "user:pass@tcp(localhost:3306)/db",
    ReadDSNs:     []string{"read1:3306", "read2:3306"},
    MaxOpenConns: 100,
    MaxIdleConns: 10,
})

// Execute query
rows, err := client.ReadDB().Query("SELECT * FROM users WHERE id = ?", userID)
```

**Features:**

- Read/write separation
- Connection pooling
- Automatic failover
- Query logging and metrics

[Learn more](internal/core/clients/mysql/README.md)

### Redis Client

The `redis` package provides Redis client with multiple connection modes:

```go
import "github.com/mel0dys0ng/song/internal/core/clients/redis"

// Create client
client, err := redis.New(&redis.Options{
    Addr:     "localhost:6379",
    Password: "secret",
    DB:       0,
})

// Use Redis
err := client.Set(ctx, "key", "value", 0).Err()
val, err := client.Get(ctx, "key").Result()
```

**Features:**

- Standalone, cluster, and sentinel modes
- Connection pooling
- Pub/Sub support
- Redis Streams support

[Learn more](internal/core/clients/redis/README.md)

### Pub/Sub Messaging

The `pubsub` package provides Redis Streams-based messaging:

```go
import "github.com/mel0dys0ng/song/internal/core/clients/pubsub"

// Create messager
messager := pubsub.NewMessager(redisClient)

// Publish message
err := messager.Publish(ctx, "stream", &pubsub.Message{
    ID:      "msg-001",
    Payload: []byte(`{"event": "user.created"}`),
})

// Subscribe to messages
err := messager.Subscribe(ctx, "stream", handler)
```

**Features:**

- Redis Streams backend
- Message acknowledgment
- Consumer groups
- Retry mechanisms

[Learn more](internal/core/clients/pubsub/README.md)

### HTTP Client (Resty)

The `resty` package provides HTTP client with signing support:

```go
import "github.com/mel0dys0ng/song/internal/core/clients/resty"

// Create client
client := resty.New(&resty.Options{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
})

// Make request
resp, err := client.R().
    SetBody(map[string]interface{}{"key": "value"}).
    Post("/endpoint")
```

**Features:**

- Request signing
- Connection pooling
- Automatic retries
- Response caching

[Learn more](internal/core/clients/resty/README.md)

### CLI Framework (Cobras)

The `cobras` package provides CLI development built on Cobra:

```go
import "github.com/mel0dys0ng/song/internal/core/cobras"

// Create command
cmd := cobras.NewCommand("myapp", "1.0.0", "My Application")

// Add subcommand
cmd.AddCommand(&cobra.Command{
    Use:   "start",
    Short: "Start the application",
    Run: func(cmd *cobra.Command, args []string) {
        // Start application
    },
})

// Execute
cmd.Execute()
```

**Features:**

- Cobra integration
- Command hierarchy
- Flag parsing
- Help generation

[Learn more](internal/core/cobras/README.md)

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Application Layer                     │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐           │
│  │  API Server│  │   Jobs    │  │   Tools   │           │
│  └───────────┘  └───────────┘  └───────────┘           │
└─────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────┐
│                   Service Layer                          │
│  ┌─────────────────────────────────────────────────┐    │
│  │           Business Logic & Orchestration         │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                        │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐           │
│  │   MySQL   │  │   Redis   │  │ External  │           │
│  │ Repository│  │  Cache    │  │   APIs    │           │
│  └───────────┘  └───────────┘  └───────────┘           │
└─────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                     │
│  ┌────────┐ ┌───────┐ ┌────────┐ ┌────────┐ ┌────────┐ │
│  │ Config │ │ Logs  │ │ Metrics│ │ Tracing│ │  Msg   │ │
│  │ Vipers │ │Erlogs │ │        │ │        │ │ PubSub │ │
│  └────────┘ └───────┘ └────────┘ └────────┘ └────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Usage Guide

### Project Initialization

1. **Create Project Structure**

```bash
mkdir myapp
cd myapp
go mod init github.com/myorg/myapp
```

2. **Add Song Framework**

```bash
go get github.com/mel0dys0ng/song
```

3. **Create Directory Structure**

```bash
mkdir -p cmd/internal configs/local internal/core
```

### Configuration Setup

Create `configs/local/app.yaml`:

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

### Application Bootstrap

Create `cmd/main.go`:

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
    // Initialize configuration
    config, err := vipers.New(
        vipers.OnProvider(vipers.ConfigProviderYaml),
        vipers.OnPath("./configs/local/app.yaml"),
    )
    if err != nil {
        panic(err)
    }

    // Initialize metadata
    metas.New(&metas.Options{
        App:  "myapp",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // Initialize logging
    erlogs.Initialize()

    // Create HTTP server
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

    // Start server
    server.Serve()
}
```

### Building and Running

```bash
# Build
go build -o myapp ./cmd/main.go

# Run
./myapp
```

## Examples

### RESTful API Service

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
    // Initialize
    metas.New(&metas.Options{
        App:  "user-api",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // Create server
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            // User handlers
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

    // Get user from database
    user := &User{
        ID:    1,
        Name:  "John Doe",
        Email: "john@example.com",
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

    // Save user to database
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

    // Update user in database

    c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
    id := c.Param("id")

    // Delete user from database

    c.JSON(http.StatusNoContent, nil)
}
```

### Background Job Processor

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
    // Initialize metadata
    metas.New(&metas.Options{
        App:  "job-processor",
        Kind: metas.KindJob,
        Mode: metas.ModeLocal,
    })

    // Initialize Redis
    redisClient, err := redis.New(&redis.Options{
        Addr: "localhost:6379",
    })
    if err != nil {
        panic(err)
    }

    // Create messager
    messager := pubsub.NewMessager(redisClient)

    // Subscribe to job queue
    ctx := context.Background()
    err = messager.Subscribe(ctx, "jobs", func(ctx context.Context, msg *pubsub.Message) error {
        // Process job
        return processJob(ctx, msg)
    })

    if err != nil {
        panic(err)
    }

    // Keep running
    select {}
}

func processJob(ctx context.Context, msg *pubsub.Message) error {
    erlogs.Info(ctx, "Processing job",
        erlogs.OptionFields(
            zap.String("job_id", msg.ID),
            zap.String("payload", string(msg.Payload)),
        ),
    )

    // Simulate job processing
    time.Sleep(1 * time.Second)

    erlogs.Info(ctx, "Job completed",
        erlogs.OptionFields(
            zap.String("job_id", msg.ID),
        ),
    )

    return nil
}
```

### Command-Line Tool

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
    "github.com/spf13/cobra"
)

func main() {
    // Create root command
    cmd := cobras.NewCommand("mytool", "1.0.0", "My CLI Tool")

    // Add commands
    cmd.AddCommand(createGreetCommand())
    cmd.AddCommand(createVersionCommand())

    // Execute
    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
    }
}

func createGreetCommand() *cobra.Command {
    var name string

    cmd := &cobra.Command{
        Use:   "greet",
        Short: "Greet someone",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Hello, %s!\n", name)
        },
    }

    cmd.Flags().StringVarP(&name, "name", "n", "World", "Name to greet")

    return cmd
}

func createVersionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Print version",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Version: 1.0.0")
        },
    }
}
```

## Best Practices

### Project Organization

1. **Separate Concerns**: Keep API, business logic, and data access layers separate
2. **Use Dependency Injection**: Pass dependencies explicitly rather than using globals
3. **Follow Conventions**: Use consistent naming and structure across the project
4. **Document Public APIs**: Document all exported functions and types

### Configuration Management

1. **Environment-Specific Configs**: Use different configs for local, test, staging, and prod
2. **Validate Configuration**: Validate required configuration at startup
3. **Use Default Values**: Provide sensible defaults for optional configuration
4. **Secure Sensitive Data**: Use environment variables or secret management for credentials

### Logging

1. **Use Structured Logging**: Always use structured logging with contextual fields
2. **Log at Appropriate Levels**: Use DEBUG for development, INFO for normal operations, ERROR for errors
3. **Include Context**: Include request IDs and user IDs in log entries
4. **Avoid Logging Sensitive Data**: Never log passwords, tokens, or PII

### Error Handling

1. **Wrap Errors**: Use error wrapping to preserve context
2. **Handle Errors Explicitly**: Don't ignore errors
3. **Return Early**: Return early on errors to reduce nesting
4. **Use Custom Error Types**: Create custom error types for domain-specific errors

### Testing

1. **Write Unit Tests**: Test individual components in isolation
2. **Use Integration Tests**: Test component interactions
3. **Mock External Dependencies**: Use mocks for databases and external services
4. **Test Edge Cases**: Test error conditions and boundary values

### Performance

1. **Use Connection Pooling**: Reuse database and Redis connections
2. **Implement Caching**: Cache frequently accessed data
3. **Optimize Queries**: Use indexes and optimize database queries
4. **Monitor Performance**: Use metrics and tracing to identify bottlenecks

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices
- Use `gofmt` to format code
- Write meaningful commit messages
- Add tests for new features
- Update documentation

### Reporting Issues

- Use GitHub Issues to report bugs
- Include steps to reproduce
- Provide environment details
- Include logs and error messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Additional Resources

- [Core Components Documentation](internal/core/README.md)
- [MySQL Client](internal/core/clients/mysql/README.md)
- [Redis Client](internal/core/clients/redis/README.md)
- [Pub/Sub Messaging](internal/core/clients/pubsub/README.md)
- [HTTP Client](internal/core/clients/resty/README.md)
- [CLI Framework](internal/core/cobras/README.md)
- [Logging](internal/core/erlogs/README.md)
- [HTTP Server](internal/core/https/README.md)
- [Metadata](internal/core/metas/README.md)
- [Configuration](internal/core/vipers/README.md)

# HTTPS Server

A comprehensive HTTP/HTTPS server built on top of the Gin framework, providing robust web service capabilities with built-in security features, client information detection, structured logging, and distributed tracing support.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating a Server](#creating-a-server)
  - [Server Configuration](#server-configuration)
  - [Middleware Support](#middleware-support)
  - [Route Registration](#route-registration)
  - [Client Information](#client-information)
- [Security Features](#security-features)
  - [CORS](#cors)
  - [CSRF Protection](#csrf-protection)
  - [Request Signing](#request-signing)
- [Lifecycle Hooks](#lifecycle-hooks)
- [Configuration Options](#configuration-options)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Gin Framework**: Built on the high-performance Gin web framework
- **Security Middleware**: CORS, CSRF, and request signature verification
- **Client Detection**: Automatic parsing of user agent, OS, browser, and device information
- **Structured Logging**: Integration with erlogs for comprehensive logging
- **Distributed Tracing**: Built-in trace ID support for distributed system monitoring
- **Graceful Shutdown**: Support for graceful server startup and shutdown
- **TLS/HTTPS Support**: Full support for HTTPS with TLS configuration
- **Timeout Management**: Configurable read, write, and idle timeouts
- **Recovery Middleware**: Automatic panic recovery with stack traces
- **Custom Middleware**: Easy integration of custom middleware

## Installation

Ensure you have the required dependencies:

```bash
go get github.com/gin-gonic/gin
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/mel0dys0ng/song/internal/core/https"
)

func main() {
    // Create server with options
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{
                    "status": "ok",
                })
            })
        }),
    })
    
    // Start server
    server.Serve()
}
```

Run the server:
```bash
curl http://localhost:8080/health
# Output: {"status":"ok"}
```

## Architecture Overview

The HTTPS server provides a layered architecture:

```
┌─────────────────┐
│  HTTP Server    │ - Net/http server wrapper
└────────┬────────┘
         │
    ┌────┴────┐
    │ Gin Engine│ - Request routing
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Middlewares    │ - Security, logging, recovery
    └─────────────────┘
```

**Key Components:**
- **Server**: Main server structure with configuration
- **Options**: Server configuration options
- **Middleware**: Security and logging middleware
- **Client**: Client information detection

## Usage Guide

### Creating a Server

Create an HTTP server instance:

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/mel0dys0ng/song/internal/core/https"
)

func main() {
    // Basic server creation
    server := https.New([]https.Option{
        https.Port(8080),
    })
    
    // Server with routes
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/api/users", getUsersHandler)
            eng.POST("/api/users", createUserHandler)
        }),
    })
    
    // Start the server
    server.Serve()
}
```

### Server Configuration

Configure the server with various options:

```go
server := https.New([]https.Option{
    // Basic configuration
    https.Port(8080),
    https.Host("0.0.0.0"),
    
    // TLS configuration
    https.TLSOpen(true),
    https.TLSKeyFile("/path/to/key.pem"),
    https.TLSCertFile("/path/to/cert.pem"),
    
    // Timeout configuration
    https.ReadTimeout(30 * time.Second),
    https.ReadHeaderTimeout(30 * time.Second),
    https.WriteTimeout(60 * time.Second),
    https.IdleTimeout(60 * time.Second),
    
    // Connection configuration
    https.KeepAlive(true),
    https.MaxHeaderBytes(1 << 20), // 1MB
    
    // Lifecycle hooks
    https.OnStart(func() {
        fmt.Println("Server started")
    }),
    https.OnShutdown(func() {
        fmt.Println("Server shutting down")
    }),
    
    // Routes
    https.Route(func(eng *gin.Engine) {
        eng.GET("/health", healthHandler)
    }),
})
```

### Middleware Support

Add custom middleware to the server:

```go
// Custom logging middleware
func CustomLogger() https.MiddlewareHandleFunc {
    return func(eng *gin.Engine) gin.HandlerFunc {
        return func(c *gin.Context) {
            start := time.Now()
            c.Next()
            duration := time.Since(start)
            fmt.Printf("Request: %s %s - %d - %v\n", 
                c.Request.Method, c.Request.URL.Path, 
                c.Writer.Status(), duration)
        }
    }
}

// Add middleware with priority
server := https.New([]https.Option{
    https.Middleware(https.Middleware{
        Priority: 1,
        Handle: CustomLogger(),
    }),
    https.Middleware(https.Middleware{
        Priority: 2,
        Handle: AuthenticationMiddleware(),
    }),
})
```

### Route Registration

Register routes using the Route option:

```go
server := https.New([]https.Option{
    https.Route(func(eng *gin.Engine) {
        // Health check
        eng.GET("/health", healthHandler)
        
        // API routes
        api := eng.Group("/api")
        {
            api.GET("/users", getUsersHandler)
            api.POST("/users", createUserHandler)
            api.GET("/users/:id", getUserHandler)
            api.PUT("/users/:id", updateUserHandler)
            api.DELETE("/users/:id", deleteUserHandler)
        }
        
        // Static files
        eng.Static("/static", "./static")
        
        // 404 handler
        eng.NoRoute(notFoundHandler)
    }),
})
```

### Client Information

Access client information in handlers:

```go
func handler(c *gin.Context) {
    // Get client info from context
    clientInfo, exists := c.Get(https.ClientInfoContextKey)
    if exists {
        info := clientInfo.(*https.ClientInfo)
        
        fmt.Printf("IP: %s\n", info.IP)
        fmt.Printf("Device ID: %s\n", info.DeviceID)
        fmt.Printf("OS: %s\n", info.OS)
        fmt.Printf("Browser: %s\n", info.Browser)
        fmt.Printf("Client Type: %d\n", info.ClientType)
        fmt.Printf("User Agent: %s\n", info.UserAgent)
    }
    
    c.JSON(200, gin.H{"status": "ok"})
}
```

## Security Features

### CORS

Configure Cross-Origin Resource Sharing:

```go
server := https.New([]https.Option{
    https.CORS(&https.Cors{
        Enable:       true,
        AllowOrigins: []string{
            "http://localhost:3000",
            "https://yourdomain.com",
        },
        AllowHeaders: []string{
            "Origin",
            "Content-Type",
            "Authorization",
        },
        AllowMethods: []string{
            "GET",
            "POST",
            "PUT",
            "DELETE",
            "OPTIONS",
        },
        ExposeHeaders: []string{
            "Content-Length",
            "Access-Control-Allow-Origin",
        },
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }),
})
```

### CSRF Protection

Enable CSRF protection:

```go
server := https.New([]https.Option{
    https.CSRF(&https.CSRF{
        Enable:         true,
        LookupType:     "header",
        LookupName:     "X-CSRF-Token",
        CookieName:     "X-CSRF-Token",
        CookieDomain:   "yourdomain.com",
        CookiePath:     "/",
        CookieMaxAge:   3600,
        CookieSecure:   true,
        CookieHttpOnly: true,
    }),
})
```

### Request Signing

Enable request signature verification:

```go
server := https.New([]https.Option{
    https.Sign(&https.Sign{
        Enable:   true,
        Secret:   "your-secret-key",
        TTL:      300, // 5 minutes
        Query:    true,
        FormData: true,
        Header:   true,
    }),
})
```

## Lifecycle Hooks

The server supports various lifecycle hooks:

```go
server := https.New([]https.Option{
    // Called when server starts
    https.OnStart(func() {
        fmt.Println("Server started successfully")
        // Initialize resources, connections, etc.
    }),
    
    // Called when server start fails
    https.OnStartFail(func(err error) {
        fmt.Printf("Server start failed: %v\n", err)
    }),
    
    // Called when server is shutting down
    https.OnShutdown(func() {
        fmt.Println("Server shutting down")
        // Cleanup resources, close connections, etc.
    }),
    
    // Called when process exits
    https.OnExit(func() {
        fmt.Println("Process exiting")
    }),
    
    // Called after each request is responded
    https.OnResponded(func(ctx context.Context, data *https.RequestResponseData) {
        fmt.Printf("Request completed: %s %s - %d\n", 
            data.Method, data.Path, data.Status)
    }),
    
    // Called when panic is recovered
    https.OnRecovered(func(ctx context.Context, data *https.RecoveredData) {
        fmt.Printf("Panic recovered: %v\n", data.Error)
    }),
})
```

## Configuration Options

### Server Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `Port` | `int` | Server port number |
| `Host` | `string` | Server host address |
| `TLSOpen` | `bool` | Enable TLS/HTTPS |
| `TLSKeyFile` | `string` | TLS private key file path |
| `TLSCertFile` | `string` | TLS certificate file path |
| `KeepAlive` | `bool` | Enable Keep-Alive |
| `ReadTimeout` | `time.Duration` | Read timeout |
| `ReadHeaderTimeout` | `time.Duration` | Read header timeout |
| `WriteTimeout` | `time.Duration` | Write timeout |
| `IdleTimeout` | `time.Duration` | Idle timeout |
| `MaxHeaderBytes` | `int` | Maximum header size |
| `TmpDir` | `string` | Temporary directory |

### Middleware Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `CORS` | `*Cors` | CORS configuration |
| `CSRF` | `*CSRF` | CSRF protection configuration |
| `Sign` | `*Sign` | Request signing configuration |
| `Middleware` | `Middleware` | Custom middleware |

### Lifecycle Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `OnStart` | `func()` | Called on server start |
| `OnStartFail` | `func(error)` | Called on start failure |
| `OnShutdown` | `func()` | Called on shutdown |
| `OnExit` | `func()` | Called on exit |
| `OnResponded` | `func(context.Context, *RequestResponseData)` | Called after response |
| `OnRecovered` | `func(context.Context, *RecoveredData)` | Called on panic recovery |

### Route Options

| Option | Parameters | Description |
|--------|------------|-------------|
| `Route` | `func(*gin.Engine)` | Route registration function |
| `Init` | `func() error` | Initialization function |
| `Defer` | `func()` | Deferred function |

## Examples

### Complete Example: REST API Server

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "go.uber.org/zap"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func main() {
    // Create server
    server := https.New([]https.Option{
        // Basic configuration
        https.Port(8080),
        https.ReadTimeout(30 * time.Second),
        https.WriteTimeout(60 * time.Second),
        
        // CORS configuration
        https.CORS(&https.Cors{
            Enable:           true,
            AllowOrigins:     []string{"http://localhost:3000"},
            AllowCredentials: true,
        }),
        
        // Lifecycle hooks
        https.OnStart(func() {
            fmt.Println("API server started on :8080")
        }),
        https.OnShutdown(func() {
            fmt.Println("API server shutting down")
        }),
        
        // Routes
        https.Route(func(eng *gin.Engine) {
            // Health check
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{"status": "ok"})
            })
            
            // API routes
            api := eng.Group("/api")
            {
                // Get all users
                api.GET("/users", getUsersHandler)
                
                // Get single user
                api.GET("/users/:id", getUserHandler)
                
                // Create user
                api.POST("/users", createUserHandler)
                
                // Update user
                api.PUT("/users/:id", updateUserHandler)
                
                // Delete user
                api.DELETE("/users/:id", deleteUserHandler)
            }
        }),
    })
    
    // Start server
    server.Serve()
}

func getUsersHandler(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Log request
    erlogs.New("Get users list").
        Options(BaseELOptions()).
        InfoLog(ctx)
    
    c.JSON(http.StatusOK, gin.H{
        "data": users,
    })
}

func getUserHandler(c *gin.Context) {
    ctx := c.Request.Context()
    id := c.Param("id")
    
    // Log request
    erlogs.New("Get user").
        Options(BaseELOptions()).
        InfoLog(ctx,
            erlogs.OptionFields(zap.String("user_id", id)),
        )
    
    // Find user
    for _, user := range users {
        if fmt.Sprintf("%d", user.ID) == id {
            c.JSON(http.StatusOK, gin.H{
                "data": user,
            })
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{
        "error": "User not found",
    })
}

func createUserHandler(c *gin.Context) {
    ctx := c.Request.Context()
    
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        erlogs.Convert(err).
            Wrap("Invalid request body").
            Options(BaseELOptions()).
            ErrorLog(ctx)
        
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }
    
    // Generate ID
    user.ID = len(users) + 1
    
    // Add to list
    users = append(users, user)
    
    // Log creation
    erlogs.New("User created").
        Options(BaseELOptions()).
        InfoLog(ctx,
            erlogs.OptionFields(
                zap.Int("user_id", user.ID),
                zap.String("name", user.Name),
            ),
        )
    
    c.JSON(http.StatusCreated, gin.H{
        "data": user,
    })
}

func updateUserHandler(c *gin.Context) {
    id := c.Param("id")
    
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }
    
    // Find and update user
    for i, u := range users {
        if fmt.Sprintf("%d", u.ID) == id {
            users[i].Name = user.Name
            users[i].Email = user.Email
            
            c.JSON(http.StatusOK, gin.H{
                "data": users[i],
            })
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{
        "error": "User not found",
    })
}

func deleteUserHandler(c *gin.Context) {
    id := c.Param("id")
    
    // Find and delete user
    for i, u := range users {
        if fmt.Sprintf("%d", u.ID) == id {
            users = append(users[:i], users[i+1:]...)
            
            c.JSON(http.StatusOK, gin.H{
                "message": "User deleted",
            })
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{
        "error": "User not found",
    })
}

func BaseELOptions() []erlogs.Option {
    return []erlogs.Option{
        erlogs.OptionKind(erlogs.KindBiz),
        erlogs.OptionCallerSkip(3),
    }
}
```

### Example: Middleware Chain

```go
// Authentication middleware
func AuthenticationMiddleware() https.MiddlewareHandleFunc {
    return func(eng *gin.Engine) gin.HandlerFunc {
        return func(c *gin.Context) {
            token := c.GetHeader("Authorization")
            
            if token == "" {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "error": "Missing authorization token",
                })
                c.Abort()
                return
            }
            
            // Validate token
            if !validateToken(token) {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "error": "Invalid token",
                })
                c.Abort()
                return
            }
            
            c.Next()
        }
    }
}

// Rate limiting middleware
func RateLimitMiddleware() https.MiddlewareHandleFunc {
    return func(eng *gin.Engine) gin.HandlerFunc {
        return func(c *gin.Context) {
            // Implement rate limiting logic
            c.Next()
        }
    }
}

// Usage
server := https.New([]https.Option{
    https.Middleware(https.Middleware{
        Priority: 1,
        Handle:   AuthenticationMiddleware(),
    }),
    https.Middleware(https.Middleware{
        Priority: 2,
        Handle:   RateLimitMiddleware(),
    }),
})
```

### Example: Graceful Shutdown

```go
func main() {
    var server *https.Server
    
    server = https.New([]https.Option{
        https.Port(8080),
        
        https.OnStart(func() {
            fmt.Println("Server started")
        }),
        
        https.OnShutdown(func() {
            fmt.Println("Cleaning up resources...")
            // Close database connections
            // Release other resources
        }),
        
        https.OnExit(func() {
            fmt.Println("Server process exiting")
        }),
        
        https.Route(func(eng *gin.Engine) {
            eng.GET("/shutdown", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "Shutting down"})
                go func() {
                    time.Sleep(1 * time.Second)
                    // Trigger shutdown
                }()
            })
        }),
    })
    
    server.Serve()
}
```

## Best Practices

1. **Use HTTPS in Production**: Always enable TLS/HTTPS for production deployments.

2. **Configure Timeouts**: Set appropriate read, write, and idle timeouts to prevent resource exhaustion.

3. **Enable CORS Carefully**: Only allow trusted origins in production environments.

4. **Implement Rate Limiting**: Protect your API from abuse with rate limiting middleware.

5. **Use Middleware Priority**: Order middleware correctly (authentication before business logic).

6. **Log Important Events**: Log requests, errors, and important business events.

7. **Handle Errors Gracefully**: Return meaningful error messages without exposing internal details.

8. **Validate Input**: Always validate and sanitize user input.

9. **Use Health Checks**: Implement health check endpoints for monitoring and load balancing.

10. **Monitor Performance**: Track request latency, error rates, and resource usage.

11. **Secure Sensitive Data**: Never log sensitive information like passwords or tokens.

12. **Use Connection Pooling**: Reuse database and other resource connections.

13. **Implement Circuit Breakers**: Add circuit breakers for external service calls.

14. **Test Thoroughly**: Write comprehensive tests for your handlers and middleware.

## Additional Resources

- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [Song Framework Documentation](../../README.md)

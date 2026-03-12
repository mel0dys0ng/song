# Erlogs

A comprehensive logging framework built on Uber's Zap logger, providing structured logging, error handling, and log level management. This package offers a fluent API for creating detailed logs with support for business logs, system logs, and distributed tracing.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating Log Entries](#creating-log-entries)
  - [Log Levels](#log-levels)
  - [Error Handling](#error-handling)
  - [Adding Custom Fields](#adding-custom-fields)
  - [Logging Methods](#logging-methods)
- [Configuration](#configuration)
- [Log Types](#log-types)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Structured Logging**: JSON-formatted logs with consistent structure
- **Multiple Log Levels**: Support for Debug, Info, Warn, Error, Panic, and Fatal levels
- **Error Wrapping**: Rich error context with stack traces
- **Log Types**: Separate business logs, system logs, and trace logs
- **Fluent API**: Chainable methods for easy log construction
- **Performance**: Built on high-performance Zap logger
- **Log Rotation**: Automatic log file rotation and compression
- **Sensitive Data Handling**: Automatic masking of sensitive information
- **Trace Integration**: Built-in support for distributed tracing
- **Color Support**: Colorized console output for local development

## Installation

Ensure you have the required dependencies:

```bash
go get go.uber.org/zap
go get go.uber.org/zap/zapcore
go get gopkg.in/natefinch/lumberjack.v2
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "go.uber.org/zap"
)

func main() {
    ctx := context.Background()
    
    // Simple info log
    erlogs.New("User logged in").Options(BaseELOptions()).InfoLog(ctx)
    
    // Error log with fields
    err := someOperation()
    if err != nil {
        erlogs.Convert(err).
            Wrap("failed to process request").
            Options(BaseELOptions()).
            ErrorLog(ctx,
                erlogs.OptionFields(
                    zap.String("user_id", "123"),
                    zap.Int("attempt", 3),
                ),
            )
    }
    
    // Business status log
    erlogs.New("Order created").
        Status(200, "Order created successfully").
        Options(BaseELOptions()).
        InfoLog(ctx)
}
```

## Architecture Overview

The Erlogs package provides a layered logging architecture:

```
┌─────────────────┐
│   ErLog Entry   │ - Log entry builder
└────────┬────────┘
         │
    ┌────┴────┐
    │ Logger  │ - Zap logger wrapper
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Log Output     │ - Console + File
    └─────────────────┘
```

**Key Components:**
- **ErLog**: Log entry builder with fluent API
- **Logger**: Wrapper around Zap logger with configuration
- **Config**: Log configuration (level, output, rotation)
- **Options**: Functional options for customization

## Usage Guide

### Creating Log Entries

Create log entries using the `New` constructor:

```go
import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
)

func main() {
    ctx := context.Background()
    
    // Basic log entry
    log := erlogs.New("Operation completed")
    
    // Log with status
    log = erlogs.New("Request processed").Status(200, "Success")
    
    // Convert error to log
    err := someFunction()
    log = erlogs.Convert(err)
    
    // Write log
    log.Options(BaseELOptions()).InfoLog(ctx)
}
```

### Log Levels

The package supports multiple log levels:

```go
// Debug level - Detailed debugging information
erlogs.New("Variable value").
    Options(BaseELOptions()).
    DebugLog(ctx, erlogs.OptionFields(zap.Any("value", data)))

// Info level - General informational messages
erlogs.New("User logged in").
    Options(BaseELOptions()).
    InfoLog(ctx)

// Warn level - Warning messages that don't stop execution
erlogs.New("Cache miss, falling back to database").
    Options(BaseELOptions()).
    WarnLog(ctx)

// Error level - Error messages that need attention
erlogs.New("Database connection failed").
    Options(BaseELOptions()).
    ErrorLog(ctx)

// Panic level - Critical errors that cause panic
erlogs.New("Configuration missing").
    Options(BaseELOptions()).
    PanicLog(ctx)

// Fatal level - Critical errors that terminate the application
erlogs.New("Unable to start server").
    Options(BaseELOptions()).
    FatalLog(ctx)
```

### Error Handling

Comprehensive error handling with context:

```go
// Convert error to log
err := someOperation()
if err != nil {
    erlogs.Convert(err).
        Options(BaseELOptions()).
        ErrorLog(ctx)
}

// Wrap error with additional context
err := processOrder(order)
if err != nil {
    erlogs.Convert(err).
        Wrap("failed to process order").
        Options(BaseELOptions()).
        ErrorLog(ctx,
            erlogs.OptionFields(
                zap.Int("order_id", order.ID),
                zap.String("customer", order.Customer),
            ),
        )
}

// Format wrap with sprintf-style formatting
err := validateUser(user)
if err != nil {
    erlogs.Convert(err).
        Wrapf("validation failed for user %d", user.ID).
        Options(BaseELOptions()).
        ErrorLog(ctx)
}

// Wrap with error and fields
err := connectToDatabase()
if err != nil {
    erlogs.Convert(err).
        WrapE(err,
            zap.String("host", dbHost),
            zap.Int("port", dbPort),
        ).
        Options(BaseELOptions()).
        ErrorLog(ctx)
}
```

### Adding Custom Fields

Add structured fields to your logs:

```go
// Add fields to log entry
erlogs.New("API request received").
    AppendFields(
        zap.String("method", "POST"),
        zap.String("path", "/api/users"),
        zap.String("client_ip", "192.168.1.1"),
        zap.Int64("duration_ms", 150),
    ).
    Options(BaseELOptions()).
    InfoLog(ctx)

// Use OptionFields in logging call
erlogs.New("User action").
    Options(BaseELOptions()).
    InfoLog(ctx,
        erlogs.OptionFields(
            zap.Int("user_id", 123),
            zap.String("action", "purchase"),
            zap.Float64("amount", 99.99),
        ),
    )

// Add stack trace
erlogs.New("Unexpected error occurred").
    AppendStack().
    Options(BaseELOptions()).
    ErrorLog(ctx)
```

### Logging Methods

Various logging methods for different scenarios:

```go
// Record log without specific level (uses default level)
erlogs.New("General message").
    Options(BaseELOptions()).
    RecordLog(ctx)

// Info log with message
erlogs.New("Server started on port 8080").
    Options(BaseELOptions()).
    InfoLog(ctx)

// Info log with formatted message
erlogs.Newf("Server started on port %d", 8080).
    Options(BaseELOptions()).
    InfoLog(ctx)

// Error log
erlogs.New("Connection timeout").
    Options(BaseELOptions()).
    ErrorLog(ctx)

// Panic log (logs and panics)
erlogs.New("Critical configuration missing").
    Options(BaseELOptions()).
    PanicLog(ctx)

// Fatal log (logs and exits)
erlogs.New("Unable to initialize").
    Options(BaseELOptions()).
    FatalLog(ctx)
```

## Configuration

Configure the logging system:

```yaml
erlog:
  level: "info"           # Log level: debug, info, warn, error
  filePath: "./logs/app"  # Log file path prefix
  maxSize: 100            # Max file size in MB
  maxAge: 30              # Max age in days
  maxBackups: 10          # Max number of old log files
  compress: true          # Enable compression
```

Programmatic configuration:

```go
import (
    "github.com/mel0dys0ng/song/internal/core/erlogs"
)

// Get base options for logging
func BaseELOptions() []erlogs.Option {
    return []erlogs.Option{
        erlogs.OptionKind(erlogs.KindBiz),
        erlogs.OptionCallerSkip(3),
    }
}
```

## Log Types

The package supports different log types for different purposes:

### Business Logs (KindBiz)

For business logic and application events:

```go
erlogs.New("Order placed").
    Options([]erlogs.Option{
        erlogs.OptionKind(erlogs.KindBiz),
    }).
    InfoLog(ctx,
        erlogs.OptionFields(
            zap.Int("order_id", 123),
            zap.String("customer_id", "CUST-456"),
        ),
    )
```

### System Logs (KindSystem)

For system-level events and infrastructure:

```go
erlogs.New("Database connection pool initialized").
    Options([]erlogs.Option{
        erlogs.OptionKind(erlogs.KindSystem),
    }).
    InfoLog(ctx,
        erlogs.OptionFields(
            zap.Int("pool_size", 20),
            zap.String("database", "users"),
        ),
    )
```

### Trace Logs (KindTrace)

For distributed tracing and request tracking:

```go
erlogs.New("Request received").
    Options([]erlogs.Option{
        erlogs.OptionKind(erlogs.KindTrace),
    }).
    InfoLog(ctx,
        erlogs.OptionFields(
            zap.String("trace_id", "abc-123-xyz"),
            zap.String("span_id", "span-456"),
        ),
    )
```

## Examples

### Complete Example: Web Service Logging

```go
package main

import (
    "context"
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "go.uber.org/zap"
)

type Order struct {
    ID       int
    UserID   int
    Amount   float64
    Status   string
}

func processOrder(ctx context.Context, order *Order) error {
    // Log order processing start
    erlogs.New("Processing order").
        Options(BaseELOptions()).
        InfoLog(ctx,
            erlogs.OptionFields(
                zap.Int("order_id", order.ID),
                zap.Int("user_id", order.UserID),
                zap.Float64("amount", order.Amount),
            ),
        )
    
    // Validate order
    if order.Amount <= 0 {
        err := fmt.Errorf("invalid order amount")
        erlogs.Convert(err).
            Wrap("order validation failed").
            Options(BaseELOptions()).
            ErrorLog(ctx,
                erlogs.OptionFields(
                    zap.Int("order_id", order.ID),
                    zap.Float64("amount", order.Amount),
                ),
            )
        return err
    }
    
    // Process payment
    err := processPayment(order)
    if err != nil {
        erlogs.Convert(err).
            Wrap("payment processing failed").
            Options(BaseELOptions()).
            ErrorLog(ctx,
                erlogs.OptionFields(
                    zap.Int("order_id", order.ID),
                ),
            )
        return err
    }
    
    // Update order status
    order.Status = "completed"
    
    // Log success
    erlogs.New("Order processed successfully").
        Status(200, "Order completed").
        Options(BaseELOptions()).
        InfoLog(ctx,
            erlogs.OptionFields(
                zap.Int("order_id", order.ID),
                zap.String("status", order.Status),
            ),
        )
    
    return nil
}

func processPayment(order *Order) error {
    // Simulate payment processing
    if order.Amount > 10000 {
        return fmt.Errorf("amount exceeds limit")
    }
    return nil
}

func BaseELOptions() []erlogs.Option {
    return []erlogs.Option{
        erlogs.OptionKind(erlogs.KindBiz),
        erlogs.OptionCallerSkip(3),
    }
}

func main() {
    ctx := context.Background()
    
    order := &Order{
        ID:     123,
        UserID: 456,
        Amount: 99.99,
    }
    
    err := processOrder(ctx, order)
    if err != nil {
        // Log fatal error
        erlogs.Convert(err).
            Wrap("failed to process order").
            Options(BaseELOptions()).
            FatalLog(ctx)
    }
}
```

### Example: Error Handling in Repository Layer

```go
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
    var user User
    
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // Log as warning for not found
            erlogs.New("User not found").
                Status(404, "User does not exist").
                Options(BaseELOptions()).
                WarnLog(ctx,
                    erlogs.OptionFields(
                        zap.Int("user_id", id),
                    ),
                )
            return nil, nil
        }
        
        // Log database error
        erlogs.Convert(err).
            Wrap("failed to query user").
            Options(BaseELOptions()).
            ErrorLog(ctx,
                erlogs.OptionFields(
                    zap.Int("user_id", id),
                    zap.String("table", "users"),
                ),
            )
        return nil, err
    }
    
    // Log successful query
    erlogs.New("User retrieved").
        Options(BaseELOptions()).
        DebugLog(ctx,
            erlogs.OptionFields(
                zap.Int("user_id", id),
            ),
        )
    
    return &user, nil
}
```

### Example: HTTP Middleware Logging

```go
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        startTime := time.Now()
        
        // Log request
        erlogs.New("HTTP request received").
            Options(BaseELOptions()).
            InfoLog(ctx,
                erlogs.OptionFields(
                    zap.String("method", c.Request.Method),
                    zap.String("path", c.Request.URL.Path),
                    zap.String("client_ip", c.ClientIP()),
                    zap.String("user_agent", c.Request.UserAgent()),
                ),
            )
        
        // Process request
        c.Next()
        
        // Calculate duration
        duration := time.Since(startTime)
        
        // Log response
        erlogs.New("HTTP response sent").
            Options(BaseELOptions()).
            InfoLog(ctx,
                erlogs.OptionFields(
                    zap.Int("status", c.Writer.Status()),
                    zap.Duration("duration_ms", duration),
                    zap.Int("body_size", c.Writer.Size()),
                ),
            )
        
        // Log errors
        if len(c.Errors) > 0 {
            for _, err := range c.Errors {
                erlogs.Convert(err.Err).
                    Wrap("HTTP handler error").
                    Options(BaseELOptions()).
                    ErrorLog(ctx,
                        erlogs.OptionFields(
                            zap.String("path", c.Request.URL.Path),
                        ),
                    )
            }
        }
    }
}
```

### Example: Batch Operation Logging

```go
func ProcessBatchOrders(ctx context.Context, orders []Order) error {
    total := len(orders)
    success := 0
    failed := 0
    
    erlogs.New("Starting batch processing").
        Options(BaseELOptions()).
        InfoLog(ctx,
            erlogs.OptionFields(
                zap.Int("total_orders", total),
            ),
        )
    
    for i, order := range orders {
        err := processOrder(ctx, &order)
        if err != nil {
            failed++
            erlogs.Convert(err).
                Wrapf("failed to process order %d/%d", i+1, total).
                Options(BaseELOptions()).
                ErrorLog(ctx,
                    erlogs.OptionFields(
                        zap.Int("order_id", order.ID),
                        zap.Int("index", i),
                    ),
                )
        } else {
            success++
        }
    }
    
    // Log batch summary
    erlogs.New("Batch processing completed").
        Status(int64(success), fmt.Sprintf("Success: %d, Failed: %d", success, failed)).
        Options(BaseELOptions()).
        InfoLog(ctx,
            erlogs.OptionFields(
                zap.Int("total", total),
                zap.Int("success", success),
                zap.Int("failed", failed),
            ),
        )
    
    if failed > 0 {
        return fmt.Errorf("batch processing completed with %d failures", failed)
    }
    
    return nil
}
```

## Best Practices

1. **Use Appropriate Log Levels**: 
   - DEBUG: Detailed debugging information
   - INFO: General operational events
   - WARN: Warning events that don't stop execution
   - ERROR: Error events that need attention
   - PANIC: Critical errors that cause panic
   - FATAL: Critical errors that terminate the application

2. **Add Context with Fields**: Always include relevant context fields like IDs, timestamps, and operation names.

3. **Use Structured Logging**: Prefer structured fields over string concatenation for better log analysis.

4. **Wrap Errors with Context**: When logging errors, add context about what operation was being performed.

5. **Avoid Logging Sensitive Data**: Never log passwords, tokens, or other sensitive information.

6. **Use Consistent Field Names**: Maintain consistency in field naming across your application.

7. **Log at the Right Level**: Don't log everything as ERROR. Use appropriate levels for different scenarios.

8. **Include Trace IDs**: In distributed systems, always include trace IDs for request correlation.

9. **Monitor Log Volume**: Be mindful of log volume in high-traffic applications. Use sampling if necessary.

10. **Review Log Output**: Regularly review log output to ensure logs are useful and not noisy.

11. **Use Log Rotation**: Configure log rotation to prevent disk space issues.

12. **Test Logging**: Include logging in your tests to ensure important events are being logged.

## Additional Resources

- [Zap Logger Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Lumberjack Log Rotation](https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2)
- [Song Framework Documentation](../../README.md)

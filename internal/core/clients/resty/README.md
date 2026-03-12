# Resty Client

A powerful HTTP client built on resty/v2 with built-in request signing, retry logic, and service-to-service communication features. This package provides a robust foundation for making HTTP requests with automatic headers, signature generation, and comprehensive error handling.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Usage Guide](#usage-guide)
  - [Creating a Client](#creating-a-client)
  - [Making HTTP Requests](#making-http-requests)
  - [Request Signing](#request-signing)
  - [Retry Configuration](#retry-configuration)
  - [Error Handling](#error-handling)
- [Configuration](#configuration)
- [Configuration Options](#configuration-options)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Built on Resty**: Leverages the powerful resty/v2 HTTP client library
- **Request Signing**: Automatic request signature generation and verification
- **Retry Logic**: Configurable retry mechanism with exponential backoff
- **Service Headers**: Automatic injection of service metadata headers
- **Trace Support**: Built-in trace ID and span ID support for distributed tracing
- **Connection Pooling**: Efficient connection management
- **Configuration Management**: Load configuration from YAML, JSON, or other sources
- **Singleton Pattern**: Automatic client reuse and lifecycle management
- **Debug Mode**: Comprehensive logging for debugging
- **Caller Information**: Automatic capture of calling location for debugging

## Installation

Ensure you have the required dependencies:

```bash
go get github.com/go-resty/resty/v2
go get github.com/samber/lo
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "context"
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/resty"
)

func main() {
    ctx := context.Background()
    
    // Create a Resty client
    client := resty.CreateClient(ctx, "", "http.client")
    
    // Make a GET request
    resp, err := client.R(ctx).
        SetQueryParam("key", "value").
        Get("https://api.example.com/users")
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %d\n", resp.StatusCode())
    fmt.Printf("Body: %s\n", string(resp.Body()))
}
```

## Architecture Overview

The Resty client enhances standard HTTP client functionality with service-to-service communication features:

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   Service   │────▶│ Resty Client │────▶│  HTTP API    │
│             │     │  + Signing   │     │   Endpoint   │
└─────────────┘     │  + Retry     │     └──────────────┘
                    │  + Headers   │
                    └──────────────┘
```

**Key Components:**
- **Client**: Enhanced resty client with automatic headers and signing
- **Signer**: Generates and verifies request signatures
- **Retry Handler**: Manages request retries with configurable backoff
- **Header Injector**: Automatically adds service metadata headers

## Usage Guide

### Creating a Client

Create a Resty client instance:

```go
import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/clients/resty"
)

func main() {
    ctx := context.Background()
    
    // Basic client creation (loads config from "http.client")
    client := resty.CreateClient(ctx, "", "http.client")
    
    // Client with custom name (for multiple instances)
    paymentClient := resty.CreateClient(ctx, "payment", "http.client")
    userClient := resty.CreateClient(ctx, "user", "http.client")
}
```

**Parameters:**
- `ctx`: Context for the operation
- `name`: Custom name for the client (optional, use empty string for default)
- `key`: Configuration key to load from config source
- `opts`: Optional configuration options

### Making HTTP Requests

The client provides a fluent interface for making HTTP requests:

```go
// GET request
resp, err := client.R(ctx).
    SetQueryParam("page", "1").
    SetQueryParam("limit", "10").
    Get("/api/users")

// POST request with JSON body
resp, err := client.R(ctx).
    SetHeader("Content-Type", "application/json").
    SetBody(map[string]interface{}{
        "name":  "John",
        "email": "john@example.com",
    }).
    Post("/api/users")

// PUT request
resp, err := client.R(ctx).
    SetPathParam("id", "123").
    SetBody(map[string]interface{}{
        "name": "Jane",
    }).
    Put("/api/users/{id}")

// DELETE request
resp, err := client.R(ctx).
    SetPathParam("id", "123").
    Delete("/api/users/{id}")

// Request with timeout
resp, err := client.R(ctx).
    SetTimeout(5 * time.Second).
    Get("/api/slow-endpoint")
```

### Request Signing

The client supports automatic request signing for secure service-to-service communication:

```go
// Client automatically signs requests when Type is "intranet"
client := resty.CreateClient(ctx, "", "http.client",
    resty.Type(resty.Intranet),
    resty.Did("payment-service"),
    resty.SignSecret("your-secret-key"),
    resty.SignTTL(300), // 5 minutes
)

// All requests from this client will be automatically signed
resp, err := client.R(ctx).
    SetBody(map[string]interface{}{
        "amount": 100,
    }).
    Post("http://payment-service/api/charge")
```

**Automatic Headers:**
When making requests, the client automatically adds:
- `X-Song-Did`: Dependency service ID
- `X-Song-Kd`: App kind
- `X-Song-Na`: App name
- `X-Song-Nd`: App node
- `X-Song-Ts`: Timestamp
- `X-Song-Rs`: Random string
- `X-Song-Fl`: Caller file location
- `X-Song-Sign`: Request signature

### Retry Configuration

Configure retry behavior for transient failures:

```go
client := resty.CreateClient(ctx, "", "http.client",
    resty.RetryCount(3),              // Retry 3 times
    resty.RetryWaitTime(100*time.Millisecond),  // Initial wait
    resty.RetryWaitMaxTime(2*time.Second),      // Max wait
)

// Requests will automatically retry on failure
resp, err := client.R(ctx).Get("/api/unstable-endpoint")
```

### Error Handling

Handle errors from HTTP requests:

```go
resp, err := client.R(ctx).Get("/api/users")

if err != nil {
    // Network error or request failed
    log.Printf("Request failed: %v", err)
    return
}

if resp.StatusCode() >= 400 {
    // HTTP error status
    log.Printf("HTTP error: %d - %s", resp.StatusCode(), string(resp.Body()))
    return
}

// Success
var user User
err = json.Unmarshal(resp.Body(), &user)
```

## Configuration

### Configuration File Example (YAML)

```yaml
http:
  client:
    debug: false
    type: "intranet"        # intranet or extranet
    trace: true
    baseUrl: "https://api.example.com"
    did: "user-service"     # Dependency service ID
    timeout: 500            # milliseconds
    retryCount: 3
    retryWaitTime: 100      # milliseconds
    retryWaitMaxTime: 2000  # milliseconds
    signSecret: "your-secret-key"
    signTTL: 300            # seconds
```

## Configuration Options

### Programmatic Configuration

You can configure the client using functional options:

```go
client := resty.CreateClient(ctx, "", "http.client",
    resty.Debug(true),
    resty.Type(resty.Intranet),
    resty.BaseURL("https://api.example.com"),
    resty.Did("payment-service"),
    resty.Trace(true),
    resty.Timeout(1000*time.Millisecond),
    resty.RetryCount(5),
    resty.RetryWaitTime(200*time.Millisecond),
    resty.RetryWaitMaxTime(5*time.Second),
    resty.SignSecret("your-secret-key"),
    resty.SignTTL(600),
)
```

### Available Options

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `Debug` | `bool` | Enable debug mode | `false` |
| `Type` | `string` | Network type (intranet/extranet) | `"intranet"` |
| `Trace` | `bool` | Enable tracing | `true` |
| `BaseURL` | `string` | Base URL for requests | `""` |
| `Did` | `string` | Dependency service ID | `""` |
| `Timeout` | `time.Duration` | Request timeout | `500ms` |
| `RetryCount` | `int` | Number of retries | `3` |
| `RetryWaitTime` | `time.Duration` | Initial retry wait time | `100ms` |
| `RetryWaitMaxTime` | `time.Duration` | Maximum retry wait time | `2s` |
| `SignSecret` | `string` | Secret key for signing | `"default_sign_secret"` |
| `SignTTL` | `int` | Signature TTL in seconds | `300` |

## Examples

### Complete Example: Service-to-Service Communication

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    "github.com/mel0dys0ng/song/internal/core/clients/resty"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    ctx := context.Background()
    
    // Create client for user service
    userClient := resty.CreateClient(ctx, "user", "http.client",
        resty.BaseURL("http://user-service:8080"),
        resty.Did("user-service"),
        resty.Type(resty.Intranet),
        resty.Timeout(2*time.Second),
        resty.RetryCount(3),
    )
    
    // Create a new user
    createReq := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    resp, err := userClient.R(ctx).
        SetBody(createReq).
        Post("/api/users")
    
    if err != nil {
        log.Fatalf("Failed to create user: %v", err)
    }
    
    if resp.StatusCode() != 201 {
        log.Fatalf("HTTP error: %d", resp.StatusCode())
    }
    
    var createResp CreateUserResponse
    json.Unmarshal(resp.Body(), &createResp)
    fmt.Printf("Created user with ID: %d\n", createResp.ID)
    
    // Get user details
    getResp, err := userClient.R(ctx).
        SetPathParam("id", fmt.Sprintf("%d", createResp.ID)).
        Get("/api/users/{id}")
    
    if err != nil {
        log.Fatalf("Failed to get user: %v", err)
    }
    
    var user User
    json.Unmarshal(getResp.Body(), &user)
    fmt.Printf("User: %+v\n", user)
}
```

### Example: External API Integration

```go
func GetWeatherData(ctx context.Context, city string) (*WeatherData, error) {
    // Create client for external API
    weatherClient := resty.CreateClient(ctx, "weather", "http.client",
        resty.BaseURL("https://api.weather.com"),
        resty.Type(resty.Extranet), // External API
        resty.Timeout(5*time.Second),
        resty.RetryCount(2),
    )
    
    resp, err := weatherClient.R(ctx).
        SetQueryParam("city", city).
        SetQueryParam("appid", "your-api-key").
        Get("/v1/current")
    
    if err != nil {
        return nil, fmt.Errorf("failed to fetch weather: %w", err)
    }
    
    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("weather API error: %d", resp.StatusCode())
    }
    
    var weather WeatherData
    err = json.Unmarshal(resp.Body(), &weather)
    if err != nil {
        return nil, fmt.Errorf("failed to parse weather data: %w", err)
    }
    
    return &weather, nil
}
```

### Example: Batch Requests

```go
func GetUserBatch(ctx context.Context, userIDs []int) ([]User, error) {
    client := resty.CreateClient(ctx, "", "http.client",
        resty.BaseURL("http://user-service:8080"),
        resty.Did("user-service"),
    )
    
    users := make([]User, 0, len(userIDs))
    
    for _, userID := range userIDs {
        resp, err := client.R(ctx).
            SetPathParam("id", fmt.Sprintf("%d", userID)).
            Get("/api/users/{id}")
        
        if err != nil {
            log.Printf("Failed to get user %d: %v", userID, err)
            continue
        }
        
        if resp.StatusCode() == 200 {
            var user User
            json.Unmarshal(resp.Body(), &user)
            users = append(users, user)
        }
    }
    
    return users, nil
}
```

### Example: File Upload

```go
func UploadFile(ctx context.Context, filePath string) error {
    client := resty.CreateClient(ctx, "", "http.client",
        resty.BaseURL("http://file-service:8080"),
        resty.Did("file-service"),
        resty.Timeout(30*time.Second), // Longer timeout for file uploads
    )
    
    resp, err := client.R(ctx).
        SetFile("file", filePath).
        SetFormData(map[string]string{
            "description": "Upload description",
        }).
        Post("/api/files/upload")
    
    if err != nil {
        return fmt.Errorf("upload failed: %w", err)
    }
    
    if resp.StatusCode() != 200 {
        return fmt.Errorf("upload error: %d", resp.StatusCode())
    }
    
    return nil
}
```

### Example: Request with Authentication

```go
func CallAuthenticatedAPI(ctx context.Context) error {
    client := resty.CreateClient(ctx, "", "http.client",
        resty.BaseURL("https://api.example.com"),
        resty.Type(resty.Extranet),
    )
    
    resp, err := client.R(ctx).
        SetHeader("Authorization", "Bearer your-token").
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "data": "value",
        }).
        Post("/api/protected/resource")
    
    if err != nil {
        return err
    }
    
    if resp.StatusCode() != 200 {
        return fmt.Errorf("API error: %d", resp.StatusCode())
    }
    
    return nil
}
```

## Best Practices

1. **Reuse Client Instances**: The client uses singleton pattern internally. Create clients once and reuse them throughout your application.

2. **Set Appropriate Timeouts**: Always configure timeouts based on the expected response time of the API you're calling.

3. **Use Retry Wisely**: Configure retry logic for transient failures but avoid excessive retries for permanent errors.

4. **Enable Tracing**: Keep tracing enabled in production for better observability and debugging.

5. **Secure Secrets**: Store sign secrets in secure configuration sources, not in code.

6. **Handle Errors Gracefully**: Always check both network errors and HTTP status codes.

7. **Use Context**: Pass context to all requests for proper timeout and cancellation handling.

8. **Set BaseURL**: Use BaseURL to avoid repeating the base URL in every request.

9. **Monitor Request Metrics**: Track request latency, error rates, and retry counts in production.

10. **Use Connection Pooling**: The underlying resty client maintains connection pools. Reuse clients to benefit from connection reuse.

11. **Validate Responses**: Always validate response data before using it in your application logic.

12. **Implement Circuit Breakers**: Add circuit breaker patterns for critical external dependencies.

## Additional Resources

- [Resty Documentation](https://github.com/go-resty/resty)
- [Song Framework Documentation](../../README.md)

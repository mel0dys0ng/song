# Resty HTTP 客户端

一个功能强大的 HTTP 客户端，基于 Resty 库构建，提供请求签名、连接池、自动重试、响应缓存等功能。该客户端简化了 HTTP 请求的编写，提供了链式 API 和丰富的配置选项。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建客户端](#创建客户端)
  - [发送请求](#发送请求)
  - [请求签名](#请求签名)
  - [自动重试](#自动重试)
  - [响应缓存](#响应缓存)
- [配置选项](#配置-options)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **链式 API**：简洁的链式调用风格
- **请求签名**：支持多种签名算法（HMAC、AWS、OAuth 等）
- **连接池**：高效的 HTTP 连接池管理
- **自动重试**：自动重试失败的请求
- **响应缓存**：支持响应缓存减少重复请求
- **JSON/XML**：内置 JSON 和 XML 支持
- **文件上传/下载**：支持文件上传和下载
- **中间件**：支持请求/响应中间件
- **调试模式**：详细的调试日志

## 安装

安装 Resty 客户端依赖：

```bash
go get github.com/go-resty/resty/v2
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/resty"
)

func main() {
    // 创建客户端
    client := resty.New(&resty.Options{
        BaseURL: "https://api.example.com",
        Timeout: 30 * 1000000000, // 30秒
    })
    
    // 发送 GET 请求
    resp, err := client.R().Get("/users/1")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("状态码: %d\n", resp.StatusCode())
    fmt.Printf("响应体: %s\n", resp.String())
}
```

## 架构概览

Resty 客户端提供了简洁易用的 HTTP 请求接口：

```
┌─────────────────┐
│  Resty Client    │ - 主客户端接口
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options │ - 配置选项
    └────┬────┘
         │
    ┌────┴────────────┐
    │  Request/Response│ - 请求/响应处理
    └─────────────────┘
         │
    ┌────┴────┐
    │ Signer  │ - 请求签名
    └────┬────┘
         │
    ┌────┴────┐
    │  Retry  │ - 自动重试
    └────┬────┘
         │
    ┌────┴────┐
    │  Cache  │ - 响应缓存
    └─────────┘
```

**核心组件：**
- **Client**：主客户端结构体
- **Options**：客户端配置选项
- **Request**：请求构建器
- **Response**：响应对象

## 使用指南

### 创建客户端

创建一个新的 HTTP 客户端实例：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/resty"

// 基本客户端
client := resty.New(&resty.Options{
    BaseURL: "https://api.example.com",
    Timeout: 30 * 1000000000,
})

// 高级配置
client := resty.New(&resty.Options{
    BaseURL:         "https://api.example.com",
    Timeout:         30 * 1000000000,
    MaxRetryCount:   3,
    RetryWaitTime:   1000000000,  // 1秒
    RetryMaxWaitTime: 30 * 1000000000, // 30秒
    
    // 连接池配置
    PoolConnections: 100,
    PoolIdleTimeout: 30 * 1000000000,
    
    // 请求头
    Headers: map[string]string{
        "User-Agent": "MyApp/1.0",
    },
})
```

### 发送请求

发送各种类型的 HTTP 请求：

```go
// GET 请求
resp, err := client.R().Get("/users")

// POST 请求（JSON）
resp, err := client.R().
    SetHeader("Content-Type", "application/json").
    SetBody(map[string]interface{}{"name": "John", "email": "john@example.com"}).
    Post("/users")

// PUT 请求
resp, err := client.R().
    SetBody(user).
    SetPathParams(map[string]string{"id": "1"}).
    Put("/users/{id}")

// PATCH 请求
resp, err := client.R().
    SetBody(map[string]interface{}{"name": "Jane"}).
    Patch("/users/1")

// DELETE 请求
resp, err := client.R().Delete("/users/1")

// HEAD 请求
resp, err := client.R().Head("/users")
```

### 请求签名

使用签名验证请求：

```go
// HMAC 签名
client := resty.New(&resty.Options{
    BaseURL: "https://api.example.com",
    Signer: &resty.HMASigner{
        Key:       "your-secret-key",
        Algorithm: "sha256",
    },
})

// 使用签名发送请求
resp, err := client.R().
    SetSignUUID("unique-request-id").
    Get("/protected")
```

### 自动重试

配置自动重试：

```go
client := resty.New(&resty.Options{
    BaseURL:       "https://api.example.com",
    MaxRetryCount: 3,              // 最大重试次数
    RetryWaitTime: 1 * 1000000000, // 初始重试等待时间（1秒）
    RetryMaxWaitTime: 30 * 1000000000, // 最大重试等待时间（30秒）
    RetryConditions: []resty.ConditionFunc{
        // 重试条件：状态码为 5xx 或 429
        func(r *resty.Response, err error) bool {
            return r.StatusCode() >= 500 || r.StatusCode() == 429
        },
    },
})

// 发送请求（会自动重试）
resp, err := client.R().Get("/api/data")
```

### 响应缓存

使用响应缓存减少重复请求：

```go
client := resty.New(&resty.Options{
    BaseURL: "https://api.example.com",
    Cache:   true,
})

// GET 请求（自动缓存响应）
resp, err := client.R().
    SetCacheTime(5 * 60 * 1000000000). // 缓存 5 分钟
    Get("/users")

// 清除缓存
client.ClearCache()

// 清除特定 URL 缓存
client.RemoveCache("/users")
```

## 配置 Options

### Options 结构体

```go
type Options struct {
    // 基础配置
    BaseURL string            // 基础 URL
    Timeout time.Duration    // 请求超时
    
    // 重试配置
    MaxRetryCount     int              // 最大重试次数
    RetryWaitTime     time.Duration    // 重试等待时间
    RetryMaxWaitTime  time.Duration    // 最大重试等待时间
    RetryConditions   []ConditionFunc  // 重试条件
    
    // 连接池配置
    PoolConnections int           // 连接池大小
    PoolIdleTimeout time.Duration // 空闲连接超时
    
    // 请求配置
    Headers         map[string]string // 默认请求头
    Auth            string             // 认证信息
    UserAgent       string             // 用户代理
    
    // 签名配置
    Signer Signer // 请求签名器
    
    // 缓存配置
    Cache bool // 启用缓存
    
    // 调试配置
    Debug bool // 调试模式
}
```

### Signer 接口

```go
type Signer interface {
    Sign(request *Request) error
}
```

### HMAC 签名器

```go
type HMASigner struct {
    Key       string // 密钥
    Algorithm string // 算法（sha256, sha512）
}
```

## 示例代码

### 完整示例：API 客户端

```go
package main

import (
    "fmt"
    "time"
    "github.com/mel0dys0ng/song/internal/core/clients/resty"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type APIClient struct {
    client *resty.Client
}

func NewAPIClient(baseURL string) *APIClient {
    client := resty.New(&resty.Options{
        BaseURL:         baseURL,
        Timeout:         30 * 1000000000,
        MaxRetryCount:   3,
        RetryWaitTime:   1 * 1000000000,
        RetryMaxWaitTime: 30 * 1000000000,
        Headers: map[string]string{
            "Content-Type": "application/json",
            "Accept":       "application/json",
        },
    })
    
    return &APIClient{client: client}
}

func (c *APIClient) GetUser(id int) (*User, error) {
    resp, err := c.client.R().Get(fmt.Sprintf("/users/%d", id))
    if err != nil {
        return nil, err
    }
    
    var user User
    err = resp.UnmarshalJSON(&user)
    return &user, err
}

func (c *APIClient) CreateUser(user *User) (*User, error) {
    resp, err := c.client.R().
        SetBody(user).
        Post("/users")
    if err != nil {
        return nil, err
    }
    
    var createdUser User
    err = resp.UnmarshalJSON(&createdUser)
    return &createdUser, err
}

func (c *APIClient) UpdateUser(id int, user *User) (*User, error) {
    resp, err := c.client.R().
        SetBody(user).
        Put(fmt.Sprintf("/users/%d", id))
    if err != nil {
        return nil, err
    }
    
    var updatedUser User
    err = resp.UnmarshalJSON(&updatedUser)
    return &updatedUser, err
}

func (c *APIClient) DeleteUser(id int) error {
    _, err := c.client.R().Delete(fmt.Sprintf("/users/%d", id))
    return err
}

func main() {
    api := NewAPIClient("https://api.example.com")
    
    // 获取用户
    user, err := api.GetUser(1)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户: %+v\n", user)
    
    // 创建用户
    newUser := &User{Name: "张三", Email: "zhangsan@example.com"}
    created, err := api.CreateUser(newUser)
    if err != nil {
        panic(err)
    }
    fmt.Printf("创建用户: %+v\n", created)
    
    // 更新用户
    user.Name = "李四"
    updated, err := api.UpdateUser(user.ID, user)
    if err != nil {
        panic(err)
    }
    fmt.Printf("更新用户: %+v\n", updated)
    
    // 删除用户
    err = api.DeleteUser(user.ID)
    if err != nil {
        panic(err)
    }
}
```

### 示例：文件上传

```go
func uploadFile(client *resty.Client, filePath string) error {
    resp, err := client.R().
        SetFile("file", filePath).
        SetField("description", "测试文件").
        Post("/upload")
    
    if err != nil {
        return err
    }
    
    if resp.StatusCode() >= 400 {
        return fmt.Errorf("上传失败: %s", resp.String())
    }
    
    return nil
}
```

### 示例：文件下载

```go
func downloadFile(client *resty.Client, url, savePath string) error {
    resp, err := client.R().Get(url)
    if err != nil {
        return err
    }
    
    if resp.StatusCode() >= 400 {
        return fmt.Errorf("下载失败: %s", resp.String())
    }
    
    return WriteFile(savePath, resp.Body())
}
```

### 示例：Webhook 签名验证

```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

func verifyWebhookSignature(payload, secret, signature string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(payload))
    expectedSignature := hex.EncodeToString(mac.Sum(nil))
    
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func sendWebhookWithSignature(client *resty.Client, url, payload, secret string) error {
    // 计算签名
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(payload))
    signature := hex.EncodeToString(mac.Sum(nil))
    
    // 发送请求
    _, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetHeader("X-Webhook-Signature", signature).
        SetBody(payload).
        Post(url)
    
    return err
}
```

### 示例：请求日志中间件

```go
func logRequests(client *resty.Client) {
    client.OnRequestLog(func(r *resty.RequestLog) error {
        fmt.Printf("请求: %s %s\n", r.Method, r.URL)
        return nil
    })
    
    client.OnResponseLog(func(r *resty.ResponseLog) error {
        fmt.Printf("响应: %d %s\n", r.StatusCode, r.Status)
        return nil
    })
}
```

## 最佳实践

1. **复用客户端**：创建一次客户端并复用，不要每次请求都创建新客户端

2. **设置超时**：始终设置合理的请求超时

3. **使用重试**：对于可能失败的请求启用自动重试

4. **错误处理**：正确处理 HTTP 错误和非 2xx 状态码

5. **使用上下文**：对于长时间运行的操作使用 context.Context

6. **连接池配置**：根据应用负载合理配置连接池大小

7. **使用签名**：对于需要认证的 API 使用请求签名

8. **响应缓存**：对于不经常变化的数据启用缓存

9. **调试模式**：在开发环境启用调试模式排查问题

10. **使用结构体**：使用结构体和 JSON 标签处理响应，提高可读性

## 相关文档

- [Song 框架文档](../../README.md)
- [HTTP 服务器](../https/README.md)
- [配置管理](../vipers/README.md)

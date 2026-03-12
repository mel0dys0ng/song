# Erlogs 日志框架

一个功能强大的结构化日志框架，基于 Zap 构建，提供结构化日志记录、错误处理、上下文传播和多种日志级别支持。该框架简化了日志记录的实现，提供了链式 API 和丰富的配置选项。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [初始化日志](#初始化日志)
  - [日志级别](#日志级别)
  - [结构化日志](#结构化日志)
  - [错误处理](#错误处理)
  - [上下文传播](#上下文传播)
- [配置选项](#配置-options)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **结构化日志**：使用 JSON 格式记录结构化日志
- **多种级别**：支持 Debug、Info、Warn、Error、Fatal、Panic 级别
- **错误包装**：支持错误包装和堆栈跟踪
- **上下文支持**：支持 context.Context 传播
- **字段记录**：支持添加自定义字段
- **日志轮转**：支持日志文件轮转
- **高性能**：基于 Zap 提供高性能日志记录
- **开发模式**：支持开发模式，更易读的日志输出

## 安装

安装日志框架依赖：

```bash
go get go.uber.org/zap
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "go.uber.org/zap"
)

func main() {
    // 初始化日志
    erlogs.Initialize()
    
    ctx := context.Background()
    
    // 记录日志
    erlogs.Info(ctx, "应用程序启动",
        zap.String("app", "myapp"),
        zap.Int("port", 8080),
    )
    
    // 记录错误
    erlogs.Error(ctx, "发生错误",
        zap.String("error", "something went wrong"),
    )
}
```

## 架构概览

Erlogs 日志框架提供了统一的日志接口：

```
┌─────────────────┐
│    Erlogs        │ - 日志主接口
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options  │ - 配置选项
    └────┬────┘
         │
    ┌────┴────┐
    │   Zap    │ - 底层日志库
    └────┬────┘
         │
    ┌────┴────┐
    │ Encoder  │ - 日志编码器
    └──────────┘
```

**核心组件：**
- **Logger**：日志记录器
- **Options**：配置选项
- **Field**：日志字段

## 使用指南

### 初始化日志

初始化日志系统：

```go
// 基本初始化
erlogs.Initialize()

// 开发模式初始化
erlogs.Initialize(
    erlogs.WithDevelopment(true),
    erlogs.WithLevel(erlogs.DebugLevel),
)

// 生产模式初始化
erlogs.Initialize(
    erlogs.WithDevelopment(false),
    erlogs.WithLevel(erlogs.InfoLevel),
    erlogs.WithOutput("/var/log/myapp.log"),
)
```

### 日志级别

使用不同的日志级别：

```go
ctx := context.Background()

// Debug 级别 - 调试信息
erlogs.Debug(ctx, "调试信息",
    zap.String("key", "value"),
)

// Info 级别 - 一般信息
erlogs.Info(ctx, "应用程序启动",
    zap.Int("port", 8080),
)

// Warn 级别 - 警告信息
erlogs.Warn(ctx, "配置缺失，使用默认值",
    zap.String("key", "timeout"),
    zap.Any("default", 30),
)

// Error 级别 - 错误信息
erlogs.Error(ctx, "数据库连接失败",
    zap.String("error", err.Error()),
)

// Fatal 级别 - 致命错误
erlogs.Fatal(ctx, "无法启动服务器",
    zap.String("error", err.Error()),
)

// Panic 级别 - 恐慌信息
erlogs.Panic(ctx, "发生严重错误",
    zap.String("error", err.Error()),
)
```

### 结构化日志

使用结构化日志记录信息：

```go
ctx := context.Background()

// 基础字段
erlogs.Info(ctx, "用户登录",
    zap.String("user_id", "12345"),
    zap.String("username", "john_doe"),
    zap.String("ip", "192.168.1.1"),
)

// 复杂类型
erlogs.Info(ctx, "请求处理",
    zap.Int("status_code", 200),
    zap.Duration("duration", time.Millisecond*150),
    zap.Float64("cpu_usage", 0.75),
)

// 结构体
erlogs.Info(ctx, "用户信息",
    zap.Stringer("user", user),
)

// 嵌套结构
erlogs.Info(ctx, "订单创建",
    zap.Int("order_id", 12345),
    zap.String("customer_name", "张三"),
    zap.Any("items", []string{"item1", "item2"}),
)

// 使用字段选项
erlogs.Info(ctx, "用户操作",
    erlogs.OptionFields(
        zap.String("user_id", "12345"),
        zap.String("action", "create"),
    ),
)
```

### 错误处理

使用错误处理功能：

```go
ctx := context.Background()

// 简单错误记录
err := someOperation()
if err != nil {
    erlogs.Error(ctx, "操作失败",
        zap.Error(err),
    )
}

// 错误包装
err := someOperation()
if err != nil {
    erlogs.Convert(err).
        Wrap("处理用户请求失败").
        ErrorLog(ctx,
            zap.String("user_id", "12345"),
        )
}

// 错误包装链式调用
err = step1()
if err != nil {
    erlogs.Convert(err).
        Wrap("步骤1失败").
        ErrorLog(ctx)
}

err = step2()
if err != nil {
    erlogs.Convert(err).
        Wrap("步骤2失败").
        Options(BaseELOptions()).
        ErrorLog(ctx,
            zap.String("user_id", "12345"),
        )
}
```

### 上下文传播

使用 context 传播日志上下文：

```go
// 从 context 创建带日志的 context
ctx := context.Background()
logCtx := erlogs.WithContext(ctx, 
    zap.String("request_id", "req-123"),
    zap.String("trace_id", "trace-456"),
)

// 在新的 context 中记录日志
erlogs.Info(logCtx, "处理请求",
    zap.String("endpoint", "/api/users"),
)
```

## 配置 Options

### Option 函数

```go
// 开发模式
erlogs.WithDevelopment(true)  // 开发模式

// 日志级别
erlogs.WithLevel(erlogs.DebugLevel)
erlogs.WithLevel(erlogs.InfoLevel)
erlogs.WithLevel(erlogs.WarnLevel)
erlogs.WithLevel(erlogs.ErrorLevel)

// 输出目标
erlogs.WithOutput("/var/log/myapp.log")
erlogs.WithStdout()

// 编码器
erlogs.WithJSONEncoder()
erlogs.WithConsoleEncoder()

// 额外字段
erlogs.AddCaller(true)
erlogs.AddStacktrace(true)
```

### 日志级别常量

```go
const (
    DebugLevel  zapcore.Level = zapcore.DebugLevel
    InfoLevel   zapcore.Level = zapcore.InfoLevel
    WarnLevel   zapcore.Level = zapcore.WarnLevel
    ErrorLevel  zapcore.Level = zapcore.ErrorLevel
    FatalLevel  zapcore.Level = zapcore.FatalLevel
    PanicLevel  zapcore.Level = zapcore.PanicLevel
)
```

## 示例代码

### 完整示例：Web 应用日志

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    "github.com/mel0dys0ng/song/internal/core/erlogs"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func main() {
    // 初始化日志
    erlogs.Initialize(
        erlogs.WithDevelopment(true),
        erlogs.WithLevel(erlogs.InfoLevel),
    )
    
    // 创建 Gin 引擎
    r := gin.Default()
    
    // 中间件：请求日志
    r.Use(func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        latency := time.Since(start)
        status := c.Writer.Status()
        
        erlogs.Info(c.Request.Context(), "HTTP 请求",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", status),
            zap.Duration("latency", latency),
            zap.String("ip", c.ClientIP()),
        )
    })
    
    // 路由
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "time":   time.Now().Format(time.RFC3339),
        })
    })
    
    r.GET("/users/:id", getUser)
    r.POST("/users", createUser)
    
    // 启动服务器
    erlogs.Info(context.Background(), "启动服务器",
        zap.Int("port", 8080),
    )
    r.Run(":8080")
}

func getUser(c *gin.Context) {
    ctx := c.Request.Context()
    id := c.Param("id")
    
    erlogs.Info(ctx, "获取用户",
        zap.String("user_id", id),
    )
    
    // 模拟错误
    err := fmt.Errorf("用户不存在")
    erlogs.Error(ctx, "获取用户失败",
        zap.Error(err),
        zap.String("user_id", id),
    )
    
    c.JSON(http.StatusNotFound, gin.H{
        "error": err.Error(),
    })
}

func createUser(c *gin.Context) {
    ctx := c.Request.Context()
    
    // 记录请求体
    erlogs.Debug(ctx, "创建用户请求",
        zap.String("body", c.Request.Body),
    )
    
    // 模拟处理
    userID := "12345"
    
    erlogs.Info(ctx, "创建用户成功",
        zap.String("user_id", userID),
    )
    
    c.JSON(http.StatusCreated, gin.H{
        "id":   userID,
        "name": "新用户",
    })
}
```

### 示例：错误处理中间件

```go
func errorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // 检查是否有错误
        if len(c.Errors) > 0 {
            for _, e := range c.Errors {
                erlogs.Error(c.Request.Context(), "请求错误",
                    zap.Error(e.Err),
                    zap.String("path", c.Request.URL.Path),
                )
            }
        }
    }
}
```

### 示例：业务日志包装

```go
type Logger struct {
    ctx context.Context
}

func NewLogger(ctx context.Context) *Logger {
    return &Logger{ctx: ctx}
}

func (l *Logger) LogUserAction(action string, userID string) {
    erlogs.Info(l.ctx, "用户操作",
        zap.String("action", action),
        zap.String("user_id", userID),
    )
}

func (l *Logger) LogOrderCreated(orderID string, amount float64) {
    erlogs.Info(l.ctx, "订单创建",
        zap.String("order_id", orderID),
        zap.Float64("amount", amount),
    )
}

func (l *Logger) LogError(err error, message string) {
    erlogs.Convert(err).
        Wrap(message).
        ErrorLog(l.ctx)
}

// 使用
func handleRequest(ctx context.Context) {
    logger := NewLogger(ctx)
    
    logger.LogUserAction("login", "user123")
    logger.LogOrderCreated("order456", 99.99)
}
```

### 示例：性能日志

```go
func logPerformance(ctx context.Context, operation string, start time.Time) {
    duration := time.Since(start)
    
    if duration > time.Second {
        erlogs.Warn(ctx, "操作耗时过长",
            zap.String("operation", operation),
            zap.Duration("duration", duration),
        )
    } else {
        erlogs.Debug(ctx, "操作完成",
            zap.String("operation", operation),
            zap.Duration("duration", duration),
        )
    }
}

// 使用
func someOperation() {
    start := time.Now()
    defer logPerformance(context.Background(), "someOperation", start)
    
    // 执行操作
    time.Sleep(100 * time.Millisecond)
}
```

### 示例：结构化错误

```go
type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

func handleError(ctx context.Context, err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        erlogs.Error(ctx, "应用错误",
            zap.String("code", appErr.Code),
            zap.String("message", appErr.Message),
            zap.Error(appErr.Err),
        )
    } else {
        erlogs.Error(ctx, "未知错误",
            zap.Error(err),
        )
    }
}
```

## 最佳实践

1. **使用结构化日志**：始终使用结构化日志，便于查询和分析

2. **添加上下文**：在日志中添加相关上下文信息

3. **选择合适的级别**：根据情况选择合适的日志级别

4. **错误包装**：使用错误包装保留错误链

5. **避免敏感信息**：不要在日志中记录密码、令牌等敏感信息

6. **使用 context**：使用 context 传播请求级别的信息

7. **开发/生产区分**：在开发环境使用更详细的日志

8. **性能考虑**：避免在日志中执行复杂操作

9. **日志轮转**：配置日志轮转，避免日志文件过大

10. **统一格式**：保持日志格式一致性

## 相关文档

- [Song 框架文档](../../README.md)
- [HTTP 服务器](../https/README.md)
- [元数据管理](../metas/README.md)

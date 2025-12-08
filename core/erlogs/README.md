# erlogs

`erlogs` 是一个集错误处理和日志记录于一体的 Go 组件库，采用链式错误处理机制，支持分布式追踪、结构化日志等功能。通过将错误处理和日志记录有机结合，提供了一套完整的错误管理和追踪解决方案。其链式设计使得错误信息更加丰富和可追溯，而多层次 API 设计则兼顾了灵活性和易用性。这种设计模式特别适合复杂的微服务架构，能够有效提升系统的可观测性和问题排查效率。

## 核心设计理念

### 1. 链式错误处理机制

`erlogs` 采用了链式错误处理的设计思想，允许开发者在现有错误基础上包装新的错误信息，形成错误链条：

```go
// 业务Code和Name

var (
    UserBiz = erlogs.Biz(10001, 'user')
    PassportBiz = erlogs.Biz(10002, 'passport')
)
```

```go
// 基础错误

var (
    // 未定义业务归属
    BaseEL = erlogs.New(erlogs.Log(true), erlogs.TypeBiz())
    // User业务基础错误
    UserEL = erlogs.New(erlogs.Log(true), erlogs.TypeBiz(), UserBiz)
    // Passport业务基础错误
    PassportEL = erlogs.New(erlogs.Log(true), erlogs.TypeBiz(), PassportBiz)
)
```

```go
// 错误码&错误提示（返回给用户的错误信息）

var (
    InvalidArguments        = BaseEL.WithStatus(40001, "请求参数错误")
	TooFrequentOperation    = BaseEL.WithStatus(40002, "操作太频繁,请稍后重试")
	Unauthorized            = BaseEL.WithStatus(40003, "请登录")
	ServerExcept            = BaseEL.WithStatus(50000, "服务异常，请稍后重试")
)
```

```go
if err = validator.New().Struct(request); err != nil {
    // 未指定业务
    // err = status.InvalidArguments.Info(ctx, erlogs.ValidateError(request, err))
    // 指定业务
	err = status.InvalidArguments.Info(ctx, PassportBiz, erlogs.ValidateError(request, err))
	return
}
```

```go
// 自定义错误内容
err0 := status.ServerExcept.Error(ctx, erlogs.Content("xxxxxxxxxxx"),
    erlogs.Fields(zap.String("request", request)),
)

// 基于已有错误内容
err1 := status.ServerExcept.Warn(ctx, erlogs.ContentError(err),
    erlogs.Fields(zap.String("request", request)),
)

// 链式错误处理
err2 := err1.Erorr(ctx, PassportBiz, erlogs.Fields(
    zap.String("request", request),
    zap.Uint64("uid", userId),
))

// 链式错误处理
err3 := err1.Erorr(ctx, UserBiz, erlogs.Fields(
    zap.String("request", request),
    zap.Uint64("uid", userId),
    zap.String("username", username),
))
```

### 2. 错误即日志对象

在 `erlogs` 中，错误对象本身就是日志记录器，具备双重身份：

- 作为错误信息载体，包含错误码、消息、调用栈等
- 作为日志记录器，使用 zap 库实现高性能结构化日志记录

### 3. 不可变性设计

每个错误对象都是不可变的，任何修改操作都会产生新的实例，确保并发安全和状态隔离。

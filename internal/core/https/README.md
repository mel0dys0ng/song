# Song 框架 HTTPS 模块

Song 框架的 HTTPS 模块是基于 Gin 框架构建的 HTTP 服务组件，提供了完整的 Web 服务功能，包括安全特性、日志记录、客户端信息识别等功能。

## 主要特性

- **基于 Gin 框架**：充分利用 Gin 框架的高性能和灵活性
- **安全中间件**：支持 CORS、CSRF 和请求签名验证
- **客户端信息识别**：自动解析用户代理、操作系统、浏览器和设备信息
- **结构化日志**：集成 erlogs 实现结构化日志记录
- **分布式追踪**：内置追踪 ID 支持分布式系统监控
- **优雅关闭**：支持服务优雅启动和关闭

## 核心功能

### 1. 服务器配置

HTTPS 模块支持多种服务器配置选项，包括：

- **基本配置**：端口、TLS 设置、超时配置等
- **连接管理**：Keep-Alive、最大请求头大小等
- **日志配置**：日志级别、存储位置、轮转策略等

#### 配置模板

```yaml
https:
  port: 8080 # 服务器监听端口
  TLSOpen: false # 是否开启TLS
  TLSKeyFile: "" # TLS私钥文件路径
  TLSCertFile: "" # TLS证书文件路径
  keepAlive: true # 是否启用Keep-Alive
  readTimeout: "30s" # 读取超时时间
  readHeaderTimeout: "30s" # 读取头部超时时间
  writeTimeout: "60s" # 写入超时时间
  idleTimeout: "60s" # 空闲超时时间
  hammerTime: "30s" # 优雅关闭等待时间
  maxHeaderBytes: 1048576 # 最大请求头大小(字节)
  tmpDir: "./tmp" # 临时目录
  loggerHeaderKeys: # 需要在日志中记录的请求头键
    - "User-Agent"
    - "Accept-Encoding"

  # 日志配置
  erlog:
    dir: "logs" # 日志存储目录
    fileName: "https.log" # 日志文件名
    level: "info" # 日志级别(debug/info/warn/error/panic/fatal)
    maxSize: 100 # 单个日志文件最大大小(MB)
    maxBackups: 10 # 日志备份文件数量
    maxAge: 30 # 日志保留天数
    compress: false # 是否压缩日志

  # CORS配置
  cors:
    enable: false # 是否启用CORS
    allowOrigins: # 允许的来源域名
      - "http://localhost:3000"
      - "https://yourdomain.com"
    allowHeaders: # 允许的请求头
      - "Origin"
      - "Content-Length"
      - "Content-Type"
      - "Authorization"
    allowMethods: # 允许的HTTP方法
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
      - "OPTIONS"
    exposeHeaders: # 允许客户端访问的响应头
      - "Content-Length"
      - "Access-Control-Allow-Origin"
    allowCredentials: true # 是否允许携带凭据(cookie等)
    allowWildcard: false # 是否允许通配符匹配
    maxAge: "12h" # 预检请求缓存时间

  # CSRF配置
  csrf:
    enable: false # 是否启用CSRF保护
    lookupType: "header" # 查找类型(header/form/query)
    lookupName: "X-Song-Csrf-Token" # 查找名称
    cookieName: "X-Song-Csrf-Token" # CSRF Cookie名称
    cookieDomain: "" # CSRF Cookie域名
    cookiePath: "/" # CSRF Cookie路径
    cookieMaxAge: 3600 # CSRF Cookie最大存活时间(秒)
    cookieSecure: false # CSRF Cookie是否只在HTTPS下传输
    cookieHttpOnly: true # CSRF Cookie是否禁用JS访问

  # 签名验证配置
  sign:
    enable: false # 是否启用签名验证
    secret: "default_sign_secret" # 签名密钥
    ttl: 300 # 签名有效期(秒)
    query: true # 是否验证查询参数
    formData: true # 是否验证表单数据
    header: true # 是否验证请求头
```

### 2. 安全中间件

#### CORS 中间件

- 控制跨域资源共享策略
- 支持自定义允许的来源、请求头、方法等
- 可配置凭证传递和通配符匹配

#### CSRF 中间件

- 防止跨站请求伪造攻击
- 支持从请求头、表单或查询参数中验证令牌
- 为安全方法（GET、HEAD、OPTIONS、TRACE）自动生成令牌

#### 签名验证中间件

- 验证请求的数字签名
- 支持查询参数、表单数据和请求头的签名验证
- 包含时间戳验证，防止重放攻击

### 3. 客户端信息识别

客户端信息中间件能够自动识别并记录以下信息：

- **IP 地址**：通过 X-Forwarded-For、X-Real-Ip 或直接连接获取客户端 IP
- **设备 ID**：支持多种设备 ID 来源（请求头、查询参数、表单字段）
- **用户代理**：解析操作系统、浏览器类型和版本
- **客户端类型**：区分 Web 应用、移动应用和桌面应用
- **版本信息**：支持解析客户端版本号和构建号

#### 操作系统检测

- Windows
- macOS
- iOS
- Android
- Linux
- Unix

#### 浏览器检测

- Chrome / Chrome Mobile
- Safari / Safari Mobile
- Firefox
- Edge
- Internet Explorer
- Opera
- UC Browser
- Samsung Browser

#### ClientInfo 结构体方法

ClientInfo 结构体提供了全面的方法来访问客户端信息，包括安全的 getter 方法，确保在 ClientInfo 为 nil 时返回零值：

- **IP 相关**：
  - `GetIP()` - 获取客户端 IP 地址
  - `GetDeviceID()` - 获取设备 ID
  - `GetUserAgent()` - 获取用户代理字符串
- **请求信息**：
  - `GetMethod()` - 获取请求方法
  - `GetPath()` - 获取请求路径
- **操作系统信息**：
  - `GetOS()` - 获取操作系统名称
  - `GetOSVersion()` - 获取操作系统版本
- **浏览器信息**：
  - `GetBrowser()` - 获取浏览器名称
  - `GetBrowserVersion()` - 获取浏览器版本
- **客户端类型和版本**：
  - `GetClientType()` - 获取客户端类型 (Web, Mobile App, Desktop App)
  - `GetClientVersion()` - 获取客户端版本号
  - `GetAppBuild()` - 获取应用构建号
- **类型判断方法**：
  - `IsWindows()` / `IsMacOS()` / `IsAndroid()` / `IsIOS()` - 操作系统类型判断
  - `IsChrome()` / `IsSafari()` / `IsFirefox()` / `IsEdge()` - 浏览器类型判断
  - `IsMobileOS()` / `IsDesktopOS()` - 操作系统分类判断
  - `IsWebClient()` / `IsMobileApp()` / `IsDesktopApp()` - 客户端类型判断
- **版本比较方法**：
  - `CompareVersion(target string)` - 比较客户端版本
  - `IsVersionEqual(target string)` - 检查版本是否相等
  - `IsVersionGreater(target string)` - 检查版本是否大于指定版本
  - `IsVersionGreaterOrEqual(target string)` - 检查版本是否大于或等于指定版本
  - `IsVersionLess(target string)` - 检查版本是否小于指定版本
  - `IsVersionLessOrEqual(target string)` - 检查版本是否小于或等于指定版本

### 4. 响应处理

模块提供了标准化的响应处理机制：

- **成功响应**：使用 `ResponseSuccess` 函数返回成功结果
- **错误响应**：使用 `ResponseError` 函数返回错误信息
- **自定义响应**：通过 `ResponseWithStatus` 实现自定义响应
- **多种格式**：支持 JSON、JSONP、HTML、流式等多种响应格式

### 5. 日志记录

- **结构化日志**：记录请求的详细信息，包括时间、IP、状态码、路径等
- **性能监控**：记录请求耗时，帮助性能分析
- **错误追踪**：集成追踪 ID，便于错误定位
- **自定义字段**：支持记录特定请求头信息

### 6. 分布式追踪

- **请求追踪**：为每个请求分配唯一的追踪 ID
- **链路追踪**：支持分布式系统的请求链路追踪

## 使用方法

### 基本服务器启动

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/https"
)

func main() {
    server := https.New([]https.Option{})
    server.Serve()
}
```

### 自定义配置启动

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/https"
)

func main() {
    server := https.New([]https.Option{
        https.Port(9090),
        https.TLS(false),
    })
    server.Serve()
}
```

### 添加路由和中间件

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/mel0dys0ng/song/internal/core/https"
)

func main() {
    server := https.New([]https.Option{
        https.Route(func(eng *gin.Engine) {
            eng.GET("/api/hello", func(ctx *gin.Context) {
                https.ResponseSuccess(ctx, "Hello, World!")
            })
        }),
        https.Middleware(https.Middleware{
            Priority: 1,
            Handle: func(eng *gin.Engine) gin.HandlerFunc {
                return func(ctx *gin.Context) {
                    // 自定义中间件逻辑
                    ctx.Next()
                }
            },
        }),
    })

    server.Serve()
}
```

### 获取客户端信息

```go
func handler(ctx *gin.Context) {
    clientInfo := https.GetClientInfo(ctx)

    // 使用客户端信息
    fmt.Println("IP:", clientInfo.GetIP())
    fmt.Println("OS:", clientInfo.GetOS())
    fmt.Println("Browser:", clientInfo.GetBrowser())
    fmt.Println("Device ID:", clientInfo.GetDeviceID())

    // 版本比较
    if clientInfo.IsVersionGreaterOrEqual("1.2.0") {
      // 特定版本逻辑
    }

    https.ResponseSuccess(ctx, "OK")
}
```

## API 文档

### 服务器配置选项

- `Port(int)` - 设置服务器端口
- `Host(string)` - 设置服务器主机
- `TLS(bool)` - 启用/禁用 TLS
- `TLSKeyFile(string)` - 设置 TLS 私钥文件
- `TLSCertFile(string)` - 设置 TLS 证书文件
- `CORS(*Cors)` - 配置 CORS
- `CSRF(*CSRF)` - 配置 CSRF
- `Sign(*Sign)` - 配置签名验证
- `Route(Route)` - 添加路由
- `Middleware(Middleware)` - 添加中间件

### 响应函数

- `Response(ctx, data, err, opts...)` - 通用响应函数
- `ResponseSuccess(ctx, data, opts...)` - 成功响应
- `ResponseError(ctx, err, opts...)` - 错误响应
- `ResponseWithStatus(ctx, status, opts...)` - 自定义状态响应

### 客户端信息方法

- `GetClientInfo(ctx)` - 获取客户端信息
- `IsWindows()` / `IsMacOS()` / `IsAndroid()` / `IsIOS()` - 操作系统检测
- `IsChrome()` / `IsSafari()` / `IsFirefox()` / `IsEdge()` - 浏览器检测
- `IsWebClient()` / `IsMobileApp()` / `IsDesktopApp()` - 客户端类型检测
- `CompareVersion(targetVersion)` - 版本比较
- `IsVersionGreater(targetVersion)` / `IsVersionLess(targetVersion)` - 版本比较便捷方法

## 性能特点

- **中间件排序**：支持中间件优先级排序
- **并发优化**：使用读写锁优化高并发场景
- **资源管理**：支持优雅关闭，确保资源正确释放
- **内存优化**：使用对象池减少 GC 压力

## 安全特性

- **CSRF 保护**：防止跨站请求伪造
- **签名验证**：确保请求完整性
- **CORS 控制**：限制跨域访问
- **输入验证**：对用户输入进行验证

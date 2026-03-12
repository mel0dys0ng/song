# HTTPS 服务器

一个基于 Gin 框架构建的 HTTP 服务器，提供简洁的 API、中间件支持、安全特性和生命周期钩子。该服务器简化了 HTTP 服务的创建和配置，提供了声明式的路由定义和丰富的配置选项。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建服务器](#创建服务器)
  - [定义路由](#定义路由)
  - [中间件](#中间件)
  - [安全特性](#安全特性)
  - [生命周期钩子](#生命周期钩子)
- [配置选项](#配置选项)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **Gin 框架集成**：基于高性能的 Gin 框架
- **中间件支持**：支持 CORS、CSRF、请求签名等中间件
- **路由分组**：支持路由分组和嵌套
- **TLS/HTTPS**：支持 TLS/HTTPS 配置
- **生命周期钩子**：支持服务器启动和停止钩子
- **请求验证**：支持请求体验证
- **上下文支持**：完全支持 context.Context
- **错误处理**：统一的错误处理机制

## 安装

安装 Gin 依赖：

```bash
go get github.com/gin-gonic/gin
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/gin-gonic/gin"
)

func main() {
    // 创建 HTTP 服务器
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

    // 启动服务器
    server.Serve()
}
```

运行：

```bash
go run main.go
# 访问 http://localhost:8080/health
```

## 架构概览

HTTPS 服务器提供了简洁的 HTTP 服务接口：

```
┌─────────────────┐
│   HTTPS Server   │ - 主服务器接口
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options │ - 配置选项
    └────┬────┘
         │
    ┌────┴────────────┐
    │     Gin Engine   │ - Gin 引擎
    └─────────────────┘
         │
    ┌────┴────────────┐
    │  Middleware      │ - 中间件
    └─────────────────┘
```

**核心组件：**

- **Server**：主服务器结构
- **Options**：服务器配置选项
- **Engine**：Gin 引擎

## 使用指南

### 创建服务器

创建一个新的 HTTP 服务器实例：

```go
import "github.com/mel0dys0ng/song/internal/core/https"

// 基本服务器
server := https.New([]https.Option{
    https.Port(8080),
})

// 高级配置
server := https.New([]https.Option{
    https.Port(8080),
    https.Host("0.0.0.0"),
    https.ReadTimeout(30),
    https.WriteTimeout(30),
    https.IdleTimeout(60),
})
```

### 定义路由

定义不同的路由：

```go
server := https.New([]https.Option{
    https.Port(8080),
    https.Route(func(eng *gin.Engine) {
        // GET 路由
        eng.GET("/users", getUsers)
        eng.GET("/users/:id", getUser)

        // POST 路由
        eng.POST("/users", createUser)

        // PUT 路由
        eng.PUT("/users/:id", updateUser)

        // DELETE 路由
        eng.DELETE("/users/:id", deleteUser)

        // 路由分组
        v1 := eng.Group("/api/v1")
        {
            v1.GET("/users", getUsers)
            v1.POST("/users", createUser)
        }
    }),
})
```

### 中间件

使用各种中间件：

```go
server := https.New([]https.Option{
    https.Port(8080),
    https.Route(func(eng *gin.Engine) {
        // 使用中间件
        eng.Use(gin.Recovery())
        eng.Use(gin.Logger())

        // CORS 中间件
        eng.Use(cors.New(cors.Config{
            AllowOrigins:     []string{"https://example.com"},
            AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
            AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
            ExposeHeaders:    []string{"Content-Length"},
            AllowCredentials: true,
        }))

        // 自定义中间件
        eng.Use(func(c *gin.Context) {
            // 处理请求
            c.Next()
            // 处理响应
        })

        // 定义路由
        eng.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })
    }),
})
```

### 安全特性

配置安全特性：

```go
server := https.New([]https.Option{
    https.Port(8080),
    https.TLSCertFile("/path/to/cert.pem"),
    https.TLSKeyFile("/path/to/key.pem"),
    https.Route(func(eng *gin.Engine) {
        // CSRF 保护
        eng.Use(csrf.New())

        // 请求签名验证
        eng.Use(signature.Verify())

        // 定义路由
        eng.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })
    }),
})
```

### 生命周期钩子

使用生命周期钩子：

```go
server := https.New([]https.Option{
    https.Port(8080),
    https.BeforeStart(func() error {
        println("服务器即将启动")
        return nil
    }),
    https.AfterStart(func() error {
        println("服务器已启动")
        return nil
    }),
    https.BeforeStop(func() error {
        println("服务器即将停止")
        return nil
    }),
    https.AfterStop(func() error {
        println("服务器已停止")
        return nil
    }),
    https.Route(func(eng *gin.Engine) {
        eng.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })
    }),
})
```

## 配置 Options

### Option 函数

```go
// 服务器配置
https.Port(8080)                    // 端口
https.Host("0.0.0.0")               // 主机地址
https.ReadTimeout(30)               // 读取超时（秒）
https.WriteTimeout(30)              // 写入超时（秒）
https.IdleTimeout(60)               // 空闲超时（秒）

// TLS 配置
https.TLSCertFile("/path/to/cert") // TLS 证书文件
https.TLSKeyFile("/path/to/key")   // TLS 密钥文件

// 路由配置
https.Route(func(eng *gin.Engine) {
    // 路由定义
})

// 生命周期钩子
https.BeforeStart(func() error { return nil }) // 启动前
https.AfterStart(func() error { return nil })  // 启动后
https.BeforeStop(func() error { return nil })  // 停止前
https.AfterStop(func() error { return nil })   // 停止后
```

## 示例代码

### 完整示例：RESTful API

```go
package main

import (
    "net/http"
    "github.com/mel0dys0ng/song/internal/core/https"
    "github.com/mel0dys0ng/song/internal/core/metas"
    "github.com/gin-gonic/gin"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "张三", Email: "zhangsan@example.com"},
    {ID: 2, Name: "李四", Email: "lisi@example.com"},
}

func main() {
    // 初始化元数据
    metas.New(&metas.Options{
        App:  "user-api",
        Kind: metas.KindAPI,
        Mode: metas.ModeLocal,
    })

    // 创建服务器
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            // 中间件
            eng.Use(gin.Recovery())
            eng.Use(gin.Logger())

            // 健康检查
            eng.GET("/health", healthCheck)

            // 用户路由
            usersGroup := eng.Group("/users")
            {
                usersGroup.GET("", listUsers)
                usersGroup.GET(":id", getUser)
                usersGroup.POST("", createUser)
                usersGroup.PUT(":id", updateUser)
                usersGroup.DELETE(":id", deleteUser)
            }
        }),
    })

    // 启动服务器
    server.Serve()
}

func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "app":    metas.Metadata().App(),
    })
}

func listUsers(c *gin.Context) {
    c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
    id := c.Param("id")

    for _, user := range users {
        if string(rune(user.ID)) == id {
            c.JSON(http.StatusOK, user)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{
        "error": "用户不存在",
    })
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    user.ID = len(users) + 1
    users = append(users, user)

    c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
    id := c.Param("id")
    var updatedUser User

    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    for i, user := range users {
        if string(rune(user.ID)) == id {
            updatedUser.ID = user.ID
            users[i] = updatedUser
            c.JSON(http.StatusOK, updatedUser)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{
        "error": "用户不存在",
    })
}

func deleteUser(c *gin.Context) {
    id := c.Param("id")

    for i, user := range users {
        if string(rune(user.ID)) == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(http.StatusNoContent, nil)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{
        "error": "用户不存在",
    })
}
```

### 示例：带中间件的 API

```go
func main() {
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            // 全局中间件
            eng.Use(gin.Recovery())
            eng.Use(authMiddleware())
            eng.Use(loggingMiddleware())

            // 公开路由
            eng.GET("/public/hello", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "你好，世界！"})
            })

            // 受保护路由
            protected := eng.Group("/api")
            protected.Use(gin.BasicAuth(gin.Accounts{
                "admin": "password123",
            }))
            {
                protected.GET("/secret", func(c *gin.Context) {
                    c.JSON(200, gin.H{"secret": "这是机密信息"})
                })
            }
        }),
    })

    server.Serve()
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "缺少认证令牌",
            })
            return
        }
        c.Next()
    }
}

func loggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 处理请求
        c.Next()
        // 记录响应
    }
}
```

### 示例：文件上传下载

```go
func main() {
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            // 静态文件服务
            eng.Static("/static", "./static")

            // 文件上传
            eng.POST("/upload", func(c *gin.Context) {
                file, err := c.FormFile("file")
                if err != nil {
                    c.JSON(http.StatusBadRequest, gin.H{
                        "error": err.Error(),
                    })
                    return
                }

                err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
                if err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{
                        "error": err.Error(),
                    })
                    return
                }

                c.JSON(http.StatusOK, gin.H{
                    "message":  "文件上传成功",
                    "filename": file.Filename,
                    "filesize": file.Size,
                })
            })

            // 文件下载
            eng.GET("/download/:filename", func(c *gin.Context) {
                filename := c.Param("filename")
                c.Header("Content-Description", "File Transfer")
                c.Header("Content-Transfer-Encoding", "binary")
                c.Header("Content-Disposition", "attachment; filename="+filename)
                c.Header("Content-Type", "application/octet-stream")
                c.File("./uploads/" + filename)
            })
        }),
    })

    server.Serve()
}
```

### 示例：WebSocket 支持

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func main() {
    server := https.New([]https.Option{
        https.Port(8080),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/ws", func(c *gin.Context) {
                conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
                if err != nil {
                    return
                }
                defer conn.Close()

                for {
                    messageType, message, err := conn.ReadMessage()
                    if err != nil {
                        break
                    }

                    if messageType == websocket.TextMessage {
                        conn.WriteMessage(messageType, message)
                    }
                }
            })
        }),
    })

    server.Serve()
}
```

### 示例：Graceful Shutdown

```go
func main() {
    server := https.New([]https.Option{
        https.Port(8080),
        https.BeforeStop(func() error {
            println("开始关闭服务器...")
            return nil
        }),
        https.AfterStop(func() error {
            println("服务器已完全关闭")
            return nil
        }),
        https.Route(func(eng *gin.Engine) {
            eng.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{"status": "ok"})
            })
        }),
    })

    server.Serve()
}
```

## 最佳实践

1. **使用中间件**：使用中间件处理通用逻辑（认证、日志、CORS 等）

2. **错误处理**：正确处理错误，返回适当的 HTTP 状态码

3. **验证输入**：始终验证用户输入，使用结构体绑定

4. **安全配置**：在生产环境中使用 HTTPS

5. **超时设置**：设置合理的超时时间

6. **优雅关闭**：使用生命周期钩子实现优雅关闭

7. **路由分组**：使用路由分组组织相关路由

8. **日志记录**：记录请求日志便于调试

9. **CORS 配置**：根据需求配置 CORS

10. **性能优化**：使用连接池和缓存提高性能

## 相关文档

- [Song 框架文档](../../README.md)
- [日志框架](../erlogs/README.md)
- [配置管理](../vipers/README.md)

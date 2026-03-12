# Cobras CLI 框架

一个基于 Cobra 构建的命令行界面（CLI）部署框架，提供命令层级、参数解析、帮助生成等功能。该框架简化了 CLI 应用的部署，提供了声明式的命令定义和丰富的配置选项。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建命令](#创建命令)
  - [添加子命令](#添加子命令)
  - [定义参数](#定义参数)
  - [添加动作](#添加动作)
- [配置选项](#配置选项)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **Cobra 集成**：基于成熟的 Cobra 库
- **命令层级**：支持多级命令嵌套
- **参数解析**：支持短选项、长选项和位置参数
- **帮助生成**：自动生成帮助信息
- **子命令支持**：支持多个子命令
- **版本信息**：内置版本命令
- **bash/zsh 补全**：支持 shell 自动补全
- **错误处理**：统一的错误处理机制

## 安装

安装 Cobra 依赖：

```bash
go get github.com/spf13/cobra
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
)

func main() {
    // 创建根命令
    cmd := cobras.NewCommand("myapp", "1.0.0", "我的应用程序")

    // 添加子命令
    cmd.AddCommand(&cobra.Command{
        Use:   "greet",
        Short: "问候某人",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("你好，世界！")
        },
    })

    // 执行
    cmd.Execute()
}
```

运行：

```bash
go run main.go greet
# 输出: 你好，世界！
```

## 架构概览

Cobras 框架提供了清晰的命令层级结构：

```
┌─────────────────┐
│  Root Command   │ - 根命令
└────────┬────────┘
         │
    ┌────┴─────────┐
    │  Subcommand  │ - 子命令
    └────┬─────────┘
         │
    ┌────┴────────────┐
    │ Flags/Args      │ - 参数和选项
    └─────────────────┘
```

**核心组件：**

- **Command**：命令结构
- **Options**：命令选项
- **Flag**：命令行标志

## 使用指南

### 创建命令

创建一个新的命令实例：

```go
import "github.com/mel0dys0ng/song/internal/core/cobras"

// 创建根命令
rootCmd := cobras.NewCommand(
    "myapp",           // 命令名称
    "1.0.0",          // 版本号
    "我的应用程序",    // 简短描述
)

// 或者使用完整选项
rootCmd := &cobra.Command{
    Use:   "myapp",
    Short: "我的应用程序",
    Long:  `一个功能强大的 CLI 应用程序`,
    Version: "1.0.0",
}
```

### 添加子命令

为命令添加子命令：

```go
// 添加简单子命令
rootCmd.AddCommand(&cobra.Command{
    Use:   "serve",
    Short: "启动服务",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("服务已启动")
    },
})

// 添加带配置的子命令
serveCmd := &cobra.Command{
    Use:   "serve",
    Short: "启动服务器",
    Long:  `启动 HTTP 服务器`,
}
serveCmd.Flags().IntP("port", "p", 8080, "端口号")
serveCmd.Flags().StringP("host", "H", "localhost", "主机地址")

serveCmd.Run = func(cmd *cobra.Command, args []string) {
    port, _ := cmd.Flags().GetInt("port")
    host, _ := cmd.Flags().GetString("host")
    fmt.Printf("启动服务器: %s:%d\n", host, port)
}

rootCmd.AddCommand(serveCmd)
```

### 定义参数

定义不同类型的参数：

```go
// 位置参数
cmd := &cobra.Command{
    Use:   "create <name>",
    Short: "创建用户",
    Args:  cobra.ExactArgs(1), // 必须提供 1 个参数
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        fmt.Printf("创建用户: %s\n", name)
    },
}

// 可选位置参数
cmd := &cobra.Command{
    Use:   "greet [name]",
    Short: "问候",
    Args:  cobra.RangeArgs(0, 1), // 0 或 1 个参数
    Run: func(cmd *cobra.Command, args []string) {
        name := "世界"
        if len(args) > 0 {
            name = args[0]
        }
        fmt.Printf("你好，%s！\n", name)
    },
}

// 选项标志
cmd.Flags().StringP("name", "n", "默认", "名称")
cmd.Flags().IntP("age", "a", 0, "年龄")
cmd.Flags().BoolP("verbose", "v", false, "详细输出")
cmd.Flags().StringSliceP("tags", "t", []string{}, "标签")
cmd.Flags().Int64P("id", "i", 0, "ID")
```

### 添加动作

为命令添加执行逻辑：

```go
// 简单动作
cmd := &cobra.Command{
    Use:   "hello",
    Short: "打招呼",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("你好！")
    },
}

// 带参数的动作
cmd := &cobra.Command{
    Use:   "greet",
    Short: "打招呼",
    PreRun: func(cmd *cobra.Command, args []string) {
        // 前置处理
        fmt.Println("准备中...")
    },
    Run: func(cmd *cobra.Command, args []string) {
        name, _ := cmd.Flags().GetString("name")
        fmt.Printf("你好，%s！\n", name)
    },
    PostRun: func(cmd *cobra.Command, args []string) {
        // 后置处理
        fmt.Println("完成！")
    },
}

// 异步动作
cmd := &cobra.Command{
    Use:   "async",
    Short: "异步任务",
    Run: func(cmd *cobra.Command, args []string) {
        done := make(chan bool)
        go func() {
            time.Sleep(2 * time.Second)
            done <- true
        }()
        <-done
        fmt.Println("任务完成")
    },
}
```

## 配置 Options

### Command 结构体

```go
type Command struct {
    Use               string              // 使用方式
    Short             string              // 简短描述
    Long              string              // 详细描述
    Example           string              // 示例
    Version           string              // 版本号
    Args              PositionalArgs      // 参数验证
    ValidArgs         []string           // 有效参数
    ArgAliases        []string           // 参数别名
    Suggestions       []string           // 建议
    SuggestFor        []string           // 建议替代
    DisableFlagParsing bool              // 禁用标志解析
    DisableAutoGenTag bool               // 禁用自动生成标签
    DisableCommands  bool                // 禁用子命令
    HideDefaultCmd   bool                // 隐藏默认命令
    Hidden           bool                // 隐藏命令
    Run              func(cmd *cobra.Command, args []string) // 执行函数
    PreRun           func(cmd *cobra.Command, args []string) // 前置函数
    PostRun          func(cmd *cobra.Command, args []string) // 后置函数
    PersistentPreRun  func(cmd *cobra.Command, args []string) // 持久前置函数
    PersistentPostRun func(cmd *cobra.Command, args []string) // 持久后置函数
    Flags            *flag.Flags         // 标志集合
    PersistentFlags  *flag.Flags         // 持久标志集合
    Commands         []*Command          // 子命令
    Parent          *Command            // 父命令
}
```

### 常用函数

```go
// 创建命令
cobras.NewCommand(name, version, short string) *cobra.Command

// 添加子命令
cmd.AddCommand(cmds ...*cobra.Command)

// 添加标志
cmd.Flags().StringP(name, shorthand, defaultValue, usage)
cmd.Flags().IntP(name, shorthand, defaultValue, usage)
cmd.Flags().BoolP(name, shorthand, defaultValue, usage)

// 绑定配置
cmd.Flags().BindPFlag(name, flag *pflag.Flag)
```

## 示例代码

### 完整示例：用户管理 CLI

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
    "github.com/spf13/cobra"
)

func main() {
    // 创建根命令
    rootCmd := cobras.NewCommand("userctl", "1.0.0", "用户管理工具")

    // 添加子命令
    rootCmd.AddCommand(
        createUserCommand(),
        listUsersCommand(),
        getUserCommand(),
        deleteUserCommand(),
        versionCommand(),
    )

    // 执行
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
    }
}

func createUserCommand() *cobra.Command {
    var name, email string
    var age int

    cmd := &cobra.Command{
        Use:   "create",
        Short: "创建用户",
        Long:  `创建一个新的用户账户`,
        Example: `
  userctl create -n John -e john@example.com -a 25
  userctl create --name Jane --email jane@example.com`,
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("创建用户:\n")
            fmt.Printf("  姓名: %s\n", name)
            fmt.Printf("  邮箱: %s\n", email)
            fmt.Printf("  年龄: %d\n", age)
        },
    }

    cmd.Flags().StringVarP(&name, "name", "n", "", "用户姓名 (必需)")
    cmd.Flags().StringVarP(&email, "email", "e", "", "用户邮箱 (必需)")
    cmd.Flags().IntVarP(&age, "age", "a", 0, "用户年龄")
    cmd.MarkFlagRequired("name")
    cmd.MarkFlagRequired("email")

    return cmd
}

func listUsersCommand() *cobra.Command {
    var all bool

    cmd := &cobra.Command{
        Use:   "list",
        Short: "列出用户",
        Long:  `列出所有用户或活跃用户`,
        Example: `
  userctl list
  userctl list --all`,
        Run: func(cmd *cobra.Command, args []string) {
            status := "活跃"
            if all {
                status = "所有"
            }
            fmt.Printf("用户列表 (%s):\n", status)
            fmt.Println("  1. 张三 - zhangsan@example.com")
            fmt.Println("  2. 李四 - lisi@example.com")
        },
    }

    cmd.Flags().BoolVarP(&all, "all", "a", false, "显示所有用户")

    return cmd
}

func getUserCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "get <id>",
        Short: "获取用户信息",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            id := args[0]
            fmt.Printf("获取用户 ID: %s\n", id)
            fmt.Println("  姓名: 张三")
            fmt.Println("  邮箱: zhangsan@example.com")
            fmt.Println("  状态: 活跃")
        },
    }

    return cmd
}

func deleteUserCommand() *cobra.Command {
    var force bool

    cmd := &cobra.Command{
        Use:   "delete <id>",
        Short: "删除用户",
        Long:  `删除指定的用户账户`,
        Example: `
  userctl delete 123
  userctl delete 123 --force`,
        Args: cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            id := args[0]
            if !force {
                fmt.Printf("确认删除用户 %s? (y/N): ", id)
                // 这里可以添加交互式确认
            }
            fmt.Printf("用户 %s 已删除\n", id)
        },
    }

    cmd.Flags().BoolVarP(&force, "force", "f", false, "强制删除")

    return cmd
}

func versionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "显示版本信息",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("userctl 版本: 1.0.0")
            fmt.Println("构建时间: 2024-01-01")
        },
    }
}
```

### 示例：服务器管理命令

```go
func serverCommands() *cobra.Command {
    serverCmd := &cobra.Command{
        Use:   "server",
        Short: "服务器管理",
    }

    serverCmd.AddCommand(
        startServerCmd(),
        stopServerCmd(),
        restartServerCmd(),
        statusServerCmd(),
    )

    return serverCmd
}

func startServerCmd() *cobra.Command {
    var port int
    var host string

    cmd := &cobra.Command{
        Use:   "start",
        Short: "启动服务器",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("启动服务器: %s:%d\n", host, port)
        },
    }

    cmd.Flags().IntVarP(&port, "port", "p", 8080, "服务器端口")
    cmd.Flags().StringVarP(&host, "host", "H", "localhost", "服务器主机")

    return cmd
}

func stopServerCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "stop",
        Short: "停止服务器",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("服务器已停止")
        },
    }
}

func restartServerCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "restart",
        Short: "重启服务器",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("重启服务器...")
            fmt.Println("服务器已重启")
        },
    }
}

func statusServerCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "status",
        Short: "查看服务器状态",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("服务器状态: 运行中")
            fmt.Println("运行时间: 2小时30分钟")
            fmt.Println("CPU使用率: 15%")
            fmt.Println("内存使用: 256MB")
        },
    }
}
```

### 示例：数据库命令

```go
func databaseCommands() *cobra.Command {
    dbCmd := &cobra.Command{
        Use:   "db",
        Short: "数据库管理",
    }

    dbCmd.AddCommand(
        migrateCmd(),
        seedCmd(),
        resetCmd(),
    )

    return dbCmd
}

func migrateCmd() *cobra.Command {
    var direction string

    cmd := &cobra.Command{
        Use:   "migrate",
        Short: "数据库迁移",
        Run: func(cmd *cobra.Command, args []string) {
            if direction == "up" {
                fmt.Println("执行数据库迁移...")
                fmt.Println("迁移完成")
            } else {
                fmt.Println("回滚数据库...")
                fmt.Println("回滚完成")
            }
        },
    }

    cmd.Flags().StringVarP(&direction, "direction", "d", "up", "迁移方向 (up/down)")

    return cmd
}

func seedCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "seed",
        Short: "填充测试数据",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("填充测试数据...")
            fmt.Println("数据填充完成")
        },
    }
}

func resetCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "reset",
        Short: "重置数据库",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("重置数据库...")
            fmt.Println("数据库已重置")
        },
    }
}
```

### 示例：配置命令

```go
func configCommands() *cobra.Command {
    configCmd := &cobra.Command{
        Use:   "config",
        Short: "配置管理",
    }

    configCmd.AddCommand(
        showConfigCmd(),
        setConfigCmd(),
    )

    return configCmd
}

func showConfigCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "show",
        Short: "显示配置",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("当前配置:")
            fmt.Println("  database.host: localhost")
            fmt.Println("  database.port: 5432")
            fmt.Println("  server.port: 8080")
        },
    }
}

func setConfigCmd() *cobra.Command {
    var key, value string

    cmd := &cobra.Command{
        Use:   "set",
        Short: "设置配置",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("设置配置: %s = %s\n", key, value)
        },
    }

    cmd.Flags().StringVarP(&key, "key", "k", "", "配置键")
    cmd.Flags().StringVarP(&value, "value", "v", "", "配置值")

    return cmd
}
```

## 最佳实践

1. **清晰的结构**：使用子命令组织相关功能

2. **一致的命名**：使用小写字母和连字符命名命令

3. **提供帮助**：始终提供简短的描述

4. **参数验证**：使用 Args 验证参数数量和类型

5. **默认值**：为可选参数提供合理的默认值

6. **长描述**：为复杂命令提供详细的使用说明

7. **示例**：提供使用示例帮助用户理解

8. **错误处理**：优雅地处理错误情况

9. **版本命令**：添加版本命令显示版本信息

10. **shell 补全**：考虑添加 shell 自动补全支持

## 相关文档

- [Song 框架文档](../../README.md)
- [HTTP 服务器](../https/README.md)
- [配置管理](../vipers/README.md)

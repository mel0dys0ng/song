# Cobras

Cobras是基于spf13/cobra的命令行工具封装包，提供了更简洁的命令定义和管理方式，支持命令继承和层级管理。

## 功能特点

- **简化命令定义**: 封装了Cobra的复杂性，提供更简洁的命令定义方式
- **命令继承**: 支持父命令向子命令传递配置和参数
- **层级管理**: 支持多级命令结构，便于组织复杂CLI应用
- **钩子函数**: 支持PreRun、PostRun、PersistentPreRun、PersistentPostRun等多种钩子函数
- **参数绑定**: 支持命令行参数和标志的便捷绑定

### 定义命令接口实现

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/cobras"
)

// 定义具体的命令实现
type HelloCommand struct{}

func (h *HelloCommand) Short() string {
    return "Say hello to someone"
}

func (h *HelloCommand) Long() string {
    return "A command that says hello to the specified person"
}

func (h *HelloCommand) PreRun(cmd *cobra.Command, args []string) {
    fmt.Println("PreRun: Validating inputs...")
}

func (h *HelloCommand) Run(cmd *cobra.Command, args []string) {
    name, _ := cmd.Flags().GetString("name")
    if name == "" {
        name = "World"
    }
    fmt.Printf("Hello, %s!\n", name)
}

func (h *HelloCommand) PostRun(cmd *cobra.Command, args []string) {
    fmt.Println("PostRun: Cleaning up...")
}

func (h *HelloCommand) PersistentPreRun(cmd *cobra.Command, args []string) {
    fmt.Println("PersistentPreRun: Setting up...")
}

func (h *HelloCommand) PersistentPostRun(cmd *cobra.Command, args []string) {
    fmt.Println("PersistentPostRun: Tearing down...")
}

func (h *HelloCommand) BindFlags(flags *flag.FlagSet) {
    flags.String("name", "", "Name to greet (default: World)")
}
```

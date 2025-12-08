# Cobras - 基于 Cobra 的命令包

## 简介

`cobras` 是一个基于 [Cobra](https://github.com/spf13/cobra) 库构建的命令管理包。它提供了更高级别的抽象和封装，使得创建和管理复杂的命令行应用程序变得更加简单。

## 核心概念

### 主要组件

- **CbassInterface**: 根命令接口，代表整个命令树的入口点
- **CommandInterface**: 命令接口，定义了命令的基本行为
- **CbaInterface**: 子命令注册接口，用于向父命令添加子命令
- **EmptyCommand**: 空命令实现，作为默认命令占位符

### 关键特性

1. **层级化命令结构**：支持嵌套的命令结构，每个命令都可以有多个子命令
2. **继承执行机制**：子命令会继承父命令的部分执行逻辑
3. **灵活的命令注册**：支持动态注册和替换命令实现

## 工作流程

### 初始化

```go
// 创建一个新的命令应用实例
app := cobras.New("myapp")
```

初始化过程会创建一个根命令，并使用 [EmptyCommand](file://d:\work\go\src\github\song\cobras\cobras.go#L11-L11) 作为默认实现。

### 命令注册

有两种方式注册命令：

1. **注册根命令**：

   ```go
   app.RegisterRoot(func(name string) CommandInterface {
       // 返回自定义的根命令实现
   })
   ```

2. **注册子命令**：
   ```go
   app.RegisterCommand("subcmd", mySubCommand)
   ```

### 命令执行

通过 `RegisterExecute` 或 `Execute` 方法启动命令解析和执行流程。

## 内部实现机制

### 命令索引系统

使用冒号(`:`)分隔符构建命令路径索引，例如：`parent:child:grandchild`，便于快速查找和管理命令层级关系。

### 继承执行模型

实现了五级执行链：

1. `PersistentPreRun` - 持续前置运行（可继承）
2. `PreRun` - 前置运行（不可继承）
3. `Run` - 主运行逻辑
4. `PostRun` - 后置运行（不可继承）
5. `PersistentPostRun` - 持续后置运行（可继承）

### 错误处理

集成了统一的日志和错误处理机制，在命令执行失败时提供详细的错误信息。

## 使用模式

典型的使用模式如下：

```go
// 创建应用
app := cobras.New("myapp")

// 注册命令
app.RegisterCommand("serve", serveCommand).
    RegisterCommand("config", configCommand)

// 执行应用
app.Execute()
```

或者使用回调方式：

```go
app.RegisterExecute(func(c CbasInterface) {
    // 在此处注册所有命令
    c.RegisterCommand("serve", serveCommand)
})
```

这个设计使开发者能够专注于业务逻辑实现，而无需关心底层的命令行解析细节。

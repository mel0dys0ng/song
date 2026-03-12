# MySQL 客户端

一个功能强大的 MySQL 数据库客户端，支持读写分离、连接池管理、自动故障转移和查询日志记录。该客户端基于 Go 的 database/sql 包构建，提供了更高级别的抽象和更便捷的 API。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [架构概览](#架构概览)
- [使用指南](#使用指南)
  - [创建客户端](#创建客户端)
  - [执行查询](#执行查询)
  - [读写分离](#读写分离)
  - [连接池配置](#连接池配置)
  - [事务处理](#事务处理)
- [配置选项](#配置选项)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 特性

- **读写分离**：自动将读操作路由到从库，写操作路由到主库
- **连接池管理**：高效管理数据库连接，支持连接复用
- **自动故障转移**：主库故障时自动切换到从库
- **查询日志**：记录慢查询和错误日志
- **指标监控**：支持 Prometheus 指标导出
- **上下文支持**：完全支持 Go 的 context.Context
- **结构体映射**：支持将查询结果直接映射到 Go 结构体

## 安装

确保已安装 MySQL 数据库，然后安装客户端依赖：

```bash
go get github.com/go-sql-driver/mysql
```

## 快速开始

以下是一个最简单的入门示例：

```go
package main

import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

func main() {
    // 创建 MySQL 客户端
    client, err := mysql.New(&mysql.Options{
        DSN:          "user:password@tcp(localhost:3306)/mydb",
        ReadDSNs:     []string{"read1:3306", "read2:3306"},
        MaxOpenConns: 100,
        MaxIdleConns: 10,
    })
    if err != nil {
        panic(err)
    }
    
    // 执行查询（读操作）
    rows, err := client.ReadDB().Query("SELECT * FROM users WHERE id = ?", 1)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    
    // 处理结果
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Printf("User: %d, %s\n", id, name)
    }
}
```

## 架构概览

MySQL 客户端提供了统一的数据库访问接口：

```
┌─────────────────┐
│   MySQL Client  │ - 主客户端接口
└────────┬────────┘
         │
    ┌────┴────┐
    │ Options │ - 配置选项
    └────┬────┘
         │
    ┌────┴────────────┐
    │   DB Pool       │ - 连接池管理
    └─────────────────┘
         │
    ┌────┴────┐
    │ Master  │ ←── 写操作
    │ Slave 1 │ ←── 读操作
    │ Slave 2 │ ←── 读操作
    └─────────┘
```

**核心组件：**
- **Client**：主客户端结构体
- **Options**：客户端配置选项
- **Repository**：数据库仓储接口

## 使用指南

### 创建客户端

创建一个新的 MySQL 客户端实例：

```go
import "github.com/mel0dys0ng/song/internal/core/clients/mysql"

// 使用主从配置创建客户端
client, err := mysql.New(&mysql.Options{
    // 主库连接信息
    DSN: "user:password@tcp(localhost:3306)/mydb",
    
    // 从库连接信息（可选，用于读写分离）
    ReadDSNs: []string{
        "user:password@tcp(read1:3306)/mydb",
        "user:password@tcp(read2:3306)/mydb",
    },
    
    // 连接池配置
    MaxOpenConns: 100,    // 最大打开连接数
    MaxIdleConns: 10,     // 最大空闲连接数
    ConnMaxLifetime: 3600, // 连接最大生命周期（秒）
    ConnMaxIdleTime: 600, // 空闲连接最大时间（秒）
    
    // 日志配置
    EnableSlowQuery: true,
    SlowQueryThreshold: 200, // 毫秒
})
```

### 执行查询

执行不同类型的数据库操作：

```go
// 执行读操作（自动路由到从库）
rows, err := client.ReadDB().QueryContext(ctx, 
    "SELECT id, name, email FROM users WHERE status = ?", "active")

// 执行写操作（自动路由到主库）
result, err := client.WriteDB().ExecContext(ctx,
    "INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com")

// 执行更新操作
result, err := client.WriteDB().ExecContext(ctx,
    "UPDATE users SET name = ? WHERE id = ?", "Jane", 1)

// 执行删除操作
result, err := client.WriteDB().ExecContext(ctx,
    "DELETE FROM users WHERE id = ?", 1)
```

### 读写分离

客户端自动处理读写分离：

```go
// 读操作会自动路由到从库
func getUserByID(client *mysql.Client, id int) (*User, error) {
    rows, err := client.ReadDB().QueryContext(ctx,
        "SELECT id, name, email FROM users WHERE id = ?", id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var user User
    if rows.Next() {
        err := rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            return nil, err
        }
    }
    
    return &user, nil
}

// 写操作会自动路由到主库
func createUser(client *mysql.Client, user *User) (int64, error) {
    result, err := client.WriteDB().ExecContext(ctx,
        "INSERT INTO users (name, email) VALUES (?, ?)",
        user.Name, user.Email)
    if err != nil {
        return 0, err
    }
    
    return result.LastInsertId()
}
```

### 连接池配置

合理配置连接池以获得最佳性能：

```go
client, err := mysql.New(&mysql.Options{
    DSN: "user:password@tcp(localhost:3306)/mydb",
    
    // 连接池大小配置
    MaxOpenConns: 100,    // 最大打开连接数
    MaxIdleConns: 10,     // 最大空闲连接数
    
    // 连接生命周期
    ConnMaxLifetime: 3600, // 连接最大生命周期（秒）
    ConnMaxIdleTime: 600,  // 空闲连接最大时间（秒）
    
    // 慢查询日志
    EnableSlowQuery: true,
    SlowQueryThreshold: 200, // 超过 200ms 记录为慢查询
})
```

### 事务处理

使用事务确保数据一致性：

```go
func transferMoney(client *mysql.Client, fromID, toID int, amount float64) error {
    // 获取事务
    tx, err := client.WriteDB().BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 扣款
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance - ? WHERE id = ?",
        amount, fromID)
    if err != nil {
        return err
    }
    
    // 充值
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance + ? WHERE id = ?",
        amount, toID)
    if err != nil {
        return err
    }
    
    // 提交事务
    return tx.Commit()
}
```

## 配置选项

### Options 结构体

```go
type Options struct {
    // 必需配置
    DSN     string   // 主库数据源名称
    ReadDSNs []string // 从库数据源名称列表（可选）
    
    // 连接池配置
    MaxOpenConns int    // 最大打开连接数（默认：100）
    MaxIdleConns int    // 最大空闲连接数（默认：10）
    ConnMaxLifetime int // 连接最大生命周期（秒）（默认：3600）
    ConnMaxIdleTime int // 空闲连接最大时间（秒）（默认：600）
    
    // 日志配置
    EnableSlowQuery bool // 是否启用慢查询日志（默认：true）
    SlowQueryThreshold int64 // 慢查询阈值（毫秒）（默认：200）
    
    // 其他配置
    Logger interface{} // 自定义日志记录器
}
```

### DSN 格式

MySQL 数据源名称格式：

```
username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
```

**参数说明：**
- `username`：数据库用户名
- `password`：数据库密码
- `host`：数据库主机地址
- `port`：数据库端口（默认：3306）
- `database`：数据库名称
- `charset`：字符编码（推荐使用 utf8mb4）
- `parseTime`：是否解析时间字段
- `loc`：时区设置

## 示例代码

### 完整示例：用户服务

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository struct {
    client *mysql.Client
}

func NewUserRepository(client *mysql.Client) *UserRepository {
    return &UserRepository{client: client}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    rows, err := r.client.ReadDB().QueryContext(ctx,
        "SELECT id, name, email, status, created_at, updated_at FROM users WHERE id = ?", id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var user User
    if rows.Next() {
        err := rows.Scan(&user.ID, &user.Name, &user.Email, 
            &user.Status, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, err
        }
    }
    
    return &user, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
    rows, err := r.client.ReadDB().QueryContext(ctx,
        "SELECT id, name, email, status, created_at, updated_at FROM users LIMIT ? OFFSET ?",
        limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email,
            &user.Status, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    
    return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *User) (int64, error) {
    result, err := r.client.WriteDB().ExecContext(ctx,
        "INSERT INTO users (name, email, status) VALUES (?, ?, ?)",
        user.Name, user.Email, user.Status)
    if err != nil {
        return 0, err
    }
    
    return result.LastInsertId()
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
    _, err := r.client.WriteDB().ExecContext(ctx,
        "UPDATE users SET name = ?, email = ?, status = ?, updated_at = ? WHERE id = ?",
        user.Name, user.Email, user.Status, time.Now(), user.ID)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    _, err := r.client.WriteDB().ExecContext(ctx,
        "DELETE FROM users WHERE id = ?", id)
    return err
}

func main() {
    // 创建客户端
    client, err := mysql.New(&mysql.Options{
        DSN:     "root:password@tcp(localhost:3306)/mydb",
        ReadDSNs: []string{"root:password@tcp(read1:3306)/mydb"},
        MaxOpenConns: 100,
        MaxIdleConns: 10,
    })
    if err != nil {
        panic(err)
    }
    
    // 创建仓储
    repo := NewUserRepository(client)
    ctx := context.Background()
    
    // 创建用户
    user := &User{
        Name:   "张三",
        Email:  "zhangsan@example.com",
        Status: "active",
    }
    
    id, err := repo.Create(ctx, user)
    if err != nil {
        panic(err)
    }
    fmt.Printf("创建用户 ID: %d\n", id)
    
    // 查询用户
    fetchedUser, err := repo.GetByID(ctx, id)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户信息: %+v\n", fetchedUser)
    
    // 更新用户
    fetchedUser.Name = "李四"
    err = repo.Update(ctx, fetchedUser)
    if err != nil {
        panic(err)
    }
    
    // 删除用户
    err = repo.Delete(ctx, id)
    if err != nil {
        panic(err)
    }
}
```

### 示例：连接池监控

```go
import (
    "fmt"
    "github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

func printPoolStats(client *mysql.Client) {
    stats := client.WriteDB().Stats()
    fmt.Printf("连接池统计:\n")
    fmt.Printf("  打开连接数: %d\n", stats.OpenConnections)
    fmt.Printf("  空闲连接数: %d\n", stats.Idle)
    fmt.Printf("  使用中的连接数: %d\n", stats.InUse)
    fmt.Printf("  等待连接数: %d\n", stats.WaitCount)
    fmt.Printf("  最大等待时间: %v\n", stats.WaitDuration)
    fmt.Printf("  关闭的连接数: %d\n", stats.Closed)
}
```

## 最佳实践

1. **使用读写分离**：读操作多使用 ReadDB()，写操作使用 WriteDB()

2. **配置合理的连接池**：
   - 根据应用负载设置 MaxOpenConns
   - 设置合适的 MaxIdleConns 以减少连接创建开销
   - 设置 ConnMaxLifetime 防止连接失效

3. **使用上下文**：始终使用 context.Context 进行超时和取消控制

4. **处理错误**：正确处理数据库错误，包括连接超时和查询错误

5. **关闭资源**：使用 defer 关闭 rows 和 stmt

6. **使用预处理语句**：对于频繁执行的查询，使用预处理语句提高性能

7. **监控慢查询**：启用慢查询日志，监控查询性能

8. **连接字符串安全**：不要在代码中硬编码密码，使用环境变量或密钥管理服务

9. **使用结构体映射**：使用 Scan 将结果映射到结构体，提高代码可读性

10. **事务使用**：对于需要原子性的操作，使用事务确保数据一致性

## 相关文档

- [Song 框架文档](../../README.md)
- [Redis 客户端](../redis/README.md)
- [配置管理](../vipers/README.md)

# MySQL Client

A powerful MySQL database client with built-in support for read/write separation, connection pooling, and load balancing. Built on top of GORM, this client provides a robust repository pattern for database operations.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage Guide](#usage-guide)
  - [Creating a Client](#creating-a-client)
  - [Basic Database Operations](#basic-database-operations)
  - [Using the Repository Pattern](#using-the-repository-pattern)
  - [Advanced Query Operations](#advanced-query-operations)
  - [Transaction Management](#transaction-management)
- [Configuration Options](#configuration-options)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Features

- **Read/Write Separation**: Automatically route write operations to master and read operations to slave databases
- **Connection Pooling**: Efficient connection management with configurable pool sizes
- **Load Balancing**: Distribute read queries across multiple slave databases
- **GORM Integration**: Full support for GORM ORM features
- **Repository Pattern**: Generic repository for type-safe database operations
- **Logging**: Comprehensive query logging with slow query detection
- **Configuration Management**: Load configuration from YAML, JSON, or other sources
- **Singleton Pattern**: Automatic client reuse and lifecycle management

## Installation

The MySQL client is part of the `song` framework. Ensure you have the required dependencies:

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/plugin/dbresolver
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

func main() {
    ctx := context.Background()

    // Create a MySQL client
    client := mysql.CreateClient(ctx, "", "database.mysql")

    // Perform a simple query
    type User struct {
        ID   uint   `gorm:"primaryKey"`
        Name string `gorm:"size:255"`
    }

    var user User
    client.Slave().First(&user, 1)
}
```

## Configuration

The MySQL client supports flexible configuration through configuration files or programmatic options.

### Configuration File Example (YAML)

```yaml
database:
  mysql:
    debug: true
    master: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    slaves:
      - "user:password@tcp(slave1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
      - "user:password@tcp(slave2:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    tablePrefix: "tbl_"
    singularTable: false
    maxIdle: 100
    maxActive: 200
    maxConnLifeTime: 3600000
    idleTimeout: 3600000
    logSlow: 1000
    logLevel: "warn"
    logDir: "/data/logs/mysql"
```

## Usage Guide

### Creating a Client

Create a MySQL client instance using the `CreateClient` function:

```go
import (
    "context"
    "github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

func main() {
    ctx := context.Background()

    // Basic client creation
    client := mysql.CreateClient(ctx, "", "database.mysql")

    // Client with custom name (for multiple instances)
    client2 := mysql.CreateClient(ctx, "analytics", "database.mysql")
}
```

**Parameters:**

- `ctx`: Context for the operation
- `name`: Custom name for the client (optional, use empty string for default)
- `key`: Configuration key to load from config source
- `opts`: Optional configuration options

### Basic Database Operations

The client provides direct access to GORM's DB instance:

```go
// Write operations (uses master database)
client.Master().Create(&user)
client.Master().Delete(&user, id)
client.Master().Updates(&user)

// Read operations (uses slave database)
client.Slave().First(&user, id)
client.Slave().Find(&users)
client.Slave().Where("status = ?", "active").Find(&users)
```

### Using the Repository Pattern

The repository pattern provides type-safe database operations:

```go
// Define your model
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"size:255"`
    Email    string `gorm:"size:255"`
    Status   string `gorm:"size:20"`
}

// Implement ModelInterface
func (User) TableName() string {
    return "users"
}

// Create repository
repo := mysql.NewRepository[User](client)

// Query single record
user, err := repo.Query(ctx, &mysql.QueryRequest{
    Query: "id = ?",
    Arguments: []interface{}{1},
})

// Query list
users, err := repo.QueryList(ctx, &mysql.QueryListRequest{
    Query: "status = ?",
    Arguments: []interface{}{"active"},
    Order: "created_at DESC",
    Limit: 10,
})

// Create record
err := repo.Create(ctx, &user)

// Update record
err := repo.Update(ctx, &user)

// Delete record
err := repo.Delete(ctx, &user)
```

### Advanced Query Operations

#### Join Operations

```go
user, err := repo.Query(ctx, &mysql.QueryRequest{
    Query: "u.id = ?",
    Arguments: []interface{}{1},
    Alias: "u",
    Joins: []*mysql.JoinRequest{
        {
            Query: "LEFT JOIN orders o ON u.id = o.user_id",
        },
    },
    Fields: []string{"u.*", "o.id as order_id"},
})
```

#### Pagination

```go
list, err := repo.QueryList(ctx, &mysql.QueryListRequest{
    Query: "status = ?",
    Arguments: []interface{}{"active"},
    Page:    1,
    Limit:   20,
    Order:   "created_at DESC",
})
```

#### For Update (Locking)

```go
user, err := repo.Query(ctx, &mysql.QueryRequest{
    Query:     "id = ?",
    Arguments: []interface{}{1},
    ForUpdate: true, // Adds FOR UPDATE clause
})
```

### Transaction Management

```go
// Begin a transaction
tx := client.Begin()

// Use the transaction in repository operations
err := repo.Create(ctx, &user, mysql.WithTx(tx))

// Commit or rollback
if err != nil {
    tx.Rollback()
} else {
    tx.Commit()
}

// Or use transaction with query
user, err := repo.Query(ctx, &mysql.QueryRequest{
    Query: "id = ?",
    Arguments: []interface{}{1},
    Tx: tx,
})
```

## Configuration Options

### Programmatic Configuration

You can configure the client using functional options:

```go
client := mysql.CreateClient(ctx, "", "database.mysql",
    mysql.Debug(true),
    mysql.Master("user:pass@tcp(master:3306)/db"),
    mysql.Slaves([]string{
        "user:pass@tcp(slave1:3306)/db",
        "user:pass@tcp(slave2:3306)/db",
    }),
    mysql.TablePrefix("tbl_"),
    mysql.MaxIdle(100),
    mysql.MaxActive(200),
    mysql.IdleTimeout(3600000 * time.Millisecond),
    mysql.LogSlow(1000 * time.Millisecond),
    mysql.LogLevel("warn"),
)
```

### Available Options

| Option                    | Type            | Description                    | Default              |
| ------------------------- | --------------- | ------------------------------ | -------------------- |
| `Debug`                   | `bool`          | Enable debug mode              | `true`               |
| `Master`                  | `string`        | Master database DSN            | `""`                 |
| `Slaves`                  | `[]string`      | Slave database DSNs            | `[]`                 |
| `TablePrefix`             | `string`        | Table name prefix              | `""`                 |
| `SingularTable`           | `bool`          | Use singular table names       | `false`              |
| `MaxIdle`                 | `int`           | Max idle connections           | `100`                |
| `MaxActive`               | `int`           | Max active connections         | `200`                |
| `MaxConnLifeTime`         | `time.Duration` | Max connection lifetime        | `1h`                 |
| `IdleTimeout`             | `time.Duration` | Idle connection timeout        | `1h`                 |
| `LogSlow`                 | `time.Duration` | Slow query threshold           | `1s`                 |
| `LogLevel`                | `string`        | Log level                      | `"warn"`             |
| `LogDir`                  | `string`        | Log directory                  | `"/data/logs/mysql"` |
| `LogIgnoreRecordNotFound` | `bool`          | Ignore record not found errors | `true`               |
| `LogColorful`             | `bool`          | Enable colorful logs           | `true`               |

## Examples

### Complete Example with Error Handling

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/mel0dys0ng/song/internal/core/clients/mysql"
)

type Product struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:255"`
    Price int    `gorm:"type:decimal(10,2)"`
}

func (Product) TableName() string {
    return "products"
}

func main() {
    ctx := context.Background()

    // Create client
    client := mysql.CreateClient(ctx, "", "database.mysql")
    if client == nil {
        log.Fatal("Failed to create MySQL client")
    }

    // Create repository
    repo := mysql.NewRepository[Product](client)

    // Create a new product
    product := Product{
        Name:  "Laptop",
        Price: 999,
    }

    err := repo.Create(ctx, &product)
    if err != nil {
        log.Printf("Failed to create product: %v", err)
        return
    }

    fmt.Printf("Created product with ID: %d\n", product.ID)

    // Query the product
    found, err := repo.Query(ctx, &mysql.QueryRequest{
        Query:     "id = ?",
        Arguments: []interface{}{product.ID},
    })
    if err != nil {
        log.Printf("Failed to query product: %v", err)
        return
    }

    fmt.Printf("Found product: %+v\n", found)

    // List products
    products, err := repo.QueryList(ctx, &mysql.QueryListRequest{
        Query: "price > ?",
        Arguments: []interface{}{500},
        Order:   "price DESC",
        Limit:   10,
    })
    if err != nil {
        log.Printf("Failed to list products: %v", err)
        return
    }

    fmt.Printf("Found %d products\n", len(products))
}
```

### Multiple Database Connections

```go
// Primary database
primaryClient := mysql.CreateClient(ctx, "", "database.mysql.primary")

// Analytics database
analyticsClient := mysql.CreateClient(ctx, "analytics", "database.mysql.analytics")

// Use different clients for different purposes
userRepo := mysql.NewRepository[User](primaryClient)
metricRepo := mysql.NewRepository[Metric](analyticsClient)
```

## Best Practices

1. **Reuse Client Instances**: The client uses singleton pattern internally. Create clients once and reuse them throughout your application.

2. **Use Repository Pattern**: For type-safe operations and better code organization, prefer the repository pattern over direct GORM usage.

3. **Read/Write Separation**: Always use `Master()` for writes and `Slave()` for reads to leverage read/write separation.

4. **Connection Pooling**: Tune `MaxIdle`, `MaxActive`, and timeout settings based on your application's load patterns.

5. **Slow Query Logging**: Set appropriate `LogSlow` threshold to identify performance bottlenecks.

6. **Error Handling**: Always check errors returned from repository operations and handle them appropriately.

7. **Context Usage**: Pass context to all operations for proper timeout and cancellation handling.

8. **Transactions**: Use transactions for operations that require atomicity. Remember to commit or rollback.

9. **Configuration Management**: Store sensitive information like database credentials in environment variables or secure configuration sources.

10. **Monitoring**: Monitor connection pool metrics and query performance in production.

## Additional Resources

- [GORM Documentation](https://gorm.io/docs/)
- [Database Resolver Plugin](https://gorm.io/docs/dbresolver.html)
- [Song Framework Documentation](../../README.md)

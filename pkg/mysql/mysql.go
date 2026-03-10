package mysql

import (
	"context"

	"github.com/mel0dys0ng/song/internal/core/clients/mysql"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

type (
	Client                          = mysql.Client
	Option                          = mysql.Option
	JoinQueryArguments              = mysql.JoinQueryArguments
	QueryRequest                    = mysql.QueryRequest
	QueryListRequest                = mysql.QueryListRequest
	QueryCountRequest               = mysql.QueryCountRequest
	CreateRequest[T ModelInterface] = mysql.CreateRequest[T]
	UpdateRequest[T ModelInterface] = mysql.UpdateRequest[T]
	DeleteRequest[T ModelInterface] = mysql.DeleteRequest[T]
	ModelInterface                  = mysql.ModelInterface
	Repository[T ModelInterface]    = mysql.Repository[T]
)

var (
	ErrInvalidParams = mysql.ErrInvalidParams
)

// NewMySQLClient 返回仅使用统一配置的MySQL数据库连接，支持读写分离、负载均衡、连接池
// @key string database config key
func NewMySQLClient(ctx context.Context, key string) *mysql.Client {
	ctx = erlogs.StartTracef(ctx, "NewMySQLClient:%s", key)
	defer erlogs.EndTrace(ctx, nil)
	return mysql.CreateClient(ctx, "", key)
}

// CustomMySQLClient 返回基于统一配置、自定义配置选项的MySQL数据库连接，支持读写分离、负载均衡、连接池
// @key string database config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @options []Option 自定义配置选项
func CustomMySQLClient(ctx context.Context, name string, key string, opts ...mysql.Option) *mysql.Client {
	ctx = erlogs.StartTracef(ctx, "CustomMySQLClient:%s:%s", name, key)
	defer erlogs.EndTrace(ctx, nil)
	return mysql.CreateClient(ctx, name, key, opts...)
}

// NewRepository 创建一个新的泛型仓库实例
// 该函数用于创建一个基于MySQL客户端的泛型仓库，提供类型安全的数据访问层
//
// 参数:
//   - client: *mysql.Client - MySQL数据库客户端实例，必须是已初始化的连接对象
//     该客户端应通过NewMySQLClient或CustomMySQLClient创建，包含完整的数据库连接配置
//
// 返回值:
//   - *mysql.Repository[T]: 指向泛型仓库的指针，其中T是ModelInterface的实现
//     返回的仓库实例提供了对指定类型数据的CRUD操作，具有类型安全保证
//
// 使用示例:
//
//	client := NewMySQLClient(ctx, "default")
//	repo := NewRepository[*Users](client)
func NewRepository[T ModelInterface](client *mysql.Client) *mysql.Repository[T] {
	return mysql.NewRepository[T](client)
}

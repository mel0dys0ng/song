package clients

import (
	"context"
	"fmt"

	"github.com/mel0dys0ng/song/core/clients/internal/mysql"
	"github.com/mel0dys0ng/song/core/erlogs"
)

type (
	MySQLClient = mysql.Client
)

// NewMySQLClient 返回仅使用统一配置的MySQL数据库连接，支持读写分离、负载均衡、连接池
// @key string database config key
func NewMySQLClient(ctx context.Context, key string) *mysql.Client {
	ctx = erlogs.StartTrace(ctx, fmt.Sprintf("NewMySQLClient:%s", key))
	defer erlogs.EndTrace(ctx, nil)
	return mysql.CreateClient(ctx, "", key)
}

// CustomMySQLClient 返回基于统一配置、自定义配置选项的MySQL数据库连接，支持读写分离、负载均衡、连接池
// @key string database config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @options []Option 自定义配置选项
func CustomMySQLClient(ctx context.Context, name string, key string, opts ...mysql.Option) *mysql.Client {
	ctx = erlogs.StartTrace(ctx, fmt.Sprintf("CustomMySQLClient:%s:%s", name, key))
	defer erlogs.EndTrace(ctx, nil)
	return mysql.CreateClient(ctx, name, key, opts...)
}

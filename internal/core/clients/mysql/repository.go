// Repository 是一个泛型数据库操作仓库，封装了常用的增删改查操作
// 使用GORM作为底层ORM框架，支持泛型实体操作
package mysql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository 定义了数据库操作的仓库结构
// T 是泛型参数，代表数据模型类型，必须是可比较类型
type Repository[T ModelInterface] struct {
	Client *Client
}

// NewRepository 创建一个新的Repository实例
// client 参数是已初始化的数据库客户端
// 返回指向Repository的指针
func NewRepository[T ModelInterface](client *Client) *Repository[T] {
	return &Repository[T]{
		Client: client,
	}
}

// Query 根据请求参数执行单条记录查询
// ctx 上下文用于控制超时和取消
// req 查询请求参数，包含查询条件、字段等
// 返回查询结果和可能发生的错误
func (r *Repository[T]) Query(ctx context.Context, req *QueryRequest) (res T, err error) {
	// 验证请求参数的有效性
	if err = req.Validate(); err != nil {
		return
	}

	// 设置查询结果的数据结构
	var data T
	query := r.Client.DB(req.Tx, req.UseSlave).Model(&data)

	// 设置表别名
	if len(req.Alias) > 0 {
		query = query.Table(fmt.Sprintf("%s as %s", data.TableName(), req.Alias))
	}

	// 添加JOIN条件到查询中
	for _, join := range req.Joins {
		if join != nil && len(join.Query) > 0 {
			query = query.Joins(join.Query, join.Arguments...)
		}
	}

	if req.ForUpdate {
		// 如果需要FOR UPDATE，则添加FOR UPDATE到查询中
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	query = query.Select(req.Fields)

	if len(req.Query) > 0 {
		// 添加WHERE条件到查询中
		query = query.Where(req.Query, req.Arguments...)
	}

	// 执行查询并将结果存储在data变量中
	err = query.First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有找到记录，将错误设为nil并返回
		err = nil
		return
	}

	if err != nil {
		// 如果发生其他错误，则直接返回
		return
	}

	// 将查询到的数据赋值给返回值
	res = data
	return
}

// QueryList 根据请求参数执行列表查询
// ctx 上下文用于控制超时和取消
// req 列表查询请求参数，包含查询条件、排序、分组等
// 返回查询结果列表和可能发生的错误
func (r *Repository[T]) QueryList(ctx context.Context, req *QueryListRequest) (res []T, err error) {
	// 验证请求参数的有效性
	if err = req.Validate(); err != nil {
		return
	}

	var model T

	query := r.Client.DB(req.Tx, req.UseSlave).Model(model)

	// 设置表别名
	if len(req.Alias) > 0 {
		query = query.Table(fmt.Sprintf("%s as %s", model.TableName(), req.Alias))
	}

	// 添加JOIN条件到查询中
	for _, join := range req.Joins {
		if join != nil && len(join.Query) > 0 {
			query = query.Joins(join.Query, join.Arguments...)
		}
	}

	if req.ForUpdate {
		// 如果需要FOR UPDATE，则添加FOR UPDATE到查询中
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	query = query.Select(req.Fields)
	if len(req.Query) > 0 {
		// 添加WHERE条件到查询中
		query = query.Where(req.Query, req.Arguments...)
	}

	if len(req.Group) > 0 {
		// 添加GROUP BY条件到查询中
		query = query.Group(req.Group)
	}

	if len(req.Order) > 0 {
		// 添加ORDER BY条件到查询中
		query = query.Order(req.Order)
	}

	// 添加LIMIT和OFFSET限制查询结果数量和偏移量
	query = query.Limit(req.Limit).Offset(req.Offset)

	var data []T
	err = query.Find(&data).Error
	if err != nil {
		// 如果发生其他错误，则直接返回
		return
	}

	res = data

	return
}

// QueryCount 根据请求参数执行计数查询
// ctx 上下文用于控制超时和取消
// req 计数查询请求参数，包含查询条件等
// 返回符合条件的记录总数和可能发生的错误
func (r *Repository[T]) QueryCount(ctx context.Context, req *QueryCountRequest) (res int64, err error) {
	// 验证请求参数的有效性
	if err = req.Validate(); err != nil {
		return
	}

	var model T

	query := r.Client.DB(req.Tx, req.UseSlave).Model(model)

	// 设置表别名
	if len(req.Alias) > 0 {
		query = query.Table(fmt.Sprintf("%s as %s", model.TableName(), req.Alias))
	}

	// 添加JOIN条件到查询中
	for _, join := range req.Joins {
		if join != nil && len(join.Query) > 0 {
			query = query.Joins(join.Query, join.Arguments...)
		}
	}

	if req.ForUpdate {
		// 如果需要FOR UPDATE，则添加FOR UPDATE到查询中
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	if len(req.Query) > 0 {
		query = query.Where(req.Query, req.Arguments...)
	}

	// 执行计数查询
	var data int64
	err = query.Count(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有找到记录，将错误设为nil并返回
		err = nil
		return
	}

	if err != nil {
		// 如果发生其他错误，则直接返回
		return
	}

	res = data

	return
}

// Create 根据请求参数创建新记录
// ctx 上下文用于控制超时和取消
// req 创建请求参数，包含待创建的数据
// 返回影响的行数和可能发生的错误
func (r *Repository[T]) Create(ctx context.Context, req *CreateRequest[T]) (rowsAffected int64, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var model T
	result := r.Client.DB(req.Tx, false).Model(model).Create(&req.Data)
	if err = result.Error; err != nil {
		return
	}

	rowsAffected = result.RowsAffected

	return
}

// Update 根据请求参数更新现有记录
// ctx 上下文用于控制超时和取消
// req 更新请求参数，包含更新条件和更新数据
// 返回影响的行数和可能发生的错误
func (r *Repository[T]) Update(ctx context.Context, req Updater[T]) (rowsAffected int64, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var model T
	query := r.Client.DB(req.GetTx(), false).Model(model)
	if len(req.GetFields()) > 0 {
		query = query.Select(req.GetFields())
	}

	if len(req.GetQuery()) > 0 {
		query = query.Where(req.GetQuery(), req.GetArguments()...)
	}

	if req.GetLimit() > 0 {
		query = query.Limit(req.GetLimit())
	}

	result := query.Updates(req.GetData())
	if err = result.Error; err != nil {
		return
	}

	rowsAffected = result.RowsAffected

	return
}

// Delete 根据请求参数删除记录
// ctx 上下文用于控制超时和取消
// req 删除请求参数，包含删除条件
// 返回影响的行数和可能发生的错误
func (r *Repository[T]) Delete(ctx context.Context, req *DeleteRequest[T]) (rowsAffected int64, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var model T
	query := r.Client.DB(req.Tx, false).Model(model)
	if len(req.Query) > 0 {
		query = query.Where(req.Query, req.Arguments...)
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	result := query.Delete(&model)
	if err = result.Error; err != nil {
		return
	}

	rowsAffected = result.RowsAffected

	return
}

// Upsert 在创建发生冲突时执行更新操作
// 使用GORM的Clauses方法配合OnConflict子句实现冲突处理逻辑
//
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期和传递额外信息
//   - req: 创建或更新请求对象，包含要操作的数据和冲突处理策略
//
// 返回值:
//   - rowsAffected: 操作影响的行数
//   - err: 操作过程中可能发生的错误
func (r *Repository[T]) Upsert(ctx context.Context, req *UpsertRequest[T]) (rowsAffected int64, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var model T
	query := r.Client.DB(req.Tx, false).Model(model)
	result := query.Clauses(req.OnConflict).Create(&req.Data)
	if err = result.Error; err != nil {
		return
	}

	rowsAffected = result.RowsAffected

	return
}

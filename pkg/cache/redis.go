package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/mel0dys0ng/song/pkg/crypto"
	"github.com/mel0dys0ng/song/pkg/result"
	"github.com/redis/go-redis/v9"
)

type (
	// redisCache Redis缓存结构体
	redisCache[T any] struct {
		client redis.UniversalClient
		ttl    time.Duration
	}
)

// newRedisCache 创建一个新的Redis缓存实例
func newRedisCache[T any](client redis.UniversalClient, ttl time.Duration) *redisCache[T] {
	return &redisCache[T]{
		client: client,
		ttl:    ttl,
	}
}

// isSet 检查Redis缓存是否已设置
func (c *redisCache[T]) isSet() bool {
	return c != nil
}

// isRetryEnable 检查Redis缓存是否支持重试
func (c *redisCache[T]) isRetryEnable() bool {
	return true
}

// getType 获取缓存类型
func (c *redisCache[T]) getType() string {
	return TypeRedis
}

// get 从Redis缓存中获取数据
func (c *redisCache[T]) get(ctx context.Context, key any, prefix string) (res result.Interface[T]) {
	var data T

	keyKey := c.genKey(prefix, key)
	idNameKey, err := c.client.Get(ctx, keyKey).Result()
	if errors.Is(err, redis.Nil) || len(idNameKey) == 0 {
		return result.Success(data)
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("GetByRedis: getIdNameKeyByKeyKey, %w", err))
	}

	value, err := c.client.Get(ctx, idNameKey).Result()
	if errors.Is(err, redis.Nil) || len(value) == 0 {
		return result.Success(data)
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("GetByRedis: %w", err))
	}

	err = json.Unmarshal([]byte(value), &data)
	if err != nil {
		return result.Error[T](fmt.Errorf("GetByRedis: json.Unmarshal, %w", err))
	}

	return result.Success(data)
}

// set 在Redis中设置一个值，并建立相关的键映射关系
// 此函数用于缓存数据并维护键之间的映射关系，以支持通过不同类型的键（如主键或ID）访问同一数据
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - key: 主键，可以通过它检索缓存的值
//   - name: 缓存名称，用于区分不同的缓存空间
//   - prefix: 键名前缀，用于组织和区分不同业务的缓存
//   - id: 数据的唯一标识符，也可以用来检索缓存值
//   - value: 实际要缓存的值
//
// 返回值:
//   - result.Interface[T]: 包含操作结果的接口，成功时包含缓存的值，失败时包含错误信息
func (c *redisCache[T]) set(ctx context.Context, key, name any, prefix string, id any, value T) (res result.Interface[T]) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: %w", err))
	}

	// 设置key->idNameKey映射
	// 根据传入的key生成Redis键，然后查找对应的idNameKey
	keyKey := c.genKey(prefix, key)
	idNameKey, err := c.client.Get(ctx, keyKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: getIdNameKeyByKeyKey, %w", err))
	}

	// 如果keyKey不存在，则生成新的idNameKey并建立key->idNameKey映射
	if len(idNameKey) == 0 {
		idNameKey = c.genKey(prefix, name, id)
		_, err = c.client.Set(ctx, keyKey, idNameKey, c.ttl).Result()
		if err != nil {
			return result.Error[T](fmt.Errorf("SetByRedis: setIdNameKeyByKeyKey, %w", err))
		}
	}

	// 设置id->idNameKey映射
	// 根据传入的id生成Redis键，然后建立id->idNameKey映射
	idKey := c.genKey(prefix, id)
	idNameKy, err := c.client.Get(ctx, idKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: getIdNameKeyByIdKey, %w", err))
	}

	// 如果idKey不存在对应的idNameKey，则建立id->idNameKey映射
	if len(idNameKy) == 0 {
		_, err = c.client.Set(ctx, idKey, idNameKey, c.ttl).Result()
		if err != nil {
			return result.Error[T](fmt.Errorf("SetByRedis: setIdNameKeyByIdKey, %w", err))
		}
	}

	// 更新关联键列表
	// 获取与idNameKey相关的所有键列表，用于后续批量删除操作
	keysRes, err := c.client.Get(ctx, idNameKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: getKeysByIdNameKey, %w", err))
	}

	var keys []string
	if len(keysRes) > 0 {
		err = json.Unmarshal([]byte(keysRes), &keys)
		if err != nil {
			return result.Error[T](fmt.Errorf("SetByRedis: json.Unmarshal, %w", err))
		}
	}

	// 将keyKey和idKey添加到关联键列表中，如果它们还不存在的话
	for _, v := range []string{keyKey, idKey} {
		if !slices.Contains(keys, v) {
			keys = append(keys, v)
		}
	}

	keysBytes, err := json.Marshal(keys)
	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: json.Marshal, %w", err))
	}

	keysKey := c.genKey(prefix, idNameKey)
	_, err = c.client.Set(ctx, keysKey, string(keysBytes), c.ttl).Result()
	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: SetKeys, %w", err))
	}

	// 设置实际值
	// 将序列化的值存储在idNameKey下，这是实际数据的存储位置
	_, err = c.client.Set(ctx, idNameKey, string(bytes), c.ttl).Result()
	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: %w", err))
	}

	return result.Success(value)
}

// del 从Redis缓存中删除数据
func (c *redisCache[T]) del(ctx context.Context, key any, prefix string) (res result.Interface[T]) {
	getRes := c.get(ctx, key, prefix)
	if getRes.Err() != nil {
		return getRes
	}

	// 获取要删除的idNameKey
	keyKey := c.genKey(prefix, key)
	idNameKey, err := c.client.Get(ctx, keyKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("DelByRedis: getIdNameKeyByKeyKey, %w", err))
	}

	var data T
	if len(idNameKey) == 0 {
		return result.Success(data)
	}

	// 获取关联键列表
	keysKey := c.genKey(prefix, idNameKey)
	keysRes, err := c.client.Get(ctx, keysKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("DelByRedis: getKeysByIdNameKey, %w", err))
	}

	var keys []string
	if len(keysRes) > 0 {
		err = json.Unmarshal([]byte(keysRes), &keys)
		if err != nil {
			return result.Error[T](fmt.Errorf("DelByRedis: json.Unmarshal, %w", err))
		}
	}

	// 删除所有相关键
	delKeys := append(keys, idNameKey, keysKey)
	_, err = c.client.Del(ctx, delKeys...).Result()
	if err != nil {
		return result.Error[T](fmt.Errorf("DelByRedis: %w", err))
	}

	return result.Success(getRes.Data())
}

// genKey 生成缓存键
func (c *redisCache[T]) genKey(prefix string, data ...any) string {
	return fmt.Sprintf("%s%s", prefix, crypto.MD5(data))
}

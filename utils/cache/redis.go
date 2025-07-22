package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/mel0dys0ng/song/utils/crypto"
	"github.com/mel0dys0ng/song/utils/result"
	"github.com/redis/go-redis/v9"
)

type (
	redisCache[T any] struct {
		client redis.UniversalClient
		ttl    time.Duration
	}
)

func newRedisCache[T any](client redis.UniversalClient, ttl time.Duration) *redisCache[T] {
	return &redisCache[T]{
		client: client,
		ttl:    ttl,
	}
}

func (c *redisCache[T]) isSet() bool {
	return c != nil
}

func (c *redisCache[T]) isRetryEnable() bool {
	return true
}

func (c *redisCache[T]) getType() string {
	return TypeRedis
}

// get 获取缓存数据
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

// set 设置缓存数据
func (c *redisCache[T]) set(ctx context.Context, key, name any, prefix string, id any, value T) (res result.Interface[T]) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: %w", err))
	}

	// -----------------------------------------------------------------------------------
	keyKey := c.genKey(prefix, key)
	idNameKey, err := c.client.Get(ctx, keyKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: getIdNameKeyByKeyKey, %w", err))
	}

	if len(idNameKey) == 0 {
		idNameKey = c.genKey(prefix, name, id)
		_, err = c.client.Set(ctx, keyKey, idNameKey, c.ttl).Result()
		if err != nil {
			return result.Error[T](fmt.Errorf("SetByRedis: setIdNameKeyByKeyKey, %w", err))
		}
	}

	// -----------------------------------------------------------------------------------
	idKey := c.genKey(prefix, id)
	idNameKy, err := c.client.Get(ctx, idKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: getIdNameKeyByIdKey, %w", err))
	}

	if len(idNameKy) == 0 {
		_, err = c.client.Set(ctx, idKey, idNameKey, c.ttl).Result()
		if err != nil {
			return result.Error[T](fmt.Errorf("SetByRedis: setIdNameKeyByIdKey, %w", err))
		}
	}

	// -----------------------------------------------------------------------------------
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

	// -----------------------------------------------------------------------------------
	_, err = c.client.Set(ctx, idNameKey, string(bytes), c.ttl).Result()
	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: %w", err))
	}

	return result.Success(value)
}

// del 删除缓存数据
func (c *redisCache[T]) del(ctx context.Context, key any, prefix string) (res result.Interface[T]) {
	getRes := c.get(ctx, key, prefix)
	if getRes.Err() != nil {
		return getRes
	}

	// -----------------------------------------------------------------------------------
	keyKey := c.genKey(prefix, key)
	idNameKey, err := c.client.Get(ctx, keyKey).Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}

	if err != nil {
		return result.Error[T](fmt.Errorf("SetByRedis: getIdNameKeyByKeyKey, %w", err))
	}

	var data T
	if len(idNameKey) == 0 {
		return result.Success(data)
	}

	// -----------------------------------------------------------------------------------
	keysKey := c.genKey(prefix, idNameKey)
	keysRes, err := c.client.Get(ctx, keysKey).Result()
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

	// -----------------------------------------------------------------------------------
	delKeys := append(keys, idNameKey, keysKey)
	_, err = c.client.Del(ctx, delKeys...).Result()
	if err != nil {
		return result.Error[T](fmt.Errorf("DelByRedis, %w", err))
	}

	return result.Success(getRes.Data())
}

func (c *redisCache[T]) genKey(prefix string, data ...any) string {
	return fmt.Sprintf("%s%s", prefix, crypto.MD5(data))
}

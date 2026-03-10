package cache

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/mel0dys0ng/song/pkg/crypto"
	"github.com/mel0dys0ng/song/pkg/result"
)

type (
	// lruCache LRU缓存结构体
	lruCache[T any] struct {
		keyClient   *expirable.LRU[string, []string]
		valueClient *expirable.LRU[string, T]
		size        int // 缓存数量
		ttl         time.Duration
	}
)

// newLRUCache 创建一个新的LRU缓存实例
func newLRUCache[T any](size int, ttl time.Duration) *lruCache[T] {
	return &lruCache[T]{
		keyClient:   expirable.NewLRU[string, []string](size, nil, ttl),
		valueClient: expirable.NewLRU[string, T](size*10, nil, ttl),
		size:        size,
		ttl:         ttl,
	}
}

// isSet 检查LRU缓存是否已设置
func (c *lruCache[T]) isSet() bool {
	return c != nil
}

// isRetryEnable 检查LRU缓存是否支持重试
func (c *lruCache[T]) isRetryEnable() bool {
	return false
}

// getType 获取缓存类型
func (c *lruCache[T]) getType() string {
	return TypeLRU
}

// get 从LRU缓存中获取数据
func (c *lruCache[T]) get(ctx context.Context, key any, prefix string) (res result.Interface[T]) {
	var data T

	keyKey := c.genKey(prefix, key)
	idNameKeys, _ := c.keyClient.Get(keyKey)
	if len(idNameKeys) > 0 {
		data, _ = c.valueClient.Get(idNameKeys[0])
	}

	return result.Success(data)
}

// set 将数据设置到LRU缓存中
func (c *lruCache[T]) set(ctx context.Context, key, name any, prefix string, id any, value T) (res result.Interface[T]) {
	keyKey := c.genKey(prefix, key)
	idNameKeys, _ := c.keyClient.Get(keyKey)
	if len(idNameKeys) == 0 {
		idNameKeys = []string{c.genKey(prefix, name, id)}
		c.keyClient.Add(keyKey, idNameKeys)
	}

	idKey := c.genKey(prefix, id)
	if v, _ := c.keyClient.Get(idKey); len(v) == 0 {
		c.keyClient.Add(idKey, idNameKeys)
	}

	idNameKey := idNameKeys[0]
	keys, _ := c.keyClient.Get(idNameKey)
	if keys == nil {
		keys = []string{}
	}
	for _, v := range []string{keyKey, idKey} {
		if !slices.Contains(keys, v) {
			keys = append(keys, v)
		}
	}

	c.keyClient.Add(idNameKey, keys)
	c.valueClient.Add(idNameKey, value)

	return result.Success(value)
}

// del 从LRU缓存中删除数据
func (c *lruCache[T]) del(ctx context.Context, key any, prefix string) (res result.Interface[T]) {
	getRes := c.get(ctx, key, prefix)
	if getRes.Err() != nil {
		return getRes
	}

	keyKey := c.genKey(prefix, key)
	defer func() {
		_ = c.keyClient.Remove(keyKey)
	}()

	idNameKeys, _ := c.keyClient.Get(keyKey)
	if len(idNameKeys) == 0 {
		var data T
		return result.Success(data)
	}

	idNameKey := idNameKeys[0]
	_ = c.valueClient.Remove(idNameKey)

	keys, _ := c.keyClient.Get(idNameKey)
	for _, v := range keys {
		_ = c.keyClient.Remove(v)
	}

	return result.Success(getRes.Data())
}

// genKey 生成缓存键
func (c *lruCache[T]) genKey(prefix string, data ...any) string {
	return fmt.Sprintf("%s%s", prefix, crypto.MD5(data))
}
package cache

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/mel0dys0ng/song/pkgs/utils/crypto"
	"github.com/mel0dys0ng/song/pkgs/utils/result"
)

type (
	lruCache[T any] struct {
		keyClient   *expirable.LRU[string, []string]
		valueClient *expirable.LRU[string, T]
		size        int // 缓存数量
		ttl         time.Duration
	}
)

func newLRUCache[T any](size int, ttl time.Duration) *lruCache[T] {
	return &lruCache[T]{
		keyClient:   expirable.NewLRU[string, []string](size, nil, ttl),
		valueClient: expirable.NewLRU[string, T](size*10, nil, ttl),
		size:        size,
		ttl:         ttl,
	}
}

func (c *lruCache[T]) isSet() bool {
	return c != nil
}

func (c *lruCache[T]) isRetryEnable() bool {
	return false
}

func (c *lruCache[T]) getType() string {
	return TypeLRU
}

// get 获取缓存数据
func (c *lruCache[T]) get(ctx context.Context, key any, prefix string) (res result.Interface[T]) {
	var data T

	keyKey := c.genKey(prefix, key)
	idNameKeys, _ := c.keyClient.Get(keyKey)
	if len(idNameKeys) > 0 {
		data, _ = c.valueClient.Get(idNameKeys[0])
	}

	return result.Success(data)
}

// get 设置缓存数据
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
	for _, v := range []string{keyKey, idKey} {
		if !slices.Contains(keys, v) {
			keys = append(keys, v)
		}
	}

	c.keyClient.Add(idNameKey, keys)
	c.valueClient.Add(idNameKey, value)

	return result.Success(value)
}

// del 删除缓存数据
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

func (c *lruCache[T]) genKey(prefix string, data ...any) string {
	return fmt.Sprintf("%s%s", prefix, crypto.MD5(data))
}

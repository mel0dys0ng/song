package cache

import (
	"context"

	"github.com/mel0dys0ng/song/pkgs/utils/result"
	"github.com/mel0dys0ng/song/pkgs/utils/retry"
	"github.com/mel0dys0ng/song/pkgs/utils/sys"
)

const (
	TypeLRU   = "lru"
	TypeRedis = "redis"
)

type (
	Cache[T any] struct {
		keyPrefix  string
		retryConf  retryConf
		lruCache   *lruCache[T]
		redisCache *redisCache[T]
		isZero     func(data T) bool
		dataId     func(data T) any
	}

	Option[T any] struct {
		apply func(c *Cache[T])
	}

	retryConf struct {
		enable       bool
		singleflight bool
		options      []retry.Option
	}

	cacheInterface[T any] interface {
		isSet() bool
		getType() string
		isRetryEnable() bool
		genKey(prefix string, data ...any) string
		set(ctx context.Context, key, name any, prefix string, id any, value T) result.Interface[T]
		get(ctx context.Context, key any, prefix string) result.Interface[T]
		del(ctx context.Context, key any, prefix string) result.Interface[T]
	}

	SetFunc[T any] func(ctx context.Context) result.Interface[T]
)

/*
New 返回缓存对象.
RedisCache和LRUCache至少设置一个，没设置则不支持对应类型的缓存.
Retry 设置重试和singleflight.
LRUCache不支持retry和singleflight.
*/
func New[T any](opts ...Option[T]) (c *Cache[T]) {
	c = &Cache[T]{retryConf: retryConf{enable: true}}
	for _, opt := range opts {
		if opt.apply != nil {
			opt.apply(c)
		}
	}

	if c.isZero == nil || c.dataId == nil {
		sys.Panic("isZero or dataId func is nil")
	}

	return
}

// GetOrSet 从多级缓存（LRU和Redis）中获取数据，若均未找到则调用set函数回源获取并更新缓存
//
//   - 1）若开启lru和redis cache，则优先级: lru > redis
//   - 1.1）若LRU缓存数据不存在，而redis缓存数据存在，则使用redis缓存数据更新LRU缓存数据；
//   - 1.2）若LRU和redis缓存数据都不存在，则使用源数据更新LRU和redis缓存数据）。
//   - 2）若lru和redis cache仅开启一个，缓存数据不存在，则使用源数据更新缓存数据。
//   - 3）retry默认开启，可关闭可配置。
//   - 4）singleflight默认关闭，可配置，开启之后singleflight key为缓存key。
//
// 参数:
//   - ctx: 上下文，用于控制请求生命周期，如超时或取消
//   - name: 缓存名称，用于区分不同业务
//   - key: 缓存键，支持任意类型
//   - set: 缓存未命中时用于生成数据的函数
//
// 返回值:
//   - result.Interface[T]: 包含获取结果或错误信息的包装对象
func (c *Cache[T]) GetOrSet(ctx context.Context, name, key any, set SetFunc[T]) (res result.Interface[T]) {
	// 初始化错误结果和缓存未命中标记
	res = result.Error[T](nil)
	isLRUCacheNotFound, isRedisCacheNotFound := false, false

	// 按缓存层级顺序尝试获取（先LRU后Redis）
	for _, cache := range []cacheInterface[T]{c.lruCache, c.redisCache} {
		cache := cache

		// 跳过未启用的缓存组件
		if !cache.isSet() {
			continue
		}

		// 定义基础获取逻辑
		handler := func(ctx context.Context) result.Interface[T] {
			return cache.get(ctx, key, c.keyPrefix)
		}

		// 根据重试配置选择执行方式
		if cache.isRetryEnable() {
			res = retryDo(ctx, &retryDoRequest[T, T]{
				key:          cache.genKey(c.keyPrefix, key, name, "getorset"),
				singleflight: c.retryConf.singleflight,
				handler:      handler,
				cache:        c,
			})
		} else {
			res = handler(ctx)
		}

		// 成功获取有效数据时的处理
		if res.Err() == nil && !c.isZero(res.Data()) {
			// 当从Redis获取到数据且LRU缓存未命中时，异步更新LRU缓存
			if c.isRedisCache(cache) && c.lruCache.isSet() && isLRUCacheNotFound {
				_ = c.set(ctx, []cacheInterface[T]{c.lruCache}, name, key,
					func(ctx context.Context) result.Interface[T] {
						return res
					},
				)
			}
			return
		}

		// 更新缓存未命中状态标记
		if c.isLRUCache(cache) {
			isLRUCacheNotFound = true
		} else {
			isRedisCacheNotFound = true
		}
	}

	// 多级缓存均未命中时的回源处理
	if c.isZero(res.Data()) || (isLRUCacheNotFound || isRedisCacheNotFound) {
		return c.set(ctx, []cacheInterface[T]{c.lruCache, c.redisCache}, name, key, set)
	}

	return
}

// set 方法用于将数据设置到缓存链中的每个可用缓存节点
// ctx: 上下文对象，用于控制请求生命周期
// name: 缓存名称，用于区分不同业务
// list: 缓存接口列表，包含多个层级的缓存实现
// key: 缓存键值，用于标识存储的数据
// set: 数据获取函数，当缓存未命中时用于获取源数据
// 返回值: 缓存操作结果对象，包含数据或错误信息
func (c *Cache[T]) set(ctx context.Context, list []cacheInterface[T], name, key any, set SetFunc[T]) (res result.Interface[T]) {
	var setRes result.Interface[T] // 源数据获取结果缓存，避免重复获取

	// 遍历缓存链中的每个缓存实现
	for _, cache := range list {
		cache := cache // 创建局部变量避免闭包问题

		// 跳过不支持设置或配置不设置的缓存类型
		if !cache.isSet() {
			continue
		}

		// 定义数据获取和缓存设置的处理闭包
		handler := func(ctx context.Context) result.Interface[T] {
			// 保证源数据只获取一次（首次需要时触发）
			if setRes == nil {
				// 执行源数据获取函数
				setRes = set(ctx)
				// 当数据有效时设置到当前缓存节点
				if setRes.Err() == nil && !c.isZero(setRes.Data()) {
					id := c.dataId(setRes.Data())
					return cache.set(ctx, key, name, c.keyPrefix, id, setRes.Data())
				}
			}
			return setRes
		}

		// 根据缓存配置选择执行策略
		if cache.isRetryEnable() {
			// 启用重试机制的执行路径（禁用 singleflight 并发控制）
			res = retryDo(ctx, &retryDoRequest[T, T]{
				key:          cache.genKey(c.keyPrefix, key, name, "set"),
				singleflight: false,
				handler:      handler,
				cache:        c,
			})
		} else {
			// 直接执行处理函数的路径
			res = handler(ctx)
		}

		// 遇到失败立即终止处理流程
		if res.Err() != nil {
			return
		}
	}

	return
}

// Del 从多级缓存中删除指定键以及对于ID相关的键，支持LRU本地缓存和Redis远程缓存两级删除
//
// 若开启lru和redis cache，则删除优先级: redis > lru
// retry默认开启，可关闭可配置，配置默认为retry的默认配置
//
// 参数:
//   - ctx: 上下文对象，用于传递超时、取消信号等
//   - key: 要删除的缓存键，支持任意类型
//
// 返回值:
//   - result.Interface[T]: 删除操作结果，包含操作状态和可能残留的数据
func (c *Cache[T]) Del(ctx context.Context, key any) (res result.Interface[T]) {
	var data T

	// 多级缓存处理：按顺序尝试从LRU缓存和Redis缓存删除
	for _, cache := range []cacheInterface[T]{c.lruCache, c.redisCache} {
		cache := cache

		if !cache.isSet() {
			continue
		}

		// 创建删除操作的闭包，携带键前缀配置
		handler := func(ctx context.Context) result.Interface[T] {
			return cache.del(ctx, key, c.keyPrefix)
		}

		// 重试逻辑处理：当启用重试且缓存键有效时，使用带重试机制的删除操作
		if cache.isRetryEnable() {
			res = retryDo(ctx, &retryDoRequest[T, T]{
				key:          cache.genKey(c.keyPrefix, key, "delete"),
				singleflight: c.retryConf.singleflight,
				handler:      handler,
				cache:        c,
			})
		} else {
			res = handler(ctx)
		}

		if res.Err() != nil {
			continue
		}

		// 保留非空数据（可能来自删除操作的响应）
		if !c.isZero(res.Data()) {
			data = res.Data()
		}
	}

	// 最终返回最近一次成功操作获取的数据（如果有）
	res.SetData(data)
	return
}

func (c *Cache[T]) isLRUCache(cache cacheInterface[T]) bool {
	return cache.getType() == TypeLRU
}

func (c *Cache[T]) isRedisCache(cache cacheInterface[T]) bool {
	return cache.getType() == TypeRedis
}

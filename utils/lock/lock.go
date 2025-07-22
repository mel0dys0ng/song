package lock

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	_TTLDefault      = 10 * time.Second
	_TimeoutDefault  = 1 * time.Second
	renewMinInterval = 200 * time.Millisecond
)

var (
	ErrorKeyEmpty                = errors.New("lock key is empty")
	ErrorRedisClientNil          = errors.New("redis client is nil")
	ErrorLockTimeout             = errors.New("lock timeout")
	ErrorLockRenewCanceledUnlock = errors.New("lock renew canceled, unlock")
	ErrorLockRenewCanceledPanic  = errors.New("lock renew canceled, panic")
	ErrorLockAcquiredByOther     = errors.New("lock has been acquired by the other process")
	ErrorLockNotHeld             = errors.New("lock not held by current client")
)

type (
	Lock struct {
		options
	}

	Core struct {
		C           chan error              // 续期取消错误通知
		key         string                  // 业务key
		value       string                  // uuid，解决错误删除锁的问题
		renewCancel context.CancelCauseFunc // 续期协程取消函数
		options
	}

	Option struct {
		apply func(l *Lock)
	}

	options struct {
		redisClient redis.UniversalClient
		TTL         time.Duration // 锁超时时间
		Timeout     time.Duration // 锁等待时间
	}
)

func RedisClient(client redis.UniversalClient) Option {
	return Option{
		apply: func(l *Lock) {
			l.redisClient = client
		},
	}
}

func TTL(t time.Duration) Option {
	return Option{
		apply: func(l *Lock) {
			l.TTL = t
		},
	}
}

func Timeout(t time.Duration) Option {
	return Option{
		apply: func(l *Lock) {
			l.Timeout = t
		},
	}
}

/*
New 创建并初始化一个分布式锁实例。

该函数通过可选参数配置锁属性，包括TTL、Timeout和Redis客户端等。
若未提供Redis客户端，将返回对应错误;若TTL或Timeout未设置或无效，则使用默认值。

参数:

	opts ...Option - 可选的配置函数，用于自定义锁属性

返回值:

	lock *Lock - 初始化后的锁实例指针（出错时为nil）
	err error - 错误信息，可能包含：
	            - ErrorRedisClientNil: Redis客户端未配置
*/
func New(opts ...Option) (lock *Lock, err error) {
	lock = &Lock{
		options{
			TTL:     _TTLDefault,
			Timeout: _TimeoutDefault,
		},
	}

	for _, opt := range opts {
		opt.apply(lock)
	}

	if lock.redisClient == nil {
		err = ErrorRedisClientNil
		return
	}

	if lock.TTL <= 0 {
		lock.TTL = _TTLDefault
	}

	if lock.Timeout <= 0 {
		lock.Timeout = _TimeoutDefault
	}

	return
}

// Lock 尝试获取分布式锁。
//
// 参数:
//
//	ctx: 上下文，用于控制请求的生命周期（如超时或取消）。
//	key: 要锁定的资源唯一标识符，不可为空。
//
// 返回值:
//
//	c: 成功时返回锁的核心控制对象，包含锁的详细信息；失败时为nil。
//	ok: 是否成功获取到锁。
//	err: 操作过程中发生的错误（如参数无效、底层存储错误等）。
func (l *Lock) Lock(ctx context.Context, key string) (c *Core, ok bool, err error) {
	if len(key) == 0 {
		err = ErrorKeyEmpty
		return
	}

	c = &Core{
		key:     key,
		value:   uuid.New().String(),
		options: l.options,
		C:       make(chan error, 1),
	}

	ok, err = c.lock(ctx)

	return
}

// Lockf 使用格式化字符串和参数生成锁键，并尝试获取锁。
//
// 参数:
//
//	ctx: 上下文，用于控制请求生命周期（如超时、取消）。
//	format: 格式化字符串模板，用于生成锁键。
//	value: 可变参数，用于填充格式化字符串的占位符。
//
// 返回值:
//
//	*Core: 成功获取锁时返回的核心对象，用于后续解锁操作。
//	bool: 表示是否成功获取锁（true=成功, false=未获取）。
//	error: 获取锁过程中发生的错误（如上下文取消、系统错误等）。
func (l *Lock) Lockf(ctx context.Context, format string, value ...any) (*Core, bool, error) {
	return l.Lock(ctx, fmt.Sprintf(format, value...))
}

func (c *Core) lock(ctx context.Context) (ok bool, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("lock panic: %v./n  stack: %+v", er, debug.Stack())
		}
	}()

	// 加入超时控制，避免获取锁超时长时间阻塞
	timeoutCtx, cancel := context.WithTimeoutCause(ctx, c.Timeout, ErrorLockTimeout)
	defer cancel()

	ok, err = c.redisClient.SetNX(timeoutCtx, c.key, c.value, c.TTL).Result()
	if err != nil {
		if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
			err = fmt.Errorf("%w: %v", context.Cause(timeoutCtx), err)
		}
		return
	}

	// 锁已被占用
	if !ok {
		err = ErrorLockAcquiredByOther
		return
	}

	// 锁获取成功，协程续期，TTL过短时禁用续期
	if c.TTL > time.Second {
		renewCtx, cancel := context.WithCancelCause(ctx)
		c.renewCancel = cancel
		go c.renew(renewCtx)
	}

	return
}

func (c *Core) renew(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			// 续期panic，通知业务
			c.C <- fmt.Errorf("panic: %v.\n stack: %v.", r, debug.Stack())
		}
	}()

	// 续期间隔：TTL-1秒 或 TTL/3（取正值）
	d := c.TTL - time.Second
	if d <= 0 {
		d = c.TTL / 3
	}

	if d <= renewMinInterval {
		d = renewMinInterval // 最小间隔
	}

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// 业务执行完毕，锁已取消, 停止续期
			return
		case <-ticker.C:
			// 续期
			if err := c.renewTTL(ctx); err != nil {
				// 续期失败，通知业务
				c.C <- err
				return
			}
		}
	}
}

func (c *Core) Close() {
	close(c.C) // 避免通道泄漏
}

// Unlock 释放锁
func (c *Core) Unlock(ctx context.Context) (err error) {
	// 取消续期协程，（如果存在）
	if c.renewCancel != nil {
		c.renewCancel(ErrorLockRenewCanceledUnlock)
	}

	script := `
if redis.call('get', KEYS[1]) == ARGV[1] then
	return redis.call('del', KEYS[1])
else
	return 0
end
`

	res, err := c.redisClient.Eval(ctx, script, []string{c.key}, c.value).Int64()
	if res != 1 {
		err = ErrorLockNotHeld
	}

	return
}

// 锁续期
func (c *Core) renewTTL(ctx context.Context) (err error) {
	script := `
if redis.call('get', KEYS[1]) == ARGV[1] then
	return redis.call('pexpire', KEYS[1], ARGV[2])
else
	return 0
end
`

	res, err := c.redisClient.Eval(ctx, script, []string{c.key}, c.value, c.TTL.Milliseconds()).Int64()
	if res != 1 {
		err = ErrorLockNotHeld
	}

	return
}

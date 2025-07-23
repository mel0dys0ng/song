package cache

import (
	"context"
	"fmt"

	"github.com/mel0dys0ng/song/pkg/utils/result"
	"github.com/mel0dys0ng/song/pkg/utils/retry"
)

type retryDoRequest[K any, V any] struct {
	key          string
	singleflight bool
	handler      retry.HandlerFunc[K]
	cache        *Cache[V]
}

func retryDo[K any, V any](ctx context.Context, request *retryDoRequest[K, V]) (res result.Interface[K]) {
	if request == nil || request.cache == nil {
		return result.Error[K](fmt.Errorf("retryDoRequest is invalid"))
	}

	// 未开启重试
	if !request.cache.retryConf.enable {
		return request.handler(ctx)
	}

	// 开启重试
	l := len(request.cache.retryConf.options)
	opts := make([]retry.Option, 0, l+1)
	if l > 0 {
		copy(opts, request.cache.retryConf.options)
	}

	if request.singleflight { // 启用singleflight
		opts = append(opts, retry.SingleflightKey(request.key))
	} else { // 未启用singleflight
		opts = append(opts, retry.SingleflightKey(""))
	}

	return retry.Do(ctx, request.handler, opts...)
}

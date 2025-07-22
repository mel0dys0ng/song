package singleton

import (
	"sync"
)

var (
	instances = &sync.Map{}
	onces     = &sync.Map{}
)

// loadOrStore 辅助函数，从sync.Map加载值，不存在则创建并store
func loadOrStore[K comparable, V any](m *sync.Map, key K, fn func() V) (V, bool) {
	if value, ok := m.Load(key); ok {
		if v, ok := value.(V); ok {
			return v, true
		}
	}
	v := fn()
	m.Store(key, v)
	return v, false
}

// GetInstance 返回指定类型的单例实例
func GetInstance[T any](key string, fn func() T) (res T) {
	res, _ = loadOrStore(instances, key, func() (res T) {
		once, _ := loadOrStore(onces, key, func() *sync.Once { return &sync.Once{} })
		once.Do(func() { res = fn() })
		return
	})
	return
}

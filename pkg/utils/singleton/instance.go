package singleton

import "sync"

// entry 存储单例实例及其初始化控制
type entry struct {
	once  sync.Once
	value any
}

var (
	// m 全局存储 key 到 entry 的映射
	m sync.Map
)

// GetInstance 返回指定 key 的单例实例
// 注意：同一 key 必须始终用于相同类型 T，否则会 panic
func GetInstance[T any](key string, fn func() T) T {
	actual, _ := m.LoadOrStore(key, &entry{})
	e := actual.(*entry)
	e.once.Do(func() { e.value = fn() })
	return e.value.(T)
}

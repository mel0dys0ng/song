package singleton

import (
	"runtime"
	"sync"
)

// entry 存储单例实例及其初始化控制
type entry struct {
	once  sync.Once
	value any
}

var (
	// m 全局存储 key 到 entry 的映射
	m sync.Map
)

// Key 返回调用者的 PC 地址，用于作为单例实例的键
// 调用链: runtime.Callers -> CallerPC -> Key -> A
// skip=3 跳过 Callers、CallerPC、Key，返回 A 中调用 Key 的位置
func Key() uintptr {
	pc := CallerPC(2)
	return pc
}

// CallerPC 返回调用栈中第 skip 层的 PC 地址
// skip=0: runtime.Callers 自身
// skip=1: CallerPC 函数
// skip=2: 调用 CallerPC 的函数
// skip=3: 调用调用 CallerPC 的函数（再上一层）
func CallerPC(skip int) uintptr {
	var pcs [1]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	if n == 0 {
		return 0
	}
	return pcs[0]
}

// Once 返回单例实例，如果不存在则创建
// 参数:
//   - key: 实例的唯一标识符，通常使用调用者的 PC 地址
//   - fn: 创建实例的函数，在第一次调用时执行
//
// 返回值:
//   - T: 类型为T的单例实例
func Once[T any](key any, fn func() T) T {
	actual, _ := m.LoadOrStore(key, &entry{})
	e := actual.(*entry)
	e.once.Do(func() { e.value = fn() })
	return e.value.(T)
}

// Get 从单例映射中获取指定键的实例
// 参数:
//   - key: 实例的唯一标识符
//
// 返回值:
//   - T: 类型为T的实例，如果不存在则返回类型的零值
//   - bool: 如果实例存在则返回true，否则返回false
func Get[T any](key any) (T, bool) {
	actual, ok := m.Load(key)
	if !ok {
		return *new(T), false
	}

	e := actual.(*entry)
	return e.value.(T), true
}

// Clear 删除指定 key 的单例实例
func Clear(key any) {
	m.Delete(key)
}

// ClearAll 清空所有单例实例
func ClearAll() {
	keys := make([]any, 0)

	m.Range(func(key, value any) bool {
		keys = append(keys, key)
		return true
	})

	for _, key := range keys {
		m.Delete(key)
	}
}

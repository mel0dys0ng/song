package sljces

import "github.com/mel0dys0ng/song/utils/aob"

// Unique 函数接受一个切片，返回一个包含唯一值的新切片。
// 它会保留原切片中第一次出现的每个元素的顺序。
//
// 参数:
//   slice: 输入切片，元素类型 T 必须满足 comparable 约束（支持 == 和 != 操作）。
//
// 返回值:
//   []T: 去重后的新切片，包含输入切片中所有不重复的元素，按首次出现顺序排列。
func Unique[T comparable](data []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(data))

	for _, v := range data {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Index 从切片中安全获取指定索引位置的元素
//
// 该函数处理了索引越界、负索引和空切片等情况：
//   - 当索引超出范围时，返回默认值
//   - 支持负索引（从末尾开始计数）
//   - 自动处理索引回绕（超出长度时取模）
//   - 空切片直接返回默认值
//
// 参数：
//   data     : 目标切片，可以是任意类型的切片
//   index    : 要获取的索引位置（支持负索引）
//   defaultV : 索引无效时返回的默认值
//
// 返回值：
//   索引有效时返回对应元素，否则返回 defaultV
func Index[T any](data []T, index int, defaultV T) T {
	item := defaultV

	if dataLen := len(data); dataLen > 0 {
		absIndex := aob.Aorb(index >= 0, index, -index)
		normalized := aob.Aorb(absIndex > dataLen, absIndex%dataLen, absIndex)
		finalIndex := aob.Aorb(index >= 0, normalized, dataLen-normalized)
		if finalIndex >= 0 && finalIndex < dataLen {
			item = data[finalIndex]
		}
	}

	return item
}

// First 从切片中获取第一个元素，如果切片为空则返回默认值。
//
// 参数:
//   data: 输入的切片，从中获取第一个元素。如果切片为空，则返回默认值。
//   defaultV: 当切片为空时返回的默认值。
//
// 返回值:
//   切片中的第一个元素；若切片为空，则返回 defaultV。
func First[T any](data []T, defaultV T) T {
	return Index(data, 0, defaultV)
}

// Last 返回切片的最后一个元素。如果切片为空，则返回默认值。
//
// 参数:
//   data: 目标切片，从中获取最后一个元素。
//   defaultV: 切片为空时返回的默认值。
//
// 返回值:
//   T: 切片的最后一个元素或默认值。
func Last[T any](data []T, defaultV T) T {
	return Index(data, -1, defaultV)
}

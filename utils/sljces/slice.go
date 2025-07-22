package sljces

import "github.com/song/utils/aob"

type Slice[T comparable] struct {
	data []T
}

func New[T comparable](data []T) *Slice[T] {
	return &Slice[T]{data: data}
}

func (s *Slice[T]) Add(item ...T) {
	s.data = append(s.data, item...)
}

func (s *Slice[T]) First(defaultValue T) (item T) {
	return s.IndexOf(0, defaultValue)
}

func (s *Slice[T]) Last(defaultValue T) (item T) {
	return s.IndexOf(-1, defaultValue)
}

func (s *Slice[T]) IndexOf(index int, defaultValue T) (item T) {
	item = defaultValue
	if dataLen := len(s.data); dataLen > 0 {
		a := aob.Aorb(index >= 0, index, -index)
		b := aob.Aorb(a > dataLen, a%dataLen, a)
		index = aob.Aorb(index >= 0, b, dataLen-b)
		if index >= 0 && index < dataLen {
			item = s.data[index]
		}
	}
	return
}

func (s *Slice[T]) Contain(item T) bool {
	return s.contain(s.data, item)
}

func (s *Slice[T]) contain(data []T, item T) bool {
	for _, v := range data {
		if v == item {
			return true
		}
	}
	return false
}

func (s *Slice[T]) Contains(items ...T) bool {
	if len(items) == 0 {
		return false
	}

	for _, item := range items {
		if !s.Contain(item) {
			return false
		}
	}

	return true
}

func (s *Slice[T]) Unique() *Slice[T] {
	dataLen := len(s.data)
	if dataLen == 0 {
		return s
	}

	tmpK := make([]int, 0, dataLen)
	tmpV := make(map[T]struct{}, dataLen)
	for k, v := range s.data {
		if _, ok := tmpV[v]; !ok {
			tmpV[v] = struct{}{}
			tmpK = append(tmpK, k)
		}
	}

	var data []T
	for k := range tmpK {
		data = append(data, s.data[k])
	}

	s.data = data

	return s
}

func (s *Slice[T]) Filter(elements ...T) *Slice[T] {
	if len(elements) == 0 || len(s.data) == 0 {
		return s
	}

	var data []T
	for _, v := range s.data {
		if !s.contain(elements, v) {
			data = append(data, v)
		}
	}

	s.data = data

	return s
}

func (s *Slice[T]) Data() []T {
	return s.data
}

func (s *Slice[T]) Length() int {
	return len(s.data)
}

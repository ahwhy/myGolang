package mset

import (
	"sync"
)

// 参考 https://allenwu.itscoder.com/set-in-go

type Set struct {
	sync.RWMutex
	Mset map[interface{}]interface{}
}

func NewSet(item ...interface{}) *Set {
	s := &Set{}
	s.Mset = make(map[interface{}]interface{})
	s.Add(item...)

	return s
}

// Add 添加元素
func (s *Set) Add(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()

	for _, item := range items {
		s.Mset[item] = item
	}

	return nil
}

// Delete 删除元素
func (s *Set) Delete(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()

	for _, item := range items {
		delete(s.Mset, item)
	}

	return nil
}

// Modify 修改元素
func (s *Set) Modify(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()

	for _, item := range items {
		_, ok := s.Mset[item]
		if !ok {
			continue
		}

		switch value := item.(type) {
		case string:
			s.Mset[item] = value + "modify"
		case int:
			s.Mset[item] = value + 10
		case bool:
			s.Mset[item] = true
		case float64:
			s.Mset[item] = value + 3.1415926535
		}
	}

	return nil
}

// Contains 包含&&查询
func (s *Set) Contains(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.Mset[item]

	return ok
}

// Size 容量
func (s *Set) Size() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.Mset)
}

// Clear 清除
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()

	s.Mset = make(map[interface{}]interface{})
}

// Equel 比较
func (s *Set) Equel(other *Set) bool {
	// 如果两者Size不相等，返回false
	if s.Size() != other.Size() {
		return false
	}

	// 迭代查询遍历
	for key := range s.Mset {
		// 若不存在返回false
		if !other.Contains(key) {
			return false
		}
	}

	return true
}

// IsSubset 子集
func (s *Set) IsSubset(other *Set) bool {
	// s的size长于other，返回false
	if s.Size() > other.Size() {
		return false
	}

	// 迭代遍历
	for key := range s.Mset {
		if !other.Contains(key) {
			return false
		}
	}

	return true
}

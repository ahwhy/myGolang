package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// 参考 https://allenwu.itscoder.com/set-in-go
// var Exist = struct{}{}

type Set struct {
	sync.RWMutex
	mset map[interface{}]interface{}
	// mset map[interface{}]struct{}
}

// 初始化set
func New(item ...interface{}) *Set { // 获取Set的地址
	s := &Set{}
	// 声明map类型的数据结构
	s.mset = make(map[interface{}]interface{})
	s.Add(item...)
	return s
}

// 添加元素
func (s *Set) Add(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.mset[item] = item
	}
	return nil
}

// 删除元素
func (s *Set) Delete(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		delete(s.mset, item)
	}
	return nil
}

// 修改元素
func (s *Set) Modify(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		_, ok := s.mset[item]
		if !ok {
			continue
		}
		switch value := item.(type) {
		case string:
			s.mset[item] = value + "modity"
		case int:
			s.mset[item] = value + 10
		case bool:
			s.mset[item] = true
		case float64:
			s.mset[item] = value + 3.1415926
		}
	}
	return nil
}

// 包含&&查询
func (s *Set) Contains(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.mset[item]
	return ok
}

// 长度
func (s *Set) Size() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.mset)
}

// 清除
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.mset = make(map[interface{}]interface{})
}

// 比较
func (s *Set) Equel(other *Set) bool {
	// 如果两者Size不相等，返回false
	if s.Size() != other.Size() {
		return false
	}
	// 迭代查询遍历
	for key := range s.mset {
		// 若不存在返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 子集
func (s *Set) IsSubset(other *Set) bool {
	// s的size长于other，返回false
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.mset {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	rand.Seed(time.Now().Unix())
	a := New()
	b := New()

	a.Add("aa", 123, false, 456.789, 1000)
	b.Add("aa", 123, false, 456.789, 1000)

	a.Delete(1000)
	b.Modify("aa", 123, false, 456.789)

	log.Printf("%T %+v ", a, a.mset)
	log.Printf("%T %#v ", b, b.mset)

	b.Clear()
	b.Add("aa", 123, false, 456.789, 1000)

	log.Println(a.Contains("aa"), a.Size(), a.Equel(b), a.IsSubset(b))

	c := New()
	for i := 0; i < 10000; i++ {
		v := fmt.Sprintf("KEY_%d", i)
		go func() {
			c.Add(v)
			for i := 0; i < 1000; i++ {
				c.Add(rand.Intn(1000))
			}
		}()

		go func() {
			c.Modify(v)
			for i := 0; i < 1000; i++ {
				c.Modify(rand.Intn(1000))
			}
		}()
	}

	for i := 0; i < 8000; i++ {
		v := fmt.Sprintf("KEY_%d", rand.Intn(1000))
		go func() {
			c.Delete(v)
			for i := 0; i < 1000; i++ {
				c.Delete(rand.Intn(1000))
			}
		}()

		go func() {
			c.Modify(v)
			for i := 0; i < 1000; i++ {
				c.Modify(rand.Intn(1000))
			}
		}()
	}

	log.Println(c.Contains(rand.Intn(1000)), c.Size(), c.Equel(b), a.IsSubset(c))
	log.Println(c.mset)
}

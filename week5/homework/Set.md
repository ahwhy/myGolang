# 代码
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

/*
# 简单的
- 作业简单的题:实现一个线程安全的集合set，元素是string
- 要求有 NewSet方法初始化
- Add方法添加元素，添加重复元素可以去重
- Del方法删除元素
- Merge方法合并另一个set
- PrintElement 方法打印所有元素
- JudgeElement方法 检测输入的string 是否存在于set中
- 总之就是set相关的方法
- https://github.com/deckarep/golang-set

*/
type SafeSet struct {
	sync.RWMutex
	m map[string]struct{}
}

func NewSet() *SafeSet {
	return &SafeSet{
		m: make(map[string]struct{}),
	}
}

func (ss *SafeSet) Add(key string) {
	ss.Lock()
	defer ss.Unlock()
	ss.m[key] = struct{}{}
}

func (ss *SafeSet) Del(key string) {
	ss.Lock()
	defer ss.Unlock()
	delete(ss.m, key)
}

func (ss *SafeSet) JudgeElement(key string) bool {
	ss.RLock()
	defer ss.RUnlock()
	_, ok := ss.m[key]
	return ok
}

func (ss *SafeSet) PrintElement() []string {
	ss.RLock()
	defer ss.RUnlock()
	res := make([]string, 0)
	for k := range ss.m {
		res = append(res, k)
	}
	return res
}

func (ss *SafeSet) Merge(set *SafeSet) {
	ss.Lock()
	defer ss.Unlock()
	keys := set.PrintElement()
	for _, k := range keys {
		// 遍历的时候还没有对写锁解锁
		// 这时不能调用ss.Add，因为Add也要加写锁 ，锁中锁的问题
		
		ss.m[k] = struct{}{}

	}
}

func setAdd(set *SafeSet, n int) {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("key_%d", i)
		set.Add(key)
	}
}

func setDel(set *SafeSet, n int) {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("key_%d", i)
		set.Del(key)
	}
}

func main() {
	// 基础测试
	s1 := NewSet()
	setAdd(s1, 10)
	res := s1.PrintElement()
	fmt.Println(res)
	setDel(s1, 5)
	fmt.Println(s1.PrintElement())
	s2 := NewSet()
	setAdd(s2, 20)
	s1.Merge(s2)
	fmt.Println(s1.PrintElement())

	//测试线程安全
	go setAdd(s1, 100)
	go setDel(s1, 100)
	go fmt.Println(s1.PrintElement())
	go setAdd(s1, 200)
	go setDel(s1, 200)
	go fmt.Println(s1.PrintElement())
	go setAdd(s1, 300)

	time.Sleep(10 * time.Hour)
}

```

## 开源set库 github.com/deckarep/golang-set
- github.com/deckarep/golang-set
- 代码
```go
package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
)

func main() {
	requiredClasses := mapset.NewSet()
	requiredClasses.Add("Cooking")
	requiredClasses.Add("English")
	requiredClasses.Add("Math")
	requiredClasses.Add("Biology")

	scienceSlice := []interface{}{"Biology", "Chemistry"}
	scienceClasses := mapset.NewSetFromSlice(scienceSlice)

	electiveClasses := mapset.NewSet()
	electiveClasses.Add("Welding")
	electiveClasses.Add("Music")
	electiveClasses.Add("Automotive")

	bonusClasses := mapset.NewSet()
	bonusClasses.Add("Go Programming")
	bonusClasses.Add("Python Programming")

	//Show me all the available classes I can take
	allClasses := requiredClasses.Union(scienceClasses).Union(electiveClasses).Union(bonusClasses)
	fmt.Println(allClasses) //Set{Cooking, English, Math, Chemistry, Welding, Biology, Music, Automotive, Go Programming, Python Programming}

	//Is cooking considered a science class?
	fmt.Println(scienceClasses.Contains("Cooking")) //false

	//Show me all classes that are not science classes, since I hate science.
	fmt.Println(allClasses.Difference(scienceClasses)) //Set{Music, Automotive, Go Programming, Python Programming, Cooking, English, Math, Welding}

	//Which science classes are also required classes?
	fmt.Println(scienceClasses.Intersect(requiredClasses)) //Set{Biology}

	//How many bonus classes do you offer?
	fmt.Println(bonusClasses.Cardinality()) //2

	//Do you have the following classes? Welding, Automotive and English?
	fmt.Println(allClasses.IsSuperset(mapset.NewSetFromSlice([]interface{}{"Welding", "Automotive", "English"}))) //true

}

```
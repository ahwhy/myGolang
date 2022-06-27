# Golang-Datastructure  Golang的数据结构

## 一、数据结构

- 逻辑结构
	- 集合结构(离散结构)
	- 线性结构
		- 线性表、栈、队列
	- 非线性结构
		- 树结构
		- 图结构或网状结构

- 存储结构
	- 顺序存储结构
	- 链式存储结构
	- 索引结构
	- 哈希结构
	
## 二、链表

### 1. 单向链表
- 在底层结构上，单向链表通过指针将一组零散的内存块串联在一起
	- 内存块称为链表的"结点"(Node)
	- 为了将所有的结点串起来，每个链表的结点除了存储数据之外，还需要记录链上的下一个结点的地址
	- 这个记录下个结点地址的指针叫作后继指针 next

- 实现代码
```go
	// 定义节点
	func NewIntNode(v int) *Node {
		return &Node{Value: v}
	}
	type Node struct {
		// 需要存储的数据
		Value interface{}
		// 下一跳
		Next  *Node
	}
	func NewIntList(headValue int) *List {
		// 链表的头
		head := &Node{Value: headValue}
		return &List{
			head: head,
		}
	}
	// 定义链表结构
	type List struct {
		head *Node
	}

	func (l *List) AddNode(n *Node) {
		// 需要找到尾节点
		next := l.head.Next
		for next.Next != nil {
			next = next.Next
		}
		// 修改为节点
		next.Next = n
	}
	func (l *List) Traverse(fn func(n *Node)) {
		n := l.head
		for n.Next != nil {
			fn(n)                // func(n *list.Node){ fmt.Println(n) }
			n = n.Next
		}
		fn(n)
		fmt.Println()
	}
	func (l *List) InsertAfter(after, current *Node) error {
		// after --> current --> afterNext
		// 保存after的下一跳 
		afterNext := after.Next
		// 插入current，修改指向
		after.Next = current
		current.Next = afterNext
		return nil
	}
```

- 缺陷
	- 单向链表的情况下(单方向的)，要获取之前Node 需要遍历整个List
	- 以下功能并没有在单向链表中实现
		- 插入到指定Node的前面
		- 删除链表中的元素
	- 要高效解决这个问题，需要一个Previos指针，直接知道前一个Node的信息，而不是遍历

### 2. 双向链表
- 支持两个方向，每个结点除一个后继指针 next 指向后面的结点外，还有一个前驱指针 prev 指向前面的结点
	- 双向链表相比于单向链表需要额外增加一个空间来保存前驱结点的地址
	- 存储同样多的数据，双向链表要比单链表占用更多的内存空间
	- 可以支持双向遍历

- 实现代码
```go
	// 重新定义节点
	type Node struct {
		// 需要存储的数据
		Value interface{}
		// 下一跳
		Next *Node
		// 上一跳
		Prev *Node
	}
	func (l *List) AddNode(n *Node) {
		// 需要找到尾节点
		next := l.head
		for next.Next != nil {
			next = next.Next
		}
		// 修改为节点
		next.Next = n
		n.Prev = next
	}
	func (l *List) Len() int {
		len := 0
		n := l.head
		if n.Prev != nil {
			return -1
		}
		for n.Next != nil {
			n = n.Next
			len++
		}
		return len + 1
	}
	func (l *List) Get(idx int) interface{} {
		index := 0
		n := l.head
		for n.Next != nil {
			n = n.Next
			index++
			if index == idx {
				return n.Value
			}
		}
		return nil
	}
	func (l *List) InsertAfter(after, current *Node) error {
		// after --> current --> afterNext
		// 保存after的下一跳 
		afterNext := after.Next
		// 插入current，修改指向
		after.Next = current
		current.Next = afterNext
		// after <-- current <-- afterNext
		current.Prev = after
		afterNext.Prev = current
		return nil
	}
	func (l *List) InsertBefore(before, current *Node) error {
		// beforePrev <-- current <-- before
		// 保存before的上一跳 
		beforePrev := before.Prev
		// 插入current，修改指向
		before.Prev = current
		current.Prev = beforePrev
		// beforePrev --> current --> before
		current.Next = before
		beforePrev.Next = current
		return nil
	}
	func (l *List) Remove(current *Node) error {
		// prev --> current --> next
		prev := current.Prev
		next := current.Next
		prev.Next, next.Prev = next, prev
		return nil
	}
```

- 应用场景
	- LRU缓存淘汰 -> LRU(Least Recently Used)最近最少使用
		- 实现思路: 缓存的key放到链表中，头部的元素表示最近刚使用
			- 如果命中缓存，从链表中找到对应的key，移到链表头部
			- 如果没命中缓存
				- 如果缓存容量没超，放入缓存，并把key放到链表头部
				- 如果超出缓存容量，删除链表尾部元素，再把key放到链表头部
```go
	var cache map[int]string
	var lst *list.List
	const CAP = 10 //定义缓存容量的上限
	func init() {
		cache = make(map[int]string, CAP)
		lst = list.New()
	}
	func readFromDisk(key int) string {
		return "china"
	}
	func read(key int) string {
		if v, exists := cache[key]; exists { //命中缓存
			head := lst.Front()
			notFound := false
			for {
				if head == nil {
					notFound = true
					break
				}
				if head.Value.(int) == key { //从链表里找到相应的key
					lst.MoveToFront(head) //把key移到链表头部
					break
				} else {
					head = head.Next()
				}
			}
			if !notFound { //正常情况下不会发生这种情况
				lst.PushFront(key)
			}
			return v
		} else { //没有命中缓存
			v = readFromDisk(key) //从磁盘中读取数据
			cache[key] = v        //放入缓存
			lst.PushFront(key) //放入链表头部
			if len(cache) > CAP { //缓存已满
				tail := lst.Back()
				delete(cache, tail.Value.(int)) //从缓存是移除很久不使用的元素
				lst.Remove(tail)                //从链表中删除最后一个元素
				fmt.Printf("remove %d from cache\n", tail.Value.(int))
			}
			return v
		}
	}
	func TraversList(lst *list.List) {
		head := lst.Front()
		for head.Next() != nil {
			fmt.Printf("%v ", head.Value)
			head = head.Next()
		}
		fmt.Println(head.Value)
	}
	func ReverseList(lst *list.List) {
		tail := lst.Back()
		for tail.Prev() != nil {
			fmt.Printf("%v ", tail.Value)
			tail = tail.Prev()
		}
		fmt.Println(tail.Value)
	}
```

### 3. 循环链表 -> ring
- 将头节点的前趋指针指向尾节点，将尾节点的后驱指针指向头节点

- 它的表现为一个环(Ring)，可以正向转和反向转

- 实现代码
```go
	// ChangeToRing 将链表头尾相连成环
	func (l *List) ChangeToRing() {
		// 需要找到尾节点
		next := l.head
		for next.Next != nil {
			next = next.Next
		}
		head, tail := l.head, next
		// head  -->  tail
		head.Prev = tail
		// head  <--  tail
		tail.Next = head
	}
	func (l *List) Traverse(fn func(n *Node)) {
		loopCount := 1
		n := l.head
		for n.Next != nil {
			// 最多循环5轮
			if loopCount > 5 {
				return
			}
			fn(n)
			n = n.Next
			if n == l.head {
				loopCount++
			}
		}
		fn(n)
		fmt.Println()
	}
```

- 应用场景 -> 基于滑动窗口的统计
	- 最近100次接口调用的平均耗时、最近10笔订单的平均值、最近30个交易日股票的最高点
	- ring的容量即为滑动窗口的大小，把待观察变量按时间顺序不停地写入ring即可

### 4. 数组和链表的比较
- 数组简单易用，在实现上使用的是连续的内存空间
	- 可以借助 CPU 的缓存机制，预读数组中的数据，所以访问效率更高
	- 而链表在内存中并不是连续存储，对 CPU 缓存不友好，没办法有效预读

- 数组的缺点是大小固定，一经声明就要占用整块连续内存空间
	- 如果声明的数组过大，系统可能没有足够的连续内存空间分配给它，导致"内存不足(out of memory)"
	- 如果声明的数组过小，则可能出现不够用的情况
	- 这时只能再申请一个更大的内存空间，把原数组拷贝进去，非常费时
	- 链表本身没有大小的限制，支持动态扩容，这也是其与数组最大的区别

- 应用场景
	- 频繁的插入和删除用list
	- 频繁的遍历查询选slice

### 5. Go语言中的标准库
- contianer/list  双向链表

- container/ring  循环链表

## 二、栈

- 栈stack，一种运算受限的线性表
	- 限定仅在表尾进行插入和删除操作的线性表
	- 这一端被称为栈顶，另一端被称为栈底，是被封死的
	
- FILO(First In Last Out)，先进后出

- Go语言中的栈
	- 2个核心方法
		- Push
		- Pop
	- 实现代码 -> 使用slice
```go
	// 定义需要存入的元素对象
	// 这里Item是范型, 指代任意类型
	type Item interface{}
	
	// 构建函数
	func NewStack() *Stack {
		return &Stack{
			items: []Item{},
		}
	}
	type Stack struct {
		items []Item
	}
	// Push adds an Item to the top of the stack
	func (s *Stack) Push(item Item) {
		s.items = append(s.items, item)
	}
	// Pop removes an Item from the top of the stack
	func (s *Stack) Pop() Item {
		if s.IsEmpty() {
			return nil
		}
		item := s.items[len(s.items)-1]
		s.items = s.items[0 : len(s.items)-1]
		return item
	}
	// Len 栈的大小 
	func (s *Stack) Len() int {
		return len(s.items)
	}
	// IsEmpty 判断是否为空 
	func (s *Stack) IsEmpty() bool {
		return len(s.items) == 0
	}
	// Peek 获取栈顶元素的值 Peek
	func (s *Stack) Peek() Item {
		if s.IsEmpty() {
			return nil
		}
		return s.items[len(s.items)-1]
	}
	// Clear 清空栈 
	func (s *Stack) Clear() {
		s.items = []Item{}
	}
	// Search 查询某个值 距离栈顶的距离 
	func (s *Stack) Search(item Item) (pos int, err error) {
		for i := range s.items {
			if item == s.items[i] {
				return i, nil
			}
		}
		return 0, fmt.Errorf("item %s not found", item)
	}
	// 遍历栈 ForEach
	func (s *Stack) ForEach(fn func(Item)) {
		for i := range s.items {
			fn(i)
		}
	}
	// Sort 插入排序 把stack的元素从大到小进行排序 插入排序
	func (s *Stack) Sort() {
		// 准备一个辅助的stack, 另一个容器
		orderdStack := NewStack()
	
		for !s.IsEmpty() {
			// 然后开始的排序流程
			current := s.Pop()
	
			// orderdStack顶端大于current，应该将orderdStack顶端移至s，直到orderdStack顶端小于current
			for !orderdStack.IsEmpty() && current.(int) > orderdStack.Peek().(int) {
				s.Push(orderdStack.Pop())
			}
	
			// 此时 当前current 一定是 <= orderdStack顶端
			orderdStack.Push(current)
		}
	
		// 倒过来
		for !orderdStack.IsEmpty() {
			s.Push(orderdStack.Pop())
		}
	}
```
		
## 四、堆

### 1. 树(Tree)
- 树是n(n>=0)个结点的有限集
	- n=0时称为空树
		- 在任意一颗非空树中，有且仅有一个特定的称为根(Root)的结点
	- 当n>1时，其余结点可分为m(m>0)个互不相交的有限集T1、T2、......、Tn
		- 其中每一个集合本身又是一棵树，并且称为根的子树

- 二叉树
	- 在计算机科学中，二叉树是每个结点最多有两个子树的树结构
		- 通常子树被称作"左子树"(left subtree)和"右子树"(right subtree)
	- 相关名词
		- 根节点: 一棵树最上方的节点称为根节点
		- 父节点、子节点: 如果一个节点下面连接多个节点，那么该节点称为父节点，它下面的节点称为子节点
		- 叶子节点: 没有任何子节点的节点称为叶子节点
		- 兄弟节点: 具有相同父节点的节点互称为兄弟节点
		- 节点度: 节点拥有的子树数
		- 树的深度: 从根节点开始(其深度为0)自顶向下逐层累加的
		- 树的高度: 从叶子节点开始(其高度为0)自底向上逐层累加的
			- 对于树中相同深度的每个节点来说，它们的高度不一定相同，这取决于每个节点下面的叶子节点的深度
	- 满二叉树
		- 一棵深度为k且有2^k - 1个结点的二叉树称为满二叉树
			- 除最后一层无任何子节点外，每一层上的所有节点都有两个子节点，最后一层都是叶子节点
			- 特点
				- 一颗树深度为h，最大层数为k，深度与最大层数相同，k=h
				- 叶子节点数(最后一层)为2k−1
				- 第 i 层的节点数是: 2i−1
				- 总节点数是: 2^k-1，且总节点数一定是奇数
	- 完全二叉树
		- 一棵深度为k的有n个结点的二叉树，对树中的结点按从上至下、从左到右的顺序进行编号，如果编号为i(1≤i≤n)的结点与满二叉树中编号为i的结点在二叉树中的位置相同
		- 每一层都是紧凑靠左排列的, 上层排满了, 才能排下层, 每层先排满左节点，才能排右节点

### 2. 堆(Heap)
- 堆通常是一个可以被看做一棵完全二叉树的数组对象
	- 堆是一个完全二叉树
	- 堆中每一个节点的值都必须大于等于(或小于等于)其子树中每个节点的值
		- 大顶堆: 堆中每一个节点的值都必须大于等于其子树中每个节点的值
		- 小顶堆: 堆中每一个节点的值都必须小于等于其子树中每个节点的值

- 应用场景
	- 堆排序(Heapsort)是指利用堆这种数据结构所设计的一种排序算法
		- 构建堆O(n)
		- 删除堆顶O(nlogn)
	- 求集合中最大的K个元素
		- 用集合的前K个元素构建小根堆
		- 逐一遍历集合的其他元素，如果比堆顶小直接丢弃
		- 否则替换掉堆顶，然后向下调整堆
	- 把超时的元素从缓存中删除
		- 按key的到期时间把key插入小根堆中
		- 周期扫描堆顶元素，如果它的到期时间早于当前时刻，则从堆和缓存中删除，然后向下调整堆
	- Go语言中的堆
		- 使用数组(array)/切片(slice)进行存储，
			- 由于数组(array)/切片(slice)的特性，在时间和空间上都很高效
			- 节点在数组中的位置index 和它的父节点以及子节点的索引之间有一个映射关系
				parent(i) = floor((i - 1)/2)
				left(i)   = 2i + 1
				right(i)  = 2i + 2 // right(i) 就是简单的 left(i) + 1，左右节点总是处于相邻的位置
				- 示例
					array:  10    7    2    5     1
					index:  0     1    2    3     4
					level:  1     2    2    3     3
				- 二叉数，指定层可以容纳的元素是固定的，2^level这个是二叉数的特性决定的
				- 把数组按照2^n次方，切分成多层，每层放对应的元素，就实现了一维到二维的映射
		- 添加/删除元素
			- 添加元素 -> 将新的元素插入到数组的尾部 -> 重新计算修复堆属性
			- 删除元素 -> 弹出堆顶的元素 -> 将数组的尾部元素，填补堆顶位置 -> 重新计算修复堆属性
		- 实现堆
			- 堆化(Init)
			- 往堆里添加元素(Push)
			- 弹出堆顶的元素(Pop)
		- Go语言中，由container/heap包实现
			- 需要实现五个方法，来定义一个堆
```go
	type Interface interface {
		sort.Interface
		Push(x interface{}) // add x as element Len()
		Pop() interface{}   // remove and return element Len() - 1.
	}
	// sort.Interface
	type Interface interface {
		Less(i, j int) bool
		Len() int		
		Swap(i, j int)
	}
```

- 官方示例 `/usr/local/go/src/container/heap/example_intheap_test.go`
```go
	type IntHeap []int
	func (h IntHeap) Len() int              { return len(h) }
	// func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }       // 大顶堆 
	func (h IntHeap) Less(i, j int) bool    { return h[i] < h[j] }       // 小顶堆
	func (h IntHeap) Swap(i, j int)         { h[i], h[j] = h[j], h[i] }
	func (h *IntHeap) Push(x interface{})   { *h = append(*h, x.(int)) }	
	func (h *IntHeap) Pop() interface{} {
		old := *h
		n := len(old)
		x := old[n-1]
		*h = old[0 : n-1]
		return x
	}
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
	// Output:
	// minimum: 1
	// 1 2 3 5
```

- 高性能定时器

	
## 五、Trie树

- trie树又叫字典权

- term集合
	- 根节点是总入口，不存储字符
	- 对于英文，第个节点有26个子节点，子节点可以存到数组里
	- 中文由于汉字很多，用数组存子节点太浪费内存，可以用map存子节点
	- 从根节点到叶节点的完整路径是一个term
	- 从根节点到某个中间节点也可能是一个term，即一个term可能是另一个term的前缀
```go
	type TrieNode struct {
		Word     rune                                                        // 当前节点存储的字符；byte只能表示英文字符，rune可以表示任意字符
		Children map[rune]*TrieNode                                          // 孩子节点，用一个map存储
		Term     string
	}
	type TrieTree struct {
		root *TrieNode
	}
	// add 把words[beginIndex:]插入到Trie树中
	func (node *TrieNode) add(words []rune, term string, beginIndex int) {
		if beginIndex >= len(words) {                                        // words已经遍历完了
			node.Term = term
			return
		}
		if node.Children == nil {
			node.Children = make(map[rune]*TrieNode)
		}
		word := words[beginIndex]                                            //把这个word放到node的子节点中
		if child, exists := node.Children[word]; !exists {
			newNode := &TrieNode{Word: word}
			node.Children[word] = newNode
			newNode.add(words, term, beginIndex+1)                           //递归
		} else {                                                             
			child.add(words, term, beginIndex+1)                             //递归
		}
	}
	// AddTerm 增加一个Term
	func (tree *TrieTree) AddTerm(term string) {
		if len(term) <= 1 {
			return
		}
		words := []rune(term)
		if tree.root == nil {
			tree.root = new(TrieNode)
		}
		tree.root.add(words, term, 0)
	}
	// walk words[0]就是当前节点上存储的字符，按照words的指引顺着树往下走，最终返回words最后一个字符对应的节点
	func (node *TrieNode) walk(words []rune, beginIndex int) *TrieNode {
		if beginIndex == len(words)-1 {
			return node
		}
		beginIndex += 1
		word := words[beginIndex]
		if child, exists := node.Children[word]; exists {
			return child.walk(words, beginIndex)
		} else {
			return nil
		}
	}
	// traverseTerms 遍历一个Node下面所有的Term，注意要传数组的指针，才能真正修改这个数组
	func (node *TrieNode) traverseTerms(terms *[]string) {
		if len(node.Term) > 0 {
			*terms = append(*terms, node.Term)
		}
		for _, child := range node.Children {
			child.traverseTerms(terms)
		}
	}
	// Retrieve 检索一个Term
	func (tree *TrieTree) Retrieve(prefix string) []string {
		if tree.root == nil || len(tree.root.Children) == 0 {
			return nil
		}
		words := []rune(prefix)
		firstWord := words[0]
		if child, exists := tree.root.Children[firstWord]; exists {
			end := child.walk(words, 0)
			if end == nil {
				return nil
			} else {
				terms := make([]string, 0, 100)
				end.traverseTerms(&terms)
				return terms
			}
		} else {
			return nil
		}
	}
```

## 六、算法的评估

### 1. 定义
- 算法(Algorithm)是指用来操作数据、解决程序问题的一组方法

### 2. 排序算法
- 冒泡排序: 两个数比较大小，较大的数下沉，较小的数冒起来
```go
	func BubbleSort(numbers []int) []int {
		for i := range numbers {
			for j := 0; j < len(numbers)-1; j++ {
				// 当前值 numbers[i], 后一个值是多少 numbers[j+1]
				fmt.Printf("数据: 当前: %d, 比对: %d\n", numbers[j], numbers[j+1])
	
				// 比较2个数, 交换顺序, 大数沉底, 小数冒出
				if numbers[j+1] < numbers[j] {
					numbers[j], numbers[j+1] = numbers[j+1], numbers[j]
				}
			}
			fmt.Printf("第%d趟: %v\n", i+1, numbers)
		}
		return numbers
	}
```

- 选择排序: 在长度为N的无序数组中，第一次遍历n-1个数，找到最小的数值与第一个元素交换，第二次遍历n-2个数，找到最小的数值与第二个元素交换...第n-1次遍历，找到最小的数值与第n-1个元素交换，排序完成
```go
	func SelectSort(numbers []int) []int {
		for i := range numbers {
			// 拿到第一个的数, 就是numbers[i], 比如 3
			fmt.Printf("第%d趟: %d\n", i+1, numbers[i])
	
			// 依次和后面相邻的数比较
			for j := i + 1; j < len(numbers); j++ {
				fmt.Printf("数据 -->  当前数据: %d, 比对数据: %d\n", numbers[i], numbers[j])
				if numbers[i] > numbers[j] {
					// 如果当前数 > 后面的数据, 则交换位置
					numbers[i], numbers[j] = numbers[j], numbers[i]
					fmt.Printf("交换 -->  当前数据: %d, 比对数据: %d\n", numbers[i], numbers[j])
				}
			}
	
			fmt.Println("结果: ", numbers)
		}
	
		fmt.Println("最终结果", numbers)
		return numbers
	}
```

- 插入排序: 在要排序的一组数中，假定前n-1个数已经排好序，现在将第n个数插到前面的有序数列中，使得这n个数也是排好顺序的。如此反复循环，直到全部排号顺序
```go
	func NewNumberStack(numbers []int) *Stack {
		items := make([]Item, 0, len(numbers))
		for i := range numbers {
			items = append(items, numbers[i])
		}
		return &Stack{
			items: items,
		}
	}
```

- 快速排序: 快速排序是对冒泡排序的一种改进，也属于交换类的排序算法

- 其他 
	- Go语言中内置排序 -> sort包，用于对象的排序, 参与排序的对象必须实现比较方法
```go
	// Sort sorts data.
	// It makes one call to data.Len to determine n and O(n*log(n)) calls to
	// data.Less and data.Swap. The sort is not guaranteed to be stable.
	func Sort(data Interface) {
		...
	}
	// 实现一个IntSlice结构
	func NewIntSlice(numbers []int) IntSlice {
		return IntSlice(numbers)
	}
	
	type IntSlice []int
	
	func (s IntSlice) Len() int { return len(s) }
	
	func (s IntSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
	
	func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
	// 比较函数
	func BuildInSort(numbers []int) []int {
		sort.Sort(IntSlice(numbers))
		return numbers
	}
```
			
### 3. 算法评估的维度
- 时间维度: 是指执行当前算法所消耗的时间，通常用「时间复杂度」来描述，可以估算出程序对处理器的使用程度
	- 时间频度: T(n) 通常，一个算法所花费的时间与代码语句执行的次数成正比，算法执行语句越多，消耗的时间也就越多
	- 渐进时间复杂度: 算法的时间复杂度函数为 T(n)=O(f(n))
	- 常见的算法时间复杂度由小到大依次为: Ο(1)＜Ο(log n)＜Ο(n)＜Ο(nlog n)＜Ο(n2)＜Ο(n3)＜…＜Ο(2^n)＜Ο(n!)

- 空间维度: 是指执行当前算法需要占用多少内存空间，通常用「空间复杂度」来描述，可以估算出程序对计算机内存的使用程度
	- 空间复杂度 O(1)
		- 不开辟额外空间，程序运行时，使用的空间是个常数
		- 冒泡排序、选择排序
	- 空间复杂度 O(n)
		- 程序使用的额外空间, 这个额外空间的大小和数据规模成线性关系
		- Map 的空间负责度 就介于0(1) ~ O(n)，不会因为一个元素就开辟一个bucket
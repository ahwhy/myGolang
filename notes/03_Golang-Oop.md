# Golang-Oop  Golang的编程思想

## 一、函数式编程

### 1. 函数式编程的定义
- 函数式编程，即 面向过程编程
	- 指忽略(通常是不允许)可变数据(以避免处理可改变的数据引发的边际效应)，忽略程序执行状态(不允许隐式的、隐藏的、不可见的状态)
	- 通过函数作为入参，函数作为返回值的方式进行计算，通过不断的推进(迭代、递归)这种计算，从而从输入得到输出的编程范式

### 2. 函数式编程模式
- 算子
	- 算子是一个函数空间到函数空间上的映射O: X -> X, 简单说来就是进行某种"操作"、动作
	- 与之对应的，就是被操作的对象，称之为操作数，业务中往往将可变部分抽象成算子，方便业务逻辑的替换
	- Map - Reduce -> 解耦数据结构和算法最常见的方式
		- Map模式
		- Filter模式
		- Reduce模式
		- Map/Reduce/Filter只是一种控制逻辑，真正的业务逻辑是在传给他们的数据和那个函数来定义的
			- 示例 -> 班级统计
				- 数据集合
				- 统计有多少学生数学成绩大于80
				- 列出数学成绩不及格的同学
				- 统计所有同学的数学平均分

```go
	// Map模式
	// item1 --map func--> new1
	// item2 --map func--> new2
	// item3 --map func--> new3
	// 使用Map函数将所有的字符串转换成大写
	func MapStrToUpper(data []string, fn func(string) string) []string {
		newData := make([]string, 0, len(data))
		for _, v := range data {
			newData = append(newData, fn(v))
		}
		return newData
	}
	list := []string{"abc", "def", "fqp"}
	out := MapStrToUpper(list, func(item string) string {
		return strings.ToUpper(item)
	})

	// Filter模式
	// item1 -- reduce func -->   x
	// item2 -- reduce func -->   itme2
	// item3 -- reduce func -->   x
	// 使用Filter函数将带"f"字段的数据筛选出
	func ReduceFilter(data []string, fn func(string) bool) []string {
		newData := []string{}
		for _, v := range data {
			if fn(v) {
				newData = append(newData, v)
			}
		}
		return newData
	}
	list := []string{"abc", "def", "fqp", "abc"}
	out := ReduceFilter(list, func(s string) bool {
		return strings.Contains(s, "f")
	})
	fmt.Println(out)

	// Reduce模式
	// item1 --|
	// item2 --|--reduce func--> new item
	// item3 --|
	// 使用Reduce函数进行求和
	func ReduceSum(data []string, fn func(string) int) int {
		sum := 0
		for _, v := range data {
			sum += fn(v)
		}
		return sum
	}
	list := []string{"abc", "def", "fqp", "abc"}
	out1 := ReduceSum(list, func(s string) int {          // 统计字符数量
		return len(s)
	})
	out2 := ReduceSum(list, func(s string) int {          // 出现过ab的字符串数量
		if strings.Contains(s, "ab") {
			return 1
		}
		return 0
	})
	
	// 班级统计
	type Class struct {
		Name     string     // 班级名称
		Number   uint8      // 班级编号
		Students []*Student // 班级学员
	}
	type Student struct {
		Name     string   // 名称
		Number   uint16   // 学号
		Subjects []string // 数学  语文  英语
		Score    []int    //  88   99   77
	}
```

- 修饰器 decorator 
	- 把一些函数装配到另外一些函数上，可以让代码更为的简单，也可以让一些"小功能型"的代码复用性更高，让代码中的函数可以像乐高玩具那样自由地拼装
	- 示例 -> HTTP中间件(net/http包)
	- 优化 -> 多个修饰器的管道 Pipeline
		- 在使用上，需要对函数一层层的套起来，特别在 decorator 比较多的话，影响代码的观感和阅读
		- 重构时，需要先写一个工具函数 -> 用来遍历并调用各个 decorator

```go
	// 修饰器
	func Hello(s string) {
		fmt.Printf("Hello, %s\n", s)
	}
	Hello("World")
	
	// HTTP中间件(net/http包)
	// - 打印AccessLog
	// - 为每个请求设置RequestID
	// - 添加BasicAuth
	func hello(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, http: %s", r.URL.Path)
	}
	http.ListenAndServe(":8848", http.HandlerFunc(hello))
	
	// 多个修饰器的管道 Pipeline
	type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc
	func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
		// A, B, C, D, E
		// E(D(C(B(A))))
		for i := range decors {
			d := decors[len(decors)-1-i] // iterate in reverse
			h = d(h)
		}
		return h
	}
```

- Functional Options
	- 不需要传入config参数, 函数签名优雅(缺省情况下 传nil很不优雅)
	- 不需要构建一个builder类来辅助
	- 示例
```go
	func NewServer(addr string, port int, options ...func(*Server)) (*Server, error) {}
	func WithProtocol(s string) func(s *Server) {}
	server := &Server{Addr: addr, Port: port}
	option(server)
```

### 3. 特点
- 高阶函数
	- 函数的参数和返回值都可以是函数

- 闭包
	- go函数编程的特点 -> 主要是闭包，简单来说，闭包就是在函数内使用外部自由变量

```go
	func Adder() func(int) int {
		// 自由变量
		num := 0
		return func(v int) int {
			num += v
			return num
		}
	}

	func callAdder() {
		addr := Adder()
		// plus 
		var res = 0
		for i := 0; i < 10; i++ {
			// 整个的累加过程作为变量放在循环的外部
			// 不断的对一个传入的数据进行加工
			res = addr(i)
			// 在进行plus的加工
			fmt.Printf("+.. %d=%d\n", i, res)
		}
	}
```

### 4. 函数式编程的应用场景
- 对于数据的长流程处理

- 类似流水线，装配模式
	- 可以随时增删流程
```go
	func fib() func() int {    // 斐波那契数列
		a, b := 0, 1
		return func() int {
			a, b = b, a+b
			return a
		}
	}
	fun := fib()
	fmt.Println(fun())
	fmt.Println(fun())
	fmt.Println(fun())
```

## 二、面向对象编程

### 1. 面向对象编程的定义
- 面向对象编程
	- Object Oriented Programming，OOP，一种程序设计思想
	- OOP把对象作为程序的基本单元，一个对象包含了数据和操作数据的函数

- 面向过程与面向对象
	- 面向过程
		- 面向过程的程序设计把计算机程序视为一系列的命令集合，即一组函数的顺序执行
		- 为了简化程序设计，面向过程把函数继续切分为子函数，即把大块函数通过切割成小块函数来降低系统的复杂度
	- 面向对象
		- 面向对象的程序设计把计算机程序视为一组对象的集合
		- 而每个对象都可以接收其他对象发过来的消息，并处理这些消息，计算机程序的执行就是一系列消息在各个对象之间传递

- 类和实例
	- 面向对象的设计思想是从自然界中来的，因为在自然界中，每一个实体都是对象(Object/Instance)，而这种实体的抽象类别就是类(Class)
		- 比如 车就是一个类, 而从你面前路过的福特汽车就是一个实例(Object)
	- 类(Class): 抽象的模板
	- 实例(Instance): 根据类创建出来的一个个具体的"对象"，每个对象都拥有相同的方法，但各自的数据可能不同

- Go语言的面向对象
	- Go没有class关键字，但是可以把Go当做面向对象的方式来编程
	- Go没有类，可以把struct作为类看待
	- 面向对象中，不希望直接访问包，而是通过开发者定义的New函数构造一个对象返回
		- 各类New函数，又称工厂函数
	- 类的方法: 给struct绑定的方法

- 三大特性: 封装、继承、多态

- 五大基本原则
	- 单一职责原则SRP(Single Responsibility Principle)         类的功能要单一
	- 开放封闭原则OCP(Open－Close Principle)                   一个模块对于拓展是开放的，对于修改是封闭的
	- 里式替换原则LSP(the Liskov Substitution Principle LSP)   子类可以替换父类出现在父类能够出现的任何地方。
	- 依赖倒置原则DIP(the Dependency Inversion Principle DIP)  高层次的模块不应该依赖于低层次的模块，他们都应该依赖于抽象。抽象不应该依赖于具体实现，具体实现应该依赖于抽象。
	- 接口分离原则ISP(the Interface Segregation Principle ISP) 设计时采用多个与特定客户类有关的接口比采用一个通用的接口要好。

- 其他
	- 结构体 -> 类
		- 构造函数 -> 创建类对象
	- 组合 -> 继承
		- 当前已有结构体 -> 扩展并使用
		- 匿名组合 -> 只匿名组合一个
	- 方法、接口 -> 多态
	
### 2. 封装
- 面向对象编程的一个重要特点就是数据封装

- 隐藏对象的属性和实现细节，仅对外提供公共访问方式，将变化隔离，便于使用，提高复用性和安全性

- 实例本身就拥有数据，要访问这些数据，没有必要从外面的函数去访问，可以直接在类的内部定义访问数据的函数

- 即将属性隐藏，提供相关接口供调用者访问和修改

### 3. 继承
- 在OOP程序设计中，当定义一个class的时候，可以从某个现有的class继承，新的class称为子类(Subclass)，而被继承的class称为基类、父类或超类(Base class、Super class)

- 提高代码复用性，继承是多态的前提

- 通过结构体的匿名嵌套，继承对应的字段和方法

### 4. 多态
- 父类或接口定义的引用变量可以指向子类或具体实现类的实例对象，提高了程序的拓展性

- 对象 在不同的条件下 有不同的行为

- Go语言中通过接口做多态
	- 通过interface{} 定义方法的集合
	- 多态 体现为: 各个结构体对象要实现 接口 中定义的所有方法
	- 统一的函数调用入口，传入的接口
	- 各个结构体对象中 绑定的方法只能多不能少，并且在中接口定义
	- 方法的签名要一致: 参数类型、参数个数，方法名称，函数返回值要一致
	- 多态的灵魂就是有个承载的容器，先把所有实现了接口的对象加进来，每个实例都要顾及的地方，直接遍历 调用方法即可
	
```go
	// 示例一
	/*
	体现多态
	告警通知的函数，根据不同的对象通知
	有共同的通知方法，每种对象自己实现
	*/
	type notifer interface {
		Init()                // 动作，定义的方法
		push()
		notify()
	}
	
	type user struct {
		name  string
		email string
	}
	type admin struct {
		name string
		age  int
	}
	
	func (u *user) Init()  {}
	func (u *admin) Init() {}
	
	func (u *user) push() {
		fmt.Printf("[普通用户][sendNotify to user %s]\n", u.name)
	}
	func (u *admin) push() {
		fmt.Printf("[管理员][sendNotify to user %s]\n", u.name)
	}
	
	func (u *user) notify() {
		fmt.Printf("[普通用户][sendNotify to user %s]\n", u.name)
	}
	func (u *admin) notify() {
		fmt.Printf("[管理员][sendNotify to user %s]\n", u.name)
	}
	
	func sendNotify(n notifer) {           // 入口 -> 多态的统一调用方法
		n.push()
		n.notify()
	}
	
	u1.push()
	a1.push()
	u1.notify()
	a1.notify()
	
	var n notifer
	n = &u1
	n.push()
	n.notify()
	n = &a1
	n.push()
	n.notify()

	ns := make([]notifer, 0)
	ns = append(ns, &u1, &a1)
	for _, n := range ns {
		sendNotify(n)
	}

	// 示例二
	/*
	1. 多个数据源
	2. QueryData方法做查数据
	3. PushData方法做写入数据
	*/
	type DataSource interface {
		PushData(data string)                 // 方法集合
		QueryData(name string) string
	}
	
	type redis struct {
		Name string
		Addr string
	}
	func (r *redis) PushData(data string) {
		log.Printf("[PushData][ds.name:%s][data:%s]", r.Name, data)
	}
	func (r *redis) QueryData(name string) string {
		log.Printf("[QueryData][ds.name:%s][data:%s]", r.Name, name)
		return name + "_redis"
	}
	
	type kafka struct {
		Name string
		Addr string
	}
	func (k *kafka) PushData(data string) {
		log.Printf("[PushData][ds.name:%s][data:%s]", k.Name, data)
	}
	func (k *kafka) QueryData(name string) string {
		log.Printf("[QueryData][ds.name:%s][data:%s]", k.Name, name)
		return name + "_kafka"
	}
	
	var Dm = make(map[string]DataSource)
	
	r := redis{
		Name: "redis",
		Addr: "1.1",
	}
	k := kafka{
		Name: "kafka",
		Addr: "2.2",
	}
	
	Dm["redis"] = &r                 // 注册数据源到承载的容器中
	Dm["kafka"] = &k
	for i := 0; i < 5; i++ {         // 推送数据
		key := fmt.Sprintf("key_%d", i)
		for _, ds := range Dm {
			ds.PushData(key)
		}
	}                                // 查询数据
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		for _, ds := range Dm {
			res := ds.QueryData(key)
			log.Println(res)
		}
	}
```

## 三、测试驱动开发

### 1. 测试驱动开发
- TDD
	- Test-Driven Development
	- 测试驱动开发，是敏捷开发中的一项核心实践和技术，也是一种设计方法论
	- TDD的原理是在开发功能代码之前，先编写单元测试用例代码，测试代码确定需要编写什么产品代码
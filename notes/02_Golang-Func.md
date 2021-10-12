# Golang-Func  Golong的函数

## 一、函数的定义
- 函数用于对代码块的逻辑封装，是提供代码复用的最基本方式，Go语言中有3种函数
	- 普通函数
	- 匿名函数(没有名称的函数)
	- 方法(定义在struct上的函数)

- Go实现了一级函数(first-class functions)，Go中的函数是高阶函数(high-order functions)
	- 函数是一个值，可以将函数赋值给变量，使得这个变量也成为函数
	- 函数可以作为参数传递给另一个函数
	- 函数的返回值可以是一个函数

- 定义语法
	- 签名
```
		func function_name( [parameter list] ) [return_types] {
			函数体
		}
```
	- func 函数由 `func` 开始声明
	- `function_name` 函数名称，函数名和参数列表一起构成了函数签名(signature)
	- `parameter list` 参数列表，参数就像一个占位符，当函数被调用时，可以将值传递给参数，这个值被称为实际参数。参数列表指定的是参数类型、顺序、及参数个数。参数是可选的，也就是说函数也可以不包含参数
	- `return_types` 返回类型，函数返回一列值。return_types 是该列值的数据类型；有些功能不需要返回值，这种情况下 return_types 不是必须的
	- 函数体 函数定义的代码集合
	
## 二、函数的参数
- 形参&&入参
	- 形参 定义函数时的参数
	- 入参 传递给函数的变量

- 类型合并
	- 在声明函数中若存在多个连续形参类型相同可只保留最后一个参数类型名
		func sum(x, y int) int {}

- 可变参数
	- 某些情况下函数需要处理形参数量可变，需要运算符 ARGS...TYPE 的方式声明可变参数函数或在调用时传递可变参数
	- 可变参数只能定义一个且只能在参数列表末端；在调用函数后，可变参数则被初始化为对应类型的切片(名为ARGS的slice,参数的数据类型都是TYPE)
		func max(a, b int, args ...int) int {}
	- 在调用函数时，也可以使用运算符 ... 将切片解包传递到可变参数函数中
		max(1, 2, slice[3]...)

- 值传递
	- 函数如果使用参数，该变量可称为函数的形参，形参就像定义在函数体内的局部变量
	- 值传递是指在调用函数时将实际参数复制一份传递到函数中，这样在函数中如果对参数进行修改，将不会影响到实际参数
	- 默认情况下，Go语言使用的是值传递(则先拷贝参数的副本，再将副本传递给函数)，即在调用过程中不会影响到实际参数

- 引用传递
	- 引用传递是指在调用函数时将实际参数的地址传递到函数中，那么在函数中对参数进行的修改，将影响到实际参数
	- 由于引用类型(slice、map、interface、channel)自身就是指针，所以这些类型的值拷贝给函数参数，函数内部的参数仍然指向它们的底层数据结构

- 值类型&引用类型     
	- 值类型和引用类型的差异在于赋值同类型新变量后，对新变量进行修改是否能够影响原来的变量，若不能影响则为值类型，若能影响则为引用类型
	- 值类型 数值、布尔、字符串、指针、数组、结构体等
	- 引用类型 切片、映射、接口等
	- 针对值类型可以借助指针修改原值
	- 针对值类型和引用类型在赋值后新旧变量的地址并不相同，只是引用类型在底层共享数据结构(其中包含指针类型元素)

- 其他
	- 变量是一个地址，也是一个引用
	- 引用表达的是关系，指针表达的是类型；  
	- A -- B 是引用关系，而 A 是指针
	
## 三、函数的返回值
- 多返回值
```
	func calcReturn(x, y int) (int, int, int, int) {
		return x + y, x - y, x  y, x  y
	}
```

- 命名返回值
```
	func calcReturnNamecalc(x, y int) (sum, difference, product, quotient int) {
		sum, difference, product, quotient = x + y, x - y, x  y, x  y
		return
	}
```

- 关于 `return`
	- return关键字中指定了参数时，返回值可以不用名称
		- 如果return省略参数，则返回值部分必须带名称
		- 但即使返回值命名了，return中也可以强制指定其它返回值的名称，也就是说return的优先级更高
	- return 中可以有表达式，但不能出现赋值表达式，这和其它语言可能有所不同
		- 例如 return a+b 是正确的，但 return c=a+b 是错误的

## 四、函数的递归

### 1. 定义
- 函数内部调用函数自身的函数称为递归函数
	- 退出条件基本上都使用退出点来定义，退出点常常也称为递归的基点，是递归函数的最后一次递归点，或者说没有东西可递归时就是退出点
	- 递归函数很可能会产生一大堆的goroutine(其它编程语言则是出现一大堆的线程、进程)，也很可能会出现栈空间内存溢出问题
	- 在其它编程语言可能只能设置最大递归深度或改写递归函数来解决这个问题，在Go中可以使用channel+goroutine设计的lazy evaluation来解决

### 2. 示例
- 阶乘
```
	// n*(n-1)...3*2*1
	func factorial(n int) int {
		if n = 0 {
				return -1
		} else if n == 1 { 
				return 1                     // 判断退出点
		} else {
				return n  factorial(n-1)     // 递归表达式
		}
	}
```

- 斐波那契数列   
```
	// f(n)=f(n-1)+f(n-2)且f(2)=f(1)=1
	func fib(n int) int {
		if n == 1  n == 2 {
				return 1
		}
		return fib(n-1) + fib(n-2)
	}
```

- 汉罗塔
	- 将所有a柱上的圆盘借助b柱移动到c柱，在移动过程中保证每个柱子的上面圆盘比下面圆盘小
	- a -> 开始 ；b -> 借助 ；c -> 终止
	- n : a -> c(b) ； 
	- n = 1 :  a -> c ； 
	- n > 1 : n - 1 (a -> b(c)) 、 a -> c ；n - 1 (b -> c(a))
```
		func tower(a, b, c string, layer int) {
			if layer = 0 {
					return
			}
			if layer == 1 {
					fmt.Printf(%s -  %sn, a, c)
					return
			}
			tower(a, c, b, layer-1)
			fmt.Printf(%s -  %sn, a, c)
			tower(b, a, c, layer-1)
		}
		tower(A, B, C, 3)
```

- 递归目录
	- 递归基点是文件，只要是文件就返回，只要是目录就进入
 
## 五、函数的类型
- 函数也可以赋值给变量，存储在数组、切片、映射中，也可作为参数传递给函数或作为函数返回值进行返回

- 声明&&初始化&&调用
```
	var callback func(n1, n2 int) (r1, r2, r3, r4 int)   // 定义函数类型变量，并使用零值nil进行初始化
	callback = calcReturn                                // 赋值为函数calcReturn
	callback(5, 2)                                       // 调用函数calcReturn
```
	- 声明函数类型变量f 为函数Add
```
		var f func(int, int) int = Add
		fmt.Println(f(4, 2))  6
```
	- 声明函数切片
```
		var fs []func(int, int) int
		fs = append(fs, Add, Sub, Mul, Div)
		fmt.Printf(%T;n%#vn, fs, fs)         // []func(int, int) int;
		                                     // []func(int, int) int{(func(int, int) int)(0xb6ef20), (func(int, int) int)(0xb6ef20), (func(int, int) int)(0xb6ef40), (func(int, int) int)(0xb6ef60)}
		for _, f = range fs{
			fmt.Println(f(4,2))
		}
```
	- 返回值为函数
```
		func genFunc() func() {
			if rand.Int()%2 == 0 {
				return sayHi
			} else {
				return sayHolle
			}
		}
		rand.Seed(time.Now().Unix())
		genFunc()
```

- 声明&&调用参数类型为函数的函数 *高阶函数
```
	func printResult(pf func(...string), list ...string) {    // 定义接收函数类型作为参数的函数
		pf(list...)
	}
	func line(list ...string) {                               // 回调函数，作为参数被传递的函数
		fmt.Print()
		for _, e = range list {
			fmt.Print(e)
			fmt.Print(t)
		}
		fmt.Println()
	}
	names = []string{aa, bb, cc}
	printResult(line, names...)
```

- 自定义函数类型&&调用参数类型为自定义函数类型的函数&&赋值变量并调用
```
	type addFunc func(x, y int) int                // 声明函数类型addFunc
	func asArg(fn addFunc) int {                   // 创建函数asArg使用声明函数类型addFunc作为参数
		return fn(2, 2)  2
	}
	ret = asArg(func(x, y int) int {               // 调用函数asArg并使用匿名函数传参
		return x + y
	}
```

## 六、匿名函数与闭包

### 1. 匿名函数
- 不需要定义名字的函数叫做匿名函数，常用做帮助函数在局部代码块中使用或作为其他函数的参数
	- 声明匿名函数并直接执行
```
		func(args){
			...
		}(parameters)
```
	- 使用匿名函数作为 printResult 的参数
```
		printResult(func(list ...string) {
			for i, v = range list{
				fmt.Printf(%d %s, i, v)
			}
		}, name...)
```
	- 声明自定义匿名函数类型
```
		type Callback func() error                      // 声明自定义匿名函数类型
		callback = map[string]Callback{}                // 赋值给map类型变量
		callback[add] = func(int string) error {        // 初始化为具体的匿名函数
			fmt.Println(add)
			return nil
		}
		callback[add]()
```

### 2. 闭包
- 闭包，匿名函数的一种，是指在函数内定义的匿名函数引用外部函数的变量(没有使用传参的方式)，只要匿名函数继续使用则外部函数赋值的变量不被自动销毁
	- 变量生成周期(内存中存在的时间)发生了变化 ，闭包不仅仅包含函数，还包含函数定义域和函数变量
```
		func addBase(base int) func(int) int {   // 定义闭包函数，返回一个匿名函数用于计算于base元素的和
			return func (num int) int {
				return base + num
			}
		}
		base2 = addBase(2)                       // 使用闭包函数
		fmt.Println(base2(3))
```

## 七、错误处理

### 1. error 接口
- error类型是个接口
```
	type error interface {
		Error() string
	}
```

- 函数调用时判断返回值
	- `if err != nil { ... }`
```
		for _, v = range [...]int{0, 1, 2, 3} {            // 处理函数返回的错误
			if r, err = division(6, v); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(r)
			}
		}
```

- Go语言通过 error 接口实现错误处理的标准模式，通过使用函数返回值列表中的最后一个值返回错误信息，将错误的处理交由程序员主动进行处理

- error接口的初始化方法
	- 通过 errors 包的 New 方法创建 `errors.New()`
```
		func division(n1, n2 int) (int, error){            // 定义除法函数，若除数为0则使用error返回错误信息
			if n2 == 0 {
				return 0, errors.New(除数为0)
			}
			return n1n2 , nil
		}
```
	- 通过通过fmt.Errorf方法创建方法创建
```
		err1, err2 = errors.New(error 1), fmt.Errorf(error %d, 2)
		fmt.Printf(%T, %T, %v, %v, err1, err2, err1, err2)               // errors.errorString, errors.errorString, error 1, error 2
```

### 2. 复杂的错误类型
- 以os包举例，其提供了 `LinkError、PathError、SyscallError` 的错误类型，这些 error 都是实现了 error接口的错误类型

- 可以用switch err.(type)判断类型
```
	file, err = os.Stat(test.txt)
	if err != nil {
		switch err.(type) {
		case os.PathError
			log.Printf(PathError)
		case os.LinkError
			log.Printf(LinkError)
		case os.SyscallError
			log.Printf(SyscallError)
		default
			log.Printf(unknow error)
		}
	} else {
		fmt.Println(file)
	}
```

### 3. 自定义error
- `errors.New()` 独立的error，基础的error

- 自定义结构体 - 原始错误的基础上再封自己的错误信息
	- 弊端 要定义很多 error结构体
```
	type MyError struct {
		err error
		msg string                      // 自定义的error字符串
	}
	func (e MyError) Error() string {
		return e.err.Error() + e.msg
	}
	err = errors.New(原始的错误 )
	newErr = MyError{
		err err,
		msg 自定义的错误,
	}
	fmt.Println(newErr.Error())
```

### 4. Error Wrapping 错误嵌套    golang 1.13
- 可以扩展error信息，使用 `fmt.ErrorF(newErrorStr %w,e)`

- 优势是不需要像上面一样定义结构体
```
	e = errors.New(原始的错误)
	w = fmt.Errorf("Wrap了一个新的错误", %w, e)
```

### 5. defer 函数
- defer 用户声明函数，不论函数是否发生错误都在函数执行最后执行(return之前)
	- 若使用defer声明多个函数，则按照声明的顺序，先声明后执行(堆)常用来做资源释放，记录日志等工作

- defer 的本质是，当在某个函数中使用了defer关键字，则创建一个独立的defer栈帧，并将该defer语句压入栈中，同时将其使用的相关变量也拷贝到该栈帧中(值拷贝的)
	- 因为栈是LIFO方式，所以先压栈的后执行
	- 因为是独立的栈帧，所以即使调用者函数已经返回或报错，也一样能在它们之后进入defer栈帧去执行
	
### 6. panic与 recover 函数
- panic
	- panic和recover函数用于处理运行时错误，当调用panic抛出错误，可以中断原有的控制流程，常用于不可修复性错误
	- 触发场景
		- 运行时错误会导致panic，比如数组越界、除0
		- 程序主动调用`panic(error)`
	- 执行顺序
		- 逆序执行当前goroutine的defer链(recover从这里介入)
		- 打印错误信息和调用堆栈
		- 调用exit(2)结束整个进程

- recover
	- recover函数用于终止错误处理流程，仅在defer语句的函数中有效，用于截取错误处理流程 recover 只能捕获到最后一个错误
		- 当未发生panic，且不存在panic，则recover函数得到的结果为nil
		- 当未发生panic，且存在panic，则recover函数得到的结果为panic传递的参数
		- recover只能获取到最后一次的panic的信息
```
			defer func() {
				fmt.Println(recover())
			}()
			var x, y int
			sum(x, y)
```
		- 在并发的场景中，需要在goroutine的启动函数里面专门编写recover，用于捕获当前goroutine的异常
# defer与异常

## defer

defer关键字可以让函数或语句延迟到函数语句块的最结尾时，即即将退出函数时执行，即便函数中途报错结束、即便已经panic()、即便函数已经return了，也都会执行defer所推迟的对象。

其实defer的本质是，当在某个函数中使用了defer关键字，则创建一个独立的defer栈帧，并将该defer语句压入栈中，同时将其使用的相关变量也拷贝到该栈帧中（显然是按值拷贝的）。因为栈是LIFO方式，所以先压栈的后执行。因为是独立的栈帧，所以即使调用者函数已经返回或报错，也一样能在它们之后进入defer栈帧去执行。

### defer的执行顺序

如果语句块内有多个defer，则defer的对象以LIFO(last in first out)的方式执行，也就是说，先定义的defer后执行

```go
fmt.Println("start")
defer fmt.Println(1)
defer fmt.Println(2)
defer fmt.Println(3)
defer fmt.Println(4)
fmt.Println("end")
```

执行结果:
```
start
end
4
3
2
1
```

### defer与匿名函数

```go
fmt.Println("func start")       
x := 10
defer func(x int) {
    fmt.Println("in defer: ", x)
}(x)
x = 30
fmt.Println("func end: ", x)
```

因为函数传参是值copy,所以x为10的值在defer定义的时候已经copy传入defer, 后面的修改并不会影响到defer中的值

我们也可以选择把变量的指针传达给defer， 这样外面的修改就是生效的, 例如

```go
fmt.Println("func start")
x := 10
defer func(x *int) {
    fmt.Println("in defer: ", *x)
}(&x)
x = 30
fmt.Println("func end: ", x)
```

### defer与闭包

当然最常用的就是直接使用闭包的方式:

```go
fmt.Println("func start")
x := 10
defer func() {
    fmt.Println("in defer: ", x)
}()

x = 30
fmt.Println("func end: ", x)
```

### defer的应用

defer有什么用呢？一般用来做善后操作，例如清理垃圾、释放资源，无论是否报错都执行defer对象。另一方面，defer可以让这些善后操作的语句和开始语句放在一起，无论在可读性上还是安全性上都很有改善，毕竟写完开始语句就可以直接写defer语句，永远也不会忘记关闭、善后等操作


## 异常处理: panic

panic()用于产生错误信息并终止当前的goroutine，一般将其看作是退出panic()所在函数以及退出调用panic()所在函数的函数

```go
func main) {
	fn()
}

func fn() {
	fmt.Println("start fn")
	panic("pannic in fn")
	fmt.Println("end fn")
}

// panic: pannic in fn
```

大部分场景下 panic都不是我们可以预判的, 比如下面
```go
var a *int
fmt.Println(*a) 
// panic: runtime error: invalid memory address or nil pointer dereference
// [signal 0xc0000005 code=0x0 addr=0x0 pc=0x4675e6]
```

由于panic会直接导致程序退出, 一般都不是我们期望的，比如:
```go
func main() {
	var x, y *int
	sum(x, y)
}

func sum(x, y *int) int {
	return *x + *y
}
```

如果不想panic直接退出程序，我们就需要捕获panic, Go语言内置的recover函数就是干这个的

## 异常捕获: recover

recover()用于捕捉panic()错误，并返回这个错误信息。但注意，即使recover()捕获到了panic()，但调用含有panic()函数的函数也会退出

比如, 如果我们放前面，由于在这个位置并未panic, 捕获为nil
```go
var x, y *int
fmt.Println(recover())
sum(x, y)
```

如果 我们放后面，则由于panic提前退出，根本执行不到我们的捕获代码
```go
var x, y *int
sum(x, y)
fmt.Println(recover())
```

而且我们写程序的时候 也完全预估不了哪里会panic, 所以正确的用法是, 函数调用后 再尝试捕获， 这时候我们就需要使用defer
```go
func main() {
	defer func() {
		fmt.Println(recover())
	}()

	var x, y *int
	sum(x, y)
}
```

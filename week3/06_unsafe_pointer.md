# 非类型安全指针

![](../../image/go-unsafe-pointer.png)

为了安全起见，Go指针在使用上相对于C指针有很多限制。 通过施加这些限制，Go指针保留了C指针的好处，同时也避免了C指针的危险性。

在基础篇中，我们介绍了Go指针是有类型限制的: *T, 这种具有类型限制的指针我们称为类型安全的指针,在使用上有些限制，比如:

1.Go指针不支持算术运算

```go
str := "pointer_test"
a := &str
// a++ a-- 是不合法的
```

2.一个指针类型的值不能被随意转换为另一个指针类型

```go
var x int64 = 20
a := &x
fmt.Println(a)
// 没有办法指针a的值进行如下转换: *int64 --> *uint64
```

3.一个指针类型的值不能被赋值给其它任意类型的指针值

```go
var x int64 = 20
a := &x
fmt.Println(a)

var y uint64 = 20
b := &y
fmt.Println(b)

// 我们不能进行 a = b
```

unsafe标准库包中提供的非类型安全指针（unsafe.Pointer）机制可以被用来打破上述Go指针的安全限制。 unsafe.Pointer类型类似于C语言中的void*。 但是，通常地，非类型安全指针机制不推荐在Go日常编程中使用

## 关于unsafe标准库包

非类型安全指针在Go中为一种特别的类型。 我们必须引入unsafe标准库包来使用非类型安全指针

非类型安全指针unsafe.Pointer被声明定义为:

```go
// ArbitraryType is here for the purposes of documentation only and is not actually
// part of the unsafe package. It represents the type of an arbitrary Go expression.
type ArbitraryType int

// Pointer represents a pointer to an arbitrary type. There are four special operations
// available for type Pointer that are not available for other types:
//	- A pointer value of any type can be converted to a Pointer.
//	- A Pointer can be converted to a pointer value of any type.
//	- A uintptr can be converted to a Pointer.
//	- A Pointer can be converted to a uintptr.
// Pointer therefore allows a program to defeat the type system and read and write
// arbitrary memory. It should be used with extreme care.
// ...
type Pointer *ArbitraryType
```

其中uintptr是被定义在build中内置类型
```go
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
```

有了Pointer类型，我们就可以进行如下转换:
+ 任意类型的指针的值 <--> Pointer
+ uintptr <--> Pointer


unsafe标准库包只提供了三个函数：
+ Alignof: 此函数用来取得一个值在内存中的地址对齐保证（address alignment guarantee）。 注意，同一个类型的值做为结构体字段和非结构体字段时地址对齐保证可能是不同的。 当然，这和具体编译器的实现有关。对于目前的标准编译器，同一个类型的值做为结构体字段和非结构体字段时的地址对齐保证总是相同的。 gccgo编译器对这两种情形是区别对待的。
+ Offsetof: 此函数用来取得一个结构体值的某个字段的地址相对于此结构体值的地址的偏移。 在一个程序中，对于同一个结构体类型的不同值的对应相同字段，此函数的返回值总是相同的。
+ Sizeof: 此函数用来取得一个值的尺寸（亦即此值的类型的尺寸）。 在一个程序中，对于同一个类型的不同值，此函数的返回值总是相同的

```go
m := Man{Name: "John", Age: 20}
fmt.Println(unsafe.Sizeof(m.Name), unsafe.Sizeof(m.Age), unsafe.Sizeof(m)) // 16 8 24

fmt.Println(unsafe.Offsetof(m.Name)) // 0
fmt.Println(unsafe.Offsetof(m.Age))  // 16
```

## 注意事项

通过使用这些转换规则，我们可以将任意两个类型安全指针转换为对方的类型，我们也可以将一个安全指针值和一个uintptr值转换为对方的类型
然而，尽管这些转换在编译时刻是合法的，但是它们中一些在运行时刻并非是合法和安全的。 这些转换摧毁了Go的类型系统（不包括非类型安全指针部分）精心设立的内存安全屏障

### Pointer与uintptr的区别

Go是一门支持垃圾回收的语言。 当一个Go程序在运行中，Go运行时（runtime）将不时地检查哪些内存块将不再被程序中的任何仍在使用中的值所引用并且回收这些内存块。 指针在这一过程中扮演着重要的角色。值与值之间和内存块与值之间的引用关系是通过指针来表征的

所以Ponter是安全的, 有应该关系时，内存块不会被GC回收, 但是一个uintptr值并不引用任何值，它被看作是一个整数，尽管常常它存储的是一个地址的数字表示, 所以uintptr存储的地址 有可能会被GC回收变成一个无效地址

总结:
+ Pointer 是安全的, 表征的是一种关系, 有引用就回不回收
+ uintptr 是一个整数, 表征的是一个值(内存地址的数字表示), 这个值表示的内存地址的值 有可能已经被GC回收

### uintptr地址被GC回收

在运行时刻，一次新的垃圾回收过程可能在一个不确定的时间启动，并且此过程可能需要一段不确定的时长才能完成。 所以一个不再被使用的内存块的回收时间点是不确定的

比如, 下面我们直接使用内存地址访问数组的其他元素:
```go
a := [3]int64{1, 2, 3}
fmt.Printf("%p\n", &a)

s1 := unsafe.Sizeof(a[0])
fmt.Printf("%d\n", s1)

// 我们把 Pointer -> uintptr (一波操作) -> Pointer, 这一系列动作是一次性完成的
p1 := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + s1))
fmt.Println(*p1)
```

如果我们把 p1 该写成2条语句:
```go
// 我们把 Pointer -> uintptr (一波操作)
p1Addr := uintptr(unsafe.Pointer(&a)) + s1 

// uintptr -> Pointer, 此时 p1Addr 的地址有可能被回收
p1 := (*int64)(unsafe.Pointer(p1Addr))
```

### Runtime的小动作

一个值的地址在程序运行中可能改变, 比如当一个协程的栈的大小改变时，开辟在此栈上的内存块需要移动，从而相应的值的地址将改变
而这个变化当中Pointer会跟随变化，但是uintptr是值 则不会

## 正确使用非类型安全指针

unsafe标准库包的文档中列出了[六种非类型安全指针的使用模式](https://golang.google.cn/pkg/unsafe/#Pointer)

### 获取地址

模式：Pointer --> uintptr 

此模式不是很有用。一般我们将最终的转换结果uintptr值输出到日志中用来调试，但是有很多其它安全并且简洁的途径也可以实现此目的

```go
type T struct{ a int }
var t1 T
fmt.Printf("%p\n", &t1)                          // 0xc0000a0200
println(&t1)                                     // 0xc0000a0200
fmt.Printf("%x\n", uintptr(unsafe.Pointer(&t1))) // c0000a0200
```

### 指针类型转换

将类型\*T1的一个值转换为非类型安全指针值，然后将此非类型安全指针值转换为类型\*T2

例如 math标准库包中的Float64bits函数。 此函数将一个float64值转换为一个uint64值。 在此转换过程中，此float64值在内存中的每个位（bit）都保持不变
```go
// 模式： *T1 --> Pointer --> *T2
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
```

再如, 将一个int8的整数转换成一个string, 同样内存中的值保持不变, 实现zero copy转换

```go
func bInt8(n int8) string {
	fmt.Println(*(*uint8)(unsafe.Pointer(&n))) // 1111 1111
	return strconv.FormatUint(uint64(*(*uint8)(unsafe.Pointer(&n))), 2)
}
```

### 直接操作内存地址

```go
// 模式: Pointer --> uintptr --> (一波计算) --> uintptr --> Pointer
p = unsafe.Pointer(uintptr(p) + offset)
```

将一个非类型安全指针转换为一个uintptr值，然后此uintptr值参与各种算术运算，再将算术运算的结果uintptr值转回非类型安全指针

例如 我们直接通过指针访问结构体的属性, 下面是直接通过地址访问y的第3个元素
```go
type T struct {
	x bool
	y [3]int16
}

const (
	N = unsafe.Offsetof(T{}.y)
	M = unsafe.Sizeof(T{}.y[0])
)

func TestUnsafePointer4() {
	t1 := T{y: [3]int16{123, 456, 789}}
	p := unsafe.Pointer(&t1)
	// "uintptr(p) + N + M + M"为t.y[2]的内存地址。
	ty2 := (*int16)(unsafe.Pointer(uintptr(p) + N + M + M))
	fmt.Println(*ty2) // 789
}
```

### 系统调用

```go
// 模式: Pointer --> uintptr --> syscall.Syscall
syscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(n))

// INVALID: uintptr cannot be stored in variable
// before implicit conversion back to Pointer during system call.
u := uintptr(unsafe.Pointer(p))
syscall.Syscall(SYS_READ, uintptr(fd), u, uintptr(n))
```


将非类型安全指针值转换为uintptr值并传递给syscall.Syscall函数调用

为啥uintptr传给Syscall的时候是安全的？

编译器针对每个syscall.Syscall函数调用中的每个被转换为uintptr类型的非类型安全指针实参添加了一些指令，从而保证此非类型安全指针所引用着的内存块在此调用返回之前不会被垃圾回收和移动


### 其他2种

由于涉及到反射这里暂不讲解

+ 将reflect.Value.Pointer或者reflect.Value.UnsafeAddr方法的uintptr返回值立即转换为非类型安全指针
+ 将一个reflect.SliceHeader或者reflect.StringHeader值的Data字段转换为非类型安全指针，以及其逆转换

## 声明与总结

go1 并不保证unsafe的兼容, 我们应该知晓当前的非类型安全机制规则和使用模式可能在以后的Go版本中完全失效, 几率很小， 因此，在实践中，请尽量保证能够将使用了非类型安全机制的代码轻松改为使用安全途径实现

从上面解释中，我们得知，对于某些情形，非类型安全机制可以帮助我们写出运行效率更高的代码。 但是，使用非类型安全指针也使得我们可能轻易地写出一些重现几率非常低的微妙的bug。 一个含有这样的bug的程序很可能在很长一段时间内都运行正常，但是突然变得不正常甚至崩溃。 这样的bug很难发现和调试





# 字符串格式化

Go语言用于控制文本输入和输出格式的库是fmt

## 格式化输入

fmt包中提供了3类读取输入的函数：

+ Scan家族：从标准输入os.Stdin中读取数据，包括Scan()、Scanf()、Scanln()
+ SScan家族：从字符串中读取数据，包括Sscan()、Sscanf()、Sscanln()
+ Fscan家族：从io.Reader中读取数据，包括Fscan()、Fscanf()、Fscanln()

其中:

+ Scanln、Sscanln、Fscanln在遇到换行符的时候停止
+ Scan、Sscan、Fscan将换行符当作空格处理
+ Scanf、Sscanf、Fscanf根据给定的format格式读取，就像Printf一样

这3家族的函数都返回读取的记录数量，并会设置报错信息，例如读取的记录数量不足、超出或者类型转换失败等

以下是他们的定义:

```sh
go doc fmt | grep -Ei "func [FS]*Scan"
func Fscan(r io.Reader, a ...interface{}) (n int, err error)
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)
func Fscanln(r io.Reader, a ...interface{}) (n int, err error)
func Scan(a ...interface{}) (n int, err error)
func Scanf(format string, a ...interface{}) (n int, err error)
func Scanln(a ...interface{}) (n int, err error)
func Sscan(str string, a ...interface{}) (n int, err error)
func Sscanf(str string, format string, a ...interface{}) (n int, err error)
func Sscanln(str string, a ...interface{}) (n int, err error)
```

因为还没介绍io.Reader，所以Fscan家族的函数暂且略过，但用法和另外两家族的scan类函数是一样的


## Scan家族

Scan家族函数从标准输入读取数据时，以空格做为分隔符分隔标准输入的内容，并将分隔后的各个记录保存到给定的变量中。

### Scan

输入数据时可以换行输入，Scan()会将换行符作为空格进行处理，直到读取到了2个记录之后自动终止读取操作

```go
func Scan() {
	fmt.Print("请输入你的姓名和年龄: ")
	fmt.Scan(&name, &age)
	fmt.Printf("姓名: %s 年龄: %d", name, age)
	fmt.Println()
}
```
注意 只能使用空格作为分隔符, 行号相当于空格

### Scanln

使用说明: 遇到换行符或EOF的时候终止读取

例如，使用Scanln函数等待用户输入数据，或从管道中读取数据。下面的代码将等待用户输入，且将读取的内容分别保存到name变量和age变量中：
```go
func Scanln() {
	var (
		name string
		age  uint
	)
	fmt.Print("请输入你的姓名和名字, 以空格分隔: ")
	fmt.Scanln(&name, &age)
	fmt.Printf("姓名: %s 年龄: %d", name, age)
	fmt.Println()
}
```
因为Scanln()遇到换行符或EOF的时候终止读取，所以在输入的时候只要按下回车键就会结束读取

### Scanf

Scanf 和 Scanln一样, 也是遇到换行符或EOF的时候终止读取, 只是Scanf只定了字符串格式，比如
```go
func Scanf() {
	var (
		name string
		age  uint
	)
	fmt.Print("请输入你的姓名和年龄, 以 ：分隔: ")
	fmt.Scanf("%s : %d", &name, &age)
	fmt.Printf("姓名: %s 年龄: %d", name, age)
	fmt.Println()
}
```

### 总结

由于Scanf可以自定义分隔符和输入格式, 灵活性上优于 其他2个，因此以Scanf使用较多

## SScan家族

Sscan家族的函数用于从给定字符串中读取数据，用法和Scan家族类似, 只是输入不再是标准输入了，方便内部使用

## Sscan

```go
input := "andy 30"
fmt.Sscan(input, &name, &age)
fmt.Println(name, age)
```

如果我们多给参数将自动忽略:
```go
input := "andy 30 40"
fmt.Sscan(input, &name, &age)
fmt.Println(name, age)
```

## Scanln

相比于Sscan而已, 我们以行作为结束符, 因此如果我们提前换行，那么这次scan结束
```go
input := "andy \n 30 40"
fmt.Sscanln(input, &name, &age)
fmt.Println(name, age)
```

## Sscanf

和Scanf一样，格式完全自己定义， 灵活性自己把握
```go
input := "andy : 50"
fmt.Sscanf(input, "%s : %d", &name, &age)
fmt.Println(name, age)
```

## 返回

我们讲Scan 和 Sscan时都没有处理函数返回, 返回值都是读取的参数的个数和错误
```
// returns the number of items successfully scanned. If that is less
// than the number of arguments, err will report why
```

由于没讲函数，这里可以简单看个例子:
```go
input := "andy : err"
ret, err := fmt.Sscanf(input, "%s : %d", &name, &age)
fmt.Println(ret, err)
```

## 扩展

我们刚才从 标准输入中读取了数据用于Scan, 现在我们可以从IO流中读取数据再配合Sscan 就完成了基于网络的Scan

这里采用bufio从标准输入读取数据, 真实案例你可以从Socket读取数据 发送过来
```go
fmt.Print("请输入你的姓名和年龄, 以空格分隔: ")
stdin := bufio.NewReader(os.Stdin)
line, _, err := stdin.ReadLine()
if err != nil {
    panic(err)
}
n, err := fmt.Sscanln(string(line), &name, &age)
if err != nil {
    panic(err)
}
fmt.Printf("read number of items: %d", n)
fmt.Println()
fmt.Printf("姓名: %s 年龄: %d", name, age)
fmt.Println()
```


# 循环语句

Go中只有一种循环结构：for

![xx](../../image/go-loops.svg)

## for循环

1.完整写法: for init; condition; post { }

+ init： 一般为赋值表达式，给控制变量赋初值；
+ condition： 关系表达式或逻辑表达式，循环控制条件；
+ post： 一般为赋值表达式，给控制变量增量或减量。

实例: 计算 1 到 10 的数字之和：

```go
var sum int
for i := 0; i <= 10; i++ {
   sum += i
}
fmt.Println(sum)
```

2.简短写法：for condition { }

init 和 post 参数是可选的，我们可以直接省略它，类似 While 语句

实例:

```go
sum := 1
for ; sum <= 10; {
         sum += sum
}
fmt.Println(sum)

// 这样写也可以，更像 While 语句形式
for sum <= 10{
         sum += sum
}
fmt.Println(sum)
```

## 无限循环

好几种方式实现for的无限循环。只要省略for的条件判断部分就可以实现无限循环

+ for i := 0;;i++
+ for ;; { }
+ for true { }
+ for { }

无限循环时，一般在循环体中加上退出语句，如break、os.Exit、return等

实例:

```go
var sum int
for {
   sum++
   fmt.Println(sum)
   if sum == 100 {
      return
   }
}
```

## for range遍历

range关键字非常好用，可以用来迭代那些可迭代的对象。比如slice、map、array，还可以迭代字符串，甚至是Unicode的字符串

```go
for index,value := range iterable {}
```

`注意`：value是从iterable中拷贝的副本, 我们直接修改value是无效的, 我们应该通过index来修改它，因此: 在循环体中应该总是让value作为一个只读变量

```go
iter := "abcdefg"
for index, value := range iter {
   fmt.Println(index, value)
   value = 'x'
}
fmt.Println(iter)
```

```go
iter := []int{1, 2, 3, 4, 5, 6}
for index, value := range iter {
   fmt.Println(index, value)
   iter[index] = 99
}
fmt.Println(iter)
```

## 嵌套循环

```go
for [condition |  ( init; condition; increment ) | Range]
{
   for [condition |  ( init; condition; increment ) | Range]
   {
      statement(s);
   }
   statement(s);
}
```

实例1: 打印九九乘法表

```go
// 1 x 1 = 1
// 1 x 2 = 2 2 x 2 = 4
// 1 x 3 = 3 2 x 3 = 6 3 x 3 = 9
// 1 x 4 = 4 2 x 4 = 8 3 x 4 = 12 4 x 4 = 16
// 1 x 5 = 5 2 x 5 = 10 3 x 5 = 15 4 x 5 = 20 5 x 5 = 25
// 1 x 6 = 6 2 x 6 = 12 3 x 6 = 18 4 x 6 = 24 5 x 6 = 30 6 x 6 = 36
// 1 x 7 = 7 2 x 7 = 14 3 x 7 = 21 4 x 7 = 28 5 x 7 = 35 6 x 7 = 42 7 x 7 = 49
// 1 x 8 = 8 2 x 8 = 16 3 x 8 = 24 4 x 8 = 32 5 x 8 = 40 6 x 8 = 48 7 x 8 = 56 8 x 8 = 64
// 1 x 9 = 9 2 x 9 = 18 3 x 9 = 27 4 x 9 = 36 5 x 9 = 45 6 x 9 = 54 7 x 9 = 63 8 x 9 = 72 9 x 9 = 81

for m := 1; m < 10; m++ {
   for n := 1; n <= m; n++ {
      fmt.Printf("%d x %d = %d ", n, m, m*n)
   }
   fmt.Println()
}
```

实例2: 寻找2~100间的素数

素数：除了能被一和本身整除不能被其它正数整除, 比如 1 ~ 10中的素数:

```go
2:  1x2
3:  1x3
4:  1x4  2x2
5:  1x5
6:  1x6  2x3
7:  1x7
8:  1x8  2x4
9:  1x9  2x? 3x3
10: 1x10 2x? 3x? 4x? 2x5 
```

所以判断一个数是不是素数只需要将比它小的数进行一个求余的计算
先把2-100都列出来,再把2的倍数划掉,接着是三,同理向下推,把素数的倍数都划掉,最后剩下的就都是素数了

```go
// 遍历 2 ~ 100
for m := 2; m <= 100; m++ {
   // 假定m是素数
   isP := true

   // 判断能不能分解
   for n := 2; n <= (m / n); n++ {
      if m%n == 0 {
         isP = false
      }
   }

   // 无法分解 为素数
   if isP {
      fmt.Println(m)
   }
}
```

## 循环中断

Go语言中用于控制循环中断语句有 break 和 continue

continue用于退出当前迭代，进入下一轮迭代。continue只能用于for循环中

### break

breake用于退出当前整个循环。如果是嵌套的循环，则退出它所在的那一层循环。
break除了可以用在for循环中，还可以用在switch结构或select结构

break 语法格式如下:

```go
for {
   break;
}

switch v {
   case xxx:
   break;
}
```

实例: 在变量 i == 10 的时候跳出循环

```go
i := 0
for {
   fmt.Println(i)
   if i == 10 {
      break
   }
   i++
}
```

以下实例有多重循环，演示了使用标记和不使用标记的区别：

```go
// 不使用标记
fmt.Println("---- break ----")
for i := 1; i <= 3; i++ {
   fmt.Printf("i: %d\n", i)
            for i2 := 11; i2 <= 13; i2++ {
                  fmt.Printf("i2: %d\n", i2)
                  break
            }
   }

// 使用标记
fmt.Println("---- break label ----")
re:
   for i := 1; i <= 3; i++ {
      fmt.Printf("i: %d\n", i)
      for i2 := 11; i2 <= 13; i2++ {
            fmt.Printf("i2: %d\n", i2)
            break re
      }
   }
```

### continue

continue用于退出当前迭代，进入下一轮迭代。continue只能用于for循环中

```go
for i := 0; i <= 10; i++ {
   if i%2 == 0 {
      continue
   }
   fmt.Println(i)
}
```

## 标签与跳转

Go 语言的 goto 语句可以无条件地转移到过程中指定的行。
goto 语句通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。
但是，在结构化程序设计中一般不主张使用 goto 语句， 以免造成程序流程的混乱，使理解和调试程序都产生困难

语法

goto 语法格式如下：

```go
goto label;
..
.
label: statement;
```

实例: 在变量 a 等于 15 的时候跳过本次循环并回到循环的开始语句 LOOP 处

```go
/* 定义局部变量 */
var a int = 10

/* 循环 */
LOOP:
for a < 20 {
   if a == 15 {
      /* 跳过迭代 */
      a++
      goto LOOP
   }
   fmt.Printf("a的值为 : %d\n", a)
   a++
}
```

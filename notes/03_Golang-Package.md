# Golang-Package  Golong的包

## 一、定义
- 包是函数和数据的集合，将有相关特性的函数和数据放在统一的文件目录进行管理，每个包都可以作为独立的单元维护并提供给其他项目进行使用
- 声明所在包，包名告知编译器哪些是包的源代码用于编译库文件，其次包名用于限制包内成员对外的可见性，最后包名用于在包外对公开成员的访问
- 在源文件中加上`package xxx`就可以声明xxx的包
 
## 二、成员可见性
- Go 语言使用名称首字母大小写来判断对象(常量、变量、函数、类型、结构体、方法等)的访问权限，首字母大写标识包外可见(公开的)，否者仅包内可访问(内部的);

## 三、main 包与 main 函数
- main包用于声明告知编译器 将包编译为二进制 可执行文件
- main包中的 main 函数是程序的入口，无返回值，无参数

## 四、init 函数
- init函数是初始化包使用，无返回值，无参数。建议每个包只定义一个； 
- init函数在import包时自动被调用(const -->var -->init)。

## 五、标准包
- Go提供了大量标准包，可查看 https://golang.google.cn/pkg/ &&  https://godoc.org 
	- `go list std` 查看所有标准包
	- `go doc packagename` 查看包的帮助信息
	- `go doc packagename.element` 查看包内成员 帮助信息

## 六、包的维护
- 包的提供者 -> 打tag  -> git tag
- 包的使用者 -> 改版本 -> go mod

## 七、关系说明
- import 导入的是路径，而非包名
- 包名和目录名不强制一致，但推荐一致
- 在代码中引用包的成员变量或者函数时，使用的包名不是目录名
- 在同一目录下，所有的源文件必须使用相同的包名
	- Multiple packages in directory: pk2, pk3 
- 文件名不限制，但不能有中文

## 八、设置 go mod 和 go proxy
- 设置两个环境变量
```
	$ go env -w GO111MODULE=on
	$ go env -w GOPROXY=https://goproxy.io,direct`
```

## 九、创建git，发布到github
- 项目目录下 go mod init github.com/ahwhy/myGolang
- git init 
- 添加 .gitignore 文件去掉一些和代码无关的文件/文件夹
- git add . && git commit -m "Record me learning golang"
- github上新建一个仓库
- 推送到远程
```
	git remote add origin https://github.com/ahwhy/myGolang.git
	git branch -M main
	git push -u origin main
```


# Golang-Path  Golang的Path包

## 一、Golang的标准库 path包

### 1. path
- path
	- path包实现了对斜杠分隔的路径的实用操作函数
```go
	// IsAbs 返回路径是否是一个绝对路径
	func IsAbs(path string) bool

	// Split 将路径从最后一个斜杠后面位置分隔为两个部分(dir和file)并返回
	// 如果路径中没有斜杠，函数返回值dir会设为空字符串，file会设为path；两个返回值满足path == dir+file
	// path.Split("static/myfile.css")
	// Output: static/ myfile.css
	func Split(path string) (dir, file string)

	// Join 可以将任意数量的路径元素放入一个单一路径里，会根据需要添加斜杠，结果是经过简化的，所有的空字符串元素会被忽略
	// path.Join("a", "b", "c")
	// Output: a/b/c
	func Join(elem ...string) string

	// Dir 返回路径除去最后一个路径元素的部分，即该路径最后一个元素所在的目录
	//在使用Split去掉最后一个元素后，会简化路径并去掉末尾的斜杠；如果路径是空字符串，会返回"."；如果路径由1到多个斜杠后跟0到多个非斜杠字符组成，会返回"/"；其他任何情况下都不会返回以斜杠结尾的路径
	// path.Dir("a", "b", "c")
	// Output: a/b
	func Dir(path string) string

	// Base 返回路径的最后一个元素。在提取元素前会求掉末尾的斜杠
	// 如果路径是""，会返回"."；如果路径是只有一个斜杆构成，会返回"/"
	// path.Base("/a/b")
	// Output: b
	func Base(path string) string

	// Ext 返回path文件扩展名
	// 返回值是路径最后一个斜杠分隔出的路径元素的最后一个'.'起始的后缀(包括'.')；如果该元素没有'.'会返回空字符串
	// path.Ext("/a/b/c/bar.css")
	// Output: .css
	func Ext(path string) string

	// Clean 通过单纯的词法操作返回和path代表同一地址的最短路径
	func Clean(path string) string

	// Match 要求匹配整个name字符串，而不是它的一部分
	// 如果name匹配shell文件名模式匹配字符串，Match函数返回真；只有pattern语法错误时，会返回ErrBadPattern
	func Match(pattern, name string) (matched bool, err error)
```

### 2. path/filepath
- path/filepath
	- filepath包实现了兼容各操作系统的文件路径的实用操作函数
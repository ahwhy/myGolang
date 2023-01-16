# Golang-Unsafe  Golang的非常规操作

## 一、Golang的标准库 unsafe包

### 1. unsafe
- unsafe
	- unsafe包提供了一些跳过go语言类型安全限制的操作
```golang
	// ArbitraryType 在本文档里表示任意一种类型，但并非一个实际存在与unsafe包的类型
	type ArbitraryType int

	// Pointer类型 用于表示任意类型的指针
	// 有4个特殊的只能用于Pointer类型的操作：
	// 1) 任意类型的指针可以转换为一个Pointer类型值
	// 2) 一个Pointer类型值可以转换为任意类型的指针
	// 3) 一个uintptr类型值可以转换为一个Pointer类型值
	// 4) 一个Pointer类型值可以转换为一个uintptr类型值
	// 因此，Pointer类型允许程序绕过类型系统读写任意内存，使用它时必须谨慎
	type Pointer *ArbitraryType

	// Sizeof 返回类型v本身数据所占用的字节数，返回值是“顶层”的数据占有的字节数
	// 例如，若v是一个切片，它会返回该切片描述符的大小，而非该切片底层引用的内存的大小
	func Sizeof(v ArbitraryType) uintptr

	// Alignof 返回类型v的对齐方式（即类型v在内存中占用的字节数）；若是结构体类型的字段的形式，它会返回字段f在该结构体中的对齐方式
	func Alignof(v ArbitraryType) uintptr

	// Offsetof 返回类型v所代表的结构体字段在结构体中的偏移量，它必须为结构体类型的字段的形式；即 它返回该结构起始处与该字段起始处之间的字节数
	func Offsetof(v ArbitraryType) uintptr
```
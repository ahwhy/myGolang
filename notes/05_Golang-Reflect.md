# Golang-Reflect  Golang的反射

## 一、Golang的反射定义

- 反射是指在运行时动态的访问和修改任意类型对象的结构和成员

- 为什么使用反射
	- 两个经典场景
		- 编写的一个函数，还不知道传给函数的类型具体是什么，可能是还没约定好，也可能是传入的类型很多
		- 希望通过用户的输入来决定调用按个函数(根据字符串调用方法)，动态执行函数

- python中的反射
	- 根据字符串执行函数
	- 根据字符串导入包

- go中的反射
	- go是静态语言。反射就是go提供一种机制，在编译时不知道类型的情况下，可以做如下的事情
		- 更新变量
		- 运行时查看值
		- 调用方法
		- 对他们的布局进行操作
	- 在go语言中提供reflect包提供反射的功能 ，每一个变量都有两个属性: 类型(Type)和值(Value)
		- reflect包提供 ValueOf 和 Typeof
		- `reflect.ValueOf` 获取输入接口中数据的值，如果为空返回 0
		- `reflect.Typeof` 获取输入接口中值的类型，如果为空返回 nil
		- Typeof 传入所有类型，因为所有的类型都实现了空接口
		
## 二、反射的Type

### 1. 定义
- `reflect.Type` 是一个接口类型，用于获取变量类型的信息，可通过 `reflect.TypeOf` 函数获取某个变量的类型信息
```go
	type Type interface { ... }

	// reflect.TypeOf
	func TypeOf(i interface{}) Type
```

### 2. ͨ通用方法
- `Name()` 类型名
- `PkgPath()` 包路径
- `Kind()` 类型枚举值
- `String()` Type字符串
- `Comparable()` 是否可进行比较
- `ImplementsType()` 是否实现某接口
- `AssignableTo(Type)` 是否可赋值给某类型
- `ConvertibleTo(Type)` 是否可转换为某类型
- `NumMethod()` 方法个数
- `Method(int)` 通过索引获取方法类型，Method结构体常用属性: `Name` 方法名; `Type` 函数类型; `Func` 方法值 Value
- `MethodByName(string)` 通过方法名字获取方法类型

### 3. 特定类型方法
- `reflect.Int*, reflect.UInt*, reflect.Float*k, reflect.Complex*`
	- `Bits()` 获取占用字节位数
- `reflect.Array`
	- `Len()` 获取数组长度
	- `Elem()` 获取数据元素类型
- `reflect.Slice`
	- `Elem()` 获取切片元素类型
- `reflect.Map`
	- `Key()` 获取映射键类型
	- `Elem()` 获取映射值类型
- `reflect.Ptr`
	- `Elem()` 获取指向值类型
- `reflect.Func`
	- `IsVariadic()` 是否具有可变参数
	- `NumIn()` 参数个数
	- `In(int)` ͨ通过索引获取参数类型
	- `NumOut`  返回值个数
	- `Out(int)` ͨ通过索引获取返回值类型
- `reflect.Struct`
	- `NumField` 属性个数
	- `Field(int)` ͨ通过索引获取属性
		- `StructField` 结构体常用属性
			- Name 属性名
			- Anonymous 是否为匿名
			- Tag 标签
				- StructTag 常用方法
					- `Get(string)`
					- `Lookup(string)`
	- `FieldByName(string)` ͨ通过属性名获取属性
	
## 三、反射的Value

### 1. 定义
- `reflect.Value` 是一个结构体类型，用于获取变量值的信息，可通过 `reflect.ValueOf` 函数获取某个变量的值信息
### 2. 通用方法
- `reflect.ValueOf`

### 3. ͨ通用方法
 - `Type()` 获取值类型
 - `CanAddr()` 是否可获取地址
 - `Addr()` 获取地址
 - `CanInterface()` 是否可以获取接口的
 - `InterfaceData()`
 - `Interface()` 将变量转换为 interface{}
 - `CanSet()` 是否可更新
 - `isValid()` 是否初始化为零值
 - `Kind()` 获取值 类型枚举值
 - `NumMethod()` 方法个数
 - `Method(int)` 通过索引获取方法值
 - `MethodByName(string)` 通过方法名字获取方法值
 - `ConvertType()` 转换为对应类型的值

### 4. 修改方法
- `Set/Set*` 设置变量值

### 5. 调用方法
- `Call()`
- `CallSlice()`

### 6. 特定类型方法
- `reflect.Int*`, `reflect.Uint*`
	- `Int()` 获取对应类型值
	- `Unit()` 获取对应类型值

- `reflect.Float*`
	- `Float()` 获取对应类型值

- `reflect.Complex*`
	- `Complex()` 获取对应类型值

- `reflact.Array`
	- `Len()` 获取数组长度
	- `Index(int)` 根据索引获取元素
	- `Slice(int, int)` 获取切片
	- `Slice3(int, int, int)` 获取切片
	
- `reflect.Slice`
	- `IsNil()` 判断是否为
	- `Len()` 获取元素数量
	- `Cap()` 获取容量
	- `Index(int)` 根据索引获取元素
	- `Slice(int, int)` 获取切片
	- `Slice3(int, int, int)` 获取切片

- `reflect.Map`
	- `IsNil()` 判断是否为
	- `Len()` 获取元素数量
	- `MapKeys()` 获取所有键
	- `MapIndex(Value)` 根据键获取值
	- `MapRange()` 获取键值组成的可迭代对象

 - `reflect.Ptr`
	- `Elem()` 获取指向值类型(解引用)

 - `reflect.Func`
	- `IsVariadic()` 是否具有可变参数
	- `NumIn()` 参数个数
	- `In(int)` 通过索引获取参数类型
	-` NumOut()` 返回值个数
	- `Out(int)` 通过索引获取返回值类型

 - `reflect.Struct`
	- `NumField()` 属性个数
	- `Field(int)` 通过索引获取属性
		- `StructField` 结构体常用属性
			- Name 属性名
			- Anonymous 是否为匿名
			- Tag标签
				- StructTag 常用方法
					- `Get(string)`
					- `Lookup(string)`
	- `FieldByName(string)` 通过属性名获取属性

## 应用

### 1. 内置类型的测试
```go
	var s interface{} = "abc"
	// TypeOf会返回模板的对象
	reflectType := reflect.TypeOf(s)
	reflectValue := reflect.ValueOf(s)
	log.Printf("[typeof:%v]", reflectType)
	log.Printf("[valueof:%v]", reflectValue)
```

### 2. 自定义struct的反射
- 生成的举例，对未知类型的进行 遍历 探测它的Field，抽象成一个函数

- Go语言里面struct成员变量小写，在反射的时候直接panic

- 结构体方法名小写不会panic，反射时也不会被查看到

- 指针方法不能被反射查看到
	- 对于成员变量
		- 先获取 intereface 的 `reflect.Type`，然后遍历 `NumField`
		- 再通过 `reflect.Type` 的 `Field` 获取字段
		- 最后通过 `Field` 的interface 获取对应的 value
	- 对于方法
		- 先获取 intereface 的 `reflect.Type`，然后遍历 `NumMethod`
		- 再分别通过 `reflect.Type` 的 `t.Method` 获取真实的方法
		- 最后通过 Name 和 Type 获取方法的类型和值

```go
	type Person struct {
		Name string
		Age  int
	}
	type Student struct {
		Person     // 匿名结构体嵌套
		StudentId  int
		SchoolName string
		IsBaoSong  bool // 是否保送
		Hobbies    []string
		// panic: reflect.Value.Interface: cannot return value obtained from unexported field or method
		// hobbies    []string
		Labels map[string]string
	}
	// func (s *Student) goHome() {
	// 		log.Printf("[回家][sid:%d]", s.StudentId)
	// }
	func (s *Student) GoHome() {
		log.Printf("[回家][sid:%d]", s.StudentId)
	}
	func (s Student) GotoSchool() {
		log.Printf("[去上学][sid:%d]", s.StudentId)
	}
	func (s Student) Baosong() {
		log.Printf("[竞赛保送][sid:%d]", s.StudentId)
	}
	func reflectProbeStruct(s interface{}) {
		// 获取目标对象
		t := reflect.TypeOf(s)
		log.Printf("[对象的类型名称: %s]", t.Name())
		// 获取目标对象的值类型
		v := reflect.ValueOf(s)
		// 遍历获取成员变量
		for i := 0; i < t.NumField(); i++ {
			// Field 代表对象的字段名
			key := t.Field(i)
			value := v.Field(i).Interface()
			//
			if key.Anonymous {
					log.Printf("[匿名字段][第:%d个字段][字段名:%s][字段的类型:%v][字段的值:%v]", i+1, key.Name, key.Type, value)
			} else {
					log.Printf("[命名字段][第:%d个字段][字段名:%s][字段的类型:%v][字段的值:%v]", i+1, key.Name, key.Type, value)
			}
		}
		// 打印方法
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			log.Printf("[第:%d个方法][方法名称:%s][方法的类型:%v]", i+1, m.Name, m.Type)
		}
	}
	s := Student{
		Person:     Person{Name: "xiaoyi", Age: 9900},
		StudentId:  123,
		SchoolName: "五道口皇家男子职业技术学院",
		IsBaoSong:  true,
		Hobbies:    []string{"唱", "跳", "Rap"},
		//hobbies:    []string{"唱", "跳", "Rap"},
		Labels: map[string]string{"k1": "v1", "k2": "v2"},
	}
	p := Person{
		Name: "李逵",
		Age:  124,
	}
	reflectProbeStruct(s)
	reflectProbeStruct(p)
	// 执行结果
	// 2021/07/16 17:09:30 [对象的类型名称: Student]
	// 2021/07/16 17:09:30 [匿名字段][第:1个字段][字段名:Person][字段的类型:main.Person][字段的值:{xiaoyi 9900}]
	// 2021/07/16 17:09:30 [命名字段][第:2个字段][字段名:StudentId][字段的类型:int][字段的值:123]
	// 2021/07/16 17:09:30 [命名字段][第:3个字段][字段名:SchoolName][字段的类型:string][字段的值:五道口皇家男子职业技术学院]
	// 2021/07/16 17:09:30 [命名字段][第:4个字段][字段名:IsBaoSong][字段的类型:bool][字段的值:true]
	// 2021/07/16 17:09:30 [命名字段][第:5个字段][字段名:Hobbies][字段的类型:[]string][字段的值:[唱 跳 Rap]]
	// 2021/07/16 17:09:30 [命名字段][第:6个字段][字段名:Labels][字段的类型:map[string]string][字段的值:map[k1:v1 k2:v2]]
	// 2021/07/16 17:09:30 [第:1个方法][方法名称:Baosong][方法的类型:func(main.Student)]
	// 2021/07/16 17:09:30 [第:2个方法][方法名称:GotoSchool][方法的类型:func(main.Student)]
	// 2021/07/16 17:09:30 [对象的类型名称: Person]
	// 2021/07/16 17:09:30 [命名字段][第:1个字段][字段名:Name][字段的类型:string][字段的值:李逵]
	// 2021/07/16 17:09:30 [命名字段][第:2个字段][字段名:Age][字段的类型:int][字段的值:124]
```

### 3. 反射修改值
- 必须是指针类型

- 调用方法 `pointer.Elem().Setxxx()`

```go
    var num float64 = 3.14
    log.Printf("[num原始值:%f]", num)
    // 通过reflect.ValueOf获取num中的value
    // 必须是指针才可以修改值
    pointer := reflect.ValueOf(&num)
    newValue := pointer.Elem()
    // 赋值
    newValue.SetFloat(5.6)
    log.Printf("[num新值:%f]", num)

    pointer = reflect.ValueOf(num)
    // reflect: call of reflect.Value.Elem on float64 Value
    newValue = pointer.Elem()
```
 
### 4. 反射调用方法
- 首先 `reflect.ValueOf(p1)` 获取，得到反射类型对象

- `reflect.ValueOf.MethodByName` 需要传入准确的方法名称，`MethodByName`代表注册
	- 名称错误 `panic: reflect: call of reflect.Value.Call on zero Value`

- `[]reflect.Value` 这是最终需要调用方法的参数，无参数传空切片
```go
	type Person struct {
		Name   string
		Age    int
		Gender string
	}
	func (p Person) ReflectCallFuncWithArgs(name string, age int) {
		log.Printf("[调用的是带参数的方法][args.name:%s][args.age:%d][[p.name:%s][p.age:%d]",
			name,
			age,
			p.Name,
			p.Age,
		)
	}
	func (p Person) ReflectCallFuncWithNoArgs() {
		log.Printf("[调用的是不带参数的方法]")
	}
	p1 := Person{
		Name:   "小乙",
		Age:    18,
		Gender: "男",
	}

	// 首先通过 reflect.ValueOf(p1)获取 得到反射值类型
	getValue := reflect.ValueOf(p1)

	// 带参数的方法调用
	methodValue := getValue.MethodByName("ReflectCallFuncWithArgs")

	// 参数是reflect.Value的切片
	args := []reflect.Value{reflect.ValueOf("李逵"), reflect.ValueOf(30)}
	methodValue.Call(args)

	// 不带参数的方法调用
	methodValue = getValue.MethodByName("ReflectCallFuncWithNoArgs")

	// 参数是reflect.Value的切片
	args = make([]reflect.Value, 0)
	methodValue.Call(args)
```


### 5. 结构体标签和反射
- json的标签解析 `json`
	- JSON 是对 JavaScript 中各种类型的值，如：字符串、数字、布尔值和对象等，进行Unicode本文编码

- yaml的标签解析 `yaml`

- xorm gorm的标签 标识db字段

- 自定义标签
	- 原理是 `t.Field.Tag.Lookup` "标签名"
```go
	type Person struct {
		Name string `json:"name" yaml:"yaml_name" mage:"name"`
		Age  int    `json:"age"  yaml:"yaml_age"  mage:"age"`
		City string `json:"-" yaml:"yaml_city" mage:"-"`
	}
	//json解析
	func jsonWork() {
		// 对象marshal成字符串
		p := Person{
			Name: "xiaoyi",
			Age:  18,
			City: "北京",
		}
		data, err := json.Marshal(p)
		if err != nil {
			log.Printf("[json.marshal.err][err:%v]", err)
			return
		}
		log.Printf("[person.marshal.res][res:%v]", string(data))
	
		// 从字符串解析成结构体
		p2Str := `
	{
		"name":"李逵",
		"age":28,
		"city":"山东"
	}`
		var p2 Person
		err = json.Unmarshal([]byte(p2Str), &p2)
		if err != nil {
			log.Printf("[json.unmarshal.err][err:%v]", err)
			return
		}
		log.Printf("[person.unmarshal.res][res:%v]", p2)
	}
	// yaml读取文件
	func yamlWork() {
		filename := "a.yaml"
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("[ioutil.ReadFile.err][err:%v]", err)
			return
		}
		p := &Person{}
		//err = yaml.Unmarshal(content, p)                         // func yaml.Unmarshal(in []byte, out interface{}) (err error)
		err = yaml.UnmarshalStrict(content, p)                     // func yaml.UnmarshalStrict(in []byte, out interface{}) (err error)
		if err != nil {
			log.Printf("[yaml.UnmarshalStrict.err][err:%v]", err)
			return
		}
		log.Printf("[yaml.UnmarshalStrict.res][res:%v]", p)
	}
	// 自定义标签
	func jiexizidingyibiaoqian(s interface{}) {
		// typeOf type类型
		r := reflect.TypeOf(s)
		value := reflect.ValueOf(s)
		for i := 0; i < r.NumField(); i++ {
			field := r.Field(i)
			key := field.Name
			if tag, ok := field.Tag.Lookup("mage"); ok {
				if tag == "-" {
					continue
				}
				log.Printf("[找到了mage标签][key:%v][value:%v][标签: mage=%s]",
					key,
					value.Field(i),
					tag,
				)
			}
		}
	}
	jsonWork()
	yamlWork()
	p := Person{
		Name: "xiaoyi",
		Age:  18,
		City: "北京",
	}
	jiexizidingyibiaoqian(p)
```

## 弊端

### 1. 代码可读性变差

### 2. 隐藏的错误躲过编译检查
- Go语言为静态语言，编译器能发现类型的错误

- 但是对于反射代码无能为力，可能运行很久才会 panic

- 反射调用方法的副作用，将 float64 参数传成 int `panic: reflect: Call using float64 as type int`

### 3. Go语言反射性能问题
- 反射比正常的代码要慢1-2个数据量级，如果是追求性能的关键模块应减少反射
```go
	// 它是一个具体的值，不是一个可复用的对象
	// 每次取出的 fieldValue类型是 reflect.value
	// 每次反射都要 malloc这个 reflect.Value结构体，还有GC
	type := reflect.value(obj)
	fieldValue := type_.FieldByName("xx")
```

# Javascript基础

需要为HTML页面上添加一些动态效果, Brendan Eich这哥们在两周之内设计出了JavaScript语言

几个公司联合ECMA（European Computer Manufacturers Association）组织定制了JavaScript语言的标准，被称为ECMAScript标准

## JavaScript 运行时

+ 浏览器
+ NodeJS

```sh
$ node -v
v14.17.1
```

## 数据类型

+ Number
+ 字符串
+ 布尔值
+ 数组
+ 对象

### 数组

JavaScript的Array可以包含任意数据类型，并通过索引来访问每个元素

```js
var arr1 = new Array(1, 2, 3); // 创建了数组[1, 2, 3]
var arr2 = [1, 2, 3.14, 'Hello', null, true];
```

越界不报错
``` js
arr1[0]  // 1
arr1[3]  // undefined

// 如果通过索引赋值时，索引超过了范围，统一可以赋值
arr1[3] = 3
arr1[3] // 3
```

#### push和pop

+ push()向Array的末尾添加若干元素
+ pop()则把Array的最后一个元素删除掉

```js
var arr = [1, 2];
arr.push('A', 'B'); // 返回Array新的长度: 4
arr; // [1, 2, 'A', 'B']
arr.pop(); // pop()返回'B'
arr; // [1, 2, 'A']
arr.pop(); arr.pop(); arr.pop(); // 连续pop 3次
arr; // []
arr.pop(); // 空数组继续pop不会报错，而是返回undefined
arr; // []
```

#### unshift和shift

+ unshift()往Array的头部添加若干元素
+ shift()方法则把Array的第一个元素删掉

```js
var arr = [1, 2];
arr.unshift('A', 'B'); // 返回Array新的长度: 4
arr; // ['A', 'B', 1, 2]
arr.shift(); // 'A'
arr; // ['B', 1, 2]
arr.shift(); arr.shift(); arr.shift(); // 连续shift 3次
arr; // []
arr.shift(); // 空数组继续shift不会报错，而是返回undefined
arr; // []
```

#### splice

splice()方法是修改Array的“万能方法”，它可以从指定的索引开始删除若干元素，然后再从该位置添加若干元素

```js
var arr = ['Microsoft', 'Apple', 'Yahoo', 'AOL', 'Excite', 'Oracle'];
// 从索引2开始删除3个元素,然后再添加两个元素:
arr.splice(2, 3, 'Google', 'Facebook'); // 返回删除的元素 ['Yahoo', 'AOL', 'Excite']
arr; // ['Microsoft', 'Apple', 'Google', 'Facebook', 'Oracle']
// 只删除,不添加:
arr.splice(2, 2); // ['Google', 'Facebook']
arr; // ['Microsoft', 'Apple', 'Oracle']
// 只添加,不删除:
arr.splice(2, 0, 'Google', 'Facebook'); // 返回[],因为没有删除任何元素
arr; // ['Microsoft', 'Apple', 'Google', 'Facebook', 'Oracle']
```

#### sort和reverse

+ sort()可以对当前Array进行排序
+ reverse()把整个Array的元素给调个个，也就是反转

```js
var arr = ['B', 'C', 'A'];
arr.sort();
arr; // ['A', 'B', 'C']
arr.reverse();
arr; // ['C', 'B', 'A']
```

#### concat和slice

+ concat()方法把当前的Array和另一个Array连接起来，并返回一个新的Array
+ slice()就是对应String的substring()版本，它截取Array的部分元素，然后返回一个新的Array

```js
var arr = ['A', 'B', 'C'];
var added = arr.concat([1, 2, 3]);
added; // ['A', 'B', 'C', 1, 2, 3]

var arr = ['A', 'B', 'C', 'D', 'E', 'F', 'G'];
arr.slice(0, 3); // 从索引0开始，到索引3结束，但不包括索引3: ['A', 'B', 'C']
arr.slice(3); // 从索引3开始到结束: ['D', 'E', 'F', 'G']

// 如果不给slice()传递任何参数，它就会从头到尾截取所有元素。利用这一点，我们可以很容易地复制一个Array
var aCopy = arr.slice();
aCopy; // ['A', 'B', 'C', 'D', 'E', 'F', 'G']
aCopy === arr; // false
```


#### vue数组
Vue 将被侦听的数组的变更方法进行了包裹，所以它们也将会触发视图更新。这些被包裹过的方法包括：

+ push()
+ pop()
+ shift()
+ unshift()
+ splice()
+ sort()
+ reverse()

### 对象

JavaScript的对象是一种无序的集合数据类型，它由若干键值对组成

```js
obj1 = new Object()
obj2 = {}
```

由于JavaScript的对象是动态类型，你可以自由地给一个对象添加或删除属性

未定义的属性不报错
```js
obj1.a = 1  
obj1.a // 1
obj1.b  // undefined

obj1.b = 2
obj1.b // 2

// 删除b属性
delete obj1.b
delete obj1.b // 删除一个不存在的school属性也不会报错
```

使用hasOwnProperty, 判断对象是否有该属性

```js
obj1.hasOwnProperty('b') // false
```

### null和undefined

null表示一个空的值，而undefined表示值未定义。事实证明，这并没有什么卵用

```js
var a = {a: 1}
a.b  // undefined
a.b = null
a.b  // null
```

### 逻辑运算符

+ &&: 与运算
+ ||: 或运算
+ !: 非运算

### 关系运算符

大于和小于没啥特别的, 要特别注意相等运算符==。JavaScript在设计时，有两种比较运算符：

+ ==: 它会自动转换数据类型再比较，很多时候，会得到非常诡异的结果；
+ ===: 它不会自动转换数据类型，如果数据类型不一致，返回false，如果一致，再比较。

```js
false == 0; // true
false === 0; // false
```

## 变量


### var申明

var：变量提升（无论声明在何处，都会被提至其所在作用于的顶部）

```js
var age = 20
function f1() {console.log(age)}
f1() // 20
```

### 局部变量声明

let：无变量提升（未到let声明时，是无法访问该变量的）

```js
{ let a1 = 20 }
a1 // a1 is not defined
```

### 申明常量

const：无变量提升，声明一个基本类型的时候为常量，不可修改；声明对象可以修改

```js
const c1 = 20
c1 = 30 // Assignment to constant variable
```

### 变量提升

JavaScript的函数定义有个特点，它会先扫描整个函数体的语句，把所有申明的变量“提升”到函数顶部

```js
function foo() {
    var x = 'Hello, ' + y;
    console.log(x);
    var y = 'Bob';
}

// JavaScript引擎看到的代码相当于
function foo() {
    var y; // 提升变量y的申明，此时y为undefined
    var x = 'Hello, ' + y;
    console.log(x);
    y = 'Bob';
}
```

所以js里面 变量都定义在顶部, 并且大量使用let来声明变量

### 解构赋值

js里面还有这种你看不懂的骚操作
```js
// 数组属性是index, 解开可以直接和变量对应
let [x, [y, z]] = ['hello', ['JavaScript', 'ES6']];

var person = {
    name: '小明',
    age: 20,
    gender: 'male',
    passport: 'G-12345678',
    school: 'No.4 middle school'
};

// 对象解开后是属性, 可以直接导出你需要的属性, 用的地方很多
var {name, age, passport} = person;
```

这就叫解构赋值


## 字符串

JavaScript的字符串就是用''或""括起来的字符表示

```js
str1 = 'str'
str2 = "str"
```

### 字符串转义

使用转义符: \

```js
'I\'m \"OK\"!';
```


### 多行字符串

```js
ml = `这是一个
多行
字符串`; 
// "这是一个\n多行\n字符串"
```

### 字符串模版

格式: 使用``表示的字符串 可以使用${var_name} 来实现变量替换

```js
var name = '小明'
var age = 20
console.log(`你好, ${name}, 你今年${age}岁了！`)// 你好, 小明, 你今年20岁了！
```

### 字符串拼接

直接使用+号

### 常用操作

+ toUpperCase: 把一个字符串全部变为大写
+ toLowerCase: 把一个字符串全部变为小写


## 错误处理

一种是程序写的逻辑不对，导致代码执行异常

```js
var s = null
s.length 
// VM1760:1 Uncaught TypeError: Cannot read property 'length' of null
//     at <anonymous>:1:3
```

如果在一个函数内部发生了错误，它自身没有捕获，错误就会被抛到外层调用函数，如果外层函数也没有捕获，该错误会一直沿着函数调用链向上抛出，直到被JavaScript引擎捕获，代码终止执行


我们可以判断s的合法性, 在保证安全的情况下，使用
```js
if (s !== null) {s.length}
```

也可以捕获异常, 阻断其往上传传递

### try catch

```js
try { s.length } catch (e) {console.log('has error, '+ e)}
// VM2371:1 has error, TypeError: Cannot read property 'length' of null
```

完整的try ... catch ... finally:
```js
try {
    ...
} catch (e) {
    ...
} finally {
    ...
}
```

+ try: 捕获代码块中的异常
+ catch: 出现异常时需要执行的语句块
+ finally: 无论成功还是失败 都需要执行的代码块


常见实用案例:  loading

### 错误类型

javaScript有一个标准的Error对象表示错误

```js
err = new Error('异常来')
err 
// <!-- Error: 异常来
//     at <anonymous>:1:7 -->

err instanceof Error
// true
```

### 抛出错误

程序也可以主动抛出一个错误，让执行流程直接跳转到catch块。抛出错误使用throw语句

```js
throw new Error('抛出异常')
// VM3447:1 Uncaught Error: 抛出异常
//     at <anonymous>:1:7
// (anonymous) @ VM3447:1
```

## 函数

```js
function abs(x) {
    if (x >= 0) {
        return x;
    } else {
        return -x;
    }
}
```

上述abs()函数的定义如下：

+ function指出这是一个函数定义；
+ abs是函数的名称；
+ (x)括号内列出函数的参数，多个参数以,分隔；
+ { ... }之间的代码是函数体，可以包含若干语句，甚至可以没有任何语句。

### 方法

比如我们定义了一个对象
```js
var person = {name: '小明', age: 23}
```

那么我们如何给这个对象添加方法喃?
```js
var person = {name: '小明', age: 23}
person.greet = function() {
    console.log(`hello, my name is ${this.name}`)
}

person.greet()
```

绑定到对象上的函数称为方法，和普通函数也没啥区别，但是它在内部使用了一个this关键字

在一个方法内部，this是一个特殊变量，它始终指向当前对象，也就是xiaoming这个变量

#### this 是个坑

注意这里的this, 如果你没有绑带在对象上, this 指的是 浏览器的window对象

```js
fn = function() {
    console.log(this)
}  
fn()
// <ref *1> Object [global] {
//   global: [Circular *1],
//   clearInterval: [Function: clearInterval],
//   clearTimeout: [Function: clearTimeout],
//   setInterval: [Function: setInterval],
//   setTimeout: [Function: setTimeout] {
//     [Symbol(nodejs.util.promisify.custom)]: [Getter]
//   },
//   queueMicrotask: [Function: queueMicrotask],
//   clearImmediate: [Function: clearImmediate],
//   setImmediate: [Function: setImmediate] {
//     [Symbol(nodejs.util.promisify.custom)]: [Getter]
//   },
//   fn: [Function: fn]
// }
```

```js
var person = {name: '小明', age: 23}
person.greetfn = function() {
    return function() {
        console.log(`hello, my name is ${this.name}`)
    }
}

person.greetfn()() // hello, my name is undefined
```


此时我们可以通过一个变量+ 闭包, 把当前this传递过去, 确保this正常传递
```js
var person = {name: '小明', age: 23}
person.greetfn = function() {
    var that = this  // 这个很多，别看不懂
    return function() {
        console.log(`hello, my name is ${that.name}`)
    }
}

person.greetfn()() // hello, my name is 小明
```

### 箭头函数(匿名函数)

在js中你会看到很多这样的语法:
```js
fn = x => x * x
console.log(fn(10)) 
```

这就是js特色的箭头函数
```js
x => x * x

// 等价于下面这个函数
function (x) {
    return x * x;
}
```

一个完整的箭头函数语法:
```js
(params ...) => { ... }
```

我们看看下面列子
```js
axios
.get('http://localhost:8050/hosts', {params: this.query})
.then(response => {
    console.log(response)
    this.tableData = response.data.data.items
    this.total = response.data.data.total
    console.log(this.tableData)
})
.catch(function (error) { // 请求失败处理
    console.log(error);
});
```

#### 请使用箭头函数

箭头函数看上去是匿名函数的一种简写，但实际上，箭头函数和匿名函数有个明显的区别：箭头函数内部的this是词法作用域，由上下文确定

```js
var person = {name: '小明', age: 23}
person.greetfn = function() {
    return () => {
        // this 继承自上层的this
        console.log(`hello, my name is ${this.name}`)
    }
}

person.greetfn()()
```

所有在js中 到处都是箭头函数


## 名字空间

不在任何函数内定义的变量就具有全局作用域。实际上，JavaScript默认有一个全局对象window，全局作用域的变量实际上被绑定到window的一个属性

### 全局作用域与window

```js
alert("hello")
// 等价于
window.alert("hello")
```

由于函数定义有两种方式，以变量方式var foo = function () {}定义的函数实际上也是一个全局变量，因此，顶层函数的定义也被视为一个全局变量，并绑定到window对象

```js
var a = 10

a // 10
window.a // 10
```

甚至我们可以覆盖掉浏览器的内置方法:

```js
alert = () => {console.log("覆盖alert方法")}
() => {console.log("覆盖alert方法")}
alert() // 覆盖alert方法
```

是不是很骚, 这要是做大项目 就是在玩火, 那如果避免这种问题喃? 使用命名空间

### 名字空间于Export

我们可以将我们的所有方法绑定到一个变量上，然后暴露出去，避免全局变量的混乱， 许多著名的JavaScript库都是这么干的：jQuery，YUI，underscore等等

```js
// 唯一的全局变量MYAPP:
var MYAPP = {};

// 其他变量:
MYAPP.name = 'myapp';
MYAPP.version = 1.0;

// 其他函数:
MYAPP.foo = function () {
    return 'foo';
};


export MYAPP
```

其他文件中
```js
import { MYAPP } from './export';
```

## 条件判断

语法格式:
```js
if (condition) {
    ...
} else if (condition) {
    ...
} else {
    ...
}
```

注意条件需要加上括号, 其他和Go语言的if一样:

```js
var age = 20;
if (age >= 6) {
    console.log('teenager');
} else if (age >= 18) {
    console.log('adult');
} else {
    console.log('kid');
}
```

## for 循环

语法格式:

```js
for (初始条件; 判断条件; 修改变量) {
    ...
}
```

注意条件需要加上括号:

```js
var x = 0;
var i;
for (i=1; i<=10000; i++) {
    x = x + i;
}
x; // 50005000
```

### for in (不推荐使用)

for循环的一个变体是for ... in循环，它可以把一个对象的所有属性依次循环出来


遍历对象: 遍历出来的属性是元素的key

```js
var o = {
    name: 'Jack',
    age: 20,
    city: 'Beijing'
};
for (var key in o) {
    console.log(key); // 'name', 'age', 'city'
}
```

遍历数组: 一个Array数组实际上也是一个对象，它的每个元素的索引被视为一个属性

```js
var a = ['A', 'B', 'C'];
for (var i in a) {
    console.log(i); // '0', '1', '2'
    console.log(a[i]); // 'A', 'B', 'C'
}
```

for in 有啥问题? 为啥不推荐使用, 我们看下面一个例子

当我们手动给Array对象添加了额外的属性后，for ... in循环将带来意想不到的意外效果
```js
var a = ['A', 'B', 'C'];
a.name = 'Hello';
for (var x in a) {
    console.log(x); // '0', '1', '2', 'name'
}
```

为什么? 这和for in的遍历机制相关: 遍历对象的属性名称

那如何解决这个问题喃? 答案是 for of


### for of

for ... of循环则完全修复了这些问题，它只循环集合本身的元素

```js
var a = ['A', 'B', 'C'];
a.name = 'Hello';
for (var x of a) {
    console.log(x); // 'A', 'B', 'C'
}
```


但是我们用for of 能遍历对象吗?

```js
var o = {
    name: 'Jack',
    age: 20,
    city: 'Beijing'
};
for (var key of o) {
    console.log(key);
}
// VM2749:6 Uncaught TypeError: o is not iterable
//     at <anonymous>:6:17
```

变通的方法是: 我们可以通过Object提供的方法获取key数组,然后遍历

```js
var o = {
    name: 'Jack',
    age: 20,
    city: 'Beijing'
};
for (var key of Object.keys(o)) {
    console.log(key); // 'name', 'age', 'city'
}
```


### forEach方法

forEach()方法是ES5.1标准引入的, 他也是遍历元素的一种常用手段, 也是能作用于可跌倒对象上, 和for of一样

```js
arr.forEach(function(item) {console.log(item )})
```

当然这还有一种简洁写法

```js
arr.forEach((item) => {console.log(item)})
```

### for循环应用

如果后端返回的数据不满足我们展示的需求, 需要修改，比如vendor想要友好显示，我们可以直接修改数据


## Promise对象

在JavaScript的世界中，所有代码都是单线程执行的。

由于这个“缺陷”，导致JavaScript的所有网络操作，浏览器事件，都必须是异步执行。Javascript通过回调函数实现异步, js的一大特色

### 单线程异步模型

```js
function callback() {
    console.log('Done');
}
console.log('before setTimeout()');
setTimeout(callback, 1000); // 1秒钟后调用callback函数
console.log('after setTimeout()');
// before setTimeout()
// after setTimeout()
// 等待一秒后
// Done
```

由此可见并不会真正阻塞1秒, 而是在1秒后调用该函数, 这就是javascript的编程范式: 基于回调的异步

### Promise与异步

我们来看一个函数, resolve是成功后的回调函数, reject是失败后的回调函数
```js
function testResultCallbackFunc(resolve, reject) {
    var timeOut = Math.random() * 2;
    console.log('set timeout to: ' + timeOut + ' seconds.');
    setTimeout(function () {
        if (timeOut < 1) {
            console.log('call resolve()...');
            resolve('200 OK');
        }
        else {
            console.log('call reject()...');
            reject('timeout in ' + timeOut + ' seconds.');
        }
    }, timeOut * 1000);
}
```

然后我们把处理成功和失败的函数 作为回调传递给该函数:
```js
function testResultCallback() {
    success = (message) => {console.log(`success ${message}`)}
    failed = (error) => {console.log(`failed ${error}`)}
    testResultCallbackFunc(success, failed)
}

// set timeout to: 0.059809346310547795 seconds.
// call resolve()...
// success 200 OK
```

可以看出，testResultCallbackFunc()函数只关心自身的逻辑，并不关心具体的resolve和reject将如何处理结果

Js把这种编程方式抽象成了一种对象: Promise

```js
interface PromiseConstructor {
    /**
     * Creates a new Promise.
     * @param executor A callback used to initialize the promise. This callback is passed two arguments:
     * a resolve callback used to resolve the promise with a value or the result of another promise,
     * and a reject callback used to reject the promise with a provided reason or error.
     */
    new <T>(executor: (resolve: (value: T | PromiseLike<T>) => void, reject: (reason?: any) => void) => void): Promise<T>;
```

下面我们将回调改为Promise对象:
```js
var p1 = new Promise(testResultCallbackFunc)
p1.then((resp) => {
    console.log(resp)
}).catch((err) => {
    console.log(err)
})
// set timeout to: 0.628561731809246 seconds.
// call resolve()...
// 200 OK
```

可见Promise最大的好处是在异步执行的流程中，把执行代码和处理结果的代码清晰地分离了:

![](./pic/promise.png)


### Async函数 + Promise组合

从回调函数，到Promise对象，再到Generator函数(不讲, 这是协程方案的一种过度形态)，JavaScript异步编程解决方案历程可谓辛酸，终于到了Async/await。很多人认为它是异步操作的最终解决方案

这些需要提到 async函数, async函数由内置执行器进行执行, 这和go func() 有异曲同工之妙

那我们如果声明一个异步函数, 其实很简单 在你函数前面加上一个 async关键字就可以了

```js
async function testWithAsync() {
    var p1 = new Promise(testResultCallbackFunc)

    try {
        var resp = await p1
        console.log(resp)
    } catch (err) {
        console.log(err)
    }
}
```

这里testWithAsync就是一个异步函数, 他执行的时候 是交给js的携程执行器处理的, 而  await关键字 就是 告诉执行器 当p1执行完成后 主动通知我下(协程的一种实现), 其实就是一个 event pool模型(简称epool模型)

我们修改下之前的demo, 使用async 来实现

## 总结

+ 推荐使用箭头函数
+ 判断使用 ===
+ 由于变量提升问题，尽量使用let声明变量，并且写在开头
+ for循环推荐forEach
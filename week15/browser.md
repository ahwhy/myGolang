# 浏览器

先抛出2个场景的问题:
+ 如何适配屏幕大小?
+ 如何兼容多种浏览器?


JavaScript可以获取浏览器提供的很多对象，并进行操作.

## 浏览器对象

我们的网页是通过浏览器加载出来的, 而浏览器本身也有很多功能:
+ 浏览器本身版本信息
+ 当前窗口大小
+ 历史记录
+ ...

而这些数据都是绑定在window对象上得, 所以你可以认为对于浏览器加载的网页, window就表示的浏览器本身

### 浏览器窗口大小

window的如下4个属性控制这个浏览器的窗口大小， IE9+、Safari、Opera和Chrome都支持这4个属性
+ innerWidth: 页面视图的宽
+ innerHeight: 页面视图的高
+ outerWidth: 浏览器整个窗口的宽
+ outerHeight: 浏览器整个窗口的高

```js
// 可以调整浏览器窗口大小试试:
console.log('window inner size: ' + window.innerWidth + ' x ' + window.innerHeight);
console.log('window outer size: ' + window.outerWidth + ' x ' + window.outerHeight);
```

### 浏览器信息

我们在页面的时候, 如何知道用户使用的那种浏览器, 好做不通的适配, 毕竟各种厂商的浏览器实现有各有差异

关于浏览器的相关信息都保存在navigator对象上面得:
```js
window.navigator
// appCodeName: "Mozilla"
// appName: "Netscape"
// appVersion: "5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36"
// ...
```

最常用的属性包括：
+ navigator.appName：浏览器名称；
+ navigator.appVersion：浏览器版本；
+ navigator.language：浏览器设置的语言；
+ navigator.platform：操作系统类型；
+ navigator.userAgent：浏览器设定的User-Agent字符串
+ navigator.userAgentData useragent相关信息, Object

如何判断用户使用的是pc还是手机?

### 屏幕信息

有时候 我们可能需要知道用户的屏幕尺寸, 方便我们做网页布局, 这是好就需要获取屏幕的数据:
```js
window.screen
// Screen {availWidth: 2560, availHeight: 1440, width: 2560, height: 1440, colorDepth: 24, …}
// availHeight: 1440
// availLeft: 1440
// availTop: 0
// availWidth: 2560
// colorDepth: 24
// height: 1440
// orientation: ScreenOrientation {angle: 0, type: 'landscape-primary', onchange: null}
// pixelDepth: 24
// width: 2560
```

### 访问网址

这个很有用, 为了保证每个用户访问该URL都呈现统一的效果, 我们的页面往往都需要读取当前Page的参数

```js
window.location
// Location
// hash: ""
// host: "www.baidu.com"
// hostname: "www.baidu.com"
// href: "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&rsv_idx=..."
// origin: "https://www.baidu.com"
// pathname: "/s"
// port: ""
// protocol: "https:"
// reload: ƒ reload()
// replace: ƒ replace()
// search: "?ie=utf-8&f=8&rsv_bp=1&rsv_idx=1&tn=..."
```

下面是一些常用属性:
+ location.protocol; // 'http'
+ location.host; // 'www.example.com'
+ location.port; // '8080'
+ location.pathname; // '/path/index.html'
+ location.search; // '?a=1&b=2'
+ location.hash; // 'TOP'

要加载一个新页面，可以调用location.assign()。如果要重新加载当前页面，调用location.reload()方法非常方便

```js
if (confirm('重新加载当前页' + location.href + '?')) {
    location.reload();
} else {
    location.assign('/'); // 设置一个新的URL地址
}
```

### 历史记录

history对象保存了浏览器的历史记录，JavaScript可以调用history对象的back()或forward ()，相当于用户点击了浏览器的“后退”或“前进”按钮

```js
// <- 浏览器的“后退”按钮
history.back();

// -> 浏览器的“前进”按钮
history.forward()
```

## AJAX

AJAX不是JavaScript的规范，它只是一个哥们“发明”的缩写：Asynchronous JavaScript and XML，意思就是用JavaScript执行异步网络请求

简单来说 就是浏览器中的http client

在现代浏览器上写AJAX主要依靠XMLHttpRequest对象

### XMLHttpRequest

```js
// 新建XMLHttpRequest对象
var request = new XMLHttpRequest(); 

// 发送请求:
request.open('GET', 'https://www.baidu.com');
request.send();
```

如何就这样我们是获取不到请求结果, 由于js都是异步的, 我们需要定义回调来处理返回

```js
request.onreadystatechange = function () {
    console.log(request.status)
    console.log(request.responseText)
}
```

对于AJAX请求特别需要注意跨域问题: CORS

#### 简单请求

满足下面3个条件的是简单请求: 
+ HTTP Method: GET、HEAD和POST
+ Content-Type: application/x-www-form-urlencoded、multipart/form-data和text/plain
+ 不能出现任何自定义头

通常能满足90%的需求

控制其跨域的关键头 来自于服务端设置的: Access-Control-Allow-Origin

对于简单请求，浏览器直接发出CORS请求。具体来说，就是在头信息之中，增加一个Origin字段

+ 客户端发送请求时, Origin: xxxx
+ 服务端响应请求时, 设置: Access-Control-Allow-Origin: xxxx / *

如果服务端不允许该域，就跨越失败

如果想允许特点的Header, 服务端也可以通过添加: Access-Control-Expose-Headers 来进行控制

#### 复杂请求

复杂请求:
+ HTTP Method: PUT、DELETE
+ Content-Type: 其他类型如application/json
+ 含自定义头, 比如我们后面携带的用于认证的X-OAUTH-TOKEN头

非简单请求的CORS请求，会在正式通信之前，增加一次HTTP查询请求，称为"预检"请求（preflight）

preflight请求:
+ Method: OPTIONS
+ Header: Access-Control-Request-Method, 列出浏览器的CORS请求会用到哪些HTTP方法
+ Header: Access-Control-Request-Headers 该字段是一个逗号分隔的字符串，指定浏览器CORS请求会额外发送的头信息字段

```
OPTIONS /cors HTTP/1.1
Origin: http://api.bob.com
Access-Control-Request-Method: PUT
Access-Control-Request-Headers: X-Custom-Header
Host: api.alice.com
Accept-Language: en-US
Connection: keep-alive
User-Agent: Mozilla/5.0...
```

preflight响应:

+ Header: Access-Control-Allow-Origin, 允许的源（域)
+ Header: Access-Control-Allow-Methods, 服务端 允许的方法(所有方法, 一般大于请求的该头)
+ Header: Access-Control-Allow-Headers, 服务队 允许的自定义Header
+ Header: Access-Control-Max-Age, 用来指定本次预检请求的有效期, 避免多次请求, 该字段可选

```
Access-Control-Allow-Methods: GET, POST, PUT
Access-Control-Allow-Headers: X-Custom-Header
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 1728000
```

### axios

很显然原生的client库，没那么好用, 于是诞生的axios


### 页面(DOM)

浏览器当前加载的Page也是一个对象: window.document

由于HTML在浏览器中以DOM形式表示为树形结构，document对象就是整个DOM树的根节点

比如我们获取当前Document中的一些属性: 比如 title
document的title属性是从HTML文档中的<title>xxx</title>读取的，但是可以动态改变

```js
document.title
// 'golang random int_百度搜索'
// 我们可以动态修改
document.title = '测试修改'
// 我们看到浏览器窗口标题发生了变化
```

很多js库比如JQuery都是通过动态操作Dom来实现很多高级功能的, 这些是上层库的基石

#### DOM查询

要查找DOM树的某个节点，需要从document对象开始查找。最常用的查找是根据ID和Tag Name以及ClassName

```html
<h1>列表</h1>
<ul id="list_menu" class="ul_class">
    <li>Coffee</li>
    <li>Tea</li>
    <li>Milk</li>
</ul>
```

```js
document.getElementById('list_menu')
// <ul data-v-7ba5bd90 id=​"list_menu">​<li data-v-7ba5bd90>​…​</li>​<li data-v-7ba5bd90>​…​</li>​<li data-v-7ba5bd90>​…​</li>​</ul>​

document.getElementsByTagName('ul')
// 这里返回所有的ul元素

document.getElementsByClassName('ul_class')
```

当然我们也可以组合使用
```js
document.getElementById('list_menu').getElementsByTagName('li')

// 上面也等价于
document.getElementById('list_menu').children

// firstElementChild
document.getElementById('list_menu').firstElementChild

// lastElementChild
document.getElementById('list_menu').lastElementChild

// parentElement 获取父元素
document.getElementById('list_menu').parentElement
```

#### 更新DOM

获取的元素后, 我们可以通过元素的如下2个方法，修改元素:
+ innerHTML: 不但可以修改一个DOM节点的文本内容，还可以直接通过HTML片段修改DOM节点内部的子树
+ innerText: 只修改文本内容

```js
var le = document.getElementById('list_menu').lastElementChild
le.innerText = '牛奶'
// 我们发现页面的内容已经修改

// html格式无法识别
le.innerText = '<span style="color:red">牛奶</span>'
// <span style="color:red">牛奶</span>

// 使用innterHTML则可以
le.innerHTML = '<span style="color:red">牛奶</span>'
```

#### 插入DOM

很多响应式框架都会根据数据新增，动态创建一个DOM元素，并插入到指定位置,

我们使用 createElement 来创建一个DOM元素, 比如创建一个a标签

```js
// 创建一个A标签
var newlink = document.createElement('a')
// <a></a>

// 修改A标签属性
newlink.href = "www.baiducom"
newlink.innerText = '跳转到百度'

// 追加到某个元素后面
var lm = document.getElementById('list_menu')
lm.appendChild(newlink)
```

如果我们想要控制元素插入的位置可以使用insertBefore

insertBefore的语法如下:
```js
parentElement.insertBefore(newElement, referenceElement);
```

```js
// 父元素
var lm = document.getElementById('list_menu')

// 子元素
var cf = document.getElementById('coffee')

// 需要插入的元素
var newlink = document.createElement('a')
newlink.href = "www.baiducom"
newlink.innerText = '跳转到百度'

// 插入到coffee之前
lm.insertBefore(newlink, cf)
```

总结: 有2种方式可以插入一个DOM元素
+ appendChild: 把一个子节点添加到父节点的最后一个子节点
+ insertBefore: 插入到某个元素之前


#### 删除DOM

删除一个DOM节点就比插入要容易得多

要删除一个节点，首先要获得该节点本身以及它的父节点，然后，调用父节点的removeChild把自己删掉

removeChild的语法如下:
```js
parent.removeChild(childNode);
```

比如我们删除刚才添加的那个元素
```js
lm.removeChild(newlink)
```

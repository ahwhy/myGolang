# Web-Base Web的基础

## 一、Protobuf介绍

### 1. Protobuf
一、web基础
 - Html 展现web内容
 - Css  内容展现的样式
 - JavaScript 给页面添加一些动作，或者对页面的一些操作


二、Html基础
1、Html 网页结构
 - 基本格式
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8" />
			<title>页面标题</title>           
		</head>
		<body>
			<h1>标题1</h1>
			<p>段落1</p>
			<script type="text/javascript">
				alert("Holle")
				confirm("World")
				prompt("请输入")
			</script>
		</body>
	</html>

2、标签与元素
 - Html标签  https://www.runoob.com/tags/ref-byfunc.html
	html
	head
	title
	body
	h1
	p
 - Html元素
	- Html元素指一个具体的标签实例，下面有2个Html元素，都是h标签
		<h1>这是一个标题</h1>
		<h1>这是另一个标题</h1>
	- 整个网页就是由这些标签组成的Html元素嵌套组成

3、元素语法
 - <tag att1=v1 attr2=v2>内容</tag>
 - 每种标签都有自己的一组属性，属性分为2类
	- 全局属性: 所有标签都有的属性  https://www.runoob.com/tags/ref-standardattributes.html
		id 定义元素的唯一id
		class 为html元素定义一个或多个类名（classname）(类名从样式文件引入)
		style 规定元素的行内样式（inline style）
		title 描述了元素的额外信息 (作为工具条使用)
	- 标签属性: 每种标签肯能还有一些该标签才特有的一些属性
		href 需要有引用的属性的标签才有这个属性, 比如 链接(a标签) 和 图片(img标签)

4、常用标签
 - 基础标签
	<h1> to <h6>  定义 HTML 标题
	<p>	          定义一个段落
	<br>	      定义简单的折行
	<hr>	      定义水平线
	<!--...-->	  定义一个注释
 - 文本标签
	<del> 定义被删除文本
	<i>   定义斜体文本
	<ins> 定义被插入文本
	<sub> 下标文字
	<sup> 上标文字
	<u>   下划线文本
 - 表单标签
	<form>
	<input>
	...
 - 常见元素
	<iframe> 嵌套外部网页
	<img>    展示图像
	<area>   标签定义图像映射内部的区域: https://www.runoob.com/try/try.php?filename=tryhtml_areamap
	<a>      链接标签
	<ul>     定义一个无序列表
	<ol>     定义一个有序列表
	<li>   定义一个列表项
 - 表格
	<table>  标签定义 HTML 表格, 一个 HTML 表格包括 <table> 元素，一个或多个 <tr>、<th> 以及 <td> 元素。
	<tr>     元素定义表格行，
	<th>     元素定义表头，
	<td>     元素定义表格单元
 - 容器元素
	<div>    标签定义 HTML 文档中的一个分隔区块或者一个区域部分，标签常用于组合块级元素，以便通过 CSS 来对这些元素进行格式化
	<span>   用于对文档中的行内元素进行组合，标签提供了一种将文本的一部分或者文档的一部分独立出来的方式

5、元素id
 - 标识该元素的唯一身份，并且可以在其他地方引用
	- 通过a标题跳转到指定的位置:
		<p>
		<a href="#C4">查看章节 4</a>
		</p>
		<h2>章节 1</h2>
		<p>这边显示该章节的内容……</p>
		<h2>章节 2</h2>
		<p>这边显示该章节的内容……</p>
		<h2>章节 3</h2>
		<p>这边显示该章节的内容……</p>
		<h2><a id="C4">章节 4</a></h2>
		<p>这边显示该章节的内容……</p>
 - id 也是js操作元素的重要依据之一
	document.getElementById('C4')
	<a id=​"C4">​章节 4​</a>

6、元素的样式
 - 通过元素的style属性可以控制该元素的样式
	- 把p元素里面的这段话的字体加大, 演示改为红色
		<p style="color:red;font-size:20px;">这边显示该章节的内容……</p>
		- 语法的格式
			key: value;
			- 分号分开的就是一个样式条目，可以为一个元素添加很多样式，详情见Css
	- 控制元素的宽和高
		<iframe src="https://www.runoob.com" sytle="height: 200px;width: 400px;">
			<p>您的浏览器不支持  iframe 标签。</p>
		</iframe>

7、脚本
 - 一个静态页面具有元素和样式，页面加载完成后，如果要动态修改里面的元素，就需要用到Javascript脚本
	<script> 标签用于定义客户端脚本，比如 JavaScript
	<script> 元素既可包含脚本语句，也可通过 src 属性指向外部脚本文件
 - 示例
	- 引入脚本
		// 通过src 网络引入
		<script src="https://cdn.jsdelivr.net/npm/vue@2.6.14"></script>
		// 通过本地文件引入
		<script>
			import axios from 'axios';
		</script>
	- JavaScript 最常用于图片操作、表单验证以及内容动态更新
		<script>
		c4 = document.getElementById('C4')
		// <a id=​"C4">​章节 4​</a>​
		c4.innerText
		// '章节 4'
		c4.innerText = '章节 5'
		// '章节 5'
		</script>


三、Css基础
1、Css的作用
 - 可以通过直接指定style属性来配置 DOM元素的样式
	- <p style="color:red;font-size:20px;">这边显示该章节的内容……</p>
	- 如果一个页面上 有多个元素都是用了这个样式
	- 这是就需要将样式独立出来，单独写到一个保存在外部的 .css 文件中，这就形成了css

2、Css的语法
 - 对所有的p标签 使用下面定义的样式
	/* 这是个注释 */
	p {
	color: red;
	}
 - 选择器，对应p，用于选择样式生效的范围，也就是对那些元素生效
 - 样式属性(style attribute)，也就是具体需要添加的样式，属性有key和value组成，多个属性用';'分开
 - 注释是用来解释你的代码，并且可以随意编辑它，浏览器会忽略它
	- CSS注释以 /* 开始, 以 */ 结束

3、CSS选择器  https://www.runoob.com/cssref/css-selectors.html
 - 基础选择器
	- 类型
		- 标签，标签名称  比如  h1
		- 类，类名称      比如  .className
		- ID，元素id      比如  #id
	- 标签选择器
		- 以标签名开头，选择所有div元素，比如前面的p标签选择器就是这种
			- 所有h1标签 都显示为红色
				<style scoped>
				h1 {
				color: red;
				}
				</style>
	- 类选择器
		- 给标签取class名，以点(.)加class名开头，选择所有该class名的元素
		- 最常用的一种选择器
			<h1 class="f12">基础标签</h1>
			<style scoped>
			.f12 {
			font-size: 12px;
			}
			</style>
	- id选择器
		- 给标签取id名，以 #加id名 开头，具有唯一性
			<h2><a id="C4">章节 4</a></h2>
			<style scoped>
			#C4 {
			color: blue;
			}
			</style>
	- 群选择器
		- 也可以组合使用，一次选择多个元素，比如 h1, .f12, #C4 都选中，以逗号','分隔多个选中的元素
			<style scoped>
			h1,.f12,#C4 {
			color: blue;
			}
			</style>
	- 全局选择器
		- 既然有多选，那就有全选，使用'*'
			<style scoped>
			* {
			color: blue;
			}
			</style>
 - 层级选择器
	- 类型
		- 如果直接选择元素的话，可能会很多，并不能精确选中需要的元素
		- 因此就可以引入层级关系来选择，比如 div元素下的p标签，而不是顶层p标签
		- 总共有4种层级关系选择器
			- 子选择器: 父元素 > 子元素
			- 包含选择器: 父元素 包含的元素
			- 兄弟选择器: 当前元素 ~ 兄弟元素
			- 相邻选择器: 当前元素 + 相邻元素
	- 子选择器
		- 用于已知父元素，选择子元素、
		- 语法: 以 > 隔开父子级元素，(模块名>模块名，修饰>前模块内的子模块)
			<ul id="list_menu" class="ul_class">
				<li id="coffee">Coffee</li>
				<li>Tea</li>
				<li>Milk</li>
			</ul>
			<style scoped>
			ul>li {
			font-weight: 600;
			}
			</style>
		- 子选择器，只能选择其当前层的子元素，如果不是就无法选中
			<ul id="list_menu" class="ul_class">
				<li id="coffee">Coffee</li>
				<li>Tea</li>
				<li>Milk</li>
				<div>
				<li>In Div</li>
				</div>
			</ul>
	- 包含选择器
		- 包含选择器 可以选择父元素下的 包含该标签的所有元素
		- 语法: 父元素 子元素
			<ul id="list_menu" class="ul_class">
				<li id="coffee">Coffee</li>
				<li>Tea</li>
				<li>Milk</li>
				<div>
				<li>In Div</li>
				</div>
			</ul>
			<style scoped>
			ul li {
				font-weight: 600;
			}
			</style>
	- 兄弟选择器
		- 兄弟选择器 可以选择从一层级的兄弟元素，比如 #coffee同层的li元素
		- 语法: 同层元素 ~ 需要选中的同层元素
			<ul id="list_menu" class="ul_class">
				<li id="coffee">Coffee</li>
				<li>Tea</li>
				<li>Milk</li>
				<div>
				<li>In Div</li>
				</div>
			</ul>
			<style scoped>
			#coffee~li {
				font-weight: 600;
			}
			</style>
	- 相邻选择器
		- 相邻弟选择器 可以选择相邻的元素，而不是同层所有的元素
		- 语法: 当前元素 + 相邻元素
			<ul id="list_menu" class="ul_class">
				<li id="coffee">Coffee</li>
				<li>Tea</li>
				<li>Milk</li>
				<div>
				<li>In Div</li>
				</div>
			</ul>
			<style scoped>
			#coffee+li {
				font-weight: 600;
			}
			</style>
 - 其他常用选择器
	- 属性选择器
		- 元素会有多属性，可以根据元素的属性来选择元素，比如下面的input框
			<form action="demo_form.php">
				<input type="submit" value="提交">
			</form>
			<style scoped>
			[type=submit] {
				font-size: 22px;
			}
			</style>
	- 伪类选择器
		- 当选中一堆元素后，如果对其中部分元素 进行选择 就可以使用伪类选择器
		- 语法: 选择的元素:函数方法
			- first-child:   第一个
			- last-child:    最后一个
			- nth-child(n):  第n个元素
			- not():         not函数，不包含
			- https://www.runoob.com/css/css-pseudo-classes.html
				<ul id="list_menu" class="ul_class">
					<li id="coffee">Coffee</li>
					<li>Tea</li>
					<li>Milk</li>
				</ul>
				<style scoped>
				ul>li:first-child {
					font-weight: 600;
				}
				</style>

4、CSS引入方式
 - 内联
	- 在标签内直接写的，style="attr:value;..."
		<el-table
			:data="tableData"
			style="width: 100%">
			...
		</el-table>
 - 内嵌
	- 通过style标签定义的样式
		<style scoped>
		ul>li:first-child {
			font-weight: 600;
		}
		</style>
 - 外联
	- 当样式需要被应用到很多页面的时候，外部样式表将是理想的选择
	- 使用外部样式表，就可以通过更改一个文件来改变整个站点的外观
		<head>
		<link rel="stylesheet" type="text/css" href="mystyle.css">
		</head>

5、样式基础
 - CSS单位
	- 在布局或者设置字体大小的时候经常看到: 22px; ，px其实是css里面长度单位
	- 绝对长度
		px * 像素 (1px = 1/96th of 1in)
		cm 厘米
		mm 毫米
		in 英寸 (1in = 96px = 2.54cm)
	- 相对长度
		em 它是描述相对于应用在当前元素的字体尺寸，所以它也是相对长度单位。一般浏览器字体大小默认为16px，则2em == 32px
		ex 依赖于英文字母小 x 的高度
		ch 数字 0 的宽度
		rem rem 是根 em（root em）的缩写，rem作用于非根元素时，相对于根元素字体大小；rem作用于根元素字体大小时，相对于其出初始字体大小
		vw viewpoint width，视窗宽度，1vw=视窗宽度的1%
		vh viewpoint height，视窗高度，1vh=视窗高度的1%
		vmin vw和vh中较小的那个。
		vmax vw和vh中较大的那个。
 - CSS颜色
	- https://www.w3school.com.cn/cssref/css_colors.asp

6、常用属性
 - 元素尺寸控制
	- height 设置元素的高度
	- width 设置元素的宽度
	- line-height 设置行高
	- min-width 设置元素的最小宽度
	- min-height 设置元素的最小高度
	- max-height 设置元素的最大高度
	- max-width 设置元素的最大宽度
		<div style="height: 220px;width:440px">
			<p>我们的内容</p>
		</div>
 - 盒子模型
	- Margin(外边距)  清除边框外的区域，外边距是透明的
	- Border(边框)    围绕在内边距和内容外的边框
	- Padding(内边距) 清除内容周围的区域，内边距是透明的
	- Content(内容)   盒子的内容，显示文本和图像
		<div style="height: 220px;width:440px">
			<p style="margin-top:22px;">我们的内容</p>
		</div>
 - Display
	- inline: 现在在一行
	- block: 块元素是一个元素，占用了全部宽度，在前后都是换行符
	- flex: flex布局，见后面参考
		<div>
			<li style="display: inline;">Tea</li>
			<li style="display: inline;">Milk</li>
		</div>
		<div>
			<span style="display: block;">span1</span>
			<span style="display: block;">span2</span>
		</div>
 - Overflow
	- overflow 属性可以控制内容溢出元素框时在对应的元素区间内添加滚动条
		- visible 默认值，内容不会被修剪，会呈现在元素框之外
		- hidden  内容会被修剪，并且其余内容是不可见的
		- scroll  内容会被修剪，但是浏览器会显示滚动条以便查看其余的内容
		- auto    如果内容被修剪，则浏览器会显示滚动条以便查看其余的内容
		- inherit 规定应该从父元素继承 overflow 属性的值
			<div id="overflowTest">
				<p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
				<p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
				<p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
				<p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
				<p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
				<p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
			</div>
			<style scoped>
			#overflowTest {
				background: #4CAF50;
				color: white;
				padding: 15px;
				width: 80%;
				height: 100px;
				overflow: scroll;
				border: 1px solid rgb(150, 18, 18);
			}
			</style>
 - 浮动
	- 控制元素左移还是右移
		- left
		- right
		- 元素浮动之后，周围的元素会重新排列，为了避免这种情况，使用 clear 属性
			<div id="overflowTest" style="clear:both">
				<span style="display: block;">span1</span>
				<span style="display: block;">span2</span>
			</div>
 - 对齐
	- 以常见的居中为例
		<div id="overflowTest" style="text-align: center">
			<span style="display: block;">span1</span>
			<span style="display: block;">垂直居中</span>
			<span style="display: block;">span2</span>
		</div>
		- 为了文本在元素内居中对齐，可以使用 text-align: center
		- 或者行高到div 就整体居中

7、常用网址
	- CSS布局教学           https://zh.learnlayout.com/
	- Flex布局教学游戏      http://flexboxfroggy.com/
	- 在线CSS代码可视工具   https://enjoycss.com/
	- 新拟态模拟工具        https://neumorphism.io/#e0e0e0
	- 渐变色方案            https://uigradients.com/#Copper


四、JavaScript基础
1、JavaScript 运行时
	- 浏览器
		NodeJS
	- $ node -v
		v14.17.1

2、数据类型
 - JavaScript中的数据类型
	- Number
	- 字符串
	- 布尔值
	- 数组 []
		- JavaScript的Array可以包含任意数据类型，并通过索引来访问每个元素
			var arr  = new Array(1, 2, 3);                          // 创建数组arr1[1, 2, 3]
			var arr2 = [1, 2, 3.14, 'Hello', null, true];
			console.log(typeof(arr), arr)                           // object (3) [1, 2, 3]
			- 越界不报错
				arr[0]  // 1
				arr[3]  // undefined
				// 如果通过索引赋值时，索引超过了范围，统一可以赋值
				arr[5] = 5
				arr     // [1, 2, 3, 4, empty, 5]
		- push和pop
			- push() 向Array的末尾添加若干元素，并返回数组长度
				arr.push('A', 'B')
			- pop()  把Array的最后一个元素删除掉，并返回删除元素的值
				arr.pop()
			- 空数组继续pop不会报错，而是返回undefined
		- unshift和shift
			- unshift() 往Array的头部添加若干元素，并返回数组长度
				arr.unshift('A', 'B')
			- shift()   把Array的第一个元素删掉
				arr.shift()
			- 空数组继续shift不会报错，而是返回undefined
		- splice
			- splice() 可以从指定的索引开始删除若干元素，然后再从该位置添加若干元素
				var arr = ['Microsoft', 'Apple', 'Yahoo', 'AOL', 'Excite', 'Oracle'];
				// 从索引2开始删除3个元素,然后再添加两个元素:
				arr.splice(2, 3, 'Google', 'Facebook');    // 返回删除的元素 ['Yahoo', 'AOL', 'Excite']
				arr;                                       // ['Microsoft', 'Apple', 'Google', 'Facebook', 'Oracle']
				// 只删除,不添加:                          
				arr.splice(2, 2);                          // ['Google', 'Facebook']
				arr;                                       // ['Microsoft', 'Apple', 'Oracle']
				// 只添加,不删除:                          
				arr.splice(2, 0, 'Google', 'Facebook');    // 返回[],因为没有删除任何元素
				arr;                                       // ['Microsoft', 'Apple', 'Google', 'Facebook', 'Oracle']
		- sort和reverse
			- concat() 把当前的Array和另一个Array连接起来，并返回一个新的Array
				var arr = ['A', 'B', 'C'];
				var added = arr.concat([1, 2, 3]);
				added;                             // ['A', 'B', 'C', 1, 2, 3]
			- slice()  截取Array的部分元素，然后返回一个新的Array (对应String的substring()版本)
				var arr = ['A', 'B', 'C', 'D', 'E', 'F', 'G'];
				arr.slice(0, 3);                   // 从索引0开始，到索引3结束，但不包括索引3: ['A', 'B', 'C']
				arr.slice(3);                      // 从索引3开始到结束: ['D', 'E', 'F', 'G']
				// 如果不给slice()传递任何参数，它就会从头到尾截取所有元素。利用这一点，我们可以很容易地复制一个Array
				var aCopy = arr.slice();
				aCopy;                             // ['A', 'B', 'C', 'D', 'E', 'F', 'G']
				aCopy === arr;                     // false
		- Vue中的数组
			- Vue 将被侦听的数组的变更方法进行了包裹，所以它们也将会触发视图更新
				- 这些被包裹过的方法包括: 
					push()
					pop()
					shift()
					unshift()
					splice()
					sort()
					reverse()	
	- 对象 字典 Object {}
		- JavaScript的对象是一种无序的集合数据类型，它由若干键值对组成
			obj  = new Object()
			obj2 = {}
		- 由于JavaScript的对象是动态类型，可以自由地给一个对象添加或删除属性
			- 未定义的属性不报错
				obj.a = 1  
				obj.a          // 1
				obj.b          // undefined
				obj.b = 2     
				obj.b          // 2
				// 删除b属性
				delete obj1.b
				delete obj1.b  // 删除一个不存在的c属性也不会报错
		- hasOwnProperty
			- hasOwnProperty() 判断对象是否有该属性
				obj.hasOwnProperty('b')  // false
 - null和undefined
	- null表示一个空的值，而undefined表示值未定义
		var a = {a: 1}
		a.b             // undefined
		a.b = null
		a.b             // null
 - 逻辑运算符
	&& 与运算
	|| 或运算
	!  非运算
 - 关系运算符
	- JavaScript在设计时，有两种比较运算符
		==    它会自动转换数据类型再比较，很多时候，会得到非常诡异的结果
		===   它不会自动转换数据类型，如果数据类型不一致，返回false，如果一致，再比较
			false == 0; // true
			false === 0; // false

3、变量
 - var申明
	- var  变量提升，即 无论声明在何处，都会被提至其所在作用域的顶部
 - 局部变量声明
	- let  无变量提升，即 未到let声明的语句时，是无法访问该变量的
		{ let a1 = 20 }
		a1                // a1 is not defined
 - 申明常量
	- const: 无变量提升，声明一个基本类型的时候为常量，不可修改
		const c1 = 20
		c1 = 30 // Assignment to constant variable
	- 声明对象可以修改
		const c1 = 30
 - 变量提升
	- JavaScript的函数定义有个特点，它会先扫描整个函数体的语句，把所有申明的变量"提升"到函数顶部
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
		- js里面 变量都定义在顶部，并且大量使用let来声明变量
 - 解构赋值
	- 数组属性是index，解开可以直接和变量对应
		let [x, [y, z]] = ['hello', ['JavaScript', 'ES6']];
	- 对象解开后是属性，可以直接导出需要的属性
		var person = {
			name: '小明',
			age: 20,
			gender: 'male',
			passport: 'G-12345678',
			school: 'No.4 middle school'
		};
		var {name, age, passport} = person;
		name                                // '小明'
		age                                 // 20
		passport                            // 'G-12345678'
		
4、字符串
 - 声明
	- JavaScript的字符串就是用''或""括起来的字符表示
		str  = 'str'
		str2 = "str"
 - 字符串转义
	- 使用转义符  \
		'I\'m \"OK\"!';     // `I'm "OK"!`
 - 多行字符串
		ml = `这是一个
		多行
		字符串`;            // '这是一个\n多行\n字符串'
 - 字符串模版
	- 格式: 使用``表示的字符串 可以使用${var_name} 来实现变量替换
		var name = '小明'
		var age = 20
		console.log(`你好, ${name}, 你今年${age}岁了！`)   // 你好, 小明, 你今年20岁了！
 - 字符串拼接
	- 直接使用+号
		'aa' + 'bb'         // 'aabb'
		'aa' + 123          // 'aa123'  注意字符串加数字 = 字符串
 - 常用操作
	str.toUpperCase() 把一个字符串全部变为大写
	str.toLowerCase() 把一个字符串全部变为小写
	str.toString()    把类型转换成字符串
	parseInt(str)     把字符串转换成Int
	parseFloat(str)   把字符串转换成Float

5、错误处理
 - 定义
	- 程序处理的逻辑有问题，导致代码执行异常
		- 示例
			var s = null
			s.length 
			// VM1760:1 Uncaught TypeError: Cannot read property 'length' of null
			//     at <anonymous>:1:3
		- 如果在一个函数内部发生了错误，它自身没有捕获，错误就会被抛到外层调用函数
		- 如果外层函数也没有捕获，该错误会一直沿着函数调用链向上抛出，直到被JavaScript引擎捕获，代码终止执行
		- 可以判断s的合法性，在保证安全的情况下，捕获异常，阻断其往上传传递
			if (s !== null) {s.length}
	- try catch
		- 示例
			try { s.length } catch (e) {console.log('has error, '+ e)}
			// VM2371:1 has error, TypeError: Cannot read property 'length' of null
		- 完整的{try ... catch ... finally}语句
			try {
				...
			} catch (e) {
				...
			} finally {
				...
			}
			try: 捕获代码块中的异常
			catch: 出现异常时需要执行的语句块
			finally: 无论成功还是失败 都需要执行的代码块
	- 常见实用案例 loading
 - 错误类型
	- javaScript有一个标准的Error对象表示错误
		err = new Error('异常情况')
		err 
		// <!-- Error: 异常情况
		//     at <anonymous>:1:7 -->
		err instanceof Error
		// true
 - 抛出错误
	- 程序也可以主动抛出一个错误，让执行流程直接跳转到catch块
	- 抛出错误使用throw语句
		throw new Error('抛出异常')
		// VM3447:1 Uncaught Error: 抛出异常
		//     at <anonymous>:1:7
		// (anonymous) @ VM3447:1

6、命名空间
 - JavaScript中的命名空间
	- JavaScript默认有一个全局对象window，全局作用域的变量实际上被绑定到window的一个属性
	- 不在任何函数内定义的变量就具有全局作用域
 - 全局作用域与window
	- 示例一
		alert("hello")
		// 等价于
		window.alert("hello")
	- 示例二
		var a = 10
		a // 10
		window.a // 10
	- 示例三
		- 由于函数定义有两种方式
			- 以变量方式var foo = function () {}定义的函数也是一个全局变量
			- 顶层函数的定义也被视为一个全局变量，并绑定到window对象
	- 示例四
		- 甚至可以覆盖掉浏览器的内置方法
			alert = () => {console.log("覆盖alert方法")}
			alert()                                         // 覆盖alert方法
 - export
	- 在js中，一个模块就是一个独立的文件
		- 该文件内部的所有变量，外部无法获取
		- 如果希望外部能够读取模块内部的某个变量，就必须使用export关键字输出该变
	- 使用export命令输出变量
		// profile.js
		export var firstName = 'Michael';
		export var lastName = 'Jackson';
		export var year = 1958;
		// profile.js  为了方便写到一行
		export {firstName, lastName, year};
 - import
	- 通过import来导入其他模块中定义的变量
		import { firstName, lastName, year } from './profile.js';
	- 如果导入的变量和本命令空间有冲突，可以使用 import as 语法来为变量进行重命名
		import { lastName as surname } from './profile.js';
	- 如果一不小心忘记了as，可能会导致变量覆盖，这回到了全局变量的问题，所以模块导出的时候最好设置模块命名空间
 - 命名空间与 export
	- 可以将所有方法绑定到一个变量上暴露出去，避免全局变量的混乱
	- 许多著名的JavaScript库都采用这种方法，如: jQuery，YUI，underscore等等
		- 示例
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
		- 示例 其他文件中
			import { MYAPP } from './export';
	- 默认导出
		- 使用import命令的时候，用户需要知道所要加载的变量名或函数名，否则无法加载
		- 如果想要直接导入模块，再看该模块下有哪些变量可用，就像golang的pkg一样，这个时候就要使用到export default
			- export default命令，为模块指定默认输出
				export default MYAPP
			- 与export命令的区别: 其他模块加载该模块时，import命令可以为该匿名函数指定任意名字
				import myPkg from './export';
				myPkg.foo(); // 'foo'
			- export default命令用于指定模块的默认输出
			- 一个模块只能有一个默认输出，因此export default命令只能使用一次
			- 所以，import命令后面才不用加大括号，因为只可能唯一对应export default命令
	- 注意版本
		- ES版本不同，有2种导入语法
			- require/exports: CommonJS/AMD 中为了解决模块化语法而引入的，ES5语法，注意NodeJs使用的该规范
			- import/export: ES6引入的新规范，因为浏览器引擎兼容问题，需要在node中用babel将ES6语法编译成ES5语法
		- CommonJs的语法
			- 一个模块想要对外暴露变量(函数也是变量)
				module.exports = variable;
			- 一个模块要引用其他模块暴露的变量
				var ref = require('module_name');    // 拿到引用模块的变量
			- 这里的机制是为每个文件准备一个module对象，把它装进去
				var module = {
					id: 'hello',
					exports: {}
				};

7、函数
 - JavaScript中的函数签名
	function abs(x) {
		if (x >= 0) {
			return x;
		} else {
			return -x;
		}
	}
	- function指出这是一个函数定义
	- abs是函数的名称
	- (x)括号内列出函数的参数，多个参数以,分隔
	- { ... }之间的代码是函数体，可以包含若干语句，甚至可以没有任何语句
 - 方法
	- 绑定方法
		var person = {name: '小明', age: 23}
		person.greet = function() {
			console.log(`hello, my name is ${this.name}`)
		}
		person.greet()
	- this 
		- 注意this的使用，没有绑带在对象上, this 指的是 浏览器的window对象
			var person = {name: '小明', age: 23}
			person.greetfn = function() {
				return function() {
					console.log(`hello, my name is ${this.name}`)
				}
			}
			person.greetfn()()                                   // hello, my name is undefined
		- 此时需要通过一个变量 + 闭包，把当前this传递过去，确保this正常传递
			var person = {name: '小明', age: 23}
			person.greetfn = function() {
				var that = this                                  // js中的特殊用法
				return function() {
					console.log(`hello, my name is ${that.name}`)
				}
			}
			person.greetfn()()                                   // hello, my name is 小明
 - 箭头函数(匿名函数)
	- 箭头函数语法
		(params ...) => { ... }
		- 示例一
			fn = x => x * x
			console.log(fn(10))  // 100
		- 示例二
			x => x * x
			// 等价于下面这个函数
			function (x) {
				return x * x;
			}
		- 示例三
			- 箭头函数看上去是匿名函数的一种简写，但实际上，箭头函数和匿名函数有个明显的区别
				- 箭头函数内部的this是词法作用域，由上下文确定
					var person = {name: '小明', age: 23}
					person.greetfn = function() {
						return () => {
							// this 继承自上层的this
							console.log(`hello, my name is ${this.name}`)
						}
					}
		- 示例四
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

8、条件判断
 - 语法格式
	if (condition) {
		...
	} else if (condition) {
		...
	} else {
		...
	}
	- 注意条件需要加上括号，其他和Go语言的if一样

9、for 循环
 - 语法格式
	for (初始条件; 判断条件; 修改变量) {
		...
	}
	- 注意条件需要加上括号
 - for in (不推荐使用)
	- for循环的一个变体是for ... in循环，它可以把一个对象的所有属性依次循环出来
		- 遍历对象: 遍历出来的属性是元素的key
			var o = {
				name: 'Jack',
				age: 20,
				city: 'Beijing'
			};
			for (var key in o) {
				console.log(key);            // 'name', 'age', 'city'
			}
		- 遍历数组: 一个Array数组实际上也是一个对象，它的每个元素的索引被视为一个属性
			var a = ['A', 'B', 'C'];
			for (var i in a) {
				console.log(i);              // '0', '1', '2'
				console.log(a[i]);           // 'A', 'B', 'C'
			}
		- 缺陷
			- 当手动给Array对象添加了额外的属性后，for ... in循环将带来意想不到的意外效果
				var a = ['A', 'B', 'C'];
				a.name = 'Hello';
				for (var x in a) {
					console.log(x);     // '0', '1', '2', 'name'
				}
			- 这和for in的遍历机制相关: 遍历对象的属性名称
			- 使用 for of 解决这个问题
 - for of
	- for ... of循环则修复了这些问题，它只循环集合本身的元素
		var a = ['A', 'B', 'C'];
		a.name = 'Hello';
		for (var x of a) {
			console.log(x);             // 'A', 'B', 'C'
		}
	- 但是for ... of循环不可用直接遍历对象，可以通过Object提供的方法获取key数组，然后遍历
		var o = {
			name: 'Jack',
			age: 20,
			city: 'Beijing'
		};
		for (var key of Object.keys(o)) {
			console.log(key);                // 'name', 'age', 'city'
		}
 - forEach方法
	- forEach()方法是ES5.1标准引入的，是遍历元素的一种常用手段，能作用于可迭代对象上，和for of一样
		arr.forEach(function(item) {console.log(item )})
		arr.forEach((item) => {console.log(item)})
 - for循环应用
	- 如果后端返回的数据不满足展示的需求，需要修改，比如vendor想要友好显示，就可以直接修改数据
	
10、Promise对象
 - 在JavaScript的世界中，所有代码都是单线程执行的
	- 由于这个"缺陷"，导致JavaScript的所有网络操作，浏览器事件，都必须是异步执行
	- Javascript通过回调函数实现异步，js的一大特色
 - 单线程异步模型
	function callback() {
		console.log('Done');
	}
	console.log('before setTimeout()');
	setTimeout(callback, 1000);         // 1秒钟后调用callback函数
	console.log('after setTimeout()');
	// before setTimeout()
	// after setTimeout()
	// 等待一秒后
	// Done
 - Promise与异步
	- Promise
		interface PromiseConstructor {
			/**
			* Creates a new Promise.
			* @param executor A callback used to initialize the promise. This callback is passed two arguments:
			* a resolve callback used to resolve the promise with a value or the result of another promise,
			* and a reject callback used to reject the promise with a provided reason or error.
			*/
			new <T>(executor: (resolve: (value: T | PromiseLike<T>) => void, reject: (reason?: any) => void) => void): Promise<T>;
		}
	- 测试一个函数，resolve是成功后的回调函数，reject是失败后的回调函数
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
		function testResultCallback() {
			success = (message) => {console.log(`success ${message}`)}
			failed = (error) => {console.log(`failed ${error}`)}
			testResultCallbackFunc(success, failed)
		}
	- 将回调改为Promise对象，而Promise的优势在于异步执行的流程中，把执行代码和处理结果的代码清晰地分离
		var p1 = new Promise(testResultCallbackFunc)  // 执行testResultCallbackFunc函数
		p1.then((resp) => {                           // 分析结果
			console.log(resp)
		}).catch((err) => {
			console.log(err)
		})
 - Async函数 + Promise组合
	- Async函数由内置执行器进行执行, 这和go func() 有异曲同工之妙
	- 如果声明一个异步函数，在函数前面加上一个 async关键字
		async function testWithAsync() {                        // async调用协程
			var p1 = new Promise(testResultCallbackFunc)        // new是将函数构造成一个对象
			try {
				var resp = await p1                             // await等待promise的结果
				console.log(resp)
			} catch (err) {
				console.log(err)
			}
		}
	- 这里testWithAsync就是一个异步函数，在执行的时候 是交给js的携程执行器处理的
	- 而 await关键字就是 告诉执行器当p1执行完成后，主动通知下(协程的一种实现)
	- 其实就是一个 event pool模型(简称epool模型)

11、关于JavaStript的总结
 - 推荐使用箭头函数
 - 判断使用 ===
 - 由于变量提升问题，尽量使用let声明变量，并且写在开头
 - for循环推荐forEach

12、Jquery
 - Jquery 是一JavaScript的一个库，极大的简化了js的开发
	- 包含
		- Html元素选取
		- Html元素操作
		- Css操作
		- 绑定事件
		- Ajax
		- 工具函数
		- 插件系统

五、浏览器基础
1、浏览器
 - 场景
	- 适配屏幕
	- 兼容多种浏览器
	- JavaScript可以获取浏览器提供的很多对象，并进行操作
 - 浏览器对象
	- 网页是通过浏览器加载出来的，而浏览器本身也有很多功能
		- 浏览器本身版本信息
		- 当前窗口大小
		- 历史记录
		- ...
	- 而这些数据都是绑定在window对象上，可以认为对于浏览器加载的网页，window就表示浏览器本身
 - 浏览器窗口大小
	- window的如下4个属性控制这个浏览器的窗口大小
		- IE9+、Safari、Opera和Chrome都支持这4个属性
			- innerWidth: 页面视图的宽
			- innerHeight: 页面视图的高
			- outerWidth: 浏览器整个窗口的宽
			- outerHeight: 浏览器整个窗口的高
		- 调整浏览器窗口大小
			- console.log('window inner size: ' + window.innerWidth + ' x ' + window.innerHeight);
			- console.log('window outer size: ' + window.outerWidth + ' x ' + window.outerHeight);
 - 浏览器信息
	- 关于浏览器的相关信息都保存在navigator对象上面
		window.navigator
		- 常用的属性
			navigator.appName: 浏览器名称；
			navigator.appVersion: 浏览器版本；
			navigator.language: 浏览器设置的语言；
			navigator.platform: 操作系统类型；
			navigator.userAgent: 浏览器设定的User-Agent字符串
			navigator.userAgentData: useragent相关信息, Object
 - 屏幕信息
	- 需要知道用户的屏幕尺寸，以便做网页布局，这是就需要获取屏幕的数据
		window.screen
 - 访问网址
	- 为了保证每个用户访问该URL都呈现统一的效果，页面往往都需要读取当前Page的参数
		window.location
		- 常用的属性
			location.protocol;     // 'http'
			location.host;         // 'www.example.com'
			location.port;         // '8080'
			location.pathname;     // '/path/index.html'
			location.search;       // '?a=1&b=2'
			location.hash;         // 'TOP'
			location.assign(url)   // 加载页面
			location.replace(url)  // 重定向页面
			location.reload()      // 重新加载当前页面
		- 示例
			if (confirm('重新加载当前页' + location.href + '?')) {
				location.reload();
			} else {
				location.assign('/'); // 设置一个新的URL地址
			}
 - 历史记录
	- history对象保存了浏览器的历史记录
	- JavaScript可以调用history对象的back()或forward ()，相当于用户点击了浏览器的"后退"或"前进"按钮
		- <- 浏览器的"后退"按钮
			history.back();
		- -> 浏览器的"前进"按钮
			history.forward()

2、AJAX
 - Asynchronous JavaScript and XML，意思就是用JavaScript执行异步网络请求
	- 简单来说就是浏览器中的Http Client
	- 在现代浏览器上写AJAX主要依靠XMLHttpRequest对象
	- 异步获取数据 + 局部刷新页面
 - XMLHttpRequest
	- 声明XMLHttpRequest对象
		var request = new XMLHttpRequest(); 
	- 发送请求
		request.open('GET', 'https://www.baidu.com');
		request.send();
	- 由于js都是异步的，需要定义回调来处理返回，用于获取请求结果
		request.onreadystatechange = function () {
			console.log(request.status)
			console.log(request.responseText)
		}
	- 对于AJAX请求特别需要注意跨域问题: CORS  https://www.ruanyifeng.com/blog/2016/04/cors.html
		- 简单请求
			- 请求特征
				- HTTP Method: GET、HEAD和POST
				- Content-Type: application/x-www-form-urlencoded、multipart/form-data和text/plain
				- 不能出现任何自定义头
			- 请求流程
				- 通常能满足90%的需求
				- 控制其跨域的关键头 来自于服务端设置的: Access-Control-Allow-Origin
				- 对于简单请求，浏览器直接发出CORS请求
					- 具体来说，就是在头信息之中，增加一个Origin字段
				- 客户端发送请求时 Origin: xxxx
				- 服务端响应请求时 Access-Control-Allow-Origin: xxxx / *
				- 如果服务端不允许该域，就跨越失败
				- 如果想允许特定的Header，服务端也可以通过添加 Access-Control-Expose-Headers 来进行控制
		- 复杂请求
			- 请求特征
				- HTTP Method: PUT、DELETE
				- Content-Type: 其他类型如application/json
				- 含自定义头，比如后面携带的用于认证的X-OAUTH-TOKEN头
			- 请求流程
				- 非简单请求的CORS请求，会在正式通信之前，增加一次HTTP查询请求，称为"预检"请求(preflight)
					- preflight请求
						- Method: OPTIONS
						- Header: Access-Control-Request-Method   列出浏览器的CORS请求会用到哪些HTTP方法
						- Header: Access-Control-Request-Headers  该字段是一个逗号分隔的字符串，指定浏览器CORS请求会额外发送的头信息字段
							OPTIONS /cors HTTP/1.1
							Origin: http://api.bob.com
							Access-Control-Request-Method: PUT
							Access-Control-Request-Headers: X-Custom-Header
							Host: api.alice.com
							Accept-Language: en-US
							Connection: keep-alive
							User-Agent: Mozilla/5.0...
					- preflight响应
						- Header: Access-Control-Allow-Origin    允许的源(域)
						- Header: Access-Control-Allow-Methods   服务端 允许的方法(所有方法，一般大于请求的该头)
						- Header: Access-Control-Allow-Headers   服务端 允许的自定义Header
						- Header: Access-Control-Max-Age         用来指定本次预检请求的有效期，避免多次请求，该字段可选
							Access-Control-Allow-Methods: GET, POST, PUT
							Access-Control-Allow-Headers: X-Custom-Header
							Access-Control-Allow-Credentials: true
							Access-Control-Max-Age: 1728000
 - axios
 - 页面(DOM)
	- DOM
		- 浏览器当前加载的Page也是一个对象: window.document
		- 由于HTML在浏览器中以DOM形式表示为树形结构，document对象就是整个DOM树的根节点
		- 获取当前Document中的一些属性
			- 比如 title document的title属性是从HTML文档中的<title>xxx</title>读取的，但是可以动态改变
				document.title
				// 'golang random int_百度搜索'
				// 动态修改
				document.title = '测试修改'
		- 很多js库比如JQuery都是通过动态操作Dom来实现很多高级功能的，这些是上层库的基石
	- DOM查询
		- 要查找DOM树的某个节点，需要从document对象开始查找
		- 最常用的查找是根据ID和Tag Name以及ClassName
			<h1>列表</h1>
			<ul id="list_menu" class="ul_class">
				<li>Coffee</li>
				<li>Tea</li>
				<li>Milk</li>
			</ul>
			document.getElementById('list_menu')
			// <ul data-v-7ba5bd90 id=​"list_menu">​<li data-v-7ba5bd90>​…​</li>​<li data-v-7ba5bd90>​…​</li>​<li data-v-7ba5bd90>​…​</li>​</ul>​
			document.getElementsByTagName('ul')
			// 返回所有的ul元素
			document.getElementsByClassName('ul_class')
		- 组合使用
			document.getElementById('list_menu').getElementsByTagName('li')
			// 上面也等价于
			document.getElementById('list_menu').children
			// firstElementChild
			document.getElementById('list_menu').firstElementChild
			// lastElementChild
			document.getElementById('list_menu').lastElementChild
			// parentElement 获取父元素
			document.getElementById('list_menu').parentElement
	- 更新DOM
		- 获取的元素后，可以通过元素的如下2个方法，修改元素
			- innerHTML: 不但可以修改一个DOM节点的文本内容，还可以直接通过HTML片段修改DOM节点内部的子树
			- innerText: 只修改文本内容
				var le = document.getElementById('list_menu').lastElementChild
				le.innerText = '牛奶'                                           // 页面的内容已经修改
				le.innerText = '<span style="color:red">牛奶</span>'            // html格式无法识别
				le.innerHTML = '<span style="color:red">牛奶</span>'            // 使用innterHTML则可以
	- 插入DOM
		- 很多响应式框架都会根据数据新增，动态创建一个DOM元素，并插入到指定位置
			- 使用 createElement 来创建一个DOM元素，比如创建一个a标签
				
				var newlink = document.createElement('a')                       // 创建一个A标签
				newlink.href = "http://www.baidu.com"                           // <a></a>
				newlink.innerText = '跳转到百度'                                // 修改A标签属性
				var lm = document.getElementById('list_menu')                   // 追加到某个元素后面
				lm.appendChild(newlink)
			- 想要控制元素插入的位置可以使用insertBefore
				parentElement.insertBefore(newElement, referenceElement);       // 父元素
				var lm = document.getElementById('list_menu')                   // 子元素
				var cf = document.getElementById('coffee')                      // 需要插入的元素
				var newlink = document.createElement('a')
				newlink.href = "http://www.baidu.com"   
				newlink.innerText = '跳转到百度'
				lm.insertBefore(newlink, cf)                                    // 插入到coffee之前
		- 有2种方式可以插入一个DOM元素
			- appendChild: 把一个子节点添加到父节点的最后一个子节点
			- insertBefore: 插入到某个元素之前
	- 删除DOM
		- 删除一个节点，首先要获得该节点本身以及它的父节点，然后调用父节点的removeChild把自己删掉
			parent.removeChild(childNode);
		- 删除刚才添加的那个元素
			lm.removeChild(newlink)


# Web-Base Web的Vue框架

## 一、Vue介绍

### 1. 

六、Vue
1、Vue的简介
 - Vue 是一套用于构建用户界面的渐进式框架
	- 与其它大型框架不同的是，Vue 被设计为可以自底向上逐层应用
 - Vue 的核心库只关注视图层，不仅易于上手，还便于与第三方库或既有项目整合，比如实现拖拽: vue + sortable.js
 - vue借鉴了很有框架优秀的部分进行了整合
	- 借鉴 angular 的模板和数据绑定技术
	- 借鉴 react 的组件化和虚拟 DOM 技术
 - 所有框架的逻辑都是通过js封装，上层概念比如MVVM，方便快速开发，因此学习任何框架前都需要具有Web基础
	- HTML
	- CSS
	- Javascript

2、Vue.js与环境的安装
1).安装Node.js
 - 一个基于 Chrome V8 引擎 的 JavaScript 运行时环境
	https://nodejs.org/zh-cn/
	$ node -v   // v14.17.1
	$ npm -v    // 6.14.13
 - 配置使用国内源
	$ npm config set registry http://registry.npm.taobao.org/
	$ npm config get registry

2).安装Vue和脚手架工具
 - 直接使用Node.js进行全局安装
	# 最新稳定版
	$ npm install -g vue            // -g 参数即为全局安装
 - 安装项目脚手架工具
	- 官方同时提供了1个cli工具，用于快速初始化一个vue工程，官方文档 https://cli.vuejs.org/
		$ npm install -g @vue/cli
		$ npm list -g               // 查看当前Node.js环境下安装的包
		
3).Vue Devtools
 - chrome商店
 - https://www.jianshu.com/p/63f09651724c
 
4).vscode 插件
 - Beautify:          js, css, html 语法高亮差距
 - ESLint:            js eslint语法风格检查
 - Auto Rename Tag:   tag rename
 - Vetur:             vue语法高亮插件
 
5).项目
 - vue项目创建
	$ vue create demo
	$ cd demo 
	$ npm run serve
 - 项目结构
	- 通过vue-cli搭建一个vue项目，会自动生成一系列文件
		├── dist/                      	# 项目构建后的产物
		├── node_module/               	# 项目中安装的依赖模块
		├── public/                    	# 纯静态资源，入口文件(index.html)也在里面
		├── src/
		│   ├── main.js                 # 程序入口文件
		│   ├── App.vue                 # 程序入口vue组件，大写字母开头，后缀 .vue
		│   ├── components/             # 组件
		│   │   └── ...
		│   └── assets/                 # 资源文件夹，一般放一些静态资源文件，比如CSS/字体/图片
		│       └── ...
		├── babel.config.js             # babel 配置文件，es6语法转换
		├── .gitignore                  # 用来过滤一些版本控制的文件，比如node_modules文件夹 
		└── package.json                # 项目文件，记载着一些命令和依赖还有简要的项目描述信息 
		└── README.md                   # 介绍创建的项目，可参照github上star多的项目
 - 部署项目
	$ npm run build ## 会在项目的dist目录下生成html文件, 使用这个静态文件部署即可
	// 比如使用python快速搭建一个http静态站点，如果是nginx copy到 对应的Doc Root位置
	$ cd dist
	$ python3 -m http.server

3、MVVM
 - MVVM模型
	- M  (Model，模型层): 模型层，主要负责业务数据相关，对应vue中的 data部分
	- V  (View，视图层) : 视图层，顾名思义，负责视图相关，细分下来就是html+css层，对应于vue中的模版部分
	- VM (ViewModel, 控制器): V与M沟通的桥梁，负责监听M或者V的修改，是实现MVVM双向绑定的要点，对应vue中双向绑定
 - Web技术的历史
	- CGI时代
	- 后端模版时代
		- ASP: 微软，C#体系
		- JSP: SUN，Java体系
		- PHP: 开源社区
	- JavaScript原生时代
	- 前端模版时代
		- ViewModel动态完成渲染
	- 虚拟DOM技术
	- 组件化时代
		- MVVM最早由微软提出来，它借鉴了桌面应用程序的MVC思想，在前端页面中，把Model用纯JavaScript对象表示，View负责显示，两者做到了最大限度的分离
		- 结合虚拟Dom技术，就可以动态生成view，在集合mvvm的思想，前端终于迎来了组件化时代
			- 页面由多个组建构成
			- 每个组件都有自己的 MVVM
	
4、Vue与MVVM
	- Model: vue中用于标识model的数据是 data对象，data 对象中的所有的 property 加入到 vue 的响应式系统中，由Vue监听变化  // Data
	- View:  vue使用模版来实现展示，但是渲染时要结合 vdom技术	                                                         // Html、Css 
	- ViewModle: vue的核心，负责视图的响应，也就是数据双向绑定
		- 监听view中的数据，如果数据有变化，动态同步到 data中
		- 监听data中的数据，如果数据有变化，通过vdom动态渲视图

5、Vue实例
 - Vue实例源码
	export interface Vue {
		// 生产的HTML元素
		readonly $el: Element; 
		// 实例的一些配置, 比如components, directives, filters ...
		readonly $options: ComponentOptions<Vue>;
		// 父Vue实例 
		readonly $parent: Vue;
		// Root Vue实例
		readonly $root: Vue;
		// 该vue实例的子vue实例，一般为 子组建
		readonly $children: Vue[];
		// 元素refs, 当元素有refs属性时才能获取
		readonly $refs: { [key: string]: Vue | Element | (Vue | Element)[] | undefined };
		// 插槽, 模版插槽
		readonly $slots: { [key: string]: VNode[] | undefined };
		readonly $scopedSlots: { [key: string]: NormalizedScopedSlot | undefined };
		readonly $isServer: boolean;
		// Model对应的数据
		readonly $data: Record<string, any>;
		// 实例props, 用于组件见消息传递
		readonly $props: Record<string, any>;
		readonly $ssrContext: any;
		// 虚拟node
		readonly $vnode: VNode;
		readonly $attrs: Record<string, string>;
		readonly $listeners: Record<string, Function | Function[]>;
		// 实例挂在到具体的HTML上
		$mount(elementOrSelector?: Element | string, hydrating?: boolean): this;
		// 强制刷新渲染, 收到刷新界面, 当有些情况下 界面没响应时
		$forceUpdate(): void;
		// 销毁实例
		$destroy(): void;
		$set: typeof Vue.set;
		$delete: typeof Vue.delete;
		// watch对象变化
		$watch(
			expOrFn: string,
			callback: (this: this, n: any, o: any) => void,
			options?: WatchOptions
		): (() => void);
		$watch<T>(
			expOrFn: (this: this) => T,
			callback: (this: this, n: T, o: T) => void,
			options?: WatchOptions
		): (() => void);
		$on(event: string | string[], callback: Function): this;
		$once(event: string | string[], callback: Function): this;
		$off(event?: string | string[], callback?: Function): this;
		// 触发事件
		$emit(event: string, ...args: any[]): this;
		// Dom更新完成后调用
		$nextTick(callback: (this: this) => void): void;
		$nextTick(): Promise<void>;
		$createElement: CreateElement;
	}
 - Vue实例的构造函数
	- Data: Model
	- Methods: 方法
	- Computed: 计算属性
	- Props: 类似于一个自定义 attribute
		new <Data = object, Methods = object, Computed = object, Props = object>(options?: ThisTypedComponentOptionsWithRecordProps<V, Data, Methods, Computed, Props>): CombinedVueInstance<V, Data, Methods, Computed, Record<keyof Props, any>>;

6、Vue实例生命周期
 - 生命周期
	https://gitee.com/infraboard/go-course/blob/master/day16/pic/lifecycle.png
 - 定义的函数钩子
	beforeCreate()
	created()
	beforeMount()
	mounted()
	beforeUpdate()
	updated()
	beforeDestroy()
	destroyed()

7、模板语法
1)、模板定义
 - 通过template标签定义的部分都是vue的模版
	- 模版会被vue-template-compiler编译后渲染
		<template>
		  ...
		</template>

2)、文本值
 - 当需要访问Model时，比如 data 这个Object
	- 直接使用 {{ attr }} 就可以访问
	- vue会根据属性是否变化，而动态渲染模版
		<template>
		<div>{{ name }}</div>              // data中的 name属性
		</template>
		<script>
		export default {
		  name: 'HelloWorld',
		  data() {
		    return {
		      name: '老喻'
		    }
		  }
		}
		</script>

3)、元素属性
 - 变量不能作用在 HTML attribute 上
	- 示例一  语法错误
		<template>
		  <div id={{ name }}>{{ name }}</div>
		</template>
	- 示例二
		- buttom有个disabled属性，用于控制当前按钮是否可以点击
		<template>
		  <button disabled="true">Button</button>
		</template>
 - 针对HTML元素的属性 vue专门提供一个 v-bind指令，这个指令就是模版引擎里面的一个函数，作用是完成HTML属性变量替换
	- v-bind:disabled="attr"   ==>  disabled="data.attr"
		<template>
		  <button v-bind:disabled="isButtomDisabled">Button</button>
		</template>
		<script>
		export default {
		  name: 'HelloWorld',
		  data() {
		    return {
		      name: '老喻',
		      isButtomDisabled: false,
		    }
		  }
		}
		</script>
	- v-binding 缩写 `:`
		<template>
		  <button :disabled="isButtomDisabled">Button</button>
		</template>

4)、元素事件
 - 给buttom元素绑定一个事件
	- HTML 原生语法  https://www.runoob.com/tags/ref-eventattributes.html
		<button onclick="copyText()">复制文本</button>
 - 对于vue的模版系统来说，copyText这个函数如何渲染，他不是一个文本，而是一个函数
	- vue针对事件专门定义了一个指令: v-on
		v-on:eventName="eventHandler"
		eventName: 事件的名称
		eventHandler: 处理这个事件的函数
	- data 是定义Model的地方，vue专门给一个属性用于定义方法: methods
		<template>
		  <button :disabled="isButtomDisabled" v-on:click="clickButtom" >Button</button>
		</template>
		<script>
		export default {
		  name: 'HelloWorld',
		  data() {
		    return {
			  name: '老喻',
			  isButtomDisabled: false,
		    }
		  },
		  methods: {
		    clickButtom() {
		      alert("别点我")  
		    }
		  }
		}
		</script>
	- v-on 缩写 `@`
		<template>
		  <button :disabled="isButtomDisabled" @click="clickButtom" >Button</button>
		</template>

4)、Vue指令
 - 指令类型
	v-model: 双向绑定的数据
	v-bind: html元素属性绑定
	v-on: html元素事件绑定
	v-if: if 渲染
	v-show: 控制是否显示
	v-for: for 循环
 - 指令
	- 语法
		v-directive:argument.modifier.modifier...
	- v-directive: 表示指令名称，如v-on
	- argument： 表示指令的参数，比如click
	- modifier:  修饰符，用于指出一个指令应该以特殊方式绑定
	- 示例 当用户按下回车时，表示用户输入完成，触发搜索
		v-directive: 需要使用绑定事件的指令: v-on
		argument:    监听键盘事件: keyup，按键弹起时
		modifier:    监听Enter建弹起时
		完整语法:    v-on:keyup.enter
			<template>
			  <input v-model="name" type="text" @keyup.enter="pressEnter">
			</template>
			<script>
			export default {
			  name: 'HelloWorld',
				data() {
				return {
				  name: '老喻',
				  isButtomDisabled: false,
				}
			  },
			  methods: {
			    clickButtom() {
			      alert("别点我")  
			    },
			    pressEnter() {
			 	  alert("点击了回车键")
			    }
			  },
			}
			</script>
	- 需要注意事件的指令的函数是可以接受参数的
		<template>
		  <input v-model="name" type="text" @keyup.enter="pressEnter(name)">
		  <button v-on:click="say('hi')">Say hi</button>
		</template>
	- 函数是直接读到model数据的，不能使用{{ }}，如果要传字符串 使用''
	- 修饰符还有很多其他用法，具体的请看官方文档
		<!-- 即使 Alt 或 Shift 被一同按下时也会触发 -->
		<button v-on:click.ctrl="onClick">A</button>
		
		<!-- 有且只有 Ctrl 被按下的时候才触发 -->
		<button v-on:click.ctrl.exact="onCtrlClick">A</button>
		
		<!-- 没有任何系统修饰符被按下的时候才触发 -->
		<button v-on:click.exact="onClick">A</button>

5)、JavaScript 表达式
	- 模版支持JavaScript的表达式, 可以在显示的动态的做一些处理
		<template>
		  <div>{{ name.split('').reverse().join('') }}</div>
		</template>
		<script>
		export default {
		  name: 'HelloWorld',
		  data() {
		    return {
		      name: '老喻'
		    }
		  }
		}
		</script>

6)、条件渲染
 - 有2个指令用于在模版中控制条件渲染
	v-if: 控制元素是否创建，创建开销较大
	v-show: 控制元素是否显示，对象无效销毁，开销较小
 - v-if 语法
	<t v-if="" /> 
	<t v-else-if="" /> 
	<t v-else="" /> 
 - v-show 语法
	<t v-show="" />
 - 示例 根据用户输入，判断当前分数的等级
	<input v-model="name" type="text" @keyup.enter="pressEnter(name)">
	<div v-if="name >= 90">
	  A
	</div>
	<div v-else-if="name >= 80">
	  B
	</div>
	<div v-else-if="name >= 60">
	  C
	</div>
	<div v-else-if="name >= 0">
	  D
	</div>
	<div v-else>
	  请输入正确的分数
	</div>
	// 使用 v-show
	<input v-model="name" type="text" @keyup.enter="pressEnter(name)">
	<div v-show="name >= 90">
	  A
	</div>
	<div v-show="name >= 80 && name < 90">
	  B
	</div>
	<div v-show="name >= 60 && name < 80">
	  C
	</div>
	<div v-show="name >= 0 && name < 60">
	  D
	</div>
	- 一般来说，v-if 有更高的切换开销，而 v-show 有更高的初始渲染开销
	- 因此，如果需要非常频繁地切换，则使用 v-show 较好
	- 如果在运行时条件很少改变，则使用 v-if 较好

7)、列表渲染
 - v-for元素的列表渲染
	<t v-for="(item, index) in items" :key="item.message">
	  {{ item.message }}
	</t>
	<!-- items: [
	  { message: 'Foo' },
	  { message: 'Bar' }
	] -->
 - 也可以省略 index
	<ul>
	  <li v-for="item in items" :key="item.message">
	    {{ item.message }}
	  </li>
	</ul>
	<script>
	export default {
	  name: 'HelloWorld',
	  data() {
	    return {
	    items: [
		  { message: 'Foo' },
		  { message: 'Bar' }
	    ]
	    }
	  },
	}
	</script>
 - v-for 除了可以遍历列表，可以遍历对象，比如嵌套2层循环，先遍历列表，再遍历对象
	<ul>
	  <li v-for="(item, index) in items" :key="item.message">
	    {{ item.message }} - {{ index}}
	    <br>
	    <span v-for="(value, key) in item" :key="key"> {{ value }} {{ key }} <br></span>
	  </li>
	</ul>
	<script>
	export default {
	  name: 'HelloWorld',
	  data() {
	    return {
	      items: [
	        { message: 'Foo', level: 'info' },
	        { message: 'Bar', level: 'error'}
	      ]
	    }
	  }
	}
	</script>
 - 可以在console界面里进行数据修改测试
	$vm._data.items.push({message: "num4", level: "pannic"})
	$vm._data.items.pop()
 - 不推荐在同一元素上使用 v-if 和 v-for，请另外单独再起一个元素进行条件判断
		<li v-for="todo in todos" v-if="!todo.isComplete">
		{{ todo }}
		</li>
	- 请改写成下面方式:
		<ul v-if="todos.length">
		  <li v-for="todo in todos">
		    {{ todo }}
		  </li>
		</ul>
		<p v-else>No todos left!</p>

8、计算属性
 - 如果model的数据并不是要直接渲染的，需要处理再展示，简单的方法是使用表达式
	- 示例
		<h2>{{ name.split('').reverse().join('') }}</h2>
	- 这种把数据处理逻辑嵌入的视图中，并不合适，不易于维护，可以把改成一个方法	
		<h2>{{ reverseData(name) }}</h2>
		<script>
		  methods: {
		    reverseData(data) {
		      return data.split('').reverse().join('')
		    }
		  }
		</script>
	- 除了函数，vue还提供了一个计算属性，这样视图可以看起来更干净
		computed: {
		  attrName: {
		    get() {
		      return value
		    },
		    set(value) {
		    // set value
		    }
		  }
		}
	- 修改为计算属性
		<h2>{{ reverseName }}</h2>
		<script>
		export default {
		  computed: {
		    reverseName: {
		      get() {
		        return this.name.split('').reverse().join('')
		      },
		      set(value) {
		        this.name = this.name = value.split('').reverse().join('')
		      }
		    }
		  },
		}
		</script>

9、侦听器
 - 一个页面有多个参数，用户可能把url copy给别人，需要不同的url看到页面内容不同，不然用户每次到这个页面都是第一个页面
	- 这个就需要监听url参数的变化，然后视图做调整，vue-router会有个全局属性: $route 可以监听它的变化
 - window本身提供一个事件回调
	window.onhashchange = function () {
	  console.log('URL发生变化了', window.location.hash);
	  this.urlHash = window.location.hash
	};
 - vue 提供的属性watch
	watch: {
	  // 如果 `urlHash` 发生改变，这个函数就会运行
	  urlHash: function (newData, oldData) {
	    this.debouncedGetAnswer()
	  }
	},
 - 先监听变化，挂载后修改vue对象，然后watch
	<script>
	export default {
	  name: 'HelloWorld',
	  data() {
	    return {
	      urlHash: '',
	    }
	  },
	  mounted() {
	    /* 来个骚操作 */
	    let that = this
	    window.onhashchange = function () {
	      that.urlHash = window.location.hash
	    };
	  },
	  watch: {
	    urlHash: function(newURL, oldURL) {
	      console.log(newURL, oldURL)
	    }
	  }
	}
	</script>
 - Vue Watch API  
	https://cn.vuejs.org/v2/api/#vm-watch

10、过滤器
 - Vue.js 允许自定义过滤器，可被用于一些常见的文本格式化，最常见的是 时间的格式化
	- 过滤器语法
		<!-- 在双花括号中 -->
		{{ message | capitalize }}
	- 可以将其等价于一个函数: capitalize(message)
	- 在当前组件的vue实例上定义一个过滤器
		filters: {
		  capitalize: function (value) {
		    /*过滤逻辑*/
		  }
		}
 - 定义 parseTime过滤器
	{{ ts | parseTime }}
	<script>
	export default {
	  name: 'HelloWorld',
	  data() {
	    return {
	      ts: Date.now()
	    }
	  },
	  filters: {
	    parseTime: function (value) {
	      let date = new Date(value)
	      return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
	    }
	  }
	}
	</script>
 - vue提供全局过滤器，再初始化vue实例的时候可以配置，找到main.js添加
	// 添加全局过滤器
	Vue.filter('parseTime', function (value) {
	  let date = new Date(value)
	  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
	})

11、自定义指令
 - 除了核心功能默认内置的指令 (v-model 和 v-show)，Vue 也允许注册自定义指令
	- 比如用户进入页面让输入框自动聚焦，方便快速输入
	- 比如登陆页面，快速聚焦到 username输入框
	- 如果是HTML元素聚焦，找到元素，调用focus
		let inputE = document.getElementsByTagName('input')
		inputE[0].focus()
	- 添加到mounted中进行测试
		mounted() {
		  let inputE = document.getElementsByTagName('input')
		  inputE[0].focus()
		}
    - 将这个功能做成一个vue的指令，比如 v-focus
		- 先注册一个局部指令，在本组件中使用
			<script>
			export default {
			  name: 'HelloWorld',
			  directives: {
			    focus: {
			      // 指令的定义
			      inserted: function (el) {
			        el.focus()
				  }
				}
			  },
			}
			</script>
		- 这里注册的指令名字叫focus，所有的指令在模版要加一个v前缀，因此注册的指令就是v-focus
			<input v-focus v-model="name" type="text" @keyup.enter="pressEnter(name)">
			// 注册一个全局自定义指令 `v-focus`
			Vue.directive('focus', {
			  // 当被绑定的元素插入到 DOM 中时……
			  inserted: function (el) {
			    // 聚焦元素
			    el.focus()
			  }
			})

12、组件
1)、组件的介绍
 - 组件是可复用的 Vue 实例，所以它们与 new Vue 接收相同的选项，例如 data、computed、watch、methods 以及生命周期钩子等，仅有的例外是像 el 这样根实例特有的选项
	- export default语法，ES6的模块语法
	- 每个vue组件，只能有1个根元素，所以不能在模版里面写2个并排的div
	- 组件的命名风格是首字母大写，导出名字最好和文件名称相同，方便阅读
	- 局部组件必须放到实例的components属性里面，代表这个组件注册到了该实例，只有注册后的组件才能使用
	- 组件的使用方式和HTML标签一样，可以认为是自定义的HTML标签，应该遵循html标签规范，小写加-链接，比如hello-world
		<template>
		  <div id="app">
		    <img alt="Vue logo" src="./assets/logo.png">
		    <hello-world msg="Welcome to Your Vue.js App"/>
		    <hello-world msg="component2"/>
		  </div>
		</template>
2)、向子组件传递数据
 - props就是用于传递数据的变量，使用Prop 可以在组件上注册的一些自定义 attribute，父组件通过绑定数据来传递消息给子组件
 - 子组件
	<script>
	export default {
	  name: 'HelloWorld',
	  props: {
	    msg: String
	  }
	}
	</script>
 - 父组件
	<template>
	  <div id="app">
	    <img alt="Vue logo" src="./assets/logo.png">
	    <hello-world :msg="msg1"/>
	  </div>
	</template>
	<script>
	import HelloWorld from './components/HelloWorld.vue'
	export default {
	name: 'App',
	  data() {
	    return {
	      msg1: 'Welcome to Your Vue.js App'
	    }
	  },
	  components: {
	    HelloWorld
	  }
	}
	</script>
3)、向父组件传递数据
 - 需要条件	 
	- 子组件使用$emit发送事件，事件名称: changeMsg
	- 父组件使用v-on 订阅子组件的事件
 - 具体步骤
	- 修改子组建，重新绑定一个属性: tmpMsg，等点击回车时发生事件给父组件
		<input v-model="tmpMsg" type="text" @keyup.enter="changeProps(tmpMsg)">
		<script>
		export default {
		  name: 'HelloWorld',
		  data() {
		    return {
		      tmpMsg: '',
		    }
		  },
		  methods: {
		    changeProps(msg) {
		      this.$emit('changeMsg', msg)
		    }
		  },
		  props: {
		    msg: String
		  }
		}
		</script>
	- 父组件修改
		<hello-world :msg="msg1" @changeMsg="msgChanged" />
		<script>
		import HelloWorld from './components/HelloWorld.vue'
		export default {
		  name: 'App',
		  data() {
		    return {
		      msg1: 'Welcome to Your Vue.js App'
		    }
		  },
		  methods: {
		    msgChanged(event) {
		      this.msg1 = event
		    }
		  },
		  components: {
		    HelloWorld
		  }
		}
		</script>
4)、使用 v-model 实现
 - 上述使用 v-bind 和 v-on 完成了数据的双向绑定
	- 但其实就是 v-model
		v-model <==>  v-bind:value + v-on:input
	- 将其 value attribute 绑定到一个名叫 value 的 prop 上 在其 input 事件被触发时，将新的值通过自定义的 input 事件抛出
 - 使用 v-model
	- 父组件使用 v-model进行值的双向传递
		<hello-world v-model="msg1" />
	- 子组件绑定value属性和抛出input事件
		<input :value="value" type="text" @input="$emit('input', $event.target.value)">
		<script>
		export default {
		  name: 'HelloWorld',
		  props: {
		    value: String,
		  }
		}
		</script>

13、插件
 - 插件通常用来为 Vue 添加全局功能
 - 插件的功能范围没有严格的限制，一般有下面几种
	- 添加全局方法或者 property，如：vue-custom-element
	- 添加全局资源：指令/过滤器/过渡等，如 vue-touch
	- 通过全局混入来添加一些组件选项，如 vue-router
	- 添加 Vue 实例方法，通过把它们添加到 Vue.prototype 上实现
	- 一个库，提供自己的 API，同时提供上面提到的一个或多个功能，如 vue-router
 - 示例 element-ui 插件
	import ElementUI from "element-ui"
	import "element-ui/lib/theme-chalk/index.css"
	Vue.use(ElementUI)
	// 组件
	import Vue, { PluginObject } from 'vue'
	import { ElementUIComponent, ElementUIComponentSize, ElementUIHorizontalAlignment } from './component'
	import { ElAlert } from './alert'
	import { ElAside } from './aside'
	import { ElAutocomplete } from './autocomplete'
	import { ElBadge } from './badge'
	import { ElBreadcrumb } from './breadcrumb'
	import { ElBreadcrumbItem } from './breadcrumb-item'
	import { ElButton } from './button'
	// 这些element的组件都被注册到vue里面
	import Vue from 'vue'
	/** ElementUI component common definition */
	export declare class ElementUIComponent extends Vue {
	  /** Install component into Vue */
	  static install (vue: typeof Vue): void
	}

14、Vue路由
1)、页面路由
 - 当前的vue配置
	// Root Vue实例
	new Vue({
	  render: h => h(App),
	}).$mount('#app')

2)、简单路由
 - window.location.pathname
	- 为每个path，定义一个组件，就实现了一个简单的路由
		// Root Vue实例
		// 添加currentRoute数据 和 浏览器的path绑定
		// 根据path 返回对应组件
		new Vue({
		  data: {
		    currentRoute: window.location.pathname
		  },
		  render(h) {
		    if (this.currentRoute === '/index') {
		      return h(App2)
		    }
		    return h(App)
		  },
		}).$mount('#app')
	- 官方示例
		https://cn.vuejs.org/v2/guide/routing.html
 - vue-router
	- Vue Router 是 Vue.js (opens new window)官方的路由管理器
		- 它和 Vue.js 的核心深度集成，让构建单页面应用变得易如反掌
		- 包含的功能
			- 套的路由/视图表
			- 块化的、基于组件的路由配置
			- 由参数、查询、通配符
			- 于 Vue.js 过渡系统的视图过渡效果
			- 粒度的导航控制
			- 有自动激活的 CSS class 的链接
			- TML5 历史模式或 hash 模式，在 IE9 中自动降级
			- 定义的滚动条行为
	- 安装vue-router
		- 使用npm
			- $ npm install vue-router
			- vue-router是vue的插件，按照插件的方式引入到vue中
				import Vue from 'vue'
				import VueRouter from 'vue-router'
				Vue.use(VueRouter)
		- 使用vue-cil
			- $ vue add router
	- 基本使用
		- 添加一个页面: Test.vue
			<template>
			  <div class="about">
			    <h1>This is an test page</h1>
			  </div>
			</template>
		- 补充到路由 ./router/index.js
			const routes = [
			  {
			    path: '/test',
			    name: 'Test',
			    // route level code-splitting
			    // this generates a separate chunk (about.[hash].js) for this route
			    // which is lazy-loaded when the route is visited.
			    component: () => import(/* webpackChunkName: "about" */ '../views/Test.vue')
			  }
			  ]
		- 在界面上添加一个跳转
			<div id="nav">
			  <router-link to="/">Home</router-link> |
			  <router-link to="/about">About</router-link> |
			  <router-link to="/test">Test</router-link>
			</div>
	- 编程式的导航
		- router-link 组件是需要用户点击生效，如果需要动态加载，或者跳转前检查用户的权限，这个时候再使用router-link就不合适了
			- window.history 和 location 可以模拟用户操作浏览器
				location.assign('/')
				location.reload()
				history.back()
				history.forward()
			- vue-router 提供了一个函数用于js来控制路由那就是 push 功能和location.assign类似
				vm.router
				router.push(location, onComplete?, onAbort?)
				// location    location参数 等价于 <router-link :to="...">, 比如<router-link :to="/home">  等价于 router.push('/home')
				// onComplete  完成后的回调
				// onAbort     取消后的回调
			- 使用a标签
				<div id="nav">
				  <a @click="jumpToHome">Home</a> |
				  <a @click="jumpToAbout">About</a> |
				  <a @click="jumpToTest">Test</a>
				</div>
				<script>
				export default {
				  name: 'App',
				  data() {
				    return {
				    }
				  },
				  methods: {
				    jumpToHome() {
				      this.$router.push('/')
				    },
				    jumpToAbout() {
				      this.$router.push('/about')
				    },
				    jumpToTest() {
				 	  this.$router.push('/test')
				    }
				  },
				}
				</script>
	- 动态路由匹配
		- 前后端路由的区别
			后端:  path --->   handler
			前端:  path --->   view
		- 后端demo里面的http-router路由
			r.GET("/hosts", api.QueryHost)
			r.POST("/hosts", api.CreateHost)
			r.GET("/hosts/:id", api.DescribeHost)
			r.DELETE("/hosts/:id", api.DeleteHost)
			r.PUT("/hosts/:id", api.PutHost)
			r.PATCH("/hosts/:id", api.PatchHost)
		- 前端vue-router的路由
			模式                            匹配路径              $route.params
			/user/:username                 /user/evan            { username: 'evan' }
			/user/:username/post/:post_id   /user/evan/post/123   { username: 'evan', post_id: '123' }
		- 修改./router/index.js
			{
			  path: '/test/:id',
			  name: 'Test',
			  // route level code-splitting
			  // this generates a separate chunk (about.[hash].js) for this route
			  // which is lazy-loaded when the route is visited.
			  component: () => import(/* webpackChunkName: "about" */ '../views/Test.vue')
			}
	- 404处理
		- 如果找不页面，需要返回一个视图，告诉用户页面不存在
		- vue-router在处理404的方式和后端不同，路由依次匹配
			- 如果都匹配不上，写一个特殊的*路由作为 404路由
				{
				  path: '*',
				  name: '404',
				  // route level code-splitting
				  // this generates a separate chunk (about.[hash].js) for this route
				  // which is lazy-loaded when the route is visited.
				  component: () => import(/* webpackChunkName: "about" */ '../views/404.vue')
				}
	- 加载数据
		- 制作一个详情页面，根据不同的id，向后端获取不同的对象，用于显示
			- 请求id对应的后端数据，通过axios
				$ npm install --save axios  // axios@0.21.4
			- 需要将ajax封装下，添加一些通用逻辑，模块位于 utils/request.js
				import axios from 'axios'
				// create an axios instance
				const service = axios.create({
				    baseURL: 'http://localhost:8050', // url = base url + request url
				    // withCredentials: true, // send cookies when cross-domain requests
				    timeout: 5000 // request timeout
				})
				// request interceptor
				service.interceptors.request.use(
				    config => {
				      return config
				    },
				    error => {
				      // do something with request error
				      console.log(error) // for debug
				      return Promise.reject(error)
				    }
				)
				// response interceptor
				service.interceptors.response.use(
				    /**
				    * If you want to get http information such as headers or status
				    * Please return  response => response
				    */
				    
				    /**
				    * Determine the request status by custom code
				    * Here is just an example
				    * You can also judge the status by HTTP Status Code
				    */
				    response => {
				      const res = response.data
				      // if the custom code is not 20000, it is judged as an error.
				      if (res.code !== 0) {
				  	    // 比如 token过期
				      } else {
				  	    // 正常
				  	    return res
				      }
				    },
				    error => {
				      console.log('err' + error) // for debug
				      // 传递出去
				      return Promise.reject(error)
				    }
				)
				export default service
			- 新增一个api目录用于存放所有的API请求，在里面新建一个模块: test.js
				import request from '../utils/request'
				export function GET_TEST_DATA(id, query) {
				  return request({
				    url: `/hosts/${id}`,
				    method: 'get',
				    params: query
				  })
				}
			- 最后在视图中使用: Test.vue
				- 关于选择在什么时候加载数据，通常有2种方案
					- 导航完成之后获取: 先完成导航，然后在接下来的组件生命周期钩子中获取数据，在数据获取期间显示"加载中"之类的指示
					- 导航完成之前获取: 导航完成前，在路由进入的守卫中获取数据，在数据获取成功后执行导航
				- 方案一
					<script>
					import { GET_TEST_DATA } from '../api/test'
					export default {
					  name: 'Test',
					  data () {
					    return {
					      loading: false,
					      post: null,
					      error: null
					    }
					  },
					  created () {
					    // 组件创建完后获取数据，
					    // 此时 data 已经被 observed 了
					    this.fetchData()
					  },
					  watch: {
					    // 如果路由有变化，会再次执行该方法
					    '$route': 'fetchData'
					  },
					  methods: {
					    async fetchData () {
					      this.error = this.post = null
					      this.loading = true
					      // replace GET_TEST_DATA with your data fetching util / API wrapper
					      try {
					        this.loading = true
					        let resp = await GET_TEST_DATA(this.$route.params.id)
					        this.post = resp.data
					      } catch (err) {
					        this.error = err.toString()
					      } finally {
					        this.loading = false
					      }
						}
					  }
					}
					</script>
				- 方案二
					- 在router的钩子中获取数据
						beforeRouteEnter: 进入路由前
						beforeRouteUpdate: 路由update前
					- https://router.vuejs.org/zh/guide/advanced/data-fetching.html#%E5%9C%A8%E5%AF%BC%E8%88%AA%E5%AE%8C%E6%88%90%E5%89%8D%E8%8E%B7%E5%8F%96%E6%95%B0%E6%8D%AE
	- Router对象
		- VueRouter源码
			export declare class VueRouter {
			  constructor(options?: RouterOptions)
			  
			  app: Vue
			  options: RouterOptions
			  mode: RouterMode
			  currentRoute: Route
			  
			  beforeEach(guard: NavigationGuard): Function
			  beforeResolve(guard: NavigationGuard): Function
			  afterEach(hook: (to: Route, from: Route) => any): Function
			  push(location: RawLocation): Promise<Route>
			  replace(location: RawLocation): Promise<Route>
			  push(
			    location: RawLocation,
			    onComplete?: Function,
			    onAbort?: ErrorHandler
			  ): void
			  replace(
			    location: RawLocation,
			    onComplete?: Function,
			    onAbort?: ErrorHandler
			  ): void
			  go(n: number): void
			  back(): void
			  forward(): void
			  match (raw: RawLocation, current?: Route, redirectedFrom?: Location): Route
			  getMatchedComponents(to?: RawLocation | Route): Component[]
			  onReady(cb: Function, errorCb?: ErrorHandler): void
			  onError(cb: ErrorHandler): void
			  addRoutes(routes: RouteConfig[]): void
			  
			  addRoute(parent: string, route: RouteConfig): void
			  addRoute(route: RouteConfig): void
			  getRoutes(): RouteRecordPublic[]
			  
			  resolve(
			    to: RawLocation,
			    current?: Route,
			    append?: boolean
			  ): {
			  location: Location
			  route: Route
			  href: string
			  // backwards compat
			  normalizedTo: Location
			    resolved: Route
			  }
			}
	- Router钩子
		- 在路由前后做一些额外的处理，可以通过路由留得钩子
		- 最常见的使用钩子的地方是认证，在访问页面的时候，判断用户是否有权限访问
		- router为提供了如下钩子
			beforeEach: 路由前处理
			beforeEnter
			beforeRouteEnter
			beforeRouteUpdate
			beforeRouteLeave
			afterEach: 路由后出来
		- 示例 设置钩子函数，做一个简单的页面加载progress bar
			- 使用最广泛的函数 beforeEach和afterEach
				router.beforeEach((to, from, next) => {
				  console.log(to, from, next)
				  next()
				})
				router.afterEach((to, from) => {
				  console.log(to, from)
				})
			- 选用nprogress库 https://www.npmjs.com/package/nprogress
				- 安装
					// nprogress@0.2.0
					$ npm install --save nprogress
				- 使用
					NProgress.start();
					NProgress.done();
					NProgress.set(0.0);     // Sorta same as .start()
					NProgress.set(0.4);
					NProgress.set(1.0);     // Sorta same as .done()
			- 代码
				- 引入库和样式
					import NProgress from 'nprogress' // progress bar
					import 'nprogress/nprogress.css'  // progress bar style
					// 路由开始时: NProgress.start();
					// 路由结束时: NProgress.done();
				- 修改router
					router.beforeEach((to, from, next) => {
					  // start progress bar
					  NProgress.start()
					  console.log(to, from, next)
					  next()
					})
					router.afterEach(() => {
					  // finish progress bar
					  NProgress.done()
					})
				- 调整颜色   找到样式，写入一个文件: styles/index.css，进行全局加载
					#nprogress .bar {
					    background:#13C2C2;
					  }
				- 在main.js加载全局样式
					// 加载全局样式
					import './styles/index.css'

15、Vue状态管理
 - token
	- 每个页面在访问后端数据的时候 都需要token，所有组件都依赖token这个数据，所以就需要找个地方存起来，让其他组件都能访问到它
	- 即 存储一些公用的东西，提供给各个组件使用，和服务器端的session功能很类似
 - 共享内存
	- 直接开辟一个变量，使全局都可以直接访问，同定义一个全局map
	- 示例 在root 实例上添加一个data，其他子实例 通过root.data来访问数据
		- root
			// Root Vue实例
			new Vue({
			  render: h => h(App),
			  data: {a: 1},
			}).$mount('#app')
		- 在父节点添加一个b属性
			<script>
			import HelloWorld from './components/HelloWorld.vue'
			export default {
			  name: 'App',
			  created() {
			    this.$root.$data.b = 2
			  },
			}
			</script>
		- 子节点上读取
			<script>
			export default {
			  name: 'HelloWorld',
			  mounted() {
			    console.log(this.$root.$data.b)
			  },
			}
			</script>
 - 本地存储
	- 使用到浏览器的存储功能，其提供了存储客户端临时信息，简称 Web Storage
		- F12 -> Application -> Stroage
			- cookie
			- sessionStorage
			- localStorage
		- cookie
			- cookie可以设置过期时间，同一个域下的页面都可以访问
			- cookie在没有设置过期时间时，系统默认浏览器关闭时失效
				- 只有设置了没到期的保存日期时，浏览器才会把cookie作为文件保存在本地上
				- 当expire到期时，cookie不会自动删除，仅在下次启动浏览器或者刷新浏览器时，浏览器会检测cookie过期时间，如已过期浏览器则会删除过期cookie
			- 数据存放大小: 4k，因为每次http请求都会携带cookie
			- 浏览器关闭时，cookie会失效
			- 注意cookie可以支持httpOnly，这个时候前端js是无法查看也无法修改的
			- 示例语法
				document.cookie  // 读取cookie, 注意读取出来的cookie是个字符串
				document.cookie.split('; ')
				document.cookie = "username=John Doe; expires=Thu, 18 Dec 2043 12:00:00 GMT; path=/"; // 直接赋值就添加了一个key-value；修改cookie和设置cookie一样，保证key相同就可以
				document.cookie = `cookieKey=;expires=Mon, 26 Aug 2019 12:00:00 UTC` // 删除cookie时，把expires 设置到过期的时间即可，比如设置个2019年的时间
		- sessionStorage
			- 存储的数据只有在同一个会话中的页面才能访问并且当会话结束后数据也随之销毁
				- 因此sessionStorage不是一种持久化的本地存储，仅仅是会话级别的存储
			- 一个标签页 就表示一个回话，当标签页关闭，会话就被清除
			- 不通标签页之间不共享数据
			- 示例语法
				// 通过setItem设置key-value
				sessionStorage.setItem('key1', 'value1')
				sessionStorage['key2']= 'value2'
				sessionStorage.key2= 'value2'
				
				// 查询sessionStorage对象
				sessionStorage
				Storage {key2: 'value2', key1: 'value1', length: 2}
				
				// 通过getItem获取key的值
				sessionStorage.getItem('key1')
				sessionStorage['key1']
				sessionStorage.key1
				
				// 修改
				sessionStorage.key1 = 'value11'
				sessionStorage['key1'] = 'value11'
				
				// 删除key
				sessionStorage.removeItem('key1')
				
				// 清空storage
				sessionStorage.clear()
		- localStorage
			- localStorage生命周期是永久，除非主动删除数据，否则数据是永远不会过期的
			- 相同浏览器的不同页面间可以共享相同的 localStorage(页面属于相同域名和端口)
			- 示例语法
				// 通过setItem设置key-value
				localStorage.setItem('key1', 'value1')
				localStorage['key2']= 'value2'
				localStorage.key2 = 'value2'
				
				// 查询sessionStorage对象
				localStorage
				Storage {key2: 'value2', key1: 'value1', length: 2}
				
				// 通过getItem获取key的值
				localStorage.getItem('key1')
				localStorage['key1']
				localStorage.key1
				
				// 修改
				localStorage.key1 = 'value11'
				localStorage['key1'] = 'value11'
				
				// 删除key
				localStorage.removeItem('key1')
				
				// 清空storage
				localStorage.clear()
		- vuex
			- vuex是一种更高级的抽象
				- state，驱动应用的数据源
				- actions，响应在 view 上的用户输入导致的状态变化
				- mutation，用于直接修改数据的方法
			- 读取数据
				- 组件通过store的getter方法 从中取数据
			- 修改数据
				- 组件通过store提供的dispatch方法触发一个action，由action提交mutation来修改数据
			- 安装
				- 使用npm
					- $ npm install vuex --save
					- vue-router是vue的插件，按照插件的方式引入到vue中
						import Vue from 'vue'
						import Vuex from 'vuex'
						Vue.use(Vuex)
				- 使用vue-cil
					- $ vue add vuex
				- 基本使用
					-  示例  用户设置分页大小，希望每个页面都能生效，
						- 在store中定义一个状态: pageSize    // 文件 store/index.js 
							export default new Vuex.Store({
							  state: {
							    /* 添加pageSize状态变量 */
							    pageSize: 20
							  },
							  getters: {
							    /* 设置获取方法 */
							    pageSize: state => {
							      return state.pageSize
							    } 
							  },
							  mutations: {
							    /* 定义修改pageSize的函数 */
							    setPageSize(state, ps) {
							      state.pageSize = ps
							    }
							  },
							  actions: {
							    /* 一个动作可以由可以提交多个mutation */
							    /* { commit, state } 这个是一个解构赋值, 正在的参数是context, 我们从中解出我们需要的变量*/
							    setPageSize({ commit }, ps) {
							    /* 使用commit 提交修改操作 */
							      commit('setPageSize', ps)
							    }
							  },
							  modules: {
							  }
							})
						- 子组件建中修改状态，采用store提供的dispatch方法来进行修改: Test.vue
							<input v-model="pageSize" type="text">
							computed: {
							  pageSize: {
							    get() {
							      return this.$store.getters.pageSize
							    },
							    set(value) {
							      this.$store.dispatch('setPageSize', value)
							    }
							  }
							},
			- vuex-persist
				- vuex的状态存储并不能持久化，存储在 Vuex 中的 store 里的数据，只要一刷新页面，数据就会丢失
				- 安装vuex-persistvuex插件，将 Vuex 的存储修改为localstorage   // https://github.com/championswimmer/vuex-persist
					- $ npm install --save vuex-persist // vuex-persist@3.1.3
				- 基本使用 配置vuex使用该插件: store/index.js
					// 1. 引入依赖
					import VuexPersistence from 'vuex-persist'
					// 2. 实例化一个插件对象, 我们使用localStorage作为存储
					const vuexLocal = new VuexPersistence({
					  storage: window.localStorage
					})
					// 3.配置store实例使用localStorage插件
					export default new Vuex.Store({
					  // ...
					  plugins: [vuexLocal.plugin],
					})
				- 通过console进行确认
					- localStorage
					- 可以发现vuex-persist，使用vuex做一个key，把所有数据都存储在这个字段里面，所以还是不要使用vuex存储太多数据，不然有性能问题

16、参考
 - VUE2官方文档                                               https://cn.vuejs.org/v2/guide/
 - 那些前端MVVM框架是如何诞生的                               https://zhuanlan.zhihu.com/p/36453279
 - MVVM设计模式                                               https://zhuanlan.zhihu.com/p/36141662 
 - vue核心之虚拟DOM(vdom)                                     https://www.jianshu.com/p/af0b398602bc
 - Vue Router文档                                             https://next.router.vuejs.org/zh/introduction.html
 - Vuex 文档                                                  https://vuex.vuejs.org/zh/
 - cookies、sessionStorage和localStorage解释及区别            https://www.cnblogs.com/pengc/p/8714475.html
 - JavaScript Cookie                                          https://www.runoob.com/js/js-cookies.html
 - JavaScript创建、读取和删除cookie                           https://www.jb51.net/article/169117.htm

七、项目实例
1、前端框架搭建
 - 初始化项目
	- $ vue create devcloud
		Vue CLI v4.5.13
		? Please pick a preset: Manually select features
		? Check the features needed for your project: Choose Vue version, Babel, Router, Vuex, CSS Pre-processors, Linter
		? Choose a version of Vue.js that you want to start the project with 2.x
		? Use history mode for router? (Requires proper server setup for index fallback in production) Yes
		? Pick a CSS pre-processor (PostCSS, Autoprefixer and CSS Modules are supported by default): Sass/SCSS (with node-sass)
		? Pick a linter / formatter config: Prettier
		? Pick additional lint features: Lint on save
		? Where do you prefer placing config for Babel, ESLint, etc.? In dedicated config files
		? Save this as a preset for future projects? Yes
		? Save preset as: devcloud
	- 确认项目依赖
		"core-js": "^3.6.5",
		"vue": "^2.6.11",
		"vue-router": "^3.2.0",
		"vuex": "^3.4.0"
		
 - 项目npm配置  https://www.npmjs.cn/files/npmrc/
	- 使用 npm config edit 查看可以配置的选项
	- 后面需要用到node-sass，这个是二进制的css预处理器(编译器)，默认的url是国外下载地址(不在npm源中)，针对这种类型的依赖，需要单独配置下载的url
		- 使用sass_binary_site将它配置到npm的配置文件中
		- 本项目npm配置文件: .npmrc
			sass_binary_site=https://npm.taobao.org/mirrors/node-sass/
			registry=https://registry.npm.taobao.org

 - 项目vue配置
	- vue.config.js 是一个可选的配置文件，如果项目的 (和 package.json 同级的) 根目录中存在这个文件，那么它会被 @vue/cli-service 自动加载  
		- 官网的配置说明 https://cli.vuejs.org/zh/config/
	- 基础配置
		- 创建一个js的模块，用于导出配置给cli使用，这里采用commonJS模块导出语法
		- 在项目根下面创建: vue.config.js
		- https://github.com/ahwhy/myGolang/blob/main/week17/devcloud/vue.config.js
			// All configuration item explanations can be find in https://cli.vuejs.org/config/
			module.exports = {
			}
	- DevServer
	- webpack 基础配置
		- webpack是一个用于现代 JavaScript 应用程序的 静态模块打包工具
			- 当 webpack 处理应用程序时，它会在内部从一个或多个入口点构建一个 依赖图(dependency graph)
			- 然后你项目中所需的每一个模块组合成一个或多个 bundles，它们均为静态资源，用于展示具体的内容
			- webpack 中文文档 https://webpack.docschina.org/concepts/
	- webpack插件配置
		- Vue CLI 内部的 webpack 配置是通过 webpack-chain 维护的
			- 这个库提供了一个 webpack 原始配置的上层抽象，使其可以定义具名的 loader 规则和具名插件，并有机会在后期进入这些规则并对它们的选项进行修改
		- 静态资加载
			- 静态资源的加载对页面性能起着至关重要的作用，浏览器提供的两个资源指令-preload/prefetch，它们能够辅助浏览器优化资源加载的顺序和时机，提升页面性能 其中rel="prefetch"被称为Resource-Hints（资源提示），也就是辅助浏览器进行资源优化的指令 类似的指令还有rel="preload"
			- 预提取(prefetch)
				- 其利用浏览器空闲时间来下载或预取用户在不久的将来可能访问的文档，网页向浏览器提供一组预取提示，并在浏览器完成当前页面的加载后开始静默地拉取指定的文档并将其存储在缓存中
					- 当用户访问其中一个预取文档时，便可以快速的从浏览器缓存中得到
					- <link rel="prefetch" href="static/img/ticket_bg.a5bb7c33.png">
				- 一个Vue CLI应用会为所有作为async chunk生成的JavaScript文件(通过动态import()按需code splitting的产物)自动生成prefetch提示
					- 这些提示会被@vue/preload-webpack-plugin注入，并且可以通过chainWebpack的config.plugin('prefetch')进行修改和删除
					- when there are many pages, it will cause too many meaningless requests
						config.plugins.delete('prefetch')
			- 预加载preload)
				- 对于页面即刻需要的资源，可能希望在页面加载的生命周期的早期阶段就开始获取，在浏览器的主渲染机制介入前就进行预加载
					- 简单来说，就是通过标签显式声明一个高优先级资源，强制浏览器提前请求资源
					- <link rel="preload" href="xxx" as="xx">
				- 一个Vue CLI应用会为所有初始化渲染需要的文件自动生成preload提示，这些提示会被@vue/preload-webpack-plugin注入，并且可以通过chainWebpack的config.plugin('preload')进行修改和删除
				- it can improve the speed of the first screen, it is recommended to turn on preload
					config.plugin('preload').tap(() => [
					{
						rel: 'preload',
						// to ignore runtime.js
						// https://github.com/vuejs/vue-cli/blob/dev/packages/@vue/cli-service/lib/config/app.js#L171
						fileBlacklist: [/\.map$/, /hot-update\.js$/, /runtime\..*\.js$/],
						include: 'initial'
					}
					])
		- vue-loader
			- Vue Loader 是一个 webpack 的 loader，它允许以一种名为单文件组件(SFCs)的格式撰写 Vue 组件
			- 作用: 解析和转换 .vue 文件，提取出其中的逻辑代码 script、样式代码 style、以及 HTML 模版 template，再分别把它们交给对应的 Loader 去处理
			- Vue Loader 还提供了很多酷炫的特性
				- 允许为 Vue 组件的每个部分使用其它的 webpack loader
				- 允许在一个 .vue 文件中使用自定义块，并对其运用自定义的 loader 链
				- 使用 webpack loader 将 style 和 template 中引用的资源当作模块依赖来处理
				- 为每个组件模拟出 scoped CSS
				- 在开发过程中使用热重载来保持状态
			- 简而言之，webpack 和 Vue Loader 的结合提供了一个现代、灵活且极其强大的前端工作流，来帮助撰写 Vue.js 应用
				- Vue Loader 官方文档  https://vue-loader.vuejs.org/zh/
			- 一般 vue-loader 提取出template后，会调用vue-template-compiler来编译模版，所以vue-loader和vue-template-compiler经常一起安装
				- 一般使用vue cli时已经完成安装，通过npm list 查看当前项目安装的vue-loader和vue-template-compiler信息
					$ npm list  | grep vue-loader
					│ ├─┬ vue-loader@15.9.8
					│ ├─┬ vue-loader-v16@npm:vue-loader@16.8.1
					$ npm list  | grep vue-template-compiler
					├─┬ vue-template-compiler@2.6.14
				- 如果没有安装, 使用下面命令安装:
					npm install -D vue-loader vue-template-compiler
			- 通过chainWebpack来配置 包含vue的文件使用 vue-loader来处理:
				// set preserveWhitespace
				config.module
					.rule('vue')
					.use('vue-loader')
					.loader('vue-loader')
					.tap(options => {
					options.compilerOptions.preserveWhitespace = true
					return options
					})
					.end()
		- svg图标处理
			- 无论使用哪个UI组件，总会遇到icon不够用的时候，这时候为了保证icon放大不失真，需要使用svg icon
				- 最常用的svg icon库 阿里巴巴矢量图标库 https://www.iconfont.cn/search/index?searchType=icon&q=gitee&page=1&fromCollection=-1&fills=&tag=
			- 简单使用
				- 直接使用img标签，把资源放到静态文件的目录下: assets/feishu.svg
				- 在App.vue中通过相对路径使用 <img alt="Feishu logo" src="./assets/feishu.svg" />
			- 优化
				- 当icon很多的时候，由于使用的img标签，所以每次都需要从服务端拉去，可以使用下面两个库去优化
					- svg-sprite-loader: 会把 svg 塞到一个个 symbol 中，合成一个大的 svg，最后将这个大的 svg 放入 body 中，通过symbol id引用，symbol的id如果不特别指定，就是文件名
					- svgo-loader: 帮助svg文件进行瘦身的库
				- 步骤
					- 安装
						npm i --dev svg-sprite-loader svgo-loader
					- webpack配置使用vg-sprite-loader
						// set svg-sprite-loader
						// 设置svg相对路径: src/icons
						config.module
						.rule('svg')
						.exclude.add(resolve('src/icons'))
						.end()
						// svg结尾的文件使用svg-sprite-loader处理
						// 在svg-sprite-loader处理之前, 使用svgo-loader提取处理
						config.module
						.rule('icons')
						.test(/\.svg$/)
						.include.add(resolve('src/icons'))
						.end()
						.use('svg-sprite-loader')
						.loader('svg-sprite-loader')
						.options({
							symbolId: 'icon-[name]'
						})
						.end()
					- 引入svg图片
					- 引入到main.js中
					- 通过svg-sprite-loader使用
					- 封装Svg Icon组件
		- 打包优化
			- 在进行webpack打包的时候，为了避免某个js库文件太大，打包成单个文件加载过慢的问题，需要对大文件进行切割，让浏览器可以并行加载，提高页面加载速度

 - 引入UI组件
	- 安装element ui      官方安装文档  https://element.eleme.cn/#/zh-CN/component/installation
		$ npm i element-ui -S
		$ npm install --save js-cookie
	- 将element ui组件库和样式库 引入到 vue项目中，入口文件:main.js
 - 项目样式配置
 - VUE全局指令
	- 基于 clipboard 进行的封装, 关于clipboard的用法可以参考: clipboard Github  https://github.com/zenorocha/clipboard.js
 - VUE全局过滤器
 - Home页面
 - 参考
	- VUE CLI 全局配置                                         https://cli.vuejs.org/zh/config/#vue-config-js
	- webpack dev-server配置                                   https://webpack.js.org/configuration/dev-server/
	- 使用 Preload&Prefetch 优化前端页面的资源加载             https://zhuanlan.zhihu.com/p/273298222
	- svg-sprite-loader 使用教程                               https://www.jianshu.com/p/70f9c9268c83
	- JetBrains svg-sprite-loader                              https://github.com/JetBrains/svg-sprite-loader
	- 使用 svg-sprite-loader、svgo-loader 优化项目中的 Icon    https://juejin.cn/post/6854573215646875655
	- Vue Loader                                               https://vue-loader.vuejs.org/zh/
	- webpack 中文文档                                         https://webpack.docschina.org/concepts/

2、登陆页面
- Login组件
	- 使用一个elemnt 的From组件来实现这个登陆表单，用法参考: element form文档 https://element.eleme.cn/#/zh-CN/component/form
- 配置路由
- 全局样式调整
- 页面样式
- 调整输入框样式
- 添加svg icon
- 绑定数据
- 绑定数据
- 修复自动填充背景颜色
- 补充查看密码
- 默认聚焦于输入框
- 登陆表单校验
- 登陆逻辑
- 登陆状态
- 登陆守卫

3、404页面
- 简陋版 404页面
- 添加路由
- 正常404页面
	- vue-element-admin 404页面  https://github.com/PanJiaChen/vue-element-admin/blob/master/src/views/error-page/404.vue

4、项目导航页面
- Layout布局
- Layout样式
- 顶部导航
- 顶部导航样式
- 侧边栏导航
- 调整main页面样式
- 主机列表页面
CSS 变量教程 // https://www.ruanyifeng.com/blog/2017/05/css-variables.html





















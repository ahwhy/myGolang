# 前端框架搭建

这里我们直接使用vue cli 初始化我们的骨架, 然后在此基础上进行修改
```sh
vue create devcloud
Vue CLI v4.5.13
? Please pick a preset: Manually select features
? Check the features needed for your project: Choose Vue version, Babel, Router, Vuex, CSS Pre-processors, Linter
? Choose a version of Vue.js that you want to start the project with 2.x
? Use history mode for router? (Requires proper server setup for index fallback in production) Yes
? Pick a CSS pre-processor (PostCSS, Autoprefixer and CSS Modules are supported by default): Sass/SCSS (with node-sass)
? Pick a linter / formatter config: Prettier
? Pick additional lint features: Lint on save
? Where do you prefer placing config for Babel, ESLint, etc.? In package.json
? Save this as a preset for future projects? Yes
? Save preset as: devcloud
```
确认下我们的项目依赖是否和预期一样
```js
"core-js": "^3.6.5",
"vue": "^2.6.11",
"vue-router": "^3.2.0",
"vuex": "^3.4.0"
```

## 项目npm配置

项目中经常会遇到 有些依赖报拉取不下来, npm除了可以配置国内淘宝源，其实还有很多可以配置的选项: [The npm config files](https://www.npmjs.cn/files/npmrc/)

我们也可以通过npm config edit来查看有哪些选项可以配置
```
npm config edit
```

由于我们后面需要用到node-sass, 由于这个是个二进制的css预处理器(编译器), 而默认的url又是国外下载地址(不在npm源中), 针对这种类型的依赖，我们往往需要单独配置下载的url, 因此我们使用sass_binary_site将它配置到npm的配置文件中

本项目npm配置文件: .npmrc
```
sass_binary_site=https://npm.taobao.org/mirrors/node-sass/
registry=https://registry.npm.taobao.org
```


## 项目vue配置

vue.config.js 是一个可选的配置文件，如果项目的 (和 package.json 同级的) 根目录中存在这个文件，那么它会被 @vue/cli-service 自动加载, 关于具体的配置选项可以参考: [配置说明](https://cli.vuejs.org/zh/config/)

我们创建一个js的模块, 用于导出配置给cli使用, 这里采用commonJS模块导出语法, 关于模块相关语法还有不清楚的请查看之前课件: [JavaScript入口](../day15/javascript.md)

我们在项目根下面创建: vue.config.js
```js
// All configuration item explanations can be find in https://cli.vuejs.org/config/
module.exports = {
}
```

### 基础配置

```js
  /**
   * You will need to set publicPath if you plan to deploy your site under a sub path,
   * for example GitHub Pages. If you plan to deploy your site to https://foo.github.io/bar/,
   * then publicPath should be set to "/bar/".
   * In most cases please use '/' !!!
   * Detail: https://cli.vuejs.org/config/#publicpath
   */
  /*
    部署应用包时的基本 URL, '/'表示部署于跟路径,
    如果不放置于根, 比如: /my-app/, 访问路径就是： https://www.my-app.com/my-app/
  */
  publicPath: '/',
  /*
  当运行 vue-cli-service build 时生成的生产环境构建文件的目录
  */
  outputDir: 'dist',
  /*
  放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录
  */
  assetsDir: 'static',
  /*
  是否在开发环境下通过 eslint-loader 在每次保存时 lint 代码。这个值会在 @vue/cli-plugin-eslint 被安装之后生效,
  这里配置在开发环境生效

  注意: process.env 是当前进程的环境变量对象, 通过它可以访问到当前进程的所有环境变量
  */
  lintOnSave: process.env.NODE_ENV === 'development',
  /*
  是否需要生产环境的 source map, 可以将其设置为 false 以加速生产环境构建, 
  */
  productionSourceMap: false,
```

### DevServer

```js
    // If your port is set to 80,
    // use administrator privileges to execute the command line.
    // For example, Mac: sudo npm run
    // You can change the port by the following method:
    // port = 9527 npm run dev OR npm run dev --port = 9527
    const port = process.env.port || process.env.npm_config_port || 9527 // dev port

  /*
    所有 webpack-dev-server 的选项都支持, 具体见: https://webpack.js.org/configuration/dev-server/
  */
  devServer: {
    /* 通过环保变量获取当前开放服务器需要监听的端口*/
    port: port,
    /* 启动后浏览器默认打开的页面, 比如 open: ['/my-page', '/another-page'], true表示默认浏览器的publicPath */
    open: true,
    /* 当有编译报错时, 直接显示在页面上, 下面配置errors显示, warnings不显示 */
    overlay: {
      warnings: false,
      errors: true
    }
    /* 配置服务端代理, 开发时临时解决跨域问题 */
    // proxy: {
    //   '/workflow/api': {
    //     target: 'http://keyauth.nbtuan.vip',
    //     ws: true,
    //     secure: false,
    //     changeOrigin: true
    //   },
    // } 
    },
```


### webpack基础配置

webpack 是一个用于现代 JavaScript 应用程序的 静态模块打包工具。当 webpack 处理应用程序时，它会在内部从一个或多个入口点构建一个 依赖图(dependency graph)，然后将你项目中所需的每一个模块组合成一个或多个 bundles，它们均为静态资源，用于展示你的内容

更多信息请查看: [webpack 中文文档](https://webpack.docschina.org/concepts/)

```js
const name = defaultSettings.title || '极乐研发云' // page title
/*
webpack 相关配置
*/
configureWebpack: {
    // provide the app's title in webpack's name field, so that
    // it can be accessed in index.html to inject the correct title.
    name: name,
    // webpack 插件配置, 具体见: https://webpack.js.org/configuration/plugins/#plugins
    plugins: [],
    // 使用import的路由别名, 比如'@/components/Tips' 会别解析成: src/components/Tips
    // 更多resovle相关配置请查看: https://webpack.js.org/configuration/resolve/#resolve
    resolve: {
        alias: {
            '@': resolve('src')
        }
    }
},
```

### webpack插件配置

Vue CLI 内部的 webpack 配置是通过 webpack-chain 维护的。
这个库提供了一个 webpack 原始配置的上层抽象，使其可以定义具名的 loader 规则和具名插件，
并有机会在后期进入这些规则并对它们的选项进行修改

#### 静态资加载

静态资源的加载对页面性能起着至关重要的作用, 浏览器提供的两个资源指令-preload/prefetch，它们能够辅助浏览器优化资源加载的顺序和时机，提升页面性能
其中rel="prefetch"被称为Resource-Hints（资源提示），也就是辅助浏览器进行资源优化的指令
类似的指令还有rel="preload"

1.预提取预prefetch

prefetch通常翻译为预提取, 其利用浏览器空闲时间来下载或预取用户在不久的将来可能访问的文档,网页向浏览器提供一组预取提示，并在浏览器完成当前页面的加载后开始静默地拉取指定的文档并将其存储在缓存中。

当用户访问其中一个预取文档时，便可以快速的从浏览器缓存中得到, 比如: 
```html
<link rel="prefetch" href="static/img/ticket_bg.a5bb7c33.png">
```

一个Vue CLI应用会为所有作为async chunk生成的JavaScript文件(通过动态import()按需code splitting的产物)自动生成prefetch提示。

这些提示会被@vue/preload-webpack-plugin注入，并且可以通过chainWebpack的config.plugin('prefetch')进行修改和删除

when there are many pages, it will cause too many meaningless requests

```js
config.plugins.delete('prefetch')
```

2.加载(preload)

preload则翻译为预加载, 对于页面即刻需要的资源，你可能希望在页面加载的生命周期的早期阶段就开始获取，在浏览器的主渲染机制介入前就进行预加载简单来说，就是通过标签显式声明一个高优先级资源，强制浏览器提前请求资源
```html
<link rel="preload" href="xxx" as="xx">
```

一个Vue CLI应用会为所有初始化渲染需要的文件自动生成preload提示, 这些提示会被@vue/preload-webpack-plugin注入，并且可以通过chainWebpack的config.plugin('preload')进行修改和删除

it can improve the speed of the first screen, it is recommended to turn on preload

```js
config.plugin('preload').tap(() => [
{
    rel: 'preload',
    // to ignore runtime.js
    // https://github.com/vuejs/vue-cli/blob/dev/packages/@vue/cli-service/lib/config/app.js#L171
    fileBlacklist: [/\.map$/, /hot-update\.js$/, /runtime\..*\.js$/],
    include: 'initial'
}
])
```

#### vue-loader

Vue Loader 是一个 webpack 的 loader，它允许你以一种名为单文件组件 (SFCs)的格式撰写 Vue 组件

作用: 解析和转换 .vue 文件，提取出其中的逻辑代码 script、样式代码 style、以及 HTML 模版 template，再分别把它们交给对应的 Loader 去处理

Vue Loader 还提供了很多酷炫的特性：
+ 允许为 Vue 组件的每个部分使用其它的 webpack loader
+ 允许在一个 .vue 文件中使用自定义块，并对其运用自定义的 loader 链；
+ 使用 webpack loader 将 style 和 template 中引用的资源当作模块依赖来处理；
+ 为每个组件模拟出 scoped CSS；
+ 在开发过程中使用热重载来保持状态。

简而言之，webpack 和 Vue Loader 的结合为你提供了一个现代、灵活且极其强大的前端工作流，来帮助撰写 Vue.js 应用

关于 更多的vue loader信息请查看: [Vue Loader 官方文档](https://vue-loader.vuejs.org/zh/)

一般 vue-loader 提取出template后, 会调用vue-template-compiler来编译模版, 所以vue-loader和vue-template-compiler经常一起安装, 一般而已使用vue cli时已经安装好了, 我们通过npm list 查看当前项目安装的vue-loader和vue-template-compiler信息

```sh
$ npm list  | grep vue-loader
│ ├─┬ vue-loader@15.9.8
│ ├─┬ vue-loader-v16@npm:vue-loader@16.8.1
$ npm list  | grep vue-template-compiler
├─┬ vue-template-compiler@2.6.14
```


如果没有安装, 使用下面命令安装:
```sh
npm install -D vue-loader vue-template-compiler
```

我们通过chainWebpack来配置 包含vue的文件使用 vue-loader来处理:
```js
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
```

#### svg图标处理

无论我们使用那个UI组件, 总会遇到icon不够用的时候, 这时候为了保证icon放大不失真，我们需要使用svg icon, 最常用的svg icon库就是:
[阿里巴巴矢量图标库](https://www.iconfont.cn/search/index?searchType=icon&q=gitee&page=1&fromCollection=-1&fills=&tag=)

最简单的使用svg icon的方法是 直接使用img标签, 因此我们把资源放到我们的静态文件的目录下: assets/feishu.svg

然后在我们的App.vue中通过相对路径使用:
```html
<img alt="Feishu logo" src="./assets/feishu.svg" />
```

这看起来并没有什么不妥, 但是当我们icon很多的时候, 由于我们使用的img标签，所以每次都需要从服务端拉去, 有没有其他优化办法喃?

这里有2个库可以用来优化我们的导入
+ svg-sprite-loader: 会把你的 svg 塞到一个个 symbol 中，合成一个大的 svg, 最后将这个大的 svg 放入 body 中, 通过symbol id引用, symbol的id如果不特别指定，就是你的文件名
+ svgo-loader: 帮助svg文件进行瘦身的库



1. 首先我们需要安装这2个库
```sh
npm i --dev svg-sprite-loader svgo-loader
```

2. webpack配置使用vg-sprite-loader
```js
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
```

3. 引入svg图片

我们在src下新建icons文件夹，文件夹里再新建svg文件夹，将svg图片放至svg文件夹里

把刚才的feishu.svg放到 icons/svg文件下面

我们在icons下创建一个index.js用来加载这些svg文件, require(相当于import)自动化加载文件夹下所有模块: 参考[webpack中require.context的作用](https://zhuanlan.zhihu.com/p/59564277)
```js
const req = require.context('./svg', false, /\.svg$/)
const requireAll = requireContext => requireContext.keys().map(requireContext)
requireAll(req)
```

引入到main.js中
```js
import './icons'
```

4. 通过svg-sprite-loader使用

加载好了，我们就可以使用 svg标签和use标签来使用, 在App.Vue中加入:
```html
<svg>
    <use xlink:href="#icon-feishu"></use>
</svg>
```

![](./images/svg-icon.jpg)

这样用有点原始，我们把它封装成一个组件

5. 封装Svg Icon组件

为了校验svg是不是外部资源, 在utils/validate模块中定义isExternal函数:
```js
/**
 * @param {string} path
 * @returns {Boolean}
 */
 export function isExternal(path) {
    return /^(https?:|mailto:|tel:)/.test(path)
  }
```

我们在components下面新建一个组件: SvgIcon

```html
<template>
  <div v-if="isExternal" :style="styleExternalIcon" class="svg-external-icon svg-icon" v-on="$listeners" />
  <svg v-else :class="svgClass" aria-hidden="true" v-on="$listeners">
    <use :xlink:href="iconName" />
  </svg>
</template>

<script>
// doc: https://panjiachen.github.io/vue-element-admin-site/feature/component/svg-icon.html#usage
import { isExternal } from '@/utils/validate'
export default {
  name: 'SvgIcon',
  props: {
    iconClass: {
      type: String,
      required: true
    },
    className: {
      type: String,
      default: ''
    }
  },
  computed: {
    isExternal() {
      return isExternal(this.iconClass)
    },
    iconName() {
      return `#icon-${this.iconClass}`
    },
    svgClass() {
      if (this.className) {
        return 'svg-icon ' + this.className
      } else {
        return 'svg-icon'
      }
    },
    styleExternalIcon() {
      return {
        mask: `url(${this.iconClass}) no-repeat 50% 50%`,
        '-webkit-mask': `url(${this.iconClass}) no-repeat 50% 50%`
      }
    }
  }
}
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}
.svg-external-icon {
  background-color: currentColor;
  mask-size: cover!important;
  display: inline-block;
}
</style>
```

修改icons/index.js 注册SvgIcon为全局组建, 这样我们就可以在模版中直接使用 svg-icon组件了
```js
import Vue from 'vue'
import SvgIcon from '@/components/SvgIcon'// svg component

// register globally
Vue.component('svg-icon', SvgIcon)

const req = require.context('./svg', false, /\.svg$/)
const requireAll = requireContext => requireContext.keys().map(requireContext)
requireAll(req)
```

6. 在App.vue中以组件的方式使用svg

```html
<svg-icon icon-class="feishu" />
```

7. 添加svgo优化命令

在package.json中添加svgo指令, 用于优化我们的svg icon
```js
"svgo": "svgo -f src/icons/svg --config=src/icons/svgo.yml"
```

最后我们执行
```js
npm run svgo
```

#### 打包优化 

在进行webpack打包的时候, 为了避免某个js库文件太大, 打包成单个文件加载过慢的问题, 需要对大文件进行切割, 让浏览器可以并行加载，提高页面加载速度

```js
config
    .when(process.env.NODE_ENV !== 'development',
    config => {
        config
        .plugin('ScriptExtHtmlWebpackPlugin')
        .after('html')
        .use('script-ext-html-webpack-plugin', [{
        // `runtime` must same as runtimeChunk name. default is `runtime`
            inline: /runtime\..*\.js$/
        }])
        .end()
        config
        .optimization.splitChunks({
            chunks: 'all',
            cacheGroups: {
            libs: {
                name: 'chunk-libs',
                test: /[\\/]node_modules[\\/]/,
                priority: 10,
                chunks: 'initial' // only package third parties that are initially dependent
            },
            elementUI: {
                name: 'chunk-elementUI', // split elementUI into a single package
                priority: 20, // the weight needs to be larger than libs and app or it will be packaged into libs or app
                test: /[\\/]node_modules[\\/]_?element-ui(.*)/ // in order to adapt to cnpm
            },
            commons: {
                name: 'chunk-commons',
                test: resolve('src/components'), // can customize your rules
                minChunks: 3, //  minimum common number
                priority: 5,
                reuseExistingChunk: true
            }
            }
        })
        // https:// webpack.js.org/configuration/optimization/#optimizationruntimechunk
        config.optimization.runtimeChunk('single')
    }
    )
```

## 引入UI组件

安装element ui: [官方安装文档](https://element.eleme.cn/#/zh-CN/component/installation)

```sh
$ npm i element-ui -S
$ npm install --save js-cookie
```


将element ui组件库和样式库 引入到我们的vue项目中, 入口文件:main.js:
```js
import Cookies from 'js-cookie'
import Element from 'element-ui'
import "element-ui/lib/theme-chalk/index.css"

Vue.use(Element, {
  size: Cookies.get('size') || 'mini', // set element-ui default size
})
```

然后我们到[element官网](https://element.eleme.cn/), 找个组件验证下

```js
<template>
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png">
    <HelloWorld msg="Welcome to Your Vue.js App"/>
    <el-row>
      <el-button>默认按钮</el-button>
      <el-button type="primary">主要按钮</el-button>
      <el-button type="success">成功按钮</el-button>
      <el-button type="info">信息按钮</el-button>
      <el-button type="warning">警告按钮</el-button>
      <el-button type="danger">危险按钮</el-button>
    </el-row>
  </div>
</template>
```

## 项目样式配置

我们创建一个styles文件夹, 用于存放项目的样式文件


### element UI 主题定制
我们定制下element的主题色: element-variables.scss

```js
/* icon font path, required */
$--font-path: '~element-ui/lib/theme-chalk/fonts';

@import "~element-ui/packages/theme-chalk/src/index";
```

然后修改样式引入我们定制后的主题样式
```js
import Cookies from 'js-cookie'
import Element from 'element-ui'
import './styles/element-variables.scss'

Vue.use(Element, {
  size: Cookies.get('size') || 'mini', // set element-ui default size
})
```

### element UI组件样式定制
我们可以按需调整全局element ui的某些样式: element-ui.scss

比如我们把button的风格改为直角
```scss
.el-button {
    border-radius: 0px;
}
```

所有定制样式通过index.scss导出
```scss
@import './element-ui.scss';
```

## VUE全局指令

我们把全局的所有指令都放在于目录: directives

我们这里先添加1个全局指令: v-clipboard

1. v-clipboard

基于 clipboard 进行的封装, 关于clipboard的用法可以参考: [clipboard Github](https://github.com/zenorocha/clipboard.js)

我们先安装clipboard这个依赖
```sh
npm install --save clipboard
```

然后是我们的指令组件封装: directives/clipboard/clipboard.js
```js
// Inspired by https://github.com/Inndy/vue-clipboard2
const Clipboard = require('clipboard')
if (!Clipboard) {
  throw new Error('you should npm install `clipboard` --save at first ')
}

export default {
  bind(el, binding) {
    if (binding.arg === 'success') {
      el._v_clipboard_success = binding.value
    } else if (binding.arg === 'error') {
      el._v_clipboard_error = binding.value
    } else {
      const clipboard = new Clipboard(el, {
        text() { return binding.value },
        action() { return binding.arg === 'cut' ? 'cut' : 'copy' }
      })
      clipboard.on('success', e => {
        const callback = el._v_clipboard_success
        callback && callback(e) // eslint-disable-line
      })
      clipboard.on('error', e => {
        const callback = el._v_clipboard_error
        callback && callback(e) // eslint-disable-line
      })
      el._v_clipboard = clipboard
    }
  },
  update(el, binding) {
    if (binding.arg === 'success') {
      el._v_clipboard_success = binding.value
    } else if (binding.arg === 'error') {
      el._v_clipboard_error = binding.value
    } else {
      el._v_clipboard.text = function() { return binding.value }
      el._v_clipboard.action = function() { return binding.arg === 'cut' ? 'cut' : 'copy' }
    }
  },
  unbind(el, binding) {
    if (binding.arg === 'success') {
      delete el._v_clipboard_success
    } else if (binding.arg === 'error') {
      delete el._v_clipboard_error
    } else {
      el._v_clipboard.destroy()
      delete el._v_clipboard
    }
  }
}
```

2. 我们将这2个指令注册到全局, 我们在directives 目录下，定义一个index.js模块，用于注册我们所有全局指令
```js
import Vue from "vue";

import Clipboard from './clipboard/clipboard'

// 注册一个全局自定义指令
Vue.directive('clipboard', Clipboard)
```

3. 最后我们在main.js中 引入

```js
// 加载全局指令
import '@/directives' 
```

4. 测试这个全局指令

在 App.vue中引入一个输入框组件:
```html
<el-input v-model="input" placeholder="请输入内容"></el-input>
<el-button v-clipboard:copy="input" v-clipboard:success="clipboardSuccess" type="text" icon="el-icon-document-copy" style="padding:0px;margin-left:12px;" />
<script>
export default {
  name: 'App',
  data() {
    return {
      input: '',
    }
  },
  methods: {
    clipboardSuccess() {
      this.$message({
        message: '复制成功',
        type: 'success'
      })
    },
  }
}
</script>
```

> 之前不是做过一个v-focus的全局组件吗? 这里你为啥也把它注册成全局组件喃?

答案是: element的 input组件 对其进行了封装, 直接使用el.focus是不可以的, 因此此时的el不是 input元素, 我们需要通过筛选才可以:
```js
el.querySelector('input').focus()
```

同学们可以自己去实践下

## VUE全局过滤器

我们把全局的所有过滤器都放在于目录: filters

1. 我们创建一个: time.js, 定义parseTime过滤函数

```js
export function parseTime (value) {
    let date = new Date(value)
    return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
  }
```

2. 然后我们创建一个: index.js 用于注册全局过滤器
```js
import Vue from "vue";

import * as timeFilters from './time' 
// register global utility filters
Object.keys(timeFilters).forEach(key => {
  Vue.filter(key, timeFilters[key])
})
```

3. main.js中引入

```js
// 加载全局过滤器
import '@/filters' 
```

4. 测试下

```js
{{ ts | parseTime}}

data() {
  return {
    input: '',
    ts: Date.now()
  }
},
```

## Home页面

在做Home页面之前，先清理掉脚手架为我们生成的页面

### 清理脚手架

在做Login页面之前, 请清理脚手架给我生出的页面: 
+ views/About.vue
+ views/Home.vue

删除路由: router/index.js
+ Home
+ About

删除多于组件: components
+ HelloWorld.vue

清理App.vue里面多于的元素:
```html
<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>
```

### Home组件

由于home页面将来涉及到多个系统的数据展示, 所以独立一个目录来存放: views/dashboard

```html
<template>
  <div class="dashboard-container">
    Home 页面
  </div>
</template>

<script>
export default {
  name: 'Dashboard',
  data() {
    return {}
  }
}
</script>
```

### 补充路由

```js
const routes = [
  {
    path: '/',
    name: "Home",
    component: () => import('../views/dashboard/index'),
  },
];
```

我们是一个后台管理系统, 需要用户登陆后才能看到Home页面, 因此我们接下来先完成登陆页面

## 参考 

+ [VUE CLI 全局配置](https://cli.vuejs.org/zh/config/#vue-config-js)
+ [webpack dev-server配置](https://webpack.js.org/configuration/dev-server/)
+ [使用 Preload&Prefetch 优化前端页面的资源加载](https://zhuanlan.zhihu.com/p/273298222)
+ [svg-sprite-loader 使用教程](https://www.jianshu.com/p/70f9c9268c83)
+ [JetBrains svg-sprite-loader](https://github.com/JetBrains/svg-sprite-loader)
+ [使用 svg-sprite-loader、svgo-loader 优化项目中的 Icon](https://juejin.cn/post/6854573215646875655)
+ [Vue Loader](https://vue-loader.vuejs.org/zh/)
+ [webpack 中文文档](https://webpack.docschina.org/concepts/)
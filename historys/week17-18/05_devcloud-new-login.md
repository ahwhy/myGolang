# 新登录页面

![](./images/new-login.png)

有好事之徒嫌之前登录页面丑, 而我也没能力设计出一个漂亮的登录页面, 因此去B站找了个过来改造: [原项目地址](https://github.com/ramostear/login-page-01), 如果打不开 可以点 [这个地址](https://hub.fastgit.org/ramostear/login-page-01)

## 搬迁素材

把该项目下的素材图片(img目录下)搬迁到我们的: src/assets/login 下面:
+ avatar.svg
+ bg.png
+ img-3.svg

## 搬迁HTML

原项目下的HTML文件为: login.html, 我们将他搬迁到我们 login组件的模板里面来: src/views/keyauth/login/new.vue

有几处需要调整:
+ img标签的src 替换为 @/assets/login/xxx.png
+ 之前的英文文案换成中文
+ input组件使用 v-model 绑定我们之前的数据
+ 替换掉之前的btn按钮, 使用之前的el-button, 但是样式使用它的样式

```html
<template>
<div class="login-container">
    <img src="@/assets/login/bg.png" alt="" class="wave">
    <div class="container">
        <div class="img">
            <img src="@/assets/login/img-3.svg" alt="">
        </div>
        <div class="login-box">
            <form action="">
                <img src="@/assets/login/avatar.svg" alt="" class="avatar">
                <h2>极乐研发云</h2>
                <div class="input-group">
                    <div class="icon">
                        <i class="fa fa-user"></i>
                    </div>
                    <div>
                        <h5>账号</h5>
                        <input v-model="loginForm.username" type="text" class="input">
                    </div>
                </div>
                <div class="input-group">
                    <div class="icon">
                        <i class="fa fa-lock"></i>
                    </div>
                    <div>
                        <h5>密码</h5>
                        <input v-model="loginForm.password" type="password" class="input">
                    </div>
                </div>
                <a href="#">忘记密码</a>
                <!-- 提交表单 -->
                <!-- 这里替换成原来的el-button, 只是样式使用该项目的样式: class login-btn --> btn
                <el-button class="btn" :loading="loading" tabindex="3" size="medium" type="primary" @click="handleLogin">
                    登录
                </el-button>
            </form>
        </div>
    </div>
</div>
</template>
```

## 搬迁CSS

原来项目的css文件放置于: css/style.css, 我们将它搬迁到 组件目录下: src/views/keyauth/login/style.scss

它的样式我们不做调整直接使用: [style.css](https://hub.fastgit.org/ramostear/login-page-01/blob/master/css/style.css)

然后在我们组件内部引入

```html
<style lang="scss" scoped>
@import './style.scss';
</style>
```

调整我们的icon
```html
<span class="svg-container">
    <svg-icon icon-class="password" />
</span>
```

调整样式:
```css
.svg-container {
    padding-top: 11px;
    color: #d9d9d9;
    vertical-align: middle;
    display: inline-block;
}
```

## 搬迁Js

原项目的Js处理很简单, 核心就是添加focus和blur的事件处理:
```js
const inputs = document.querySelectorAll(".input");

function focusFunction(){
    let parentNode = this.parentNode.parentNode;
    parentNode.classList.add('focus');
}
function blurFunction(){
    let parentNode = this.parentNode.parentNode;
    if(this.value == ''){
        parentNode.classList.remove('focus');
    }
}

inputs.forEach(input=>{
    input.addEventListener('focus',focusFunction);
    input.addEventListener('blur',blurFunction);
});
```

我们将他搬迁到 new.vue里面的一个方法里面
```js
addEventHandler() {
    const inputs = document.querySelectorAll(".input");

    function focusFunction(){
        let parentNode = this.parentNode.parentNode;
        parentNode.classList.add('focus');
    }
    function blurFunction(){
        let parentNode = this.parentNode.parentNode;
        if(this.value == ''){
            parentNode.classList.remove('focus');
        }
    }

    inputs.forEach(input=>{
        input.addEventListener('focus',focusFunction);
        input.addEventListener('blur',blurFunction);
    });  
},
```

最后再页面加载完成后调用:
```js
mounted() {
    this.addEventHandler()
},
```

## 切换路由到新登录页面

修改我们路由: router/index.js, 让其指向新的登录视图
```js
{
path: "/login",
name: "Login",
component: () =>
    import("../views/keyauth/login/new.vue"),
},
```


到此我们基础的搬迁工作就完成了, 剩下适配我们的登录逻辑了

## 前端输入验证改造

之前的表单验证:
```js
this.$refs.loginForm.validate(async (valid) => {})
```

我们自己写一个check函数来完成验证, 如果没用输入用户名或者密码 我们就晃动下输入框

网上去找个晃动的样式 放到我们 login的css文件 style.scss中:
```css
.shake {
    animation: shake 800ms ease-in-out;
}
@keyframes shake {
    10%, 90% { transform: translate3d(-1px, 0, 0); }
    20%, 80% { transform: translate3d(+2px, 0, 0); }
    30%, 70% { transform: translate3d(-4px, 0, 0); }
    40%, 60% { transform: translate3d(+4px, 0, 0); }
    50% { transform: translate3d(-4px, 0, 0); }
}
```

然我们写一个动态添加晃动样式的函数和check的函数:
```js
shake(elemId) {
    let elem = document.getElementById(elemId)
    if (elem) {
        elem.classList.add('shake')
        setTimeout(()=>{ elem.classList.remove('shake') }, 800)
    }
},
check() {
    if (this.loginForm.username === '') {
        this.shake('username')
        return false
    }
    if (this.loginForm.password === '') {
        this.shake('password')
        return false
    }
    return true
},
```

我们为我们需要添加样式的元素添加上id:
```html
<!-- 账号输入框组 -->
<div id="username" class="input-group">
<!-- 密码输入框组 -->
<div id="password" class="input-group">
```

最后调整我们的验证函数: handleLogin
```js
async handleLogin() {
    if (this.check()) {
        // 省略...
    }
}
```

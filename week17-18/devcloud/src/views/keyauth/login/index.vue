<template>
  <div class="login-container">
    <el-form class="login-form" ref="loginForm" :model="loginForm" :rules="loginRules">
      <!-- 切换 -->
      <div>
        <el-tabs v-model="loginForm.grant_type">
        <el-tab-pane label="普通登陆" name="password" />
        <el-tab-pane label="LDAP登陆" name="ldap" />
        </el-tabs>
      </div>

      <!-- 账号输入框 
        key: 元素的key, vue做数据绑定时，更新数据的标识符, vm实例内需要唯一
        ref: 添加引用, 通过vm.$refs中 ref名字可以找到该元素
        name/autocomplete: 一起使用, 自动填充功能
        type: 输入框类型, text 文本框, password 密码框
        tabindex: 使用tab按键进行切换时的顺序控制 -->
      <el-form-item prop="username">
      <span class="svg-container">
        <svg-icon icon-class="user" />
      </span>
      <el-input key="username" placeholder="账号" ref="username" v-model="loginForm.username" name="username" type="text" tabindex="1" autocomplete="on" />
      </el-form-item>

      <!-- 密码输入框 -->
      <el-form-item prop="password">
      <span class="svg-container">
        <svg-icon icon-class="password" />
      </span>
      <el-input key="password" placeholder="密码" ref="password" v-model="loginForm.password" name="password" :type="passwordType" tabindex="2" autocomplete="on" />
      <span class="show-pwd" @click="showPwd">
        <svg-icon :icon-class="passwordType === 'password' ? 'eye-close' : 'eye-open'" />
      </span>
      </el-form-item>

      <!-- 登陆按钮 -->
      <el-button  class="login-btn" :loading="loading" size="medium" type="primary" tabindex="3" @click="handleLogin">
        登陆
      </el-button>
    </el-form>
  </div>
</template>

<script>
const validateUsername = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入账号'))
  } else {
    callback()
  }
}

const validatePassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    callback()
  }
}

export default{
  name: 'Login',
  data(){
    return{
      loading: false,
      passwordType: 'password',
      loginForm: {
          grant_type: 'password',
          username: '',
          password: ''
      },
      loginRules: {
      // required 是否必填
      // trigger  合适触发校验， change/blur
      // message  校验失败信息
        username: [{ trigger: 'blur', validator: validateUsername }],
        password: [{ trigger: 'blur', validator: validatePassword }]
      }
    }
  },  
  methods: {
    handleLogin() {
      this.$refs.loginForm.validate(async valid => {
        if (valid) {
          this.loading = true
          try {
            // 调用后端接口进行登录, 状态保存到vuex中
            await this.$store.dispatch('user/login', this.loginForm)

            // 调用后端接口获取用户profile, 状态保存到vuex中
            const user = await this.$store.dispatch('user/getInfo')
            console.log(user)
          } catch (err) {
            // 如果登陆异常, 中断登陆逻辑
            console.log(err)
            return
          } finally {
            this.loading = false
          }

          // 登陆成功, 重定向到Home或者用户指定的URL
          this.$router.push({ path: this.$route.query.redirect || '/', query: this.otherQuery })
        }
  })
    },
    showPwd(){
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
          this.passwordType = 'password'
      }
      this.$nextTick(() => {
          this.$refs.password.focus()
      })
    }
  },
  mounted() {
    this.$refs.username.focus()
  },
}
</script>

<style lang="scss" scoped>
/* 页面样式 */
.login-container{
  height: 100%;
  width: 100%;
  background-image: linear-gradient(to top, #3584A7 0%, #473B7B 100%);
  .login-form {
    width: 520px;
    padding: 160px 35px 0;
    margin: 0 auto;
    .login-btn {
      width:100%;
    }
  }
  .svg-container {
    padding: 6px 5px 6px 15px;
    color: #889aa4;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: #889aa4;
    cursor: pointer;
    user-select: none;
  }
}

/* 调整输入框样式 */
/* reset element-ui css */
.login-container ::v-deep .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;
    input {
      background: transparent;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      height: 47px;
      caret-color: #fff;
      color: #fff;
    }
  }

.login-container ::v-deep .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #454545;
}

.login-container ::v-deep .el-tabs__item {
  color: white;
  font-size: 18px;
}

.login-container ::v-deep .is-active {
  color:#13C2C2;
}
</style>
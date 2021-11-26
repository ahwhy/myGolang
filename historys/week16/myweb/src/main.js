import Vue from 'vue'
import App from './App.vue'
import Appsec from './Appsec.vue'
import router from './router'

// 加载全局样式
import './styles/index.css'

import store from './store'


Vue.config.productionTip = false

new Vue({
  router,
  data: {
    currentRoute: window.location.pathname
  },
  store,
  render(h) {
    if (this.currentRoute === '/sec') {
      return h(Appsec)
    }
    return h(App)
  },
}).$mount('#app')

// 注册一个全局自定义指令 `v-focus`
Vue.directive('focus', {
  // 当被绑定的元素插入到 DOM 中时……
  inserted: function (el) {
    // 聚焦元素
    el.focus()
  }
})
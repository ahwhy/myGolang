import Vue from 'vue'
import App from './App.vue'
import router from './router'

Vue.config.productionTip = false

// 加载全局样式
import './styles/index.css'


Vue.filter('parseTime', function (value) {
  let date = new Date(value)
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
})

Vue.directive('focus', {
  // 当被绑定的元素插入到 DOM 中时……
  inserted: function (el) {
    // 聚焦元素
    el.focus()
  }
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')

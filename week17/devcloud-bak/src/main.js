import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import "./icons"; // icon

import Cookies from 'js-cookie'
import Element from 'element-ui'
import './styles/element-variables.scss'

// 引入全局样式
import '@/styles/index.scss' // global css

// 加载全局指令
import '@/directives' 

// 加载全局过滤器
import '@/filters' 


Vue.use(Element, {
  size: Cookies.get('size') || 'mini', // set element-ui default size
})

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");

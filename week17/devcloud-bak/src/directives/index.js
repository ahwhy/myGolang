import Vue from "vue";

import Clipboard from './clipboard/clipboard'

// 注册一个全局自定义指令
Vue.directive('clipboard', Clipboard)
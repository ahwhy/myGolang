import Vue from "vue";

import * as timeFilters from './time' 
// register global utility filters
Object.keys(timeFilters).forEach(key => {
  Vue.filter(key, timeFilters[key])
})
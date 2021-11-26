import Vue from 'vue'
import Vuex from 'vuex'
// 1. 引入依赖 vuex-persist@3.1.3
import VuexPersistence from 'vuex-persist'

Vue.use(Vuex)

// 2. 实例化一个插件对象, 使用localStorage作为存储 vuex-persist@3.1.3
const vuexLocal = new VuexPersistence({
  storage: window.localStorage
})

export default new Vuex.Store({
  state: {
    /* 添加pageSize状态变量 $vm0.pageSize  */ 
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
  },
  // 3.配置store实例使用localStorage插件 vuex-persist@3.1.3
  plugins: [vuexLocal.plugin],
})

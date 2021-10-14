import Vue from "vue";
import Vuex from "vuex";
import VuexPersistence from 'vuex-persist'
import user from './modules/user'
import getters from './getters'

const vuexLocal = new VuexPersistence({
  storage: window.localStorage
})

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {user: user},
  getters,
  plugins: [vuexLocal.plugin]
});
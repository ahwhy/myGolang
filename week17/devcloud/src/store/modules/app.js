const state = {
    sidebar: {
      opened: true,
    },
    size: 'medium',
    system: 'dashboard'
  }
  
  const mutations = {
    TOGGLE_SIDEBAR: state => {
      state.sidebar.opened = !state.sidebar.opened
    },
    CLOSE_SIDEBAR: (state) => {
      state.sidebar.opened = false
    },
    SET_SIZE: (state, size) => {
      state.size = size
    },
    SET_SYSTEM: (state, system) => {
      state.system = system
    }
  }
  
  const actions = {
    toggleSideBar({ commit }) {
      commit('TOGGLE_SIDEBAR')
    },
    closeSideBar({ commit }, { withoutAnimation }) {
      commit('CLOSE_SIDEBAR', withoutAnimation)
    },
    setSize({ commit }, size) {
      commit('SET_SIZE', size)
    },
    setSystem({ commit }, system) {
      commit('SET_SYSTEM', system)
    }
  }

  export default {
    namespaced: true,
    state,
    mutations,
    actions
  }
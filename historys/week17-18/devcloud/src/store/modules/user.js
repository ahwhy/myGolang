import { LOGIN, GET_PROFILE } from '@/api/keyauth/token'

const state = {
    accessToken: '',
    namespace: '',
    account: '',
    type: '',
    name: '',
    avatar: '',
}

const mutations = {
    SET_TOKEN: (state, token) => {
        state.accessToken = token.access_token
        state.namespace = token.namespace
    },
    CLEAN_TOKEN: (state) => {
        state.accessToken = ''
    },
    SET_PROFILE: (state, user) => {
        state.type = user.type
        state.account = user.account
        state.name = user.profile.real_name
        state.avatar = user.profile.avatar
    },
}

const actions = {
    // 用户登陆接口
    login({ commit }, loginForm) {
        return new Promise((resolve, reject) => {
            const resp = LOGIN(loginForm)
            commit('SET_TOKEN', resp.data)
            resolve(resp)
        })
    },
    // 退出登录
    logout({ commit }) {
        return new Promise((resolve, reject) => {
            commit('CLEAN_TOKEN')
            resolve()
        })
    },
    // 获取用户Profile
    getInfo({ commit }) {
        return new Promise((resolve, reject) => {
            const resp = GET_PROFILE()
            commit('SET_PROFILE', resp.data)
            resolve(resp)
        })
    },
}

export default {
    // 这个是独立模块, 每个模块独立为一个namespace
    namespaced: true,
    state,
    mutations,
    actions,
}
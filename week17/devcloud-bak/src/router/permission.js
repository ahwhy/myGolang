import store from '@/store'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style

NProgress.configure({ showSpinner: false }) // NProgress Configuration

// 不需认证的页面
const whiteList = ['/login']

// 路由前钩子, 权限检查
export function beforeEach(to, from, next) {
    // 路由开始
    NProgress.start()

    // 取出token
    const hasToken = store.getters.accessToken

    // 判断用户是否登陆
    if (hasToken) {
        // 已经登陆得用户, 再次访问登陆页面, 直接跳转到Home页面
        if (to.path === '/login') {
            next({ path: '/' })
            NProgress.done()
        } else {
            next()
        }
    } else {
        // 未登录的用户, 重定向到登录页面
        if (whiteList.indexOf(to.path) !== -1) {
            // in the free login whitelist, go directly
            next()
        } else {
            // other pages that do not have permission to access are redirected to the login page.
            next(`/login?redirect=${to.path}`)
            NProgress.done()
        }
    }
}

// 路由后构造
export function afterEach() {
    // 路由完成
    NProgress.done()
}
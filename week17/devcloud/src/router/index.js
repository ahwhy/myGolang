import Vue from "vue";
import VueRouter from "vue-router";
import {beforeEach, afterEach} from './permission'

/* Layout */
import Layout from '@/layout'


Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        name: 'Dashboard',
      }
    ]
  },
  {
    path: "/cmdb",
    component: Layout,
    redirect: '/cmdb/host',
    children: [
      {
        path: 'host',
        component: () => import('@/views/cmdb/host/index.vue'),
        name: 'Cmdb-Host',
      }
    ]
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/keyauth/login/index.vue"),
  },

  {
    path: '/404',
    component: () => import('@/views/common/error-page/404.vue'),
    hidden: true
  },
  // 如果前面所有路径都没有匹配到页面 就跳转到404页面
  { path: '*', redirect: '/404', hidden: true }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach(beforeEach)
router.afterEach(afterEach)

export default router;

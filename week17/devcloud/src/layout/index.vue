<template>
  <div>
    <!-- 顶部导航栏 -->
    <div class="navbar-container">
        <navbar />
    </div>

    <!-- 主内容区 -->
    <div class="app-wrapper">
        <!-- 侧边栏导航 -->
        <div class="sidebar-container">
            <sidebar />
        </div>
        <!-- 内容页面区 -->
        <div class="main-container">
            <app-main />
        </div>
    </div>
  </div>
</template>

<script>
import { AppMain, Navbar, Sidebar } from './components'

export default {
  name: 'Layout',
  components: {
    AppMain,
    Navbar,
    Sidebar
  },
}
</script>
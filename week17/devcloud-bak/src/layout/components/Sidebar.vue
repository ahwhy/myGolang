<template>
  <div class="sidebar" :style="{'--sidebar-width': sidebarWidth}">
    <hamburger id="hamburger-container" :is-active="sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu default-active="1-4-1" class="sidebar-el-menu" :collapse="isCollapse" router>
        <el-submenu index="/host">
          <!-- 添加个title -->
          <template slot="title">
            <i class="el-icon-location"></i>
            <span slot="title">基础资源</span>
          </template>
          <!-- 导航条目 -->
          <el-menu-item index="/cmdb/search">资源检索</el-menu-item>
          <el-menu-item index="/cmdb/host" >主机</el-menu-item>
        </el-submenu>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Hamburger from '@/components/Hamburger'

export default {
  name: 'Sidebar',
  components: { Hamburger },
  computed: {
    // 使用对象展开运算符将 getter 混入 computed 对象中
    ...mapGetters([
      'sidebar'
    ]),
    isCollapse() {
      return !this.sidebar.opened
    }
  },
  data() {
    return {
      sidebarWidth: '210px',
    }
  },
  mounted() {
    if (!this.sidebar.opened) {
      this.sidebarWidth='65px'
    } else {
      this.sidebarWidth='210px'
    }
  },
  methods: {
    toggleSideBar() {
      if (!this.isCollapse) {
        this.sidebarWidth='65px'
      } else {
        this.sidebarWidth='210px'
      }
      this.$store.dispatch('app/toggleSideBar')
    },
  }
}
</script>

<style lang="scss" scoped>
.sidebar{
  width: var(--sidebar-width);
}

.sidebar-el-menu {
  height: calc(100vh - 70px);
  
}
</style>
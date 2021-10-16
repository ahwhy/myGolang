<template>
  <div class="sidebar" :style="{'width': sidebarWidth}">
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
export default {
  name: 'Sidebar',
  data() {
    return {
      sidebarWidth: '',
    }
  },
  watch: {
    isCollapse: {
      handler(newV) {
        if (newV) {
          this.sidebarWidth = '65px'
        } else {
          this.sidebarWidth = '210px'
        }
      },
      immediate: true
    }
  },
  computed: {
    isCollapse() {
      return this.$store.getters.sidebar.opened
    }
  },
}
</script>

<style lang="scss" scoped>
.sidebar-el-menu {
  height: calc(100vh - 50px);
}
</style>
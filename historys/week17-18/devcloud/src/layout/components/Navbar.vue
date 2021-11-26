<template>
  <div class="navbar">
    <!-- logo -->
    <div class="logo-container">
      <hamburger id="hamburger-container" :is-active="isCollapse" class="hamburger-container" @toggleClick="toggleSideBar" />
      <span>极乐研发云</span>
    </div>
    <!-- 主导航栏 -->
    <div class="navbar-main">
      <span class="navbar-item" @click="changeSystem('dashboard')" :class="{ active: activeIndex === 'dashboard' }">首页</span>
      <span class="navbar-item" @click="changeSystem('product')" :class="{ active: activeIndex === 'product' }">产品运营</span>
      <span class="navbar-item" @click="changeSystem('cmdb')" :class="{ active: activeIndex === 'cmdb' }">资源管理</span>
      <span class="navbar-item" @click="changeSystem('workflow')" :class="{ active: activeIndex === 'workflow' }">研发交付</span>
      <span class="navbar-item" @click="changeSystem('monitor')" :class="{ active: activeIndex === 'monitor' }">监控告警</span>
    </div>
    <!-- 用户信息区 -->
    <div class="navbar-user">
      <el-dropdown>
        <el-button type="text" style="color:white">
          <span>Atlantis</span>
          <i class="el-icon-arrow-down el-icon--right dropdown-color" />
        </el-button>
        <el-dropdown-menu slot="dropdown">
          <!-- 个人信息 -->
          <el-dropdown-item>
            <svg-icon icon-class="person" />
            <span class="dropdown-item-text">基本信息</span>
          </el-dropdown-item>
          <!-- 项目地址 -->
          <el-dropdown-item>
            <svg-icon icon-class="github" />
            <span class="dropdown-item-text">Github地址</span>
          </el-dropdown-item>
          <!-- 退出系统 -->
          <el-dropdown-item @click.native="logout" divided>
            <span style="display:block;">退出登录</span>
          </el-dropdown-item>       
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import Hamburger from '@/components/Hamburger'

export default {
  name: 'Navbar',
  components: { Hamburger },
  data() {
    return {
      activeIndex: 'dashboard',
      // isCollapse: false
    }
  },
  computed: {
    isCollapse() {
      return this.$store.getters.sidebar.opened
    },
    activeSystem() {
      return this.$store.getters.system
    },
  },
  methods: {
    changeSystem(index) {
      this.activeIndex = index
      // this.$store.dispatch('app/setSystem', system)
    },
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push({path: '/login'})
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  display: flex;
  width: 100%;
  align-content: center;
  justify-content: space-between;
  align-items: center;
  padding: 0 12px;
}

.navbar-main {
  display: flex;
  color: rgba(255, 255, 255, 0.8);
  font-size: 12px;
  .navbar-item {
    width: 68px;
    display: flex;
    justify-content: center;
    cursor: pointer;
    line-height: 32px;
    margin: 0 3px;
  }

  .navbar-item:hover {
    background-color: rgba(255, 255, 255, 0.3);
    border-radius: 4px;
  }

  .active {
    background-color: rgba(255, 255, 255, 0.3);
    border-radius: 4px;    
  }
  
}

.logo-container {
  display: flex;
  align-items: center;
  box-sizing: border-box;
  width: 200px;
  color: rgba(255, 255, 255, 0.8);

  .title {
    font-size: 16px;
  }
}

.navbar-user {
  margin-left: auto;
}
</style>
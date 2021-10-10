<template>
  <div class="navbar">
    <!-- logo -->
    <div class="logo-container">
      <span>极乐研发云</span>
    </div>
    <!-- 主导航栏 -->
    <div class="navbar-main">
      <el-menu :default-active="activeIndex" mode="horizontal" menu-trigger="click">
        <el-menu-item index="dashboard">首页</el-menu-item>
        <el-menu-item index="product">产品运营</el-menu-item>
        <el-menu-item index="resource">资源管理</el-menu-item>
        <el-menu-item index="develop">研发交付</el-menu-item>
        <el-menu-item index="eventbox">监控告警</el-menu-item>
      </el-menu>
    </div>
    <!-- 用户信息区 -->
    <div class="navbar-user">
      <el-dropdown>
        <el-button type="text">
          <span>老喻</span>
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
          <el-dropdown-item divided>
            <span style="display:block;">退出登录</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      activeIndex: 'dashboard'
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  display: flex;
  align-content: center;
  justify-content: flex-start;
  align-items: center;
}

.logo-container {
  box-sizing: border-box;
  width: 210px;
  padding-left: 28px;
  color: #909398;
}

.navbar ::v-deep .el-menu-item {
  height: 50px;
  line-height: 50px;
}

.navbar-user {
  margin-left: auto;
  padding-right: 12px;
}

.dropdown-item-text {
  margin-left: 12px;
  margin-right: 12px;
}
</style>
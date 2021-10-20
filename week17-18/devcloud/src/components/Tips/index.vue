<template>
  <div v-if="!hidden" class="tips">
    <div class="titile">
        <div class="tips-icon">
            <svg-icon icon-class="tips-info"></svg-icon>
        </div>
        <span class="title-content">
            {{ title }}
        </span>
        <span class="close-btn">
          <i class="el-icon-close" @click="handleClose" />
        </span>
    </div>
    <div class="content">
        <li v-for="(item, i) in tips" :key="i" class="tip-item"> {{item}} </li>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Tips',
  props: {
    title: {
      type: String,
      default: '温馨提示',
    },
    tips: {
      type: Array,
      default() {
        return []
      }
    }
  },
  data() {
    return {
      hidden: false,
    }
  },
  methods: {
    handleClose() {
      this.hidden = true
    }
  }
}
</script>

<style lang="scss" scoped>
.tips {
  width: 100%;
  background-color: rgba(48, 210, 190, 0.42);
  color: rgb(20, 105, 105);
  font-size: 12px;
  padding: 8px 16px;   // 左右16px 上下8px
}

.titile {
  display: flex;

  .tips-icon {
    font-size: 16px;
  }

  .title-content {
    margin-left: 12px;
    font-size: 13px;
    font-weight: 600;
  }

  .close-btn {
    margin-left: auto;
    cursor: pointer;
  }
}

.content {
  font-size: 12px;
  padding: 6px 26px 0px 26px;
}
</style>
<template>
  <div class="main-container">
    <tips :tips="tips" />
    CMDB Host
  </div>
</template>

<script>
import Tips from '@/components/Tips'

const tips = [
  '现在仅同步了阿里云主机资产'
]

export default {
  name: 'Host',
  components: { Tips },
  data() {
    return {
      tips: tips
    }
  }
}
</script>
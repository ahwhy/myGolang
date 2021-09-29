<template>
  <div class="detail">
    <h1>This is an detail page</h1>   
    <span>{{ $route.params }}</span>
  </div>
</template>

<script>
import { GET_TEST_DATA } from '../api/detail'

export default {
  name: 'Detail',
  data () {
    return {
      loading: false,
      post: null,
      error: null
    }
  },
  created () {
    // 组件创建完后获取数据，
    // 此时 data 已经被 observed
    this.fetchData()
  },
  mounted() {
    console.log(this.$root.$data.currentRoute)
  },
  watch: {
    // 如果路由有变化，会再次执行该方法
    '$route': 'fetchData'
  },
  methods: {
      async fetchData () {
        this.error = this.post = null
        this.loading = true
        // replace GET_TEST_DATA with your data fetching util / API wrapper
        try {
          this.loading = true
          let resp = await GET_TEST_DATA(this.$route.params.id)
          this.post = resp.data
        } catch (err) {
          this.error = err.toString()
        } finally {
          this.loading = false
        }
      }
    }
}
</script>
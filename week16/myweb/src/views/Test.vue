<template>
  <div class="test">
    <h1>This is an test page</h1>
    <input v-model="pageSize" type="text">
  </div>
</template>

<script>
export default {
  name: 'Test',
  computed: {
    pageSize: {
      get() {
        return this.$store.getters.pageSize
      },
      set(value) {
        this.$store.dispatch('setPageSize', value)
      }
    }
  },
}
</script>
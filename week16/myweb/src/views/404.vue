<template>
  <div class="404">
    <h1>This is an 404 page</h1>   
  </div>
</template>

<script>
export default {
  name: '404',
  mounted() {
    console.log(this.$root.$data.currentRoute)
  },
}
</script>
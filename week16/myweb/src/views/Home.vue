<template>
  <div id="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <!--  <- <HelloWorld :msg="message"/>  --> 
    <hello-world v-model="message" :msg="message" />
  </div>
</template>

<script>
import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'Home',
  data() {
    return {
      message: "Welcome to Your Vuesec.js App",
      urlHash: '',
    }
  },
  components: {
    HelloWorld
  },
  mounted() {
    let that = this
    window.onhashchange = function () {
      that.urlHash = window.location.hash
    };
    console.log(this.$root.$data.currentRoute)
  },
  watch: {
    urlHash: function(newURL, oldURL) {
      console.log(newURL, oldURL)
    }
  }
}
</script>

<style>
#home {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>

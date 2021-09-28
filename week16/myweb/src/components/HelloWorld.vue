<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <h2>{{ reverseName }}</h2>
    <input :value="value" type="text" @input="$emit('input', $event.target.value)">
    <input v-focus v-model="name" type="text" @keyup.enter="pressEnter(name)">
    <button :disabled="isButtomDisable" @click="clickButtom">你点呀</button>
    <br>
    <div v-if="name >= 90">
      A
    </div>
    <div v-else-if="name >= 80">
      B
    </div>
    <div v-else-if="name >= 60">
      C
    </div>
    <div v-else-if="name >= 0">
      D
    </div>
    <div v-else>
      请输入正确的分数
    </div>
    <ul>
      <li v-for="(item, index) in items" :key="item.message">
        {{ item.message }} - {{ index}}
        <br>
        <span v-for="(value, key) in item" :key="key"> {{ value }} {{ key }} <br></span>
      </li>
    </ul>
    <br>
    <p>
      For a guide and recipes on how to configure / customize this project,<br>
      check out the
      <a href="https://cli.vuejs.org" target="_blank" rel="noopener">vue-cli documentation</a>.
    </p>
    <h3>Installed CLI Plugins</h3>
    <ul>
      <li><a href="https://github.com/vuejs/vue-cli/tree/dev/packages/%40vue/cli-plugin-babel" target="_blank" rel="noopener">babel</a></li>
      <li><a href="https://github.com/vuejs/vue-cli/tree/dev/packages/%40vue/cli-plugin-eslint" target="_blank" rel="noopener">eslint</a></li>
    </ul>
    <h3>Essential Links</h3>
    <ul>
      <li><a href="https://vuejs.org" target="_blank" rel="noopener">Core Docs</a></li>
      <li><a href="https://forum.vuejs.org" target="_blank" rel="noopener">Forum</a></li>
      <li><a href="https://chat.vuejs.org" target="_blank" rel="noopener">Community Chat</a></li>
      <li><a href="https://twitter.com/vuejs" target="_blank" rel="noopener">Twitter</a></li>
      <li><a href="https://news.vuejs.org" target="_blank" rel="noopener">News</a></li>
    </ul>
    <h3>Ecosystem</h3>
    <ul>
      <li><a href="https://router.vuejs.org" target="_blank" rel="noopener">vue-router</a></li>
      <li><a href="https://vuex.vuejs.org" target="_blank" rel="noopener">vuex</a></li>
      <li><a href="https://github.com/vuejs/vue-devtools#vue-devtools" target="_blank" rel="noopener">vue-devtools</a></li>
      <li><a href="https://vue-loader.vuejs.org" target="_blank" rel="noopener">vue-loader</a></li>
      <li><a href="https://github.com/vuejs/awesome-vue" target="_blank" rel="noopener">awesome-vue</a></li>
    </ul>
    <br>
      {{ urlHash }}
    <br>
    <div>
     {{ ts | parseTime }}
    </div>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  data() {
    return {
      value: "Holle World.",
      name: 'This is first Web.',
      isButtomDisabled: false,
      items:[
        { message: 'Foo' },
        { message: 'Bar' }
      ],
      urlHash: '',
      ts: Date.now(),
    }
  },
  methods: {
    pressEnter() {
      alert("别敲回车键呀")
    },
    clickButtom() {
      alert("别点我")
    }
  },
  directives: {
    focus: {
      // 指令的定义
      inserted: function (el) {
        el.focus()
      }
    }
  },
  beforeCreate() {
    console.log('beforeCreate')
  },
  created() {
    console.log('created')
  },
  beforeMount() {
    console.log('beforeMount')
  },
  mounted() {
    console.log('mounted')
    let that = this
    window.onhashchange = function () {
      that.urlHash = window.location.hash
    };
  },
  beforeUpdate() {
    console.log('beforeUpdate')
  },
  updated() {
    console.log('updated')
  },
  beforeDestroy() {
    console.log('beforeDestroy')
  },
  destroyed() {
    console.log('destroyed')
  },
  computed: {
    reverseName: {
      get() {
        return this.name.split('').reverse().join('')
      },
      set(value){
        this.name = this.name = value.split('').reverse().join('')
      }
    }
  },
  props: {
    msg: String
  },
  watch: {
    urlHash: function(newURL, oldURL) {
      console.log(newURL, oldURL)
    }
  },
  filters: {
    parseTime: function (value) {
      let date = new Date(value)
      return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 20px;
}
a {
  color: #fd0000;
}
</style>

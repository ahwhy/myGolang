<template>
  <div class="wscn-http404-container">
    <p>对不起，您请求的页面不存在、或已被删除、或暂时不可用</p>
    <p>请点击以下链接继续浏览网页</p>
    <p> >> <a style="cursor:pointer" @click="$router.go(-1)">返回上一页面</a></p>
    <p> >> <a style="cursor:pointer" @click="$router.push({path: '/'})">返回网站首页</a></p>
  </div>
</template>

<style lang="scss" scoped>
.wscn-http404-container {
    height: 120px;
    width: 420px;
    margin: 0 auto;
    padding-top: 220px;
}
</style>
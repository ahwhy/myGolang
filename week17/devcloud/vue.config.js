"use strict";
const path = require("path");
const defaultSettings = require("./src/settings.js");

function resolve(dir) {
  return path.join(__dirname, dir);
}

const name = defaultSettings.title || '极乐研发云'; // page title

// If your port is set to 80,
// use administrator privileges to execute the command line.
// For example, Mac: sudo npm run
// You can change the port by the following method:
// port = 9527 npm run dev OR npm run dev --port = 9527
const port = process.env.port || process.env.npm_config_port || 9527; // dev por

// All configuration item explanations can be find in https://cli.vuejs.org/config/
module.exports = {
  /**
   * You will need to set publicPath if you plan to deploy your site under a sub path,
   * for example GitHub Pages. If you plan to deploy your site to https://foo.github.io/bar/,
   * then publicPath should be set to "/bar/".
   * In most cases please use '/' !!!
   * Detail: https://cli.vuejs.org/config/#publicpath
   */
  /*
    部署应用包时的基本 URL, '/'表示部署于跟路径,
    如果不放置于根, 比如: /my-app/, 访问路径就是： https://www.my-app.com/my-app/
  */
    publicPath: '/',
    /*
    当运行 vue-cli-service build 时生成的生产环境构建文件的目录
    */
    outputDir: 'dist',
    /*
    放置生成的静态资源 (js、css、img、fonts) 的 (相对于 outputDir 的) 目录
    */
    assetsDir: 'static',
    /*
    是否在开发环境下通过 eslint-loader 在每次保存时 lint 代码。这个值会在 @vue/cli-plugin-eslint 被安装之后生效,
    这里配置在开发环境生效
  
    注意: process.env 是当前进程的环境变量对象, 通过它可以访问到当前进程的所有环境变量
    */
    lintOnSave: process.env.NODE_ENV === 'development',
    /*
    是否需要生产环境的 source map, 可以将其设置为 false 以加速生产环境构建, 
    */
    productionSourceMap: false,
    /*
    webpack 相关配置
    */
    configureWebpack: {
      // provide the app's title in webpack's name field, so that
      // it can be accessed in index.html to inject the correct title.
      name: name,
      // webpack 插件配置, 具体见: https://webpack.js.org/configuration/plugins/#plugins
      plugins: [],
      // 使用import的路由别名, 比如'@/components/Tips' 会别解析成: src/components/Tips
      // 更多resovle相关配置请查看: https://webpack.js.org/configuration/resolve/#resolve
      resolve: {
        alias: {
          '@': resolve('src')
        }
      }
    },
}
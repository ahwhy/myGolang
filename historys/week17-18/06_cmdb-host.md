# 主机列表页面

我们先完成主机的列表页面, 列表页面大致布局如下:

![](./images/list-page-layout.jpg)

我们列表页面布局大概如下:
+ 提示信息
+ 表格操作: 左边: 搜索区域， 右边: 操作区域
+ 表格数据: 一个Box存放, 底部是分页信息


当前的页面如下:
```html
<template>
  <div class="host-container">
    Host 页面
  </div>
</template>

<script>
export default {
  name: 'CmdbHost',
  data() {
    return {}
  }
}
</script>
```

## Tips组件

tips大概布局:
+ 一个box 容器, 带背景色, Box有一个关闭按钮
+ 第一行: icon + title文字描述
+ 下面是列表文字说明

先写组件模版:
```html
<template>
  <div v-if="!hidden" class="tips">
    <div class="titile">
        <div class="tips-icon">
            <svg-icon icon-class="tips-info"></svg-icon>
        </div>
        <span class="title-content">
            温馨提示
        </span>
        <span class="close-btn">
          <i class="el-icon-close" @click="handleClose" />
        </span>
        
    </div>
    <div class="content">
        <li class="tip-item"> 具体提示</li>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Tips',
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
```

然后我们在Host列表页面 应用Tips组件:
```html
<template>
  <div class="host-container">
    <tips />
    Host 页面
  </div>
</template>

<script>
import Tips from '@/components/Tips'

export default {
  name: 'CmdbHost',
  components: { Tips },
  data() {
    return {}
  }
}
</script>
```


然后开始补充样式:
```scss
<style lang="scss" scoped>
.tips {
  width: 100%;
  background-color: rgba(48, 210, 190, 0.42);
  color: rgb(20, 105, 105);
  font-size: 12px;
  padding: 8px 16px;
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
```

然后我们来定义这个组件需要传入的变量: 
```js
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
  // ...
}
</script>
```

最后给组件传递参数 进行测试:
```html
<tips :tips="['主机列表页面']" />
```

为了代码干净, 我们把tips定义在 data里面:
```html
<template>
  <div class="host-container">
    <tips :tips="tips" />
    Host 页面
  </div>
</template>

<script>
import Tips from '@/components/Tips'

const tips = [
  '现在仅同步了阿里云主机资产'
]

export default {
  name: 'CmdbHost',
  components: { Tips },
  data() {
    return {
      tips: tips
    }
  }
}
</script>
```

![](./images/tips.jpg)

## 表格数据

表格数据，我们采用Element UI 的Table组件: [Table 表格](https://element.eleme.cn/#/zh-CN/component/table)

表格我们采用卡片风格, [Border 边框](https://element.eleme.cn/#/zh-CN/component/border)里提供了box-shadow的样式:
```
基础投影 box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04)
浅色投影 box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
```

我们采用浅色投影, 因此补充一个全局样式: box-shadow
```css
.box-shadow {
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}
```

然后我们copy一个样例过来看看效果:
```html
<template>
  <div class="host-container">
    <tips :tips="tips" />
    <div class="box-shadow">
      <el-table
        :data="tableData"
        style="width: 100%">
        <el-table-column
          prop="date"
          label="日期"
          width="180">
        </el-table-column>
        <el-table-column
          prop="name"
          label="姓名"
          width="180">
        </el-table-column>
        <el-table-column
          prop="address"
          label="地址">
        </el-table-column>
      </el-table>
    </div>
    Host 页面
  </div>
</template>

<script>
import Tips from '@/components/Tips'

const tips = [
  '现在仅同步了阿里云主机资产'
]

export default {
  name: 'CmdbHost',
  components: { Tips },
  data() {
    return {
      tips: tips,
      tableData: [{
          date: '2016-05-02',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1518 弄'
        }, {
          date: '2016-05-04',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1517 弄'
        }, {
          date: '2016-05-01',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1519 弄'
        }, {
          date: '2016-05-03',
          name: '王小虎',
          address: '上海市普陀区金沙江路 1516 弄'
        }]
    }
  }
}
</script>
```

我们之前是把获取数据的代码直接写在该模块内部: 
```js
import axios from 'axios';

// ...

getHosts() {
  axios
    .get('http://localhost:8050/hosts', {params: this.query})
    .then(response => {
      this.tableData = response.data.data.items
      this.total = response.data.data.total
      console.log(this.tableData)
    })
    .catch(function (error) { // 请求失败处理
      console.log(error);
    });
},
```

这样并不方便互用和维护, 因此我们把这个路径都放到API模块下: api/cmdb/host.js

如果每个API都直接使用 axios, 那么后面做一些 中间件处理异常的逻辑就没发做了, 因此我们需要构造一个

### 封装axios实例

安装引入的依赖:
```
npm i axios
```

api/client.js
```js
import axios from 'axios'

// http client
const client = axios.create({
    // API请求的base URL
    baseURL: "http://localhost:8050",
    // 超时时间
    timeout: 5000,
})

// request中间件
client.interceptors.request.use(
    // 成功的处理逻辑
    request => {
        return request
    },
    // 错误时的处理逻辑
    err => {
        console.log(err)
        return Promise.reject(err)
    }
)

// response中间件
client.interceptors.response.use(
    response => {
        const resp = response.data
        // 判断返回的error code是否为0, 如果为0请求成功
        if (resp.code === 0) {
            return resp
        }

        // 如果不为0, 请求失败
        console.log(resp)
    },
    err => {
        // 错误时的处理逻辑
        console.log(err)
        return Promise.reject(err)
    }
)

export default client
```

### 编写API

api/cmdb/host.js
```js
import request from '@/api/client'

export function LIST_HOST(params) {
    return request({
        url: '/cmdb/api/v1/hosts',
        method: 'get',
        params: params
    })
}
```

### 视图数据请求

启动我们的demo-api，测试
```
go run main.go start -f etc/demo-api.toml
```

可以看到hosts的数据已经请求回来了:

![](./images/hosts_data.png)

调整下我们页面视图显示的数据和描述:
```html
<template>
  <div class="host-container">
    <tips :tips="tips" />
    <div class="box-shadow">
      <el-table
        :data="hosts"
        style="width: 100%">
        <el-table-column
          prop="name"
          label="名称"
          width="180">
        </el-table-column>
        <el-table-column
          prop="sync_at"
          label="同步时间"
          width="180">
          <template slot-scope="scope">
            {{ scope.row.sync_at | parseTime}}
          </template>
          
        </el-table-column>
        <el-table-column
          prop="description"
          label="描述">
        </el-table-column>
      </el-table>
    </div>
    Host 页面
  </div>
</template>
```


### 后端分页

我们的API并没显示 总条数和分页, 这里需要适配后端分页, 最简单的做法是直接使用element的分页组件: [Pagination 分页](https://element.eleme.cn/#/zh-CN/component/pagination)

在表格下面补充分页组件: 
```html
<el-pagination
  @size-change="handleSizeChange"
  @current-change="handleCurrentChange"
  :current-page="query.page_number"
  :page-sizes="[2,10, 20, 30, 50]"
  :page-size="query.page_size"
  layout="total, sizes, prev, pager, next, jumper"
  :total="total">
</el-pagination>
```

data里面补充相关参数:
```js
data() {
  return {
    tips: tips,
    query: {page_size: 20, page_number: 1},
    total: 0,
    hosts: []
  }
```

最后是实现分页切换时的methods:
```js
handleSizeChange(val) {
  this.query.page_size = val
  this.get_hosts()
},
handleCurrentChange(val) {
  this.query.page_number = val
  this.get_hosts()
}
```

定制分页组件: Pagination

```html
<template>
  <div :class="{'hidden':hidden}" class="pagination-container">
    <el-pagination
      :background="background"
      :current-page.sync="currentPage"
      :page-size.sync="pageSize"
      :layout="layout"
      :page-sizes="pageSizes"
      :total="total"
      v-bind="$attrs"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
export default {
  name: 'Pagination',
  props: {
    total: {
      required: true,
      type: Number
    },
    page: {
      type: Number,
      default: 1
    },
    limit: {
      type: Number,
      default: 20
    },
    pageSizes: {
      type: Array,
      default() {
        return [10, 20, 30, 50]
      }
    },
    layout: {
      type: String,
      default: 'total, sizes, prev, pager, next, jumper'
    },
    background: {
      type: Boolean,
      default: true
    },
    autoScroll: {
      type: Boolean,
      default: true
    },
    hidden: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    currentPage: {
      get() {
        return this.page
      },
      set(val) {
        this.$emit('update:page', val)
      }
    },
    pageSize: {
      get() {
        return this.limit
      },
      set(val) {
        this.$emit('update:limit', val)
      }
    }
  },
  methods: {
    handleSizeChange(val) {
      this.$emit('pagination', { page: this.currentPage, limit: val })
    },
    handleCurrentChange(val) {
      this.$emit('pagination', { page: val, limit: this.pageSize })
    }
  }
}
</script>

<style scoped>
.pagination-container {
  background: #fff;
  padding: 16px 8px;
}
.pagination-container.hidden {
  display: none;
}
</style>
```

最后修改我们的host列表页使用该组件
```html
<pagination 
  v-show="total>0" 
  :total="total" 
  :page.sync="query.page_number" 
  :limit.sync="query.page_size" 
  @pagination="get_hosts" 
/>
```

## 搜索框

我们采用Element的 输入框组件: [Input 输入框](https://element.eleme.cn/#/zh-CN/component/input)

```html
<el-input placeholder="请输入内容" v-model="input3" class="input-with-select">
  <el-select v-model="select" slot="prepend" placeholder="请选择">
    <el-option label="餐厅名" value="1"></el-option>
    <el-option label="订单号" value="2"></el-option>
    <el-option label="用户电话" value="3"></el-option>
  </el-select>
  <el-button slot="append" icon="el-icon-search"></el-button>
</el-input>
```





# CMDB搜索页面

![](./images/search-ui.png)

搜索页面由2个部分组成:
+ 搜索框
+ 搜索内容

```html
<template>
  <div class="cmdb-search-container">
        <!-- 搜索框 -->
    <div class="box-shadow search-box">
      搜索页面
    </div>

        <!-- 搜索结果 -->
    <div class="box-shadow content-box">
      搜索结果
    </div>
  </div>
</template>

<script>
export default {
  name: "Search",
  data() {
    return {};
  },
};
</script>

<style scoped>
.search-box {
}

.content-box {
  padding-top: 12px;
}
</style>
```

## 搜索框

我们使用input输入框作为搜索框

```html
<!-- 搜索框 -->
<div class="box-shadow search-box">
    <div class="search-input">
    <el-input
        v-model="query.keywords"
        placeholder="请输入 资源名称|ID|IP 敲回车进行搜索"
    ></el-input>
    </div>
</div>
```

对应样式
```css
.search-box {
  height: 80px;
  background-color: white;
  display: flex;
  align-items: center;
  border-radius: 4px;
}

.search-input {
  width: 100%;
  margin-left: 12px;
  margin-right: 12px;
}
```

## 搜索结果

搜索的结果我们使用表格展示
```html
<!-- 搜索结果 -->
<div class="content-box">
    <el-table
    :data="resources"
    v-loading="fetchResourceLoading"
    style="width: 100%"
    >
    <el-table-column prop="name" label="名称">
        <template slot-scope="{ row }">
        {{ row.name }}
        </template>
    </el-table-column>
    <el-table-column prop="name" label="资产来源">
        <template slot-scope="{ row }">
        {{ row.vendor }} <br />
        {{ row.region }}
        </template>
    </el-table-column>
    <el-table-column prop="name" label="内网IP/外网IP">
        <template slot-scope="{ row }">
        {{ row.private_ip }} <br />
        {{ row.public_ip }}
        </template>
    </el-table-column>
    <el-table-column prop="name" label="同步时间">
        <template slot-scope="{ row }">
        {{ row.os_name }}
        </template>
    </el-table-column>
    <el-table-column prop="sync_at" label="创建时间">
        <template slot-scope="scope">
        {{ scope.row.create_at | parseTime }}
        </template>
    </el-table-column>
    <el-table-column prop="expire_at" label="过期时间">
        <template slot-scope="scope">
        {{ scope.row.expire_at | parseTime }}
        </template>
    </el-table-column>
    <el-table-column prop="name" label="状态">
        <template slot-scope="{ row }">
        {{ row.status }}
        </template>
    </el-table-column>
    </el-table>
</div>

<script>
export default {
  name: "Search",
  data() {
    return {
      fetchResourceLoading: false,
      resources: [],
      query: {
        keywords: "",
      },
    };
  },
};
</script>
```

## Search API

业务模块开发的步骤:
+ 数据结构与接口定义
+ 实现接口
+ 暴露HTTP服务
+ 注册接口实例
+ 添加路由

完成Search 模块, 具体参考 [resource 服务](https://github.com/infraboard/cmdb/tree/master/pkg/resource)

完成后测试接口是否正常:

![](./images/search-api.png)


## 数据展示

封装Seach API
```js
import request from "@/api/client";

export function Search(params) {
  return request({
    url: "/search",
    method: "get",
    params: params,
  });
}
```

引入到组件中:
```js
import { Search } from "@/api/cmdb/resource";
import Pagination from "@/components/Pagination";

export default {
  name: "Search",
  components: { Pagination },
  data() {
    return {
      fetchResourceLoading: false,
      resources: [],
      total: 0,
      query: {
        page_size: 20,
        page_number: 1,
        keywords: "",
      },
    };
  },
  methods: {
    async search() {
      this.fetchResourceLoading = true;
      try {
        const resp = await Search(this.query);
        this.resources = resp.data.items;
        this.total = resp.data.total;
      } catch (error) {
        this.$notify.error({
          title: "搜索资源异常",
          message: error,
        });
      } finally {
        this.fetchResourceLoading = false;
      }
    },
  },
};
</script>
```

最后添加搜索框和分页
```html
<template>
  <div class="cmdb-search-container">
    <!-- 搜索框 -->
    <div class="box-shadow search-box">
      <div class="search-input">
        <el-input
          v-model="query.keywords"
          placeholder="请输入 资源名称|ID|IP 敲回车进行搜索"
          @keyup.enter.native="search"
        ></el-input>
      </div>
    </div>

    <!-- 搜索结果 -->
    <div class="content-box">
      <el-table
        :data="resources"
        v-loading="fetchResourceLoading"
        style="width: 100%"
      >
        <el-table-column prop="name" label="名称">
          <template slot-scope="{ row }">
            {{ row.name }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="资产来源">
          <template slot-scope="{ row }">
            {{ row.vendor }} <br />
            {{ row.region }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="内网IP/外网IP">
          <template slot-scope="{ row }">
            {{ row.private_ip }} <br />
            {{ row.public_ip }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="同步时间">
          <template slot-scope="{ row }">
            {{ row.sync_at | parseTime }}
          </template>
        </el-table-column>
        <el-table-column prop="sync_at" label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.create_at | parseTime }}
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="过期时间">
          <template slot-scope="scope">
            {{ scope.row.expire_at | parseTime }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="状态">
          <template slot-scope="{ row }">
            {{ row.status }}
          </template>
        </el-table-column>
      </el-table>

      <pagination
        v-show="total > 0"
        :total="total"
        :page.sync="query.page_number"
        :limit.sync="query.page_size"
        @pagination="search"
      />
    </div>
  </div>
</template>
```

## 测试

+ 测试name/id/ip 是否可以正常使用
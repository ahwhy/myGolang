# CMDB资产同步页面

![](./images/secret-list-page.png)

## 页面与路由

先添加一个空页面: resource/sync.vue
```html
<template>
    <div>
        资产同步
    </div>
</template>

<script>
export default {
    name: "ResourceSync",
}
</script>
```

然后添加路由:
```js
  {
    path: "/cmdb",
    component: Layout,
    redirect: "/cmdb/search",
    children: [
        ...
      {
        path: "sync/list",
        component: () => import("@/views/cmdb/sync/index"),
        name: "ResourceSync",
      },
    ],
  },
```

补充menu:
```html
<el-submenu index="/sync">
    <!-- 添加个title -->
    <template slot="title">
    <i class="el-icon-s-tools"></i>
    <span slot="title">资源同步</span>
    </template>
    <!-- 导航条目 -->
    <el-menu-item index="/cmdb/sync/list">凭证管理</el-menu-item>
</el-submenu>
```

## 页面数据

API接口，我们上次课已经完成:

![](./images/secret-list.png)

添加前端API： cmdb/secret.js
```js
import request from "@/api/client";

export function LIST_SECRET(params) {
  return request({
    url: "/secrets",
    method: "get",
    params: params,
  });
}
```

通过接口获取数据:

```html
<script>
import { LIST_SECRET } from "@/api/cmdb/secret";
export default {
  name: "ResourceSync",
  data() {
    return {
      secrets: [],
      fetchSecretLoading: false,
      total: 0,
      query: {
        page_size: 20,
        page_number: 1,
      },
    };
  },
  mounted() {
    this.get_secrets();
  },
  methods: {
    async get_secrets() {
      this.fetchSecretLoading = true;
      try {
        const resp = await LIST_SECRET(this.query);
        this.secrets = resp.data.items;
        this.total = resp.data.total;
      } catch (error) {
        this.$notify.error({
          title: "获取凭证异常",
          message: error,
        });
      } finally {
        this.fetchSecretLoading = false;
      }
    },
  },
};
</script>
```

## 数据展示

使用Table展示数据
```html
<template>
  <div>
    <tips :tips="tips" />

    <!-- 凭证表格 -->
    <div class="box-shadow secret-box">
      <el-table
        :data="secrets"
        v-loading="fetchSecretLoading"
        style="width: 100%"
      >
        <el-table-column prop="name" label="描述">
          <template slot-scope="{ row }">
            {{ row.description }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="厂商">
          <template slot-scope="{ row }">
            {{ row.vendor }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="类型">
          <template slot-scope="{ row }">
            {{ row.crendential_type }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="凭证">
          <template slot-scope="{ row }">
            {{ row.api_key }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="同步地域">
          <template slot-scope="{ row }">
            {{ row.allow_regions.join(",") }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="创建时间">
          <template slot-scope="{ row }">
            {{ row.create_at | parseTime }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="速率限制">
          <template slot-scope="{ row }">
            {{ row.request_rate }}
          </template>
        </el-table-column>
        <el-table-column prop="操作" align="center" label="状态">
          <template>
            <el-button type="text" disabled>删除</el-button>
            <el-button type="text" disabled>禁用</el-button>
            <el-button type="text" disabled>测试</el-button>
          </template>
        </el-table-column>
      </el-table>

      <pagination
        v-show="total > 0"
        :total="total"
        :page.sync="query.page_number"
        :limit.sync="query.page_size"
        @pagination="get_secrets"
      />
    </div>
  </div>
</template>
```

## 关键字搜索

我们在table上方添加 搜索和创建

```html
<!-- 表格功能区 -->
<div class="table-op">
    <div class="search">
    <el-input
        v-model="query.keywords"
        placeholder="请输入凭证描述|API KEY|用户名 敲回车进行搜索"
        @keyup.enter.native="get_secrets"
    ></el-input>
    </div>
    <div class="op">
    <el-button icon="el-icon-plus" type="primary">添加凭证</el-button>
    </div>
</div>
```

样式调试
```html
<style scoped>
.secret-box {
  margin-top: 8px;
}

.table-op {
  margin-top: 12px;
  display: flex;
  align-items: center;
}

.op {
    margin-left: auto;
}
</style>
```

后端 API 添加关键字搜索

```go
```

## 添加凭证组件

我们采用抽屉组件 作为添加凭证的弹窗: [Drawer 抽屉](https://element.eleme.cn/#/zh-CN/component/drawer)


添加组件: AddSecret.vue
```html
<template>
  <div>
    <!-- 添加凭证表单 -->
    <el-drawer
      title="添加凭证"
      :visible.sync="isVisible"
      :before-close="handleClose"
      size="40%"
    >
      <span>我来啦!</span>
    </el-drawer>
  </div>
</template>

<script>
export default {
  name: "AddSecret",
  props: {
    visible: {
      default: false,
      type: Boolean,
    },
  },
  methods: {
    handleClose() {
      this.$emit("update:visible", false);
    },
  },
  computed: {
    isVisible: {
      get() {
        return this.visible;
      },
      set(val) {
        this.$emit("update:visible", val);
      },
    },
  },
};
</script>
```

然后在页面引入该组件
```html
<!-- 添加secret -->
<add-secret :visible.sync="showAddSecretDrawer" />
```

```js
import AddSecret from "./components/AddSecret";
export default {
  name: "ResourceSync",
  components: { Tips, Pagination, AddSecret },
```

点击添加组件是，显示我们的抽屉
```html
<el-button icon="el-icon-plus" type="primary" @click="handleAddSecret"
    >添加凭证</el-button
>
```

补充函数
```js
handleAddSecret() {
    this.showAddSecretDrawer = true;
},
```

测试该组件能否正常显示

## 添加凭证表单

使用Form组件来 显示表单: [Form 表单](https://element.eleme.cn/#/zh-CN/component/form)

### 基础骨架

我们在抽屉组件里面添加一个添加凭证的表单:
```html
<template>
  <div>
    <!-- 添加凭证表单 -->
    <el-drawer
      title="添加凭证"
      :visible.sync="isVisible"
      :before-close="handleClose"
      size="40%"
    >
      <div class="add-secret-form">
        <el-form ref="addSecretForm" :model="form">
          <el-form-item label="厂商" :label-width="formLabelWidth">
          </el-form-item>
          <el-form-item label="类型" :label-width="formLabelWidth">
          </el-form-item>
          <div>
            <el-form-item label="API Key" :label-width="formLabelWidth">
            </el-form-item>
            <el-form-item label="API Secret" :label-width="formLabelWidth">
            </el-form-item>
            <el-form-item label="Region" :label-width="formLabelWidth">
            </el-form-item>
          </div>
          <div>
            <el-form-item label="Address" :label-width="formLabelWidth">
            </el-form-item>
            <el-form-item label="Username" :label-width="formLabelWidth">
            </el-form-item>
            <el-form-item label="Password" :label-width="formLabelWidth">
            </el-form-item>
          </div>
          <el-form-item label="速率限制" :label-width="formLabelWidth">
          </el-form-item>
        </el-form>

        <div class="drawer-footer">
          <el-button @click="cancelForm">取 消</el-button>
          <el-button
            type="primary"
            :loading="addSecretLoading"
            @click="submit"
            >{{ addSecretLoading ? "提交中 ..." : "确 定" }}</el-button
          >
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script>
export default {
  name: "AddSecret",
  props: {
    visible: {
      default: false,
      type: Boolean,
    },
  },
  data() {
    return {
      form: {},
      formLabelWidth: "80px",
      addSecretLoading: false,
    };
  },
  methods: {
    handleClose() {
      this.$emit("update:visible", false);
    },
    cancelForm() {
      this.addSecretLoading = false;
      this.$emit("update:visible", false);
    },
  },
  computed: {
    isVisible: {
      get() {
        return this.visible;
      },
      set(val) {
        this.$emit("update:visible", val);
      },
    },
  },
};
</script>
```

最终骨架样子:

![](./images/add-secret-raw.png)


### 枚举类型API

可以看到有2个地方需要有枚举的数据:
+ 厂商
+ 类型

后端 sync模块 为该枚举提供API接口

utils下面补充专门用于返回枚举的数据结构
```go
type EnumDescribe struct {
	Value    string `json:"value"`
	Describe string `json:"describe"`
}
```

厂商: vednor, resource 模块
```go
func RegistAPI(r *httprouter.Router) {
    ...
	r.GET("/vendors", api.ListVendor)
}
```

handler
```go
func (h *handler) ListVendor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := []utils.EnumDescribe{
		{Value: resource.VendorAliYun.String(), Describe: "阿里云"},
		{Value: resource.VendorTencent.String(), Describe: "腾讯云"},
		{Value: resource.VendorHuaWei.String(), Describe: "华为云"},
		{Value: resource.VendorVsphere.String(), Describe: "Vsphere"},
	}
	response.Success(w, resp)
}
```


类型: crendential_type
```go
func RegistAPI(r *httprouter.Router) {
    ...
	r.GET("/crendential_types", api.ListCrendentialType)
}
```

handler
```go
func (h *handler) ListCrendentialType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := []utils.EnumDescribe{
		{Value: syncer.CrendentialAPIKey.String(), Describe: "API凭证"},
		{Value: syncer.CrendentialPassword.String(), Describe: "用户名密码"},
	}
	response.Success(w, resp)
}
```

测试API
+ http://127.0.0.1:8050/vendors
+ http://127.0.0.1:8050/crendential_types

前端添加枚举的API: cmdb/enum.js
```js
import request from "@/api/client";

export function LIST_VENDOR(params) {
  return request({
    url: "/vendors",
    method: "get",
    params: params,
  });
}

export function LIST_CRENDENTIAL_TYPE(params) {
  return request({
    url: "/crendential_types",
    method: "get",
    params: params,
  });
}
```

### 列表页枚举转换

由于我们列表页的枚举显示的API 值，不太友好, 可以将其转换为对应的描述显示:

secret/index.vue
```js
import { LIST_VENDOR, LIST_CRENDENTIAL_TYPE } from "@/api/cmdb/enum.js";

// 补充枚举值
data() {
    return {
        ...
        vendors: [],
        crendentialTypes: [],
    };
}

// 编写请求方法
async getVendors() {
    try {
    const resp = await LIST_VENDOR();
    this.vendors = resp.data.items;
    } catch (error) {
    this.$notify.error({
        title: "获取厂商异常",
        message: error,
    });
    }
},
async getCrendentialTypes() {
    try {
    const resp = await LIST_CRENDENTIAL_TYPE();
    this.crendentialTypes = resp.data.items;
    } catch (error) {
    this.$notify.error({
        title: "获取凭证类型异常",
        message: error,
    });
    }
},

// 页面加载时 获取枚举
  mounted() {
    this.getVendors();
    this.getCrendentialTypes();
    this.get_secrets();
  },
```

添加转换方法:
```js
getEnumDescribe(t, v) {
    let showVendor = v
    switch (t) {
    case "vendor":
        this.vendors.forEach((item) => {
        if (item.value === v) {
            showVendor = item.describe;
        }
        });
        break;
    case "cendential":
        this.crendentialTypes.forEach((item) => {
        if (item.value === v) {
            showVendor = item.describe;
        }
        });
        break;
    }
    return showVendor;
},
```

页面展示是 使用函数转换
```html
<el-table-column prop="name" label="厂商">
    <template slot-scope="{ row }">
    {{ getEnumDescribe("vendor", row.vendor) }}
    </template>
</el-table-column>
<el-table-column prop="name" label="类型">
    <template slot-scope="{ row }">
    {{ getEnumDescribe("cendential", row.crendential_type) }}
    </template>
</el-table-column>
```

![](./images/transfer-vendor.png)

### 表单绑定数据

我们将vendor 和 cendential types传入给 组件内部使用

添加2个属性:
```js
vendors: {
    default() {
    return [];
    },
    type: Array,
},
cendentialTypes: {
    default() {
    return [];
    },
    type: Array,
},
```

然后 index页面传入参数:

```html
<!-- 添加secret -->
<add-secret
    :visible.sync="showAddSecretDrawer"
    :vendors="vendors"
    :cendentialTypes="crendentialTypes"
/>
```

1. 厂商

这里使用 [Radio 单选框](https://element.eleme.cn/#/zh-CN/component/radio)

动态选择枚举列表:

```html
<el-form-item label="厂商" :label-width="formLabelWidth">
    <el-radio-group v-model="form.vendor">
        <el-radio-button
        v-for="item in vendors"
        :key="item.value"
        :label="item.value"
        >
        {{ item.describe }}
        </el-radio-button>
    </el-radio-group>
</el-form-item>
```

通过watch 设置初始值: 注意, 由于数据是挂载后加载的, 所以页面早于数据, 需要watch动态修改(或者父页面使用created钩子)
```js
watch: {
    visible: {
        handler(newV) {
        if (newV) {
            this.form.vendor = this.vendors[0].value;
        }
        },
        immediate: true,
    },
},
```

2. 绑定其他数据

定义我们要提交的表单:
```js
data() {
    return {
        form: {
        description: "",
        vendor: "",
        crendential_type: "",
        api_key: "",
        api_secret: "",
        request_rate: 0.2,
        address: "",
        username: "",
        password: "",
        },
        formLabelWidth: "80px",
        addSecretLoading: false,
    };
},
```

绑定数据
```html
<template>
  <div>
    <!-- 添加凭证表单 -->
    <el-drawer
      title="添加凭证"
      :visible.sync="isVisible"
      :before-close="handleClose"
      size="40%"
    >
      <div class="add-secret-form">
        <el-form ref="addSecretForm" :model="form">
          <el-form-item label="厂商" :label-width="formLabelWidth">
            <el-radio-group v-model="form.vendor">
              <el-radio-button
                v-for="item in vendors"
                :key="item.value"
                :label="item.value"
              >
                {{ item.describe }}
              </el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="类型" :label-width="formLabelWidth">
            <el-radio-group v-model="form.crendential_type">
              <el-radio-button
                v-for="item in cendentialTypes"
                :key="item.value"
                :label="item.value"
              >
                {{ item.describe }}
              </el-radio-button>
            </el-radio-group>
          </el-form-item>
          <div>
            <el-form-item label="API Key" :label-width="formLabelWidth">
              <el-input
                v-model="form.api_key"
                placeholder="请输入内容"
              ></el-input>
            </el-form-item>
            <el-form-item label="API Secret" :label-width="formLabelWidth">
              <el-input
                v-model="form.api_secret"
                placeholder="请输入内容"
              ></el-input>
            </el-form-item>
            <el-form-item label="Region" :label-width="formLabelWidth">
              <el-input
                v-model="form.allow_regions"
                placeholder="请输入内容"
              ></el-input>
            </el-form-item>
          </div>
          <div>
            <el-form-item label="Address" :label-width="formLabelWidth">
              <el-input
                v-model="form.address"
                placeholder="请输入内容"
              ></el-input>
            </el-form-item>
            <el-form-item label="Username" :label-width="formLabelWidth">
              <el-input
                v-model="form.username"
                placeholder="请输入内容"
              ></el-input>
            </el-form-item>
            <el-form-item label="Password" :label-width="formLabelWidth">
              <el-input
                v-model="form.password"
                placeholder="请输入内容"
              ></el-input>
            </el-form-item>
          </div>
          <el-form-item label="速率限制" :label-width="formLabelWidth">
            <el-input
              v-model="form.request_rate"
              placeholder="请输入内容"
            ></el-input>
          </el-form-item>
        </el-form>

        <div class="drawer-footer">
          <el-button @click="cancelForm">取 消</el-button>
          <el-button
            type="primary"
            :loading="addSecretLoading"
            @click="submit"
            >{{ addSecretLoading ? "提交中 ..." : "确 定" }}</el-button
          >
        </div>
      </div>
    </el-drawer>
  </div>
</template>
```

添加样式
```css
<style scoped>
/* form表单添加边距 */
.add-secret-form {
  padding: 20px 20px;
}

/* 提交表单居中 */
.drawer-footer {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
```

![](./images/add-secret-form.png)


### 适配密码

使用v-if做判断:

```html
<div v-if="form.crendential_type === 'CRENDENTIAL_API_KEY'">
  <el-form-item label="API Key" :label-width="formLabelWidth">
    <el-input
      v-model="form.api_key"
      placeholder="请输入内容"
    ></el-input>
  </el-form-item>
  <el-form-item label="API Secret" :label-width="formLabelWidth">
    <el-input
      v-model="form.api_secret"
      placeholder="请输入内容"
    ></el-input>
  </el-form-item>
  <el-form-item label="Region" :label-width="formLabelWidth">
    <el-input
      v-model="form.allow_regions"
      placeholder="请输入内容"
    ></el-input>
  </el-form-item>
</div>
<div v-if="form.crendential_type === 'CRENDENTIAL_PASSWORD'">
  <el-form-item label="Address" :label-width="formLabelWidth">
    <el-input
      v-model="form.address"
      placeholder="请输入内容"
    ></el-input>
  </el-form-item>
  <el-form-item label="Username" :label-width="formLabelWidth">
    <el-input
      v-model="form.username"
      placeholder="请输入内容"
    ></el-input>
  </el-form-item>
  <el-form-item label="Password" :label-width="formLabelWidth">
    <el-input
      v-model="form.password"
      placeholder="请输入内容"
    ></el-input>
  </el-form-item>
</div>
```

### 自动适配API KEY 和密码

厂商绑定changeVendor方法
```html
<el-form-item label="厂商" :label-width="formLabelWidth">
  <el-radio-group @change="changeVendor" v-model="form.vendor">
    <el-radio-button
      v-for="item in vendors"
      :key="item.value"
      :label="item.value"
    >
      {{ item.describe }}
    </el-radio-button>
  </el-radio-group>
</el-form-item>
```

当有变化是, 修正crendential type
```js
changeVendor(val) {
  if (val === "VENDOR_VSPHERE") {
    this.form.crendential_type = "CRENDENTIAL_PASSWORD";
  } else {
    this.form.crendential_type = "CRENDENTIAL_API_KEY";
  }
},
```

### 添加校验规则






### 提交表单




<template>
  <div class="main-container">
    <tips :tips="tips" />
    <h3>CMDB Host</h3>

    <div class="table-op">
      <div class="search">
        <el-input placeholder="请输入内容" v-model="filterValue" class="input-with-select">
          <el-select v-model="filterKey" slot="prepend" placeholder="请选择">
            <el-option label="名称" value="1"></el-option>
            <el-option label="同步时间" value="2"></el-option>
            <el-option label="描述" value="3"></el-option>
          </el-select>
          <el-button slot="append" icon="el-icon-search"></el-button>
        </el-input>
      </div>
      <div class="op">
      </div>
    </div>

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

      <pagination 
        v-show="total>0" 
        :total="total" 
        :page.sync="query.page_number" 
        :limit.sync="query.page_size" 
        @pagination="get_hosts" 
      />
    </div>
  </div>
</template>

<script>
import Tips from '@/components/Tips'
import Pagination from '@/components/Pagination'
import { LIST_HOST } from '@/api/cmdb/host.js'

const tips = [
  '现在仅同步了阿里云主机资产'
]

export default {
  name: 'Host',
  components: { Tips, Pagination },
  data() {
    return {
      tips: tips,
      filterKey: '',
      filterValue: '',
      query: {page_size: 20, page_number: 1},
      total: 0,
      hosts: []
    }
  },
  created() {
    this.get_hosts()
  },
  methods: {
    async get_hosts() {
      const resp = await LIST_HOST(this.query)
      console.log(resp)
      this.hosts = resp.data.items
      this.total = resp.data.total
    },
    handleSizeChange(val) {
      this.query.page_size = val
      this.get_hosts()
    },
    handleCurrentChange(val) {
      this.query.page_number = val
      this.get_hosts()
    }
  }
}
</script>

<style lang="scss" scoped>
.box-shadow {
  margin: 12px 0;
}
</style>
<template>
  <div>
    <el-input v-model="query.keywords" @keyup.enter.native="getHosts" placeholder="请输入内容"></el-input>
    <el-table
      :data="tableData"
      style="width: 100%">
      <el-table-column
        prop="id"
        label="资产ID"
        width="180">
      </el-table-column>
      <el-table-column
        prop="name"
        label="主机名称"
        width="180">
      </el-table-column>
      <el-table-column
        prop="region"
        label="地址">
      </el-table-column>
    </el-table>
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="query.page_number"
      :page-sizes="[2, 10, 20, 30, 50]"
      :page-size="query.page_size"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total">
    </el-pagination>
  </div>

  </template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      query: {
        keywords: '',
        page_size: 20,
        page_number: 1,
      },
      total: 0,
      tableData: []
    }
  },
  mounted() {
    this.getHosts()
  },
  methods: {
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
    handleSizeChange(val) {
      this.query.page_size = val
      this.getHosts()
    },
    handleCurrentChange(val) {
      this.query.page_number = val
      this.getHosts()
    }
  },
}
</script>
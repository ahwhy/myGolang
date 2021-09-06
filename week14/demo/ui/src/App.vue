  <template>
    <div>
      <el-table
        :data="tableData"
        style="width: 100%">
        <el-table-column
          prop="id"
          label="ID"
          width="180">
        </el-table-column>
        <el-table-column
          prop="name"
          label="姓名"
          width="180">
        </el-table-column>
        <el-table-column
          prop="region"
          label="Region">
        </el-table-column>
      </el-table>
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="query.page_number"
        :page-sizes="[2,10, 20, 30, 50]"
        :page-size="query.page_size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total">
      </el-pagination>
      
      <h1>这是一个标题</h1>
      <h2>这是一个标题</h2>
      <h3>这是一个标题</h3>
      <h4>这是一个标题</h4>
      
      <iframe src="//www.runoob.com">
        <p>您的浏览器不支持  iframe 标签。</p>
      </iframe>
    </div>
  </template>

<script>
  import axios from 'axios';

  export default {
    data() {
      return {
        query: {
          page_size: 20,
          page_number: 1,
        },
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
              console.log(response)
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
    }
  }
</script>

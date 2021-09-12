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
      
      <h1 class="f12" id="css-start">基础标签</h1>
      <h1>这是一个H1标题</h1>
      <h6>这是一个H6标签</h6>

      <p> 某某: <br>
        我是谁谁
      </p>
      <hr>

      <p>这是一个段落</p>
      <p>段落3</p>

      <p>这是一个 <br />  段落</p>


      <h1>文本标签</h1>
        <!-- <del> 和 <ins> 一起使用，描述文档中的更新和修正。浏览器通常会在已删除文本上添加一条删除线，在新插入文本下添加一条下划线 -->
      <p></p>

      
      <u><del>被删除的文字</del></u>
      <p>上标<sup>1</sup></p>



      <p>This is a <u>parragraph</u>.</p>

      <h1>表单标签</h1>
      <form action="demo_form.php">
        <label for="male">Male</label>
        <input type="radio" name="sex" id="male" value="male"><br>
        <label for="female">Female</label>
        <input type="radio" name="sex" id="female" value="female"><br><br>
        <input type="submit" value="提交">
      </form>

      <h1>内联框架</h1>
      <iframe src="//www.runoob.com">
        <p>您的浏览器不支持  iframe 标签。</p>
      </iframe>

      <hr>

      <img src="https://img0.baidu.com/it/u=3805813364,2396910389&fm=11&fmt=auto&gp=0.jpg" alt="">

      <h1>列表</h1>
      <ul id="list_menu" class="ul_class">
          <li id="coffee">Coffee</li>
          <li>Tea</li>
          <li>Milk</li>
          <div>
            <li>In Div</li>
          </div>
      </ul>

      <h1>表格</h1>
      <table>

      </table>
      <table border="1">
      <tr>
      <th>Month</th>
      <th>Savings</th>
      </tr>
      <tr>
      <td>January</td>
      <td>$100</td>
      </tr>
      </table>

      <div>
        <p>第一段, <span>span演示</span> <span>span演示</span> </p>
        <p>第二段</p>
      </div>


      <p>
      <a href="#p4">查看章节 4</a>
      </p>

      <h2>章节 1</h2>
      <p>这边显示该章节的内容……</p>

      <h2>章节 2</h2>
      <p>这边显示该章节的内容……</p>

      <h2>章节 3</h2>
      <p>这边显示该章节的内容……</p>

      <h2><a>章节 4</a></h2>
      <p id="p4">这边显示该章节的内容……</p>

      <div style="height: 220px;width:440px">
          <p style="margin-top:22px;">我们的内容</p>
      </div>

      <hr>
      <div >
        <p>我是垂直居中的。</p>
      </div>

      <div>
          <li>Tea</li>
          <li>Milk</li>
      </div>

      <div>
          <span>span1</span>
          <span>span2</span>
      </div>

      <div id="overflowTest" style="clear:both">
        <p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
        <p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
        <p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
        <p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
        <p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
        <p>这里的文本内容是可以滚动的，滚动条方向是垂直方向。</p>
      </div>
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
        total: 0,
        tableData: []
      }
    },
    methods: {
      getHosts() {
        // loading
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

<style scoped>

.f12 {
  font-size: 12px;
}

#css-start {
  color: #4CAF50;
}

ul>li:nth-child(2) {
  font-weight: 600;
}

#overflowTest {
    background: #4c53af;
    color: white;
    padding: 15px;
    width: 80%;
    height: 100px;
    overflow: scroll;
    border: 1px solid rgb(150, 18, 18);
}


</style>
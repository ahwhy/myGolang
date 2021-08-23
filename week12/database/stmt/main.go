package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ahwhy/myGolang/week12/database"

	_ "github.com/go-sql-driver/mysql"
)

const TIME_LAYOUT = "2006-01-02"

var (
	loc *time.Location
)

func init() {
	loc, _ = time.LoadLocation("Asia/Shanghai")
}

//insert 通过stmt插入数据
func insert(db *sql.DB) {
	//一条sql，插入2行记录
	stmt, err := db.Prepare("insert into student (name,province,city,enrollment) values (?,?,?,?), (?,?,?,?)")
	database.CheckError(err)
	//字符串解析为时间。注意要使用time.ParseInLocation()函数指定时区，time.Parse()函数使用默认的UTC时区
	date1, err := time.ParseInLocation(TIME_LAYOUT, "2021-03-18", loc)
	database.CheckError(err)
	date2, err := time.ParseInLocation(TIME_LAYOUT, "2021-03-26", loc)
	database.CheckError(err)
	//执行修改操作通过stmt.Exec，执行查询操作通过stmt.Query
	res, err := stmt.Exec("小明", "深圳", "深圳", date1, "小红", "上海", "上海", date2)
	database.CheckError(err)
	lastId, err := res.LastInsertId() //ID自增，用过的id（即使对应的行已delete）不会重复使用
	database.CheckError(err)
	fmt.Printf("after insert last id %d\n", lastId)
	rows, err := res.RowsAffected() //插入2行，所以影响了2行
	database.CheckError(err)
	fmt.Printf("insert affect %d row\n", rows)
}

//replace 通过stmt插入(覆盖)数据
func replace(db *sql.DB) {
	//由于name字段上有唯一索引，insert重复的name会报错。而使用replace会先删除，再插入
	stmt, err := db.Prepare("replace into student (name,province,city,enrollment) values (?,?,?,?), (?,?,?,?)")
	database.CheckError(err)
	//字符串解析为时间。注意要使用time.ParseInLocation()函数指定时区，time.Parse()函数使用默认的UTC时区
	date1, err := time.ParseInLocation(TIME_LAYOUT, "2021-04-18", loc)
	database.CheckError(err)
	date2, err := time.ParseInLocation(TIME_LAYOUT, "2021-04-26", loc)
	fmt.Printf("day of 2021-04-26 is %d\n", date2.Local().Day())
	//执行修改操作通过stmt.Exec，执行查询操作通过stmt.Query
	res, err := stmt.Exec("小明", "深圳", "深圳", date1, "小红", "上海", "上海", date2)
	database.CheckError(err)
	lastId, err := res.LastInsertId() //ID自增，用过的id（即使对应的行已delete）不会重复使用
	database.CheckError(err)
	fmt.Printf("after insert last id %d\n", lastId)
	rows, err := res.RowsAffected() //先删除，后插入，影响了4行
	database.CheckError(err)
	fmt.Printf("insert affect %d row\n", rows)
}

//update 通过stmt修改数据
func update(db *sql.DB) {
	//不同的city加不同的分数
	stmt, err := db.Prepare("update student set score=score+? where city=?")
	database.CheckError(err)
	//执行修改操作通过stmt.Exec，执行查询操作通过stmt.Query
	res, err := stmt.Exec(10, "上海") //上海加10分
	database.CheckError(err)
	res, err = stmt.Exec(9, "深圳") //深圳加9分
	database.CheckError(err)
	lastId, err := res.LastInsertId() //0, 仅插入操作才会给LastInsertId赋值
	database.CheckError(err)
	fmt.Printf("after update last id %d\n", lastId)
	rows, err := res.RowsAffected() //where city=?命中了几行，就会影响几行
	database.CheckError(err)
	fmt.Printf("update affect %d row\n", rows)
}

//query 通过stmt查询数据
func query(db *sql.DB) {
	stmt, err := db.Prepare("select id,name,city,score from student where id>?")
	database.CheckError(err)
	//执行修改操作通过stmt.Exec，执行查询操作通过stmt.Query
	rows, err := stmt.Query(2) //查询得分大于2的记录
	database.CheckError(err)
	for rows.Next() { //没有数据或发生error时返回false
		var id int
		var score float32
		var name, city string
		err = rows.Scan(&id, &name, &city, &score) //通过scan把db里的数据赋给go变量
		database.CheckError(err)
		fmt.Printf("id=%d, score=%.2f, name=%s, city=%s \n", id, score, name, city)
	}
}

//delete 通过stmt删除数据
func delete(db *sql.DB) {
	stmt, err := db.Prepare("delete from student where id>?")
	database.CheckError(err)
	//执行修改操作通过stmt.Exec，执行查询操作通过stmt.Query
	res, err := stmt.Exec(13) //删除得分大于13的记录
	database.CheckError(err)
	rows, err := res.RowsAffected() //where id>?命中了几行，就会影响几行
	database.CheckError(err)
	fmt.Printf("delete affect %d row\n", rows)
}

func main() {
	//连接数据库的规范格式：user:password@tcp(localhost:5555)/dbname?charset=utf8
	//如果是本地，且采用默认的3306端口，可简写为：user:password@/dbname
	// db, err := sql.Open("mysql", "root:@/test")
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	database.CheckError(err)
	// insert(db)
	// replace(db)
	// update(db)
	// query(db)
	delete(db)
}

package main

import (
	"database/sql"
	"fmt"

	"github.com/ahwhy/myGolang/week12/database"

	_ "github.com/go-sql-driver/mysql"
)

// const TIME_LAYOUT = "2006-01-02"

// var (
// 	loc *time.Location
// )

// func init() {
// 	loc, _ = time.LoadLocation("Asia/Shanghai")
// }

// insert 插入数据
func insert(db *sql.DB) {
	// 一条sql，插入2行记录
	res, err := db.Exec("insert into student (name,province,city,enrollment) values ('小明', '深圳', '深圳', '2021-04-18'), ('小红', '上海', '上海', '2021-04-26')")
	database.CheckError(err)

	lastId, err := res.LastInsertId() // ID自增，用过的id(即使对应的行已delete)不会重复使用
	database.CheckError(err)
	fmt.Printf("After insert last id %d\n", lastId)

	rows, err := res.RowsAffected() // 插入2行，所以影响2行
	database.CheckError(err)
	fmt.Printf("Insert affect %d row\n", rows)
}

// replace 插入(覆盖)数据
func replace(db *sql.DB) {
	// 由于name字段上有唯一索引，insert重复的name会报错；而使用replace会先删除，再插入
	res, err := db.Exec("replace into student (name,province,city,enrollment) values ('小明', '深圳', '深圳', '2021-04-18'), ('小红', '上海', '上海', '2021-04-26')")
	database.CheckError(err)

	lastId, err := res.LastInsertId() // ID自增，用过的id（即使对应的行已delete）不会重复使用
	database.CheckError(err)
	fmt.Printf("After insert last id %d\n", lastId)

	rows, err := res.RowsAffected() // 先删除，后插入，影响了4行
	database.CheckError(err)
	fmt.Printf("Insert affect %d row\n", rows)
}

// update 修改数据
func update(db *sql.DB) {
	// 不同的city加不同的分数
	res, err := db.Exec("update student set score=score+10 where city='上海'") // 上海加10分
	database.CheckError(err)

	lastId, err := res.LastInsertId() // 0, 仅插入操作才会给LastInsertId赋值
	database.CheckError(err)
	fmt.Printf("After update last id %d\n", lastId)

	rows, err := res.RowsAffected() // where city=?命中了几行，就会影响几行
	database.CheckError(err)
	fmt.Printf("Update affect %d row\n", rows)
}

// query 查询数据
func query(db *sql.DB) {
	rows, err := db.Query("select id,name,city,score from student where id>2") //查询得分大于2的记录
	database.CheckError(err)

	for rows.Next() { // 没有数据或发生error时返回false
		var id int
		var score float32
		var name, city string
		err = rows.Scan(&id, &name, &city, &score) // 通过scan把db里的数据赋给go变量
		database.CheckError(err)
		fmt.Printf("id=%d, score=%.2f, name=%s, city=%s \n", id, score, name, city)
	}
}

// delete 删除数据
func delete(db *sql.DB) {
	res, err := db.Exec("delete from student where id>13") // 删除得分大于13的记录
	database.CheckError(err)

	rows, err := res.RowsAffected() // where id>13命中了几行，就会影响几行
	database.CheckError(err)
	fmt.Printf("delete affect %d row\n", rows)
}

func main() {
	/**
	DSN(data surce name)格式：[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	例如user:password@tcp(localhost:5555)/dbname?charset=utf8
	如果是本地MySQl，且采用默认的3306端口，可简写为：user:password@/dbname
	*/

	// db, err := sql.Open("mysql", "root:@/test")
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")
	database.CheckError(err)

	insert(db)
	replace(db)
	update(db)
	query(db)
	delete(db)
}

package main

import (
	"database/sql"
	"fmt"
	"go-course/database"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

func query() {
	where := map[string]interface{}{
		"city":     []string{"北京", "上海", "杭州"},
		"score<":   30,
		"addr":     builder.IsNotNull,
		"_orderby": "score desc",
		"_groupby": "province",
	}
	table := "student"
	fields := []string{"id", "name", "city", "score"}
	template, values, err := builder.BuildSelect(table, where, fields)
	database.CheckError(err)
	fmt.Println(template)  //包含占位符的sql模板
	fmt.Println(values...) //占位符的具体值
}

func insert() {
	data := []map[string]interface{}{
		{"name": "王五", "province": "河南", "city": "郑州", "enrollment": "2021-05-01"},
		{"name": "大王", "province": "浙江", "city": "杭州", "enrollment": "2021-04-01"},
	}
	table := "student"
	template, values, err := builder.BuildInsert(table, data)
	database.CheckError(err)
	fmt.Println(template)  //包含占位符的sql模板
	fmt.Println(values...) //占位符的具体值
}

func update() {
	where := map[string]interface{}{
		"city": []string{"北京", "上海", "杭州"},
	}
	data := map[string]interface{}{
		"score": 25,
	}
	table := "student"
	template, values, err := builder.BuildUpdate(table, where, data)
	database.CheckError(err)
	fmt.Println(template)  //包含占位符的sql模板
	fmt.Println(values...) //占位符的具体值
}

func delete() {
	where := map[string]interface{}{
		"city": "杭州",
	}
	table := "student"
	template, values, err := builder.BuildDelete(table, where)
	database.CheckError(err)
	fmt.Println(template)  //包含占位符的sql模板
	fmt.Println(values...) //占位符的具体值
}

func query2(db *sql.DB) {
	where := map[string]interface{}{
		"city":     []string{"北京", "上海", "杭州"},
		"score<":   30,
		"addr":     builder.IsNotNull,
		"_orderby": "score desc",
	}
	table := "student"
	fields := []string{"id", "name", "city", "score"}
	template, values, err := builder.BuildSelect(table, where, fields)
	database.CheckError(err)
	rows, err := db.Query(template, values...)
	database.CheckError(err)
	for rows.Next() {
		var id int
		var name, city string
		var score float32
		err := rows.Scan(&id, &name, &city, &score)
		database.CheckError(err)
		fmt.Printf("%d %s %s %.2f\n", id, name, city, score)
	}
}

func insert2(db *sql.DB) {
	data := []map[string]interface{}{
		{"name": "王五", "province": "河南", "city": "郑州", "enrollment": "2021-05-01"},
		{"name": "大王", "province": "浙江", "city": "杭州", "enrollment": "2021-04-01"},
	}
	table := "student"
	template, values, err := builder.BuildReplaceInsert(table, data) //使用replace
	database.CheckError(err)
	res, err := db.Exec(template, values...)
	database.CheckError(err)
	rows, err := res.RowsAffected()
	database.CheckError(err)
	fmt.Printf("insert affect %d rows\n", rows)
}

func update2(db *sql.DB) {
	where := map[string]interface{}{
		"city": []string{"北京", "上海", "杭州"},
	}
	data := map[string]interface{}{
		"score": 25,
	}
	table := "student"
	template, values, err := builder.BuildUpdate(table, where, data)
	database.CheckError(err)
	res, err := db.Exec(template, values...)
	database.CheckError(err)
	rows, err := res.RowsAffected()
	database.CheckError(err)
	fmt.Printf("update affect %d rows\n", rows)
}

func delete2(db *sql.DB) {
	where := map[string]interface{}{
		"city": "杭州",
	}
	table := "student"
	template, values, err := builder.BuildDelete(table, where)
	database.CheckError(err)
	res, err := db.Exec(template, values...)
	database.CheckError(err)
	rows, err := res.RowsAffected()
	database.CheckError(err)
	fmt.Printf("delete affect %d rows\n", rows)
}

func main() {
	query()
	insert()
	update()
	delete()
	fmt.Println("==================")

	dbName := "test"
	user := "root"
	password := ""
	host := "localhost"
	db, err := manager.New(dbName, user, password, host).Set(
		manager.SetCharset("utf8"),
	).Port(3306).Open(true)
	database.CheckError(err)
	query2(db)
	insert2(db)
	update2(db)
	delete2(db)
}

package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ahwhy/myGolang/week12/database"

	_ "github.com/go-sql-driver/mysql"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("Asia/Shanghai")
}

func insert(db *sql.DB) {
	stmt, err := db.Prepare("insert into student (name,province,city,enrollment) values (?,?,?,?)")
	database.CheckError(err)
	date, err := time.ParseInLocation("20060102", "20210618", loc)

	for i := 0; i < 100000; i++ {
		stmt.Exec("宋江"+strconv.Itoa(i), "山西", "大同", date)
	}
}

func query(db *sql.DB) {
	stmt, err := db.Prepare("select id,name from student where province=?")
	database.CheckError(err)

	rows, err := stmt.Query("山西")
	database.CheckError(err)

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("id=%d name=%s\n", id, name)
	}
}

func traverse(db *sql.DB) {
	var offset int
	begin := time.Now()
	stmt, _ := db.Prepare("select id,name,province from student limit ?,100")

	for i := 0; i < 100; i++ {
		t0 := time.Now()
		rows, _ := stmt.Query(offset)
		offset += 100
		fmt.Println(i, time.Since(t0))

		for rows.Next() {
			var id int
			var name string
			var province string
			rows.Scan(&id, &name, &province)
		}
	}
	fmt.Println("total", time.Since(begin))
}

func traverse2(db *sql.DB) {
	var maxid int
	begin := time.Now()
	stmt, _ := db.Prepare("select id,name,province from student where id>? limit 100")

	for i := 0; i < 100; i++ {
		t0 := time.Now()
		rows, _ := stmt.Query(maxid)
		fmt.Println(i, time.Since(t0))

		for rows.Next() {
			var id int
			var name string
			var province string
			rows.Scan(&id, &name, &province)
			if id > maxid {
				maxid = id
			}
		}
	}
	fmt.Println("total", time.Since(begin))
}

func main() {
	db, err := sql.Open("mysql", "root:@/test")
	database.CheckError(err)
	insert(db)
	query(db)

	fmt.Println("============")
	traverse(db)
	fmt.Println("============")
	traverse2(db)
}

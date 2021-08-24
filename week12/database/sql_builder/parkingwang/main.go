package main

import (
	"fmt"

	gsb "github.com/parkingwang/go-sqlbuilder"
)

func insert() {
	sql := gsb.NewContext().Insert("student").
		Columns("name", "province", "city", "enrollment").
		Values("赵四", "江苏", "南京", "2021-02-18").
		ToSQL()
	fmt.Println(sql)
}

func insert2() {
	sql := gsb.NewContext().Insert("student").
		Columns("name", "province", "city", "enrollment").
		ToSQL()
	fmt.Println(sql)
}

func update() {
	ctx := gsb.NewContext()
	sql := ctx.Update("student").
		Columns("name", "province", "city", "enrollment"). //statment占位符
		Where(ctx.EqTo("province", "河南").
			And().In("city", "郑州", "洛阳")).
		ToSQL()
	fmt.Println(sql)
}

func update2() {
	ctx := gsb.NewContext()
	sql := ctx.Update("student").
		Columns("name", "province", "city", "enrollment"). //statment占位符
		Where(ctx.Eq("province")).
		ToSQL()
	fmt.Println(sql)
}

func query() {
	sql := gsb.NewContext().Select("id", "name", "score", "city").
		From("student").
		OrderBy("score").DESC().
		Column("name").ASC().
		Limit(10).Offset(20).
		ToSQL()
	fmt.Println(sql)
}

func delete() {
	ctx := gsb.NewContext()
	sql := ctx.Delete("student").
		Where(ctx.GEtTo("score", 10)).
		ToSQL()
	fmt.Println(sql)
}

func delete2() {
	ctx := gsb.NewContext()
	sql := ctx.Delete("student").
		Where(ctx.GEt("score")).
		ToSQL()
	fmt.Println(sql)
}

func main() {
	insert()
	insert2()
	update()
	update2()
	query()
	delete()
	delete2()
}

/*
INSERT INTO `student`(`name`, `province`, `city`, `enrollment`) VALUES ('赵四', '江苏', '南京', '2021-02-18');
INSERT INTO `student`(`name`, `province`, `city`, `enrollment`) VALUES (?, ?, ?, ?);
UPDATE `student` SET `name`=?, `province`=?, `city`=?, `enrollment`=? WHERE `province` = '河南' AND `city` IN ('郑州', '洛阳');
UPDATE `student` SET `name`=?, `province`=?, `city`=?, `enrollment`=? WHERE `province` = ?;
SELECT `id`, `name`, `score`, `city` FROM `student` ORDER BY `score` DESC, `name` ASC LIMIT 10 OFFSET 20;
DELETE FROM `student` WHERE `score` >= 10;
DELETE FROM `student` WHERE `score` >= ?;
*/

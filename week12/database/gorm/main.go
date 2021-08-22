package main

import (
	"fmt"
	"go-course/database"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//默认情况下，GORM 使用 ID 作为主键，使用结构体名的 蛇形复数 作为表名，字段名的 蛇形 作为列名
type Student struct {
	Id         int    `gorm:"column:id;primaryKey"`
	Name       string `gorm:"column:name"`
	Province   string
	City       string    `gorm:"column:city"`
	Address    string    `gorm:"column:addr"`
	Score      float32   `gorm:"column:score"`
	Enrollment time.Time `gorm:"column:enrollment;type:date"`
}

//使用TableName()来修改默认的表象
func (Student) TableName() string {
	return "student"
}

func query(db *gorm.DB) {
	/**
	普通的where查询
	*/
	//返回一条记录
	var student Student
	db.Where("city=?", "郑州").First(&student) //有First就有Last
	fmt.Println(student.Name)
	fmt.Println()
	//返回多条记录
	var students []Student
	db.Where("city=?", "郑州").Find(&students)
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println()
	students = []Student{}
	db.Where("city in ?", []string{"郑州", "北京"}).Find(&students)
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println("============where end============")
	//根据主键查询
	student = Student{} //清空student，防止前后影响
	students = []Student{}
	db.First(&student, 1)
	fmt.Println(student.Name)
	fmt.Println()
	db.Find(&students, []int{1, 2, 3})
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println("============primary key end============")
	//根据map查询
	student = Student{}
	students = []Student{}
	db.Where(map[string]interface{}{"city": "郑州", "score": 0}).Find(&students)
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println("============map end============")
	//OR查询
	student = Student{}
	students = []Student{}
	db.Where("city=?", "郑州").Or("city=?", "北京").Find(&students)
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println("============or end============")
	//order by
	student = Student{}
	students = []Student{}
	db.Where("city=?", "郑州").Order("score").Find(&students)
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println("============order end============")
	//limit
	student = Student{}
	students = []Student{}
	db.Where("city=?", "郑州").Order("score").Limit(1).Offset(0).Find(&students)
	for _, ele := range students {
		fmt.Printf("id=%d, name=%s\n", ele.Id, ele.Name)
	}
	fmt.Println("============limit end============")
	//选择特定的字段
	student = Student{}
	db.Select("name").Take(&student)                                     //Take从结果中取一个，不保证是第一个或最后一个
	fmt.Printf("name=%s, province=%s\n", student.Name, student.Province) //只select了name，所以province是空的
	fmt.Println("============select specified column end============")
}

func update(db *gorm.DB) {
	//根据where更新一列
	db.Model(&Student{}).Where("city=?", "北京").Update("score", 10)
	//更新多列
	db.Model(&Student{}).Where("city=?", "北京").Updates(map[string]interface{}{"score": 3, "addr": "海淀区"})
	//where里加入object的ID
	student := Student{Id: 2, City: "太原"}
	db = db.Model(&student).Where("city=?", "郑州").Updates(map[string]interface{}{"score": 10, "addr": "中原区"})
	fmt.Printf("update %d rows\n", db.RowsAffected)
	fmt.Println("=============update end=============")
}

func create(db *gorm.DB) {
	//插入一条记录
	student := Student{Name: "光绪", Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
	db.Create(&student)
	//一次性插入多条
	students := []Student{{Name: "无极", Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}, {Name: "小王", Province: "上海", City: "上海", Score: 12, Enrollment: time.Now()}, {Name: "小亮", Province: "北京", City: "北京", Score: 20, Enrollment: time.Now()}}
	db.Create(students)
	//量太大时分批插入
	students = []Student{{Name: "大壮", Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}, {Name: "刘二", Province: "上海", City: "上海", Score: 12, Enrollment: time.Now()}, {Name: "文明", Province: "北京", City: "北京", Score: 20, Enrollment: time.Now()}}
	db = db.CreateInBatches(students, 2) //一次插入2条
	fmt.Printf("insert %d rows\n", db.RowsAffected)
	fmt.Println("=============insert end=============")
}

func delete(db *gorm.DB) {
	//用where删除
	db = db.Where("city in ?", []string{"常州", "成都"}).Delete(&Student{})
	fmt.Printf("delete %d rows\n", db.RowsAffected)
	//用主键删除
	db = db.Delete(&Student{}, 27)
	fmt.Printf("delete %d rows\n", db.RowsAffected)
	db = db.Delete(&Student{}, []int{28, 29, 30})
	fmt.Printf("delete %d rows\n", db.RowsAffected)
	fmt.Println("=============delete end=============")
}

func transaction(db *gorm.DB) {
	tx := db.Begin()
	for i := 0; i < 10; i++ {
		student := Student{Name: "学生" + strconv.Itoa(i), Province: "北京", City: "北京", Score: 38, Enrollment: time.Now()}
		db.Create(&student)
	}
	tx.Commit()
	fmt.Println("=============transaction end=============")
}

func main() {
	//想要正确的处理time.Time ，您需要带上parseTime参数
	//要支持完整的UTF-8编码，您需要将charset=utf8更改为charset=utf8mb4
	dsn := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	database.CheckError(err)
	db.AutoMigrate(&Student{})
	// query(db)
	// update(db)
	// create(db)
	// delete(db)
	transaction(db)
	db.Where("name like ?", "学生%").Delete(&Student{}) //删除name以"学生"为前缀的记录
}

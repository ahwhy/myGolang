package main

import (
	"fmt"
)

// 使用结构体表示班级和学生，请计算每个班级学科平均分
// Student 名称(Name) 学号(Number) 科目(Subjects) 成绩(Scores)
// Class   名称(Name) 编号(Number) 学员(Students)
// Class   实现一个平均值的方法
type Student struct {
	Name     string
	Number   int
	Subjects []string
	Scores   []int
}

type Class struct {
	Name     string
	Number   int
	Students []*Student
}

func main() {
	Class1_avgscores()
	Class2_avgscores()
	Class3_avgscores()
}

// 一班信息录入及各科目平均分
func Class1_avgscores() {
	var class1 Class = Class{
		Name:   "一班",
		Number: 1,
		Students: []*Student{
			{
				Name:     "aa",
				Number:   1,
				Subjects: []string{"语文", "数学", "英语"},
				Scores:   []int{85, 78, 97},
			},
			{
				Name:     "bb",
				Number:   2,
				Subjects: []string{"语文", "数学", "英语"},
				Scores:   []int{79, 87, 81},
			},
		},
	}

	class1Scires := make(map[string][]int)
	for _, v := range class1.Students {
		class1Scires["语文"] = append(class1Scires["语文"], v.Scores[0])
		class1Scires["数学"] = append(class1Scires["数学"], v.Scores[1])
		class1Scires["英语"] = append(class1Scires["英语"], v.Scores[2])
	}

	fmt.Println("一班各科目平均分")
	for k, v := range class1Scires {
		tol := 0
		for _, j := range v {
			tol += j
		}
		fmt.Printf("科目: %s，平均分: %d\n", k, tol/2)
	}
}

// 二班信息录入及各科目平均分
func Class2_avgscores() {
	var class2 = Class{
		Name:   "二班",
		Number: 2,
		Students: []*Student{
			{
				Name:     "cc",
				Number:   1,
				Subjects: []string{"语文", "数学", "英语"},
				Scores:   []int{67, 72, 78},
			},
			{
				Name:     "dd",
				Number:   2,
				Subjects: []string{"语文", "数学", "英语"},
				Scores:   []int{82, 75, 76},
			},
		},
	}

	class2Scires := make(map[string][]int)
	for _, v := range class2.Students {
		class2Scires["语文"] = append(class2Scires["语文"], v.Scores[0])
		class2Scires["数学"] = append(class2Scires["数学"], v.Scores[1])
		class2Scires["英语"] = append(class2Scires["英语"], v.Scores[2])
	}

	fmt.Println("二班各科目平均分")
	for k, v := range class2Scires {
		tol := 0
		for _, j := range v {
			tol += j
		}
		fmt.Printf("科目: %s，平均分: %d\n", k, tol/2)
	}
}

// 三班信息录入及各科目平均分
func Class3_avgscores() {
	class3 := Class{
		Name:   "三班",
		Number: 3,
		Students: []*Student{
			{
				Name:     "ee",
				Number:   1,
				Subjects: []string{"语文", "数学", "英语"},
				Scores:   []int{85, 92, 95},
			},
			{
				Name:     "ff",
				Number:   2,
				Subjects: []string{"语文", "数学", "英语"},
				Scores:   []int{90, 99, 98},
			},
		},
	}

	class3Scires := make(map[string][]int)
	for _, v := range class3.Students {
		class3Scires["语文"] = append(class3Scires["语文"], v.Scores[0])
		class3Scires["数学"] = append(class3Scires["数学"], v.Scores[1])
		class3Scires["英语"] = append(class3Scires["英语"], v.Scores[2])
	}

	fmt.Println("三班各科目平均分")
	for k, v := range class3Scires {
		tol := 0
		for _, j := range v {
			tol += j
		}
		fmt.Printf("科目: %s，平均分: %d\n", k, tol/2)
	}
}

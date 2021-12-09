package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewStudent(t *testing.T) {
	Convey("start test new", t, func() {
		stu, err := NewStudent("")
		Convey("空的name初始化错误", func() {
			So(err, ShouldBeError)
		})
		Convey("stu对象为nil", func() {
			So(stu, ShouldBeNil)
		})
	})
}

func TestScore(t *testing.T) {
	stu, _ := NewStudent("小乙")
	Convey("不设置分数可能出错", t, func() {
		sc, err := stu.GetAvgScore()
		Convey("获取分数出错了", func() {
			So(err, ShouldBeError)
		})
		Convey("分数为0", func() {
			So(sc, ShouldEqual, 0)
		})
	})
	Convey("正常情况", t, func() {
		stu.ChiScore = 60
		stu.EngScore = 70
		stu.MathScore = 80
		score, err := stu.GetAvgScore()
		Convey("获取分数出错了", func() {
			So(err, ShouldBeNil)
		})
		Convey("平均分大于60", func() {
			So(score, ShouldBeGreaterThan, 600)
		})
	})
}
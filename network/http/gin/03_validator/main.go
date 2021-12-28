package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10" // 注意要用新版本v10
)

type Student struct {
	Name       string    `form:"name" binding:"required"`                                                                // required:必须上传name参数
	Score      int       `form:"score" binding:"gt=0"`                                                                   // score必须为正数
	Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`       // 自定义验证before_today，日期格式东8区
	Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"` // 毕业时间要晚于入学时间
}

var beforeToday validator.Func = func(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if date.Before(today) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func processErr(err error) {
	if err == nil {
		return
	}

	// 给Validate.Struct()函数传了一个非法的参数
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		fmt.Println("param error:", invalid)
		return
	}

	// ValidationErrors是一个错误切片，它保存了每个字段违反的每个约束信息
	validationErrs := err.(validator.ValidationErrors)
	for _, validationErr := range validationErrs {
		fmt.Printf("field %s 不满足条件 %s\n", validationErr.Field(), validationErr.Tag())
	}
}

func main() {
	engine := gin.Default()

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("before_today", beforeToday)
	}

	// http://localhost:5656?name=zcy&score=1&enrollment=2021-08-23&graduation=2021-09-23
	engine.GET("/", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBind(&stu); err != nil {
			processErr(err) // 校验不符合时，打印出哪时不符合
			ctx.String(http.StatusBadRequest, "parse parameter failed")
		} else {
			ctx.JSON(http.StatusOK, stu)
		}
	})

	engine.Run(":5656")
}

package validation

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

type RegistRequest struct {
	UserName   string `validate:"gt=0"`             // >0  长度大于0
	PassWord   string `validate:"min=6,max=12"`     // 密码长度[6, 12]
	PassRepeat string `validate:"eqfield=PassWord"` // 跨字段相等校验
	Email      string `validate:"email"`            // 需要满足email的格式
}

func ValidateEmail(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	if pass, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,})\.([a-z]{2,4})$`, input); pass {
		return true
	}

	return false
}

type InnerRequest struct {
	Pass  string `validate:"min=6,max=12"` // 密码长度[6, 12]
	Email string `validate:"my_email"`
}

type OutterRequest struct {
	PassWord   string `validate:"eqcsfield=Nest.Pass"` // 跨结构体相等校验
	PassRepeat string `validate:"eqfield=PassWord"`    // 跨字段相等校验
	Nest       InnerRequest
}

func ProcessErr(err error) {
	if err == nil {
		return
	}

	// 给Validate.Struct()函数传了一个非法的参数
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		fmt.Println("param error:", invalid)
		return
	}

	// ValidationErrors是一个错误切片，保存了每个字段违反的每个约束信息
	validationErrs := err.(validator.ValidationErrors)
	for _, validationErr := range validationErrs {
		fmt.Printf("field %s 不满足条件 %s\n", validationErr.Field(), validationErr.Tag())
	}
}

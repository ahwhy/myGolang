package validation_test

import (
	"testing"

	"github.com/ahwhy/myGolang/network/http/validation"
	"github.com/go-playground/validator"
)

var val = validator.New()

func TestValidation(t *testing.T) {
	req := validation.RegistRequest{
		UserName:   "zcy",
		PassWord:   "12345",
		PassRepeat: "1234568",
		Email:      "123qq.com",
	}

	validation.ProcessErr(val.Struct(req)) // Struct()返回的error分为两种类型: InvalidValidationError和ValidationErrors
	validation.ProcessErr(val.Struct(3))
}

func TestCustValidation(t *testing.T) {
	val.RegisterValidation("my_email", validation.ValidateEmail) // 注册一个自定义的validator

	inreq := validation.InnerRequest{
		Pass:  "123456",
		Email: "123qq.com",
	}
	outreq := validation.OutterRequest{
		PassWord:   "123456",
		PassRepeat: "1234568",
		Nest:       inreq,
	}

	validation.ProcessErr(val.Struct(inreq))
	validation.ProcessErr(val.Struct(outreq))
}

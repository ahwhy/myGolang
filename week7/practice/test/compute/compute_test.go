package main

import "testing"

func TestAdd(t *testing.T) {
	a := 10
	b := 20
	want := 30

	actual := Add(a, b)
	if want != actual {
		t.Errorf("[Add 函数参数 :%d %d][期望：%d][实际：%d]", a, b, want, actual)
	}
}

func TestMul(t *testing.T) {
	a := 10
	b := 20
	want := 300

	actual := Mul(a, b)
	if want != actual {
		t.Errorf("[Mul 函数参数 :%d %d][期望：%d][实际：%d]", a, b, want, actual)
	}
}

func TestDiv(t *testing.T) {
	a := 10
	b := 20
	want := 2

	actual := Div(a, b)
	if want != actual {
		t.Errorf("[Div 函数参数 :%d %d][期望：%d][实际：%d]", a, b, want, actual)
	}
}

func TestMul2(t *testing.T) {
	t.Run("正数", func(t *testing.T) {
		if Mul(4, 5) != 20 {
			t.Fatal("muli.zhengshu.error")
		}
	})

	t.Run("负数", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("muli.fusshu.error")
		}
	})

}

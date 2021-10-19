package custmath

import "fmt"

var version string = "v1.0"

func init() {
	fmt.Println("Custmath Version:", version)
}

func Add(a int, b int) int {
	return a + b
}

func Sub(a int, b int) int {
	return a - b
}

func Mul(a int, b int) int {
	return a * b
}

func Div(a int, b int) int {
	return a / b
}

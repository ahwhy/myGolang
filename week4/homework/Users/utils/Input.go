package utils

import (
	"fmt"
)

func Input(prompt string) string {
	var text string
	fmt.Print(prompt)
	fmt.Scanln(&text)
	return text
}

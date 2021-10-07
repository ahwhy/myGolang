package utils_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/utils"
)

func TestBinaryFormat(t *testing.T) {
	fmt.Println(utils.BinaryFormat(4095))
}

func TestHumanBytesLoaded(t *testing.T) {
	fmt.Println(utils.HumanBytesLoaded(1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024 * 1024 * 1024))
}

func TestMd5(t *testing.T) {
	fmt.Println(utils.Md5("Atlantis"))
}

func TestMd5Salt(t *testing.T) {
	fmt.Println(utils.Md5Salt("Atlantis", ""))
}

func TestSplitMd5Salt(t *testing.T) {
	salt, md5 := utils.SplitMd5Salt(utils.Md5Salt("Atlantis", ""))
	fmt.Printf("Salt: %s \nMd5: %s\n", salt, md5)
}

func TestRandString(t *testing.T) {
	fmt.Println(utils.RandString(10))
}

func TestSnake(t *testing.T) {
	fmt.Println(utils.Snake("qwer1234ASDF"))
}

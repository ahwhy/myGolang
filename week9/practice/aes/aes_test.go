package aes_test

import (
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/aes"
)

// The text need to be encrypt. ==> 4ff9b0b21962a62da8678efebd33e37eee2caad9d69933939d93dad1
// 4ff9b0b21962a62da8678efebd33e37eee2caad9d69933939d93dad1 ==> The text need to be encrypt.
func TestAES(t *testing.T) {
	plain := "The text need to be encrypt."
	// AES 规定有3种长度的key: 16, 24, 32分别对应AES-128, AES-192, or AES-256
	key := "abcdefgehjhijkmlkjjwwoew"

	// 加密
	cipherByte, err := aes.Encrypt(plain, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %x\n", plain, cipherByte)

	// 解密
	plainText, err := aes.Decrypt(cipherByte, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%x ==> %s\n", cipherByte, plainText)
}

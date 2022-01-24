package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"log"
)

//XOR 异或运算，要求plain和key的长度相同
func XOR(plain string, key []byte) string {
	bPlain := []byte(plain)
	bCipher := make([]byte, len(key))
	for i, k := range key {
		bCipher[i] = k ^ bPlain[i]
	}
	cipher := string(bCipher)
	return cipher
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0填充
	return append(ciphertext, padtext...)
}

//ZeroUnPadding 这种方法不严谨，末尾的0不一定全是padding出来的
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0) //截掉尾部连续的0
		})
}

//DesEncrypt DES加密
//密钥必须是64位，所以key必须是长度为8的byte数组
func DesEncrypt(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key) //用des创建一个加密器cipher
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()    //分组的大小，blockSize=8
	src = ZeroPadding(src, blockSize) //填充

	out := make([]byte, len(src)) //密文和明文的长度一致
	dst := out
	for len(src) > 0 {
		//分组加密
		block.Encrypt(dst, src[:blockSize]) //对src进行加密，加密结果放到dst里
		//移到下一组
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	return hex.EncodeToString(out), nil
}

//DesDecrypt DES解密
//密钥必须是64位，所以key必须是长度为8的byte数组
func DesDecrypt(text string, key []byte) (string, error) {
	src, err := hex.DecodeString(text) //转成[]byte
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		//分组解密
		block.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	out = ZeroUnPadding(out) //反填充
	return string(out), nil
}

func DesEncryptCBC(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key) //用des创建一个加密器cipher
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()    //分组的大小，blockSize=8
	src = ZeroPadding(src, blockSize) //填充

	out := make([]byte, len(src))                   //密文和明文的长度一致
	encrypter := cipher.NewCBCEncrypter(block, key) //CBC分组模式加密
	encrypter.CryptBlocks(out, src)
	return hex.EncodeToString(out), nil
}

func DesDecryptCBC(text string, key []byte) (string, error) {
	src, err := hex.DecodeString(text) //转成[]byte
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	out := make([]byte, len(src))                   //密文和明文的长度一致
	encrypter := cipher.NewCBCDecrypter(block, key) //CBC分组模式解密
	encrypter.CryptBlocks(out, src)
	out = ZeroUnPadding(out) //反填充
	return string(out), nil
}

func main1() {
	plain := "ABCD"
	key := []byte{1, 2, 3, 4}
	cipher := XOR(plain, key)
	plain = XOR(cipher, key)
	fmt.Printf("plain: %s\n", plain)
	fmt.Println("-------------------------------------")

	key = []byte("ir489u58") //key必须是长度为8的byte数组
	plain = "因为我们没有什么不同"
	cipher, err := DesEncrypt(plain, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("密文：%s\n", cipher)

	plain, err = DesDecrypt(cipher, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("明文：%s\n", plain)
	fmt.Println("-------------------------------------")

	cipher, _ = DesEncryptCBC(plain, key)
	fmt.Printf("密文：%s\n", cipher)
	plain, err = DesDecryptCBC(cipher, key)
	fmt.Printf("明文：%s\n", plain)
}

//go  encryption/des.go

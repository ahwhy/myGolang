package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto/ecies" // 以太坊加密库，要求go版本升级到1.15
)

func genPrivateKey() (*ecies.PrivateKey, error) {
	pubkeyCurve := elliptic.P256() // 初始化椭圆曲线
	// 随机挑选基点,生成私钥
	p, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // 用golang标准库生成公私钥对
	if err != nil {
		return nil, err
	} else {
		return ecies.ImportECDSA(p), nil // 转换成以太坊的公私钥对
	}
}

//ECCEncrypt 椭圆曲线加密
func ECCEncrypt(plain string, pubKey *ecies.PublicKey) ([]byte, error) {
	src := []byte(plain)
	return ecies.Encrypt(rand.Reader, pubKey, src, nil, nil)
}

//ECCDecrypt 椭圆曲线解密
func ECCDecrypt(cipher []byte, prvKey *ecies.PrivateKey) (string, error) {
	if src, err := prvKey.Decrypt(cipher, nil, nil); err != nil {
		return "", err
	} else {
		return string(src), nil
	}
}

func main() {
	prvKey, err := genPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	pubKey := prvKey.PublicKey
	plain := "因为我们没有什么不同"
	cipher, err := ECCEncrypt(plain, &pubKey)
	if err != nil {
		log.Fatal(err)
	}
	plain, err = ECCDecrypt(cipher, prvKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("明文：%s\n", plain)
}

//go get github.com/ethereum/go-ethereum 要求go版本升级到1.15才能安装成功
//go run encryption/ecc.go

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

/*
生成1024位的RSA私钥：
openssl genrsa -out data/rsa_private_key.pem 1024
根据私钥生成公钥：
openssl rsa -in data/rsa_private_key.pem -pubout -out data/rsa_public_key.pem

pem是一种标准格式，它通常包含页眉和页脚
*/

var (
	publicKey  []byte
	privateKey []byte
)

func ReadFile(keyFile string) ([]byte, error) {
	if f, err := os.Open(keyFile); err != nil {
		return nil, err
	} else {
		content := make([]byte, 4096)
		if n, err := f.Read(content); err != nil {
			return nil, err
		} else {
			return content[:n], nil
		}
	}
}

func ReadRSAKey(publicKeyFile, privateKeyFile string) {
	publicKey, _ = ReadFile(publicKeyFile)
	privateKey, _ = ReadFile(privateKeyFile)
}

// RSA加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //目前的数字证书一般都是基于ITU（国际电信联盟）制定的X.509标准
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RSA解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
func main5() {
	ReadRSAKey("data/rsa_public_key.pem", "data/rsa_private_key.pem")

	plain := "因为我们没有什么不同"
	cipher, _ := RsaEncrypt([]byte(plain))
	fmt.Printf("密文：%s\n", hex.EncodeToString(cipher))
	bPlain, _ := RsaDecrypt(cipher)
	fmt.Printf("明文：%s\n", string(bPlain))
}

//go run encryption/rsa.go

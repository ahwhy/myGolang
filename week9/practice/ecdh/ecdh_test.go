package ecdh_test

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ahwhy/myGolang/week9/practice/ecdh"
)

func test(e ecdh.ECDH) {
	var privKey1, privKey2 crypto.PrivateKey
	var pubKey1, pubKey2 crypto.PublicKey
	var pubKey1Buf, pubKey2Buf []byte
	var err error
	var ok bool
	var secret1, secret2 []byte
	// 准备2对秘钥对,A: privKey1,pubKey1 B:privKey2,pubKey2
	privKey1, pubKey1, err = e.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	privKey2, pubKey2, err = e.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	pubKey1Buf = e.Marshal(pubKey1)
	pubKey2Buf = e.Marshal(pubKey2)
	pubKey1, ok = e.Unmarshal(pubKey1Buf)
	if !ok {
		fmt.Println("Unmarshal does not work")
	}
	pubKey2, ok = e.Unmarshal(pubKey2Buf)
	if !ok {
		fmt.Println("Unmarshal does not work")
	}
	// A 通过B给的公钥协商共享密码
	secret1, err = e.GenerateSharedSecret(privKey1, pubKey2)
	if err != nil {
		fmt.Println(err)
	}
	// B 通过A给的公钥协商共享密码
	secret2, err = e.GenerateSharedSecret(privKey2, pubKey1)
	if err != nil {
		fmt.Println(err)
	}
	// A B在没暴露直接的私钥的情况下, 协商出了一个共享密码
	fmt.Printf("The secret1 shared keys: %x\n", secret1)
	fmt.Printf("The secret2 shared keys: %x\n", secret2)
}
func TestECDH(t *testing.T) {
	e1 := ecdh.NewEllipticECDH(elliptic.P521())
	e2 := ecdh.NewCurve25519ECDH()
	test(e1)
	test(e2)
}

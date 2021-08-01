# 秘钥交换算法

一种密钥交换协议，注意该算法只能用于密钥的交换，而不能进行消息的加密和解密。双方确定要用的密钥后，要使用其他对称密钥操作加密算法实际加密和解密消息。它可以让双方在不泄漏密钥的情况下协商出一个密钥来, 常用于保证对称加密的秘钥的安全, TLS就是这样做的。
在这个领域应该2种

+ DH：ECDH是DH的加强版
+ ECDH: DH算法的加强版, 常用的是NIST系列,但是后面curve25519
+ curve25519: 实质上也是一种ECDH,但是其实现更为优秀,表现的更为安全,可能是下一代秘钥交换算法的标准。


DH全称是:Diffie-Hellman, 是一种确保共享KEY安全穿越不安全网络的方法，它是OAKLEY的一个组成部分。Whitefield与Martin Hellman在1976年提出了一个奇妙的密钥交换协议，称为Diffie-Hellman密钥交换协议/算法(Diffie-Hellman Key Exchange/Agreement Algorithm).这个机制的巧妙在于需要安全通信的双方可以用这个方法确定对称密钥。然后可以用这个密钥进行加密和解密。
DH依赖于计算离散对数的难度, 大概过程如下:

```
可以如下定义离散对数：首先定义一个素数p的原根，为其各次幂产生从1 到p-1的所有整数根，也就是说，如果a是素数p的一个原根，那么数值 a mod p,a2 mod p,…,ap-1 mod p 是各不相同的整数，并且以某种排列方式组成了从1到p-1的所有整数. 对于一个整数b和素数p的一个原根a，可以找到惟一的指数i，使得 b = a^i mod p 其中0 ≤ i ≤ （p-1） 指数i称为b的以a为基数的模p的离散对数或者指数.该值被记为inda,p(b).
```


ECDH 全称是Elliptic Curve Diffie-Hellman, 是DH算法的加强版, 基于椭圆曲线难题加密, 现在是主流的密钥交换算法

这里需要指出下golang的标准库的crypto里的椭圆曲线实现了这4种(elliptic文档): P224/P256/P384/P521, 而curve25519是单独实现的, 他不在标准库中: golang.org/x/crypto/curve25519

```go
package main
import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"golang.org/x/crypto/curve25519"
)
// ECDH 秘钥交换算法的主接口
type ECDH interface {
	GenerateKey(io.Reader) (crypto.PrivateKey, crypto.PublicKey, error)
	Marshal(crypto.PublicKey) []byte
	Unmarshal([]byte) (crypto.PublicKey, bool)
	GenerateSharedSecret(crypto.PrivateKey, crypto.PublicKey) ([]byte, error)
}
type ellipticECDH struct {
	ECDH
	curve elliptic.Curve
}
type ellipticPublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}
type ellipticPrivateKey struct {
	D []byte
}
// NewEllipticECDH 指定一种椭圆曲线算法用于创建一个ECDH的实例
// 关于椭圆曲线算法标准库里面实现了4种: 见crypto/elliptic
func NewEllipticECDH(curve elliptic.Curve) ECDH {
	return &ellipticECDH{
		curve: curve,
	}
}
// GenerateKey 基于标准库的NIST椭圆曲线算法生成秘钥对
func (e *ellipticECDH) GenerateKey(rand io.Reader) (crypto.PrivateKey, crypto.PublicKey, error) {
	var d []byte
	var x, y *big.Int
	var priv *ellipticPrivateKey
	var pub *ellipticPublicKey
	var err error
	d, x, y, err = elliptic.GenerateKey(e.curve, rand)
	if err != nil {
		return nil, nil, err
	}
	priv = &ellipticPrivateKey{
		D: d,
	}
	pub = &ellipticPublicKey{
		Curve: e.curve,
		X:     x,
		Y:     y,
	}
	return priv, pub, nil
}
// Marshal用于公钥的序列化
func (e *ellipticECDH) Marshal(p crypto.PublicKey) []byte {
	pub := p.(*ellipticPublicKey)
	return elliptic.Marshal(e.curve, pub.X, pub.Y)
}
// Unmarshal用于公钥的反序列化
func (e *ellipticECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	var key *ellipticPublicKey
	var x, y *big.Int
	x, y = elliptic.Unmarshal(e.curve, data)
	if x == nil || y == nil {
		return key, false
	}
	key = &ellipticPublicKey{
		Curve: e.curve,
		X:     x,
		Y:     y,
	}
	return key, true
}
// GenerateSharedSecret 通过自己的私钥和对方的公钥协商一个共享密码
func (e *ellipticECDH) GenerateSharedSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	priv := privKey.(*ellipticPrivateKey)
	pub := pubKey.(*ellipticPublicKey)
	x, _ := e.curve.ScalarMult(pub.X, pub.Y, priv.D)
	return x.Bytes(), nil
}
// NewCurve25519ECDH 使用密码学家Daniel J. Bernstein的椭圆曲线算法:Curve25519来创建ECDH实例
// 因为Curve25519独立于NIST之外, 没在标准库实现, 需要单独为期实现一套接口来支持ECDH
func NewCurve25519ECDH() ECDH {
	return &curve25519ECDH{}
}
type curve25519ECDH struct {
	ECDH
}
// GenerateKey 基于curve25519椭圆曲线算法生成秘钥对
func (e *curve25519ECDH) GenerateKey(rand io.Reader) (crypto.PrivateKey, crypto.PublicKey, error) {
	var pub, priv [32]byte
	var err error
	_, err = io.ReadFull(rand, priv[:])
	if err != nil {
		return nil, nil, err
	}
	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64
	curve25519.ScalarBaseMult(&pub, &priv)
	return &priv, &pub, nil
}
// 实现公钥的序列化
func (e *curve25519ECDH) Marshal(p crypto.PublicKey) []byte {
	pub := p.(*[32]byte)
	return pub[:]
}
// 实现公钥的反序列化
func (e *curve25519ECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	var pub [32]byte
	if len(data) != 32 {
		return nil, false
	}
	copy(pub[:], data)
	return &pub, true
}
// 实现秘钥协商接口
func (e *curve25519ECDH) GenerateSharedSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	var priv, pub, secret *[32]byte
	priv = privKey.(*[32]byte)
	pub = pubKey.(*[32]byte)
	secret = new([32]byte)
	curve25519.ScalarMult(secret, priv, pub)
	return secret[:], nil
}
func test(e ECDH) {
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
func main() {
	e1 := NewEllipticECDH(elliptic.P521())
	e2 := NewCurve25519ECDH()
	test(e1)
	test(e2)
}
```
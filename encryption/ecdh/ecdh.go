package ecdh

import (
	"crypto"
	"crypto/elliptic"
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

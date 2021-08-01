package hash_test

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	hashData = "create mode 100645 image/cbc-encrypto.jpg"
)

// 4ca3c744679506f7ffa9c7b9fd81ea9d
// 319ea33c0dfd741de0370f5c9418f9d7
func TestMD5Hash(t *testing.T) {
	m := md5.New()
	_, err := m.Write([]byte("create mode 100645 image/cbc-encrypto.jpg"))
	if err != nil {
		panic(err)
	}

	hd := m.Sum(nil)
	fmt.Printf("%x\n", hd)
}

// 5d8f6833d1bf27ce0e2e1f72e9047cf13f07d682
func TestSHA1Hash(t *testing.T) {
	s := sha1.New()
	// []byte -->  string  置换表<> code_point --->  字符， 编码表
	//

	// map[hashV]plainText   7c4a8d09ca3762af61e59520943dc26494f8941b -->  123456
	// 26 *2 + spec
	_, err := io.WriteString(s, "123456")
	if err != nil {
		panic(err)
	}

	hd := s.Sum(nil)
	fmt.Printf("%x", hd)
}

// abc --> c2b520d835e6a2e3f8f1ecc9fd57d794b3a9d2ae
// key：c2b520d835e6a2e3f8f1ecc9fd57d794b3a9d2ae
// md5
// sha1
// raw:  create mode 100645 image/cbc-encrypto.jpg
// sha1: 5d8f6833d1bf27ce0e2e1f72e9047cf13f07d682
// +key：c2b520d835e6a2e3f8f1ecc9fd57d794b3a9d2ae   A公司
// +key: 4b983170ca50abad2e4af6d72c23c612405b9b59   B公司

//  abc c2b520d835e6a2e3f8f1ecc9fd57d794b3a9d2ae
//  meq 4b983170ca50abad2e4af6d72c23c612405b9b59

// 
func TestHMAC(t *testing.T) {
	h := hmac.New(sha1.New, []byte("meq"))
	io.WriteString(h, hashData)
	fmt.Printf("%x", h.Sum(nil))
}

// PasswordCrypto
//  raw:  create mode 100645 image/cbc-encrypto.jpg
//        $2a$10$nKFnu5NCNrkruBCSNMoBQudD5PZ2FKzngs2WhIjQTznFdvkGEXhN.
//        $2a$10$D7xpsVGjK06zINxcAggFoOkCHT1pK0UmzNnqFmQPeggj3x.t.E8oG
func TestBcrypto(t *testing.T) {
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now).Milliseconds())
	}()

	hd, _ := bcrypt.GenerateFromPassword([]byte(hashData), 28)
	fmt.Printf("%s\n", hd)
}

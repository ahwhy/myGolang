package cbc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ahwhy/myGolang/week9/practice/cbc"
)

// 79e04f5f8817c77123b4efb53beb83628294627e9a840dff5ce9985bd57267d333a0670d1e183737798ad18f9943c3b9
func TestAESCBC(t *testing.T) {
	data := []byte("this is my plainText")
	key := []byte("123456")

	should := require.New(t)

	cipherData, err := cbc.Encrypt(data, key)
	should.NoError(err)
	t.Logf("cipher data: %x\n", cipherData)

	rawData, err := cbc.Decrypt(cipherData, key)
	should.NoError(err)
	t.Logf("raw data: %s", rawData)

	should.Equal(data, []byte(rawData))
}

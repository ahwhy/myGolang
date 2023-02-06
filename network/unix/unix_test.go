package unix_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ahwhy/myGolang/network/unix"
)

func TestDecryptFromProxySSL(t *testing.T) {
	// 初始化client
	cli := unix.InitProxySSLClient()
	secret := os.Getenv("testing")

	decryptData, err := cli.DecryptFromProxySSL([]byte(secret))
	if err != nil {
		panic(err)
	}
	fmt.Print(decryptData)
}

package unix

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	ProxySSLAddress string = "/var/pssl.sockrest"
	ProxySSLURL     string = "http://proxyssl/configs"
)

// ProxySSL-Client
var (
	transportOnce     sync.Once
	initOnceTransport *http.Transport
	DialerTimeout     = 60 * time.Second
	ResponseTimeout   = 60 * time.Second
	ClientMaxIdle     = 5
)

type client struct {
	proxySSLClient *http.Client
}

// InitProxySSLClient 初始化 ProxySSL 客户端
func InitProxySSLClient() *client {
	// 创建 http客户端, 负责与 proxySSL 进行交互
	address := ProxySSLAddress
	transportOnce.Do(func() {
		initOnceTransport = defaultTransport(address)
	})
	httpCilent := &http.Client{
		Transport: initOnceTransport,
		Timeout:   ResponseTimeout,
	}

	return &client{
		proxySSLClient: httpCilent,
	}
}

func defaultTransport(address string) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: DialerTimeout,
			}
			return d.Dial("unix", address)
		},
		MaxIdleConns:          ClientMaxIdle,
		ResponseHeaderTimeout: ResponseTimeout,
	}
}

func (c *client) DecryptFromProxySSL(data []byte) ([]byte, error) {
	// 创建 http 请求
	req, err := http.NewRequest("POST", ProxySSLURL, bytes.NewBuffer(data))
	if err != nil {
		return []byte{}, err
	}

	// 发送 http 请求，获取 response
	resp, err := c.proxySSLClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	// 读取解密后的内容
	decryptedData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode/100 != 2 {
		return []byte{}, err
	}

	return decryptedData, nil
}
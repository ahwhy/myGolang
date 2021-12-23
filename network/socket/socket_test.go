package socket_test

import (
	"testing"

	"github.com/ahwhy/myGolang/network/socket/client"
	"github.com/ahwhy/myGolang/network/socket/server"
)

func TestTCPServer_SimpleMessage(t *testing.T) {
	server.TCPServer_SimpleMessage()
}

func TestTCPClient_SimpleMessage(t *testing.T) {
	client.TCPClient_SimpleMessage()
}

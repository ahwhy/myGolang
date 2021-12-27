package client_test

import (
	"testing"

	"github.com/ahwhy/myGolang/network/http/client"
)

func TestSimpleGet(t *testing.T) {
	client.SimpleGet()
}

func TestSimplePost(t *testing.T) {
	client.SimplePost()
}

func TestPostForm(t *testing.T) {
	client.PostForm()
}

func TestHead(t *testing.T) {
	client.Head()
}

func TestComplexRequest(t *testing.T) {
	client.ComplexRequest()
}

func TestRestful(t *testing.T) {
	client.Restful()
}

func TestRequestPanic(t *testing.T) {
	client.RequestPanic()
}

func TestRequestBook(t *testing.T) {
	client.RequestBook()
}

func TestRequestFood(t *testing.T) {
	client.RequestFood()
}

func TestMiddleWare(t *testing.T) {
	client.MiddleWare()
}

func TestAuthLogin(t *testing.T) {
	client.AuthLogin()
}

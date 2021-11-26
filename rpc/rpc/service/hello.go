package service

import (
	"bytes"
	"encoding/gob"
)

const HelloServiceName = "HelloService"

// 定义hello service的接口
type HelloService interface {
	Hello(request string, reply *string) error
}

func GobEncode(val interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(val); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func GobDecode(data []byte, value interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(value)
}

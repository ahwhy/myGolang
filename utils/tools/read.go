package tools

import (
	"encoding/json"
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

func ReadFile(path string) (string, error) {
	fs, err := os.Open(path)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(fs)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ReadContentFile(filepath string) ([]byte, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	payload, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func MustReadContentFile(filepath string) string {
	content, err := ReadContentFile(filepath)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func ReadJsonFile(filepath string, v interface{}) error {
	content, err := ReadContentFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, v)
}

func MustReadJsonFile(filepath string, v interface{}) {
	err := ReadJsonFile(filepath, v)
	if err != nil {
		panic(err)
	}
}

func ReadYamlFile(filepath string, v interface{}) error {
	content, err := ReadContentFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, v)
}

func MustReadYamlFile(filepath string, v interface{}) {
	err := ReadYamlFile(filepath, v)
	if err != nil {
		panic(err)
	}
}

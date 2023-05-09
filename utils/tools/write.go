package tools

import (
	"encoding/json"

	"sigs.k8s.io/yaml"
)

func MustToYaml(v interface{}) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func Prettify(i interface{}) string {
	resp, _ := json.MarshalIndent(i, "", "   ")
	
	return string(resp)
}

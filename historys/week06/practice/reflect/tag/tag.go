package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Person struct {
	Name string `json:"name" yaml:"yaml_name" mage:"name"`
	Age  int    `json:"age"  yaml:"yaml_age"  mage:"age"`
	City string `json:"-" yaml:"yaml_city" mage:"-"`
}

//json解析
func jsonWork() {
	// 对象marshal成字符串
	p := Person{
		Name: "xiaoyi",
		Age:  18,
		City: "北京",
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Printf("[json.marshal.err][err:%v]", err)
		return
	}
	log.Printf("[person.marshal.res][res:%v]", string(data))

	// 从字符串解析成结构体
	p2Str := `
   {
    "name":"李逵",
    "age":28,
    "city":"山东"
}
`
	var p2 Person
	err = json.Unmarshal([]byte(p2Str), &p2)
	if err != nil {
		log.Printf("[json.unmarshal.err][err:%v]", err)
		return
	}
	log.Printf("[person.unmarshal.res][res:%v]", p2)

}

// yaml读取文件
func yamlWork() {
	filename := "a.yaml"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("[ioutil.ReadFile.err][err:%v]", err)
		return
	}
	p := &Person{}
	//err = yaml.Unmarshal(content, p)
	err = yaml.UnmarshalStrict(content, p)
	if err != nil {
		log.Printf("[yaml.UnmarshalStrict.err][err:%v]", err)
		return
	}
	log.Printf("[yaml.UnmarshalStrict.res][res:%v]", p)
}

// 自定义标签
func jiexizidingyibiaoqian(s interface{}) {
	// typeOf type类型
	r := reflect.TypeOf(s)
	value := reflect.ValueOf(s)
	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		key := field.Name
		if tag, ok := field.Tag.Lookup("mage"); ok {
			if tag == "-" {
				continue
			}
			log.Printf("[找到了mage标签][key:%v][value:%v][标签：mage=%s]",
				key,
				value.Field(i),
				tag,
			)
		}
	}
}

func main() {
	jsonWork()
	yamlWork()
	p := Person{
		Name: "xiaoyi",
		Age:  18,
		City: "北京",
	}
	jiexizidingyibiaoqian(p)
}

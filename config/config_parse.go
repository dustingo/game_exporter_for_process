package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// MyConfig config结构体 ，对应yaml的process_name
type MyConfig struct {
	Processnames []Info `yaml:"process_names"`
}

// Info 结构体，对应process_names下的-name和cmdline
type Info struct {
	Name    string   `yaml:"name"`
	Cmdline []string `yaml:"cmdline"`
}

// GetConfig 解析yaml，返回myconfig结构体指针
func GetConfig(f *string) *MyConfig {
	var myconfig = new(MyConfig)
	yamlInfo, _ := ioutil.ReadFile(*f)
	err := yaml.Unmarshal(yamlInfo, myconfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	return myconfig
}

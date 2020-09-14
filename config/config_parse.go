package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
func GetConfig() (*MyConfig, error) {
	var fileName string
	for _,para := range os.Args{
		if strings.HasPrefix(para,"--config.path"){
			fmt.Println(para)
			fmt.Printf("Type:%T\n",para)
			fileName = string(strings.Split(para,"=")[1])
		}
	}
	if fileName == ""{
		fileName = "./gameprocess.yaml"
	}
	var myconfig = new(MyConfig)
	yamlInfo, _ := ioutil.ReadFile(fileName)
	err := yaml.Unmarshal(yamlInfo, myconfig)
	return myconfig, err
}

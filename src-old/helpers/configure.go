package helpers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	WFGroups []int64 `yaml:"WFGroups"`
	XDGroups []int64 `yaml:"XDGroups"`
}

func LoadConfig() Config {
	funcName := "configure.go: LoadConfig"
	config := Config{}
	configFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		AddLog(funcName, "readfile", err)
		return config
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		AddLog(funcName, "unmarshal yaml", err)
		return config
	}
	return config
}

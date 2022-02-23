package config

import (
	"embed"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

var (
	CONFIG *config.Config
)

func init() {
	CONFIG = LoadConfig()
}

//go:embed *.yaml
var f embed.FS

func LoadConfig() *config.Config {
	conf := config.New("config")
	conf.AddDriver(yamlv3.Driver)

	data, err := f.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var confData interface{}
	err = yamlv3.Decoder(data, &confData)
	if err != nil {
		panic(err)
	}

	err = conf.LoadData(confData)
	if err != nil {
		panic(err)
	}

	return conf
}

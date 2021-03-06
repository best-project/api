package config

import (
	"github.com/vrischmann/envconfig"
)

// Config provide helm broker configuration
// Supported tags:
//	- json: 		github.com/ghodss/yaml
//	- envconfig: 	github.com/vrischmann/envconfig
//	- default: 		github.com/mcuadros/go-defaults
//	- valid         github.com/asaskevich/govalidator
// Example of valid tag: `valid:"alphanum,required"`
// Combining many tags: tags have to be separated by WHITESPACE: `json:"port" default:"8080" valid:"required"`
type Config struct {
	Port string

	FbAppKey    string
	FbAppSecret string

	DbHost string
	DbUser string
	DbPass string
	DbName string
	DbPort string

	PassPercent float32 `envconfig:"optional"`
	InitDB      bool    `envconfig:"optional" default:"false"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Init(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

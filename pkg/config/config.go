package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/configor"
)

type Config struct {
	ServerURL         string `yaml:"server_url" default:"https://notes.example.com" validate:"url"`
	Listen            string `yaml:"listen" default:"127.0.0.1:5000" validate:"hostname_port"`
	AllowCreateApiKey bool   `yaml:"allow_create_api_key" default:"true"`
	AllowChangeApiKey bool   `yaml:"allow_change_api_key" default:"false"`
	Data              string `yaml:"data" default:"./data"`
}

func LoadFromFile(filename string) (*Config, error) {
	cfg := Config{}

	var err error

	if err = configor.Load(&cfg, filename); err != nil {
		return nil, err
	}
	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

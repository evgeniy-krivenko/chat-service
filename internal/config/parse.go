package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	var conf Config
	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("decode file: %v", err)
	}

	err = validator.Validator.Struct(conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

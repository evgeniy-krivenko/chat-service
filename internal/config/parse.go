package config

import (
	"github.com/BurntSushi/toml"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
	"os"
)

func ParseAndValidate(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var conf Config
	_, err = toml.NewDecoder(file).Decode(&conf)
	if err != nil {
		return Config{}, err
	}

	err = validator.Validator.Struct(conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

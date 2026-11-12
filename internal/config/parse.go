package config

import (
	"github.com/BurntSushi/toml"

	"github.com/dvadert/chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return Config{}, err
	}

	if err := validator.Validator.Struct(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

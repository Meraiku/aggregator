package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {

	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c Config

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (cfg *Config) SetUser(userName string) error {

	cfg.CurrentUserName = userName

	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}

package config

import (
	"encoding/json"
	"os"
)

const configName = "/.gatorconfig.json"

func getConfigFilePath() (string, error) {

	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path + configName, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

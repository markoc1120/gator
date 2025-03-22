package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const (
	configFilePath = "/.gatorconfig.json"
	filePermission = 0644
)

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't get working directory: %v", err)
	}
	return filepath.Join(homeDir, configFilePath), nil
}

func write(config *Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("couldn't parse Go struct to JSON: %v", err)
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, data, filePermission)
	if err != nil {
		return fmt.Errorf("couldn't update/write to the config JSON: %v", err)
	}
	return nil
}

func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error in reading %v: %v", configPath, err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("error in parsing the json: %v", err)
	}
	return &config, nil
}

func SetUser(username string, config *Config) error {
	config.CurrentUser = username
	err := write(config)
	return err
}

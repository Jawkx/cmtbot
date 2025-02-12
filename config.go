package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	ApiBase   string `toml:"api_base"`
	ApiKeyEnv string `toml:"api_key_env"`
	ModelName string `toml:"model_name"`
	NumOfMsg  int    `toml:"num_of_msg"`
	Prompt    string `toml:"prompt"`
}

func LoadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "cmtbot", "cmtbot.toml")

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return Config{
			ApiBase:   "https://openrouter.ai/api/v1/chat/completions",
			ApiKeyEnv: "OPENROUTER_API_KEY",
			ModelName: "google/gemini-flash-1.5",
			NumOfMsg:  5,
		}, nil
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return config, nil
}

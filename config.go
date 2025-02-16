package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ConfigFile struct {
	ApiBase    string `toml:"api_base"`
	ApiKeyEnv  string `toml:"api_key_env"`
	ModelName  string `toml:"model_name"`
	NumOfMsg   int    `toml:"num_of_msg"`
	PromptFile string `toml:"prompt_filename"`
}

type Config struct {
	ApiBase   string
	ApiKeyEnv string
	ModelName string
	NumOfMsg  int
	Prompt    string
}

func LoadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "cmtbot", "config.toml")

	_, err = os.Stat(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("please provide a config file")
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var configToml ConfigFile
	err = toml.Unmarshal(configFile, &configToml)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	promptPath := filepath.Join(
		homeDir,
		".config",
		"cmtbot",
		configToml.PromptFile,
	) // Construct prompt path

	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading prompt file: %w", err)
	}

	config := Config{
		ApiBase:   configToml.ApiBase,
		ApiKeyEnv: configToml.ApiKeyEnv,
		ModelName: configToml.ModelName,
		NumOfMsg:  configToml.NumOfMsg,
		Prompt:    string(promptContent),
	}

	return config, nil
}

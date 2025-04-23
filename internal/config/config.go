package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName     string `yaml:"AppName"`
	AppPort     string `yaml:"AppPort"`
	DbHost      string `yaml:"DbHost"`
	DbPort      string `yaml:"DbPort"`
	DbAdminUser string `yaml:"DbAdminUser"`
	DbAdminDB   string `yaml:"DbAdminDB"`
	DbUser      string `yaml:"DbUser"`
	DbPassword  string `yaml:"DbPassword"`
	DbName      string `yaml:"DbName"`
}

func LoadConfig(filename string) (*Config, error) {
	// Get the directory of the caller (server.go)
	_, callerFile, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("failed to get caller information")
	}
	dir := filepath.Dir(callerFile)
	configPath := filepath.Join(dir, filename)

	// Read the config file
	config := &Config{}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %v", configPath, err)
	}

	// Unmarshal YAML data
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Log the loaded configuration
	log.Printf("Loaded configuration: %+v", config)

	// Validate required fields
	if config.DbAdminUser == "" {
		return nil, fmt.Errorf("DbAdminUser is required but was empty")
	}
	if config.DbAdminDB == "" {
		return nil, fmt.Errorf("DbAdminDB is required but was empty")
	}

	return config, nil
}

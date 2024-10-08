package config

import (
	"orchestrator/internal/constants"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Logging struct {
		Level     string `yaml:"level"`
		File      string `yaml:"file"`
		ToConsole bool   `yaml:"to_console"`
		DebugMode bool   `yaml:"debug_mode"`
	} `yaml:"logging"`
}

func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	// Set default values
	config.Server.Port = constants.DefaultServerPort
	config.Database.Path = constants.DefaultDBPath
	config.Logging.Level = constants.DefaultLogLevel
	config.Logging.File = constants.DefaultLogFile
	config.Logging.ToConsole = constants.DefaultLogToConsole
	config.Logging.DebugMode = constants.DefaultDebugMode

	// Read the config file
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Parse the config file
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

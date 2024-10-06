package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Jobs struct {
		DefaultTimeout time.Duration `yaml:"default_timeout"`
		MaxRetries     int           `yaml:"max_retries"`
	} `yaml:"jobs"`
	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Lire le fichier de configuration
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier de configuration: %v", err)
	}

	// Décoder le YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage du fichier de configuration: %v", err)
	}

	// Valider et définir les valeurs par défaut
	err = validateAndSetDefaults(config)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la validation de la configuration: %v", err)
	}

	return config, nil
}

func validateAndSetDefaults(config *Config) error {
	// Valider et définir les valeurs par défaut pour le serveur
	if config.Server.Port == 0 {
		config.Server.Port = 8080 // Port par défaut
	}

	// Valider le chemin de la base de données
	if config.Database.Path == "" {
		return fmt.Errorf("le chemin de la base de données doit être spécifié")
	}

	// Valider et définir les valeurs par défaut pour les jobs
	if config.Jobs.DefaultTimeout == 0 {
		config.Jobs.DefaultTimeout = 5 * time.Minute // Timeout par défaut
	}
	if config.Jobs.MaxRetries == 0 {
		config.Jobs.MaxRetries = 3 // Nombre maximal de tentatives par défaut
	}

	// Valider et définir les valeurs par défaut pour le logging
	if config.Logging.Level == "" {
		config.Logging.Level = "info" // Niveau de log par défaut
	}
	if config.Logging.File == "" {
		config.Logging.File = "orchestrator.log" // Fichier de log par défaut
	}

	return nil
}

package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name             string `yaml:"name"`
	Endpoint         string `yaml:"endpoint"`
	Threshold        int    `yaml:"threshold"`
	Heal             string `yaml:"heal"`
	ContainerName    string `yaml:"container_name"`
	FailureThreshold int    `yaml:"failure_threshold"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read config file %s: %v", path, err)
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("Failed to unmarshal config: %v", err)
		return nil, err
	}
	for _, svc := range cfg.Services {
		log.Printf("Loaded service: name=%s, endpoint=%s, threshold=%d, heal=%s, container_name=%s",
			svc.Name, svc.Endpoint, svc.Threshold, svc.Heal, svc.ContainerName)
	}
	return &cfg, nil
}

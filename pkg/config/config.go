package config

import "github.com/spf13/viper"

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name      string `yaml:"name"`
	Endpoint  string `yaml:"endpoint"`
	Threshold int    `yaml:"threshold"`
	Heal      string `yaml:"heal"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

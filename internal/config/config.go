package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Name     string `yaml:"name" json:"name"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"-"`
	Password string `yaml:"password" json:"-"`
}

type APIConfig struct {
	Port   int    `yaml:"port"`
	APIKey string `yaml:"api_key"`
}

type Config struct {
	API     APIConfig `yaml:"api"`
	Servers []Server  `yaml:"servers"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.API.Port == 0 {
		cfg.API.Port = 8080
	}

	for i, s := range cfg.Servers {
		if s.Port == 0 {
			cfg.Servers[i].Port = 623
		}
	}

	return &cfg, nil
}

func (c *Config) FindServer(name string) (*Server, error) {
	for _, s := range c.Servers {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("server %q not found", name)
}

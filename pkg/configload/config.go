package configload

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	cfg := &Config{}
	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		log.Fatalf("Error decoding config: %v", err)
	}

	return cfg
}

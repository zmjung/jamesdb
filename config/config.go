package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" envconfig:"HOST"`
		Port int    `yaml:"port" envconfig:"PORT"`
	} `yaml:"server"`

	Database struct {
		RootPath string `yaml:"rootPath" envconfig:"ROOT_PATH"`
	} `yaml:"database"`
}

func LoadConfig(cfg *Config, configPath string) error {
	// Load configuration from a file and environment variables.
	err := readFile(cfg, configPath)
	if err != nil {
		return err
	}

	return readEnv(cfg)
}

func readFile(cfg *Config, configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(cfg)
}

func readEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Printf("Error reading environment variables: %v\n", err)
	}
	return err
}

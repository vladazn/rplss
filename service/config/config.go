package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func New(path string) (*Config, error) {
	conf := Config{}
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	err = yaml.Unmarshal(dat, &conf)

	return &conf, err
}

type Config struct {
	Choices map[string][]string `yaml:"choices"`
	Redis   RedisConfigs        `yaml:"redis"`
	JWT     JWTConfigs          `yaml:"jwt"`
	Grpc    GrpcConfigs         `yaml:"grpc"`
}

type RedisConfigs struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Db   int    `yaml:"db"`
}

type JWTConfigs struct {
	Key string `yaml:"key"`
}

type GrpcConfigs struct {
	Host string
	Port int
}

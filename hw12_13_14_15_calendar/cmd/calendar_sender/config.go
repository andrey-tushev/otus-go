package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf
	RabbitMQ RabbitMQ
}

type LoggerConf struct {
	Level string
}

type RabbitMQ struct {
	URI      string
	Queue    string
	Consumer string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Parse(path string) error {
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return fmt.Errorf("config read error: %w", err)
	}
	return nil
}

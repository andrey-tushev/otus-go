package sender

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf
	RabbitMQ RabbitMQConf
}

type LoggerConf struct {
	Level string
}

type RabbitMQConf struct {
	URI      string
	Queue    string
	Consumer string
}

func New() *Config {
	return &Config{}
}

func (c *Config) Parse(path string) error {
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return fmt.Errorf("config read error: %w", err)
	}
	return nil
}

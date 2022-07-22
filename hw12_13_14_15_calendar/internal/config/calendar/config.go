package calendar

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	SQL     SQLConf
	Web     WebConf
	GRPC    GRPCConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Storage string
}

type SQLConf struct {
	DSN string
}

type WebConf struct {
	Host string
	Port string
}

type GRPCConf struct {
	Port string
}

func New() Config {
	return Config{}
}

func (c *Config) Parse(path string) error {
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return fmt.Errorf("config read error: %w", err)
	}
	return nil
}

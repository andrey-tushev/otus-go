package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	SQL     SQLConf
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

func NewConfig() Config {
	return Config{}
}

func (c *Config) Parse(path string) error {
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return fmt.Errorf("config read error: %w", err)
	}
	return nil
}

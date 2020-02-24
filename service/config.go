package service

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
	"os"
)

var (
	newFile *os.File
	err     error
)

type Config struct {
	*flag.FlagSet `json:"-"`

	Port string `toml:"port" json:"port"`

	NginxPath string `toml:"nginx-path" json:"nginx-path"`

	NginxConfig string `toml:"nginx-config" json:"nginx-config"`

	SwitchConfig string `toml:"switch-config" json:"switch-config"`

	MonitorDb DBConfig `toml:"monitor-db" json:"monitor-db"`

	SwitchDb DBConfig `toml:"switch-db" json:"switch-db"`
}

type DBConfig struct {
	Host string `toml:"host" json:"host"`

	User string `toml:"user" json:"user"`

	Password string `toml:"password" json:"password"`

	Name string `toml:"name" json:"name"`

	Port int `toml:"port" json:"port"`
}

// NewConfig creates a new config.
func NewConfig() *Config {
	cfg := &Config{}

	return cfg
}

// configFromFile loads config from file.
func (c *Config) ConfigFromFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	return errors.Trace(err)
}

package config

import (
	"flag"
	"os"
	"sync"

	"github.com/caarlos0/env"
)

const (
	defaultServerAddress = "localhost:8080"
)

type (
	Config struct {
		ServerAddress string `env:"ADDRESS"`
	}

	Option func(*Config)
)

var (
	config *Config
	once   sync.Once
)

func NewConfig(options ...Option) *Config {
	once.Do(
		func() {
			config = &Config{
				ServerAddress: defaultServerAddress,
			}
			for _, option := range options {
				option(config)
			}
		})
	return config
}

func WithParseConfig() Option {
	return func(c *Config) {
		env.Parse(c)
		c.ParseFlags()
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "ADDRESS")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
}

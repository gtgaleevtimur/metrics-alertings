package config

import (
	"flag"
	"os"
	"sync"
)

const (
	defaultServerAddress = "localhost:8080"
)

type (
	Config struct {
		ServerAddress string
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
		c.ParseFlags()
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "SERVER_ADDRESS")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
}

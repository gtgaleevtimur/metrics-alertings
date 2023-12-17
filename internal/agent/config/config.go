package config

import (
	"flag"
	"github.com/caarlos0/env"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

type (
	Config struct {
		ServerAddress  string        `env:"ADDRESS"`
		PollInterval   time.Duration `env:"POLL_INTERVAL"`
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`
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
				ServerAddress:  defaultServerAddress,
				PollInterval:   defaultPollInterval,
				ReportInterval: defaultReportInterval,
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
	reportInterval := flag.Int("r", 10, "REPORT_INTERVAL")
	pollInterval := flag.Int("p", 2, "POLL_INTERVAL")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	if !strings.Contains(c.ServerAddress, "http://") && !strings.Contains(c.ServerAddress, "https://") {
		c.ServerAddress = strings.Join([]string{"http://", c.ServerAddress}, "")
	}
	c.ReportInterval = time.Second * time.Duration(*reportInterval)
	c.PollInterval = time.Second * time.Duration(*pollInterval)
}

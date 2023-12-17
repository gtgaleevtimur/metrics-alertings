package config

import (
	"flag"
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
		ServerAddress  string
		PollInterval   time.Duration
		ReportInterval time.Duration
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
		c.ParseFlags()
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "SERVER_ADDRESS")
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "REPORT_INTERVAL")
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "POLL_INTERVAL")
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	if !strings.Contains(c.ServerAddress, "http://") && !strings.Contains(c.ServerAddress, "https://") {
		c.ServerAddress = strings.Join([]string{"http://", c.ServerAddress}, "")
	}
}

package app

import (
	"github.com/gtgaleevtimur/metrics-alertings/internal/agent/config"
	"github.com/gtgaleevtimur/metrics-alertings/internal/agent/repository"
	"log"
	"strings"
	"time"
)

func Run() {
	storager := repository.NewAgentMemStorage()
	conf := config.NewConfig(config.WithParseConfig())

	go func() {
		for {
			time.Sleep(conf.PollInterval)
			storager.UpdateMemStorage()
		}
	}()
	for {
		time.Sleep(conf.ReportInterval)
		if err := storager.SendMetrics(strings.Join([]string{conf.ServerAddress, "/update/"}, "")); err != nil {
			log.Fatal(err)
		}
	}
}

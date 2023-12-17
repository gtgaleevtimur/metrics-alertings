package app

import (
	"github.com/gtgaleevtimur/metrics-alertings/internal/agent/repository"
	"log"
	"time"
)

const (
	pollInterval               = 2 * time.Second
	reportInterval             = 10 * time.Second
	updateServerAddress string = "http://localhost:8080/update/"
)

func Run() {
	storager := repository.NewAgentMemStorage()

	go func() {
		for {
			time.Sleep(pollInterval)
			storager.UpdateMemStorage()
		}
	}()
	for {
		time.Sleep(reportInterval)
		if err := storager.SendMetrics(updateServerAddress); err != nil {
			log.Fatal(err)
		}
	}
}

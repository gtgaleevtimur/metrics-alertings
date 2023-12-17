package app

import (
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/config"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/handler"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/repository"
	"log"
	"net/http"
)

func RunServer() {
	conf := config.NewConfig(config.WithParseConfig())
	storager := repository.NewServerMemStorage()
	server := http.Server{
		Addr:    conf.ServerAddress,
		Handler: handler.NewServerRouter(storager),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

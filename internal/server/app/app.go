package app

import (
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/handler"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/repository"
	"log"
	"net/http"
)

func RunServer() {
	storager := repository.NewServerMemStorage()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: handler.NewServerRouter(storager),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

package app

import (
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/handler"
	"github.com/gtgaleevtimur/metrics-alertings/internal/server/repository/memory"
	"log"
	"net/http"
)

func RunServer() {
	repository := memory.NewServerMemStorage()
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: handler.NewServerRouter(repository),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

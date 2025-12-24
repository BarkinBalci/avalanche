package main

import (
	"log"
	"net/http"

	"github.com/BarkinBalci/event-analytics-service/internal/api"
	"github.com/BarkinBalci/event-analytics-service/internal/config"
)

func main() {
	_ = config.Load()

	handler := api.NewHandler()

	addr := ":8080"
	log.Printf("API starting on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
}

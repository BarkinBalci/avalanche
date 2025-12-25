package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/BarkinBalci/event-analytics-service/docs"
	"github.com/BarkinBalci/event-analytics-service/internal/api"
	"github.com/BarkinBalci/event-analytics-service/internal/clients"
	"github.com/BarkinBalci/event-analytics-service/internal/config"
	"github.com/BarkinBalci/event-analytics-service/internal/service"
)

// @title Event Analytics Service API
// @version 1.0
// @description API for publishing and managing analytics events
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Configure Swagger host dynamically
	docs.SwaggerInfo.Host = cfg.ServiceHost

	// Initialize SQS client
	sqsClient, err := clients.NewSQSClient(context.Background(), cfg.SQSEndpoint, cfg.SQSQueueURL, cfg.SQSRegion)
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	// Initialize event service
	eventService := service.NewEventService(sqsClient)

	// Initialize handler
	handler := api.NewHandler(eventService)

	addr := fmt.Sprintf(":%s", cfg.ServiceAPIPort)
	log.Printf("API starting on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
}

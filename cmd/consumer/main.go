package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BarkinBalci/event-analytics-service/internal/config"
	"github.com/BarkinBalci/event-analytics-service/internal/consumer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := consumer.New(cfg)

	log.Println("Consumer starting")

	go func() {
		if err := c.Start(ctx); err != nil {
			log.Fatalf("Consumer error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer")
	cancel()
}

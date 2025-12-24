package consumer

import (
	"context"
	"log"

	"github.com/BarkinBalci/event-analytics-service/internal/config"
)

type Consumer struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Consumer {
	return &Consumer{
		cfg: cfg,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Consumer running")
	<-ctx.Done()
	return nil
}

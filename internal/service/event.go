package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BarkinBalci/event-analytics-service/internal/clients"
	"github.com/BarkinBalci/event-analytics-service/internal/models"
	"github.com/google/uuid"
)

// EventService represents event service
type EventService struct {
	sqsClient *clients.SQSClient
}

// NewEventService creates a new event service
func NewEventService(sqsClient *clients.SQSClient) *EventService {
	return &EventService{
		sqsClient: sqsClient,
	}
}

// ProcessEvent processes a single event
func (s *EventService) ProcessEvent(event *models.PublishEventRequest) (string, error) {
	ctx := context.Background()

	currentTime := time.Now().Unix()
	if event.Timestamp > currentTime+1 {
		return "", fmt.Errorf("timestamp cannot be in the future: %d > %d", event.Timestamp, currentTime)
	}

	eventID := uuid.New().String()

	err := s.sqsClient.PublishEvent(ctx, event, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to publish event to SQS: %w", err)
	}

	return eventID, nil
}

// ProcessBulkEvents validates and processes multiple events
func (s *EventService) ProcessBulkEvents(events []models.PublishEventRequest) ([]string, []string, error) {
	var eventIDs []string
	var errors []string

	for i, event := range events {
		eventID, err := s.ProcessEvent(&event)
		if err != nil {
			errors = append(errors, err.Error())
			log.Printf("Failed to process event at index %d: %v", i, err)
			continue
		}
		eventIDs = append(eventIDs, eventID)
	}

	return eventIDs, errors, nil
}

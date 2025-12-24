package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/BarkinBalci/event-analytics-service/internal/models"
)

// SQSClient represents an SQS client
type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSClient creates a new SQS client
func NewSQSClient(ctx context.Context, endpoint, queueURL, region string) (*SQSClient, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("queueURL is required")
	}
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}

	configOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	var clientOpts []func(*sqs.Options)

	// Configure for local development with ElasticMQ
	if endpoint != "" {
		configOpts = append(configOpts,
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")))

		clientOpts = append(clientOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(cfg, clientOpts...)

	return &SQSClient{
		client:   sqsClient,
		queueURL: queueURL,
	}, nil
}

// PublishEvent publishes an event to SQS
func (c *SQSClient) PublishEvent(ctx context.Context, event *models.PublishEventRequest, eventID string) error {
	messageBody := map[string]interface{}{
		"event_id":    eventID,
		"event_name":  event.EventName,
		"channel":     event.Channel,
		"campaign_id": event.CampaignID,
		"user_id":     event.UserID,
		"timestamp":   event.Timestamp,
		"tags":        event.Tags,
		"metadata":    event.Metadata,
	}

	bodyJSON, err := json.Marshal(messageBody)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = c.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(bodyJSON)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"EventName": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.EventName),
			},
			"Channel": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.Channel),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return nil
}

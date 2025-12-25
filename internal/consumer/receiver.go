package consumer

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.uber.org/zap"

	"github.com/BarkinBalci/event-analytics-service/internal/queue/sqs"
)

// ReceiverConfig configures the SQS receiver
type ReceiverConfig struct {
	MaxMessages     int32
	WaitTimeSeconds int32
	BufferSize      int
}

// Receiver handles receiving messages from SQS
type Receiver struct {
	sqsClient *sqs.Client
	config    ReceiverConfig
	log       *zap.Logger
}

// NewReceiver creates a new SQS receiver
func NewReceiver(sqsClient *sqs.Client, config ReceiverConfig, log *zap.Logger) *Receiver {
	return &Receiver{
		sqsClient: sqsClient,
		config:    config,
		log:       log,
	}
}

// Start begins receiving messages and sends them to the output channel
func (r *Receiver) Start(ctx context.Context, out chan<- types.Message) {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			r.log.Info("Receiver shutting down")
			return
		default:
			result, err := r.sqsClient.Client().ReceiveMessage(ctx, &awssqs.ReceiveMessageInput{
				QueueUrl:              aws.String(r.sqsClient.QueueURL()),
				MaxNumberOfMessages:   r.config.MaxMessages,
				WaitTimeSeconds:       r.config.WaitTimeSeconds,
				MessageAttributeNames: []string{"All"},
			})

			if err != nil {
				r.log.Error("Error receiving messages from SQS", zap.Error(err))
				time.Sleep(1 * time.Second)
				continue
			}

			if len(result.Messages) == 0 {
				continue
			}

			r.log.Info("Received messages from SQS", zap.Int("message_count", len(result.Messages)))

			// Send messages to the next stage
			for _, msg := range result.Messages {
				select {
				case <-ctx.Done():
					r.log.Info("Receiver shutting down while sending messages")
					return
				case out <- msg:
					// Message sent to next stage
				}
			}
		}
	}
}

package config

import (
	"fmt"

	"github.com/BarkinBalci/envconfig"
)

type Config struct {
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" required:"true"`
	ServiceAPIPort     string `envconfig:"SERVICE_API_PORT" required:"true"`
	ValkeyHost         string `envconfig:"VALKEY_HOST" required:"true"`
	ValkeyPort         string `envconfig:"VALKEY_PORT" required:"true"`
	SQSEndpoint        string `envconfig:"SQS_ENDPOINT"`
	SQSQueueURL        string `envconfig:"SQS_QUEUE_URL" required:"true"`
	SQSRegion          string `envconfig:"SQS_REGION" required:"true"`
	ClickHouseHost     string `envconfig:"CLICKHOUSE_HOST" required:"true"`
	ClickHousePort     string `envconfig:"CLICKHOUSE_PORT" required:"true"`
	ClickHouseDB       string `envconfig:"CLICKHOUSE_DB" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	return &cfg, nil
}

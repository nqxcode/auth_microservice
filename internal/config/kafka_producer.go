package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/sarama"
)

const (
	requiredAcksEnvName    = "KAFKA_REQUIRED_ACKS"
	retryMaxEnvName        = "KAFKA_RETRY_MAX"
	returnSuccessesEnvName = "KAFKA_RETURN_SUCCESSES"
)

type kafkaProducerConfig struct {
	brokers         []string
	requiredAcks    sarama.RequiredAcks
	retryMax        int
	returnSuccesses bool
}

// NewKafkaProducerConfig new kafka producer config
func NewKafkaProducerConfig() (*kafkaProducerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	requiredAcksStr := os.Getenv(requiredAcksEnvName)
	if len(requiredAcksStr) == 0 {
		return nil, errors.New("kafka required acks not found")
	}
	requiredAcks, err := strconv.ParseInt(requiredAcksStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to parse required acks: %s", requiredAcksStr)
	}

	retryMaxStr := os.Getenv(retryMaxEnvName)
	if len(retryMaxStr) == 0 {
		return nil, errors.New("kafka retry max not found")
	}
	retryMax, err := strconv.Atoi(retryMaxStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse retry max: %s", retryMaxStr)
	}

	returnSuccessesStr := os.Getenv(returnSuccessesEnvName)
	if len(returnSuccessesStr) == 0 {
		return nil, errors.New("kafka return successes not found")
	}
	returnSuccesses, err := strconv.ParseBool(returnSuccessesStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse return successes: %s", retryMaxStr)
	}

	return &kafkaProducerConfig{
		brokers:         brokers,
		requiredAcks:    sarama.RequiredAcks(int16(requiredAcks)), //nolint: gosec
		retryMax:        retryMax,
		returnSuccesses: returnSuccesses,
	}, nil
}

func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaProducerConfig) RequiredAcks() sarama.RequiredAcks {
	return cfg.requiredAcks
}

func (cfg *kafkaProducerConfig) RetryMax() int {
	return cfg.retryMax
}

func (cfg *kafkaProducerConfig) ReturnSuccesses() bool {
	return cfg.returnSuccesses
}

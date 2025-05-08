package messaging

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/lioarce01/real-time-order-processing-system/internal/ports"
	"go.uber.org/zap"
)

type KafkaPublisher struct {
	producer sarama.SyncProducer
	logger   *zap.Logger
}

// NewKafkaPublisher creates a new Kafka publisher
func NewKafkaPublisher(brokers []string, logger *zap.Logger) (ports.EventPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Idempotent = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Net.MaxOpenRequests = 1

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaPublisher{producer: producer, logger: logger}, nil
}

// Publish send a message converted to JSON to specific kafka topic
func (p *KafkaPublisher) Publish(topic string, message interface{}) error {
	// serialize the message to json
	bytes, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal message", zap.Error(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bytes),
	}

	// sendMessage returns partition and offset if needed
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("failed to send message", zap.Error(err))
		return err
	}

	p.logger.Debug("message published to kafka", zap.String("topic", topic),
		zap.Int32("partition", partition), zap.Int64("offset", offset))
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.producer.Close()
}

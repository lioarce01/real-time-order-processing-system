package messaging

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/lioarce01/real-time-order-processing-system/internal/core/domain"
	"github.com/lioarce01/real-time-order-processing-system/internal/core/service"
	"go.uber.org/zap"
)

type OrderConsumer struct {
	handler *service.OrderService
	logger  *zap.Logger
}

func (consumer *OrderConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *OrderConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// consumeClaim processes messages from kafka
func (consumer *OrderConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		consumer.logger.Info("message received", zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition), zap.Int64("offset", msg.Offset))

		// decode message (assuming its an order json)
		var order domain.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			consumer.logger.Info("failed to unmarshal message", zap.Error(err))
			session.MarkMessage(msg, "")
			continue
		}

		//proccess the order
		consumer.logger.Info("proccessing order", zap.Uint("order_id", order.ID))
		// call the order service to process the order
		// this is where the business logic goes

		// mark message as processed
		session.MarkMessage(msg, "")

	}
	return nil
}

// new kafka consumer group initializes and runs the consumer group
func NewKafkaConsumerGroup(brokers []string, groupID string, topics []string, handler *service.OrderService, logger *zap.Logger, config *sarama.Config) error {
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}

	consumer := &OrderConsumer{handler: handler, logger: logger}
	ctx := context.Background()

	go func() {
		for {
			if err := group.Consume(ctx, topics, consumer); err != nil {
				logger.Error("error from consumer group", zap.Error(err))
			}
		}
	}()
	return nil
}

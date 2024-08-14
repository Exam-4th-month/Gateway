package msgbroker

import (
	"context"
	"gateway-service/internal/items/config"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type MsgBroker struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func NewMsgBroker(writer *kafka.Writer, logger *slog.Logger) *MsgBroker {
	return &MsgBroker{
		writer: writer,
		logger: logger,
	}
}

func (b *MsgBroker) TransactionCreated(ctx context.Context, body []byte) error {
	return b.publishMessage(ctx, "transaction_created", body)
}

func (b *MsgBroker) BudgetUpdated(ctx context.Context, body []byte) error {
	return b.publishMessage(ctx, "budget_updated", body)
}

func (b *MsgBroker) GoalProgressUpdated(ctx context.Context, body []byte) error {
	return b.publishMessage(ctx, "goal_progress_updated", body)
}

func (b *MsgBroker) NotificationCreated(ctx context.Context, body []byte) error {
	return b.publishMessage(ctx, "notification_created", body)
}

func (b *MsgBroker) publishMessage(ctx context.Context, topic string, body []byte) error {
	err := b.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: body,
	})
	if err != nil {
		b.logger.Error("Failed to publish message", "topic", topic, "error", err.Error())
		return err
	}

	b.logger.Info("Message published", "topic", topic)
	return nil
}

func CreateTopics(config *config.Config, logger *slog.Logger) error {
	topics := []string{
		"transaction_created",
		"budget_updated",
		"goal_progress_updated",
		"notification_created",
	}

	conn, err := kafka.DialContext(context.Background(), "tcp", config.Kafka.Brokers)
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, topic := range topics {
		partitions, err := conn.ReadPartitions(topic)
		if err == nil && len(partitions) > 0 {
			logger.Info("Topic already exists", "topic", topic)
			continue
		}

		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err != nil {
			logger.Error("Failed to create topic", "topic", topic, "error", err.Error())
			return err
		}

		logger.Info("Topic created successfully", "topic", topic)
	}

	return nil
}

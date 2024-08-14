package msgbroker

import (
	"context"
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

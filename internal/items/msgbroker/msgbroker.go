package msgbroker

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MsgBroker struct {
	channel *amqp.Channel
	logger  *slog.Logger
}

func NewMsgBroker(channel *amqp.Channel, logger *slog.Logger) *MsgBroker {
	return &MsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *MsgBroker) TransactionCreated(body []byte) error {
	return b.publishMessage("transaction_created", body)
}

func (b *MsgBroker) BudgetUpdated(body []byte) error {
	return b.publishMessage("budget_updated", body)
}

func (b *MsgBroker) GoalProgressUpdated(body []byte) error {
	return b.publishMessage("goal_progress_updated", body)
}

func (b *MsgBroker) NotificationCreated(body []byte) error {
	return b.publishMessage("notification_created", body)
}

// publishMessage is a helper function to publish messages to a specified queue.
func (b *MsgBroker) publishMessage(queueName string, body []byte) error {
	err := b.channel.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		b.logger.Error("Failed to publish message", "queue", queueName, "error", err.Error())
		return err
	}

	b.logger.Info("Message published", "queue", queueName)
	return nil
}

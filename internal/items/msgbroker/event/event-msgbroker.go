package event

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	EventMsgBroker struct {
		channel *amqp.Channel
		logger  *slog.Logger
	}
)

func NewEventMsgBroker(channel *amqp.Channel, logger *slog.Logger) *EventMsgBroker {
	return &EventMsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *EventMsgBroker) CreateEvent(body []byte) error {
	return b.publishMessage("create_event", body)
}

func (b *EventMsgBroker) UpdateEvent(body []byte) error {
	return b.publishMessage("update_event", body)
}

func (b *EventMsgBroker) DeleteEvent(body []byte) error {
	return b.publishMessage("delete_event", body)
}

// publishMessage is a helper function to publish messages to a specified queue.
func (b *EventMsgBroker) publishMessage(queueName string, body []byte) error {
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

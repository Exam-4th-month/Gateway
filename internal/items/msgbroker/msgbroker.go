package msgbroker

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	MsgBroker struct {
		channel *amqp.Channel
		logger  *slog.Logger
	}
)

func NewMsgBroker(channel *amqp.Channel, logger *slog.Logger) *MsgBroker {
	return &MsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *MsgBroker) CreateAthlete(body []byte) error {
	return b.publishMessage("create_athlete", body)
}

func (b *MsgBroker) UpdateAthlete(body []byte) error {
	return b.publishMessage("update_athlete", body)
}

func (b *MsgBroker) DeleteAthlete(body []byte) error {
	return b.publishMessage("delete_athlete", body)
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

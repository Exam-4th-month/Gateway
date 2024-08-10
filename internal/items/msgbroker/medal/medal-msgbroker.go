package medal

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	MedalMsgBroker struct {
		channel *amqp.Channel
		logger  *slog.Logger
	}
)

func NewMedalMsgBroker(channel *amqp.Channel, logger *slog.Logger) *MedalMsgBroker {
	return &MedalMsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *MedalMsgBroker) CreateMedal(body []byte) error {
	return b.publishMessage("create_medal", body)
}

func (b *MedalMsgBroker) UpdateMedal(body []byte) error {
	return b.publishMessage("update_medal", body)
}

func (b *MedalMsgBroker) DeleteMedal(body []byte) error {
	return b.publishMessage("delete_medal", body)
}

// publishMessage is a helper function to publish messages to a specified queue.
func (b *MedalMsgBroker) publishMessage(queueName string, body []byte) error {
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

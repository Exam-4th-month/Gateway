package athlete

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	AthleteMsgBroker struct {
		channel *amqp.Channel
		logger  *slog.Logger
	}
)

func NewAthleteMsgBroker(channel *amqp.Channel, logger *slog.Logger) *AthleteMsgBroker {
	return &AthleteMsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *AthleteMsgBroker) CreateAthlete(body []byte) error {
	return b.publishMessage("create_athlete", body)
}

func (b *AthleteMsgBroker) UpdateAthlete(body []byte) error {
	return b.publishMessage("update_athlete", body)
}

func (b *AthleteMsgBroker) DeleteAthlete(body []byte) error {
	return b.publishMessage("delete_athlete", body)
}

// publishMessage is a helper function to publish messages to a specified queue.
func (b *AthleteMsgBroker) publishMessage(queueName string, body []byte) error {
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

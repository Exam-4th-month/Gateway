package auth

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	AuthMsgBroker struct {
		channel *amqp.Channel
		logger  *slog.Logger
	}
)

func NewAthleteMsgBroker(channel *amqp.Channel, logger *slog.Logger) *AuthMsgBroker {
	return &AuthMsgBroker{
		channel: channel,
		logger:  logger,
	}
}

func (b *AuthMsgBroker) Register(body []byte) error {
	return b.publishMessage("register_user", body)
}

func (b *AuthMsgBroker) CreateAdmin(body []byte) error {
	return b.publishMessage("set_admin", body)
}

func (b *AuthMsgBroker) UpdateUser(body []byte) error {
	return b.publishMessage("update_user", body)
}

func (b *AuthMsgBroker) DeleteUser(body []byte) error {
	return b.publishMessage("delete_user", body)
}

// publishMessage is a helper function to publish messages to a specified queue.
func (b *AuthMsgBroker) publishMessage(queueName string, body []byte) error {
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

package handler

import (
	"log/slog"

	"gateway-service/internal/items/config"

	"gateway-service/internal/items/http/handler/auth"
	"gateway-service/internal/items/http/handler/budgeting"
	msgbroker "gateway-service/internal/items/msgbroker"

	amqp "github.com/rabbitmq/amqp091-go"

	"gateway-service/internal/items/redisservice"
)

type Handler struct {
	AuthRepo      *auth.AuthHandler
	BudgetingRepo *budgeting.BudgetingHandler
}

func New(redis *redisservice.RedisService, logger *slog.Logger, config *config.Config, channel *amqp.Channel) *Handler {
	msgbroker := msgbroker.NewMsgBroker(channel, logger)

	return &Handler{
		AuthRepo:      auth.NewAuthHandler(logger, redis, config),
		BudgetingRepo: budgeting.NewBudgetingHandler(logger, msgbroker, config),
	}
}

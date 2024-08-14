package handler

import (
	"log/slog"

	"gateway-service/internal/items/config"
	"gateway-service/internal/items/redisservice"

	"gateway-service/internal/items/http/handler/auth"
	"gateway-service/internal/items/http/handler/budgeting"
	msgbroker "gateway-service/internal/items/msgbroker"

	"github.com/segmentio/kafka-go"
)

type Handler struct {
	AuthRepo      *auth.AuthHandler
	BudgetingRepo *budgeting.BudgetingHandler
}

func New(redis *redisservice.RedisService, logger *slog.Logger, config *config.Config, writer *kafka.Writer) *Handler {
	msgbroker := msgbroker.NewMsgBroker(writer, logger)

	return &Handler{
		AuthRepo:      auth.NewAuthHandler(logger, config),
		BudgetingRepo: budgeting.NewBudgetingHandler(redis, logger, msgbroker, config),
	}
}

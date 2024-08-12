package budgeting

import (
	pb "gateway-service/genproto/budget"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type BudgetHandler struct {
	budget    pb.BudgetServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
}

func NewBudgetHandler(budget pb.BudgetServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *BudgetHandler {
	return &BudgetHandler{
		budget:    budget,
		logger:    logger,
		msgbroker: msgbroker,
	}
}

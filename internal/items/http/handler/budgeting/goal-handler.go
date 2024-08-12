package budgeting

import (
	pb "gateway-service/genproto/goal"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type GoalHandler struct {
	goal pb.GoalServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
}

func NewGoalHandler(goal pb.GoalServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *GoalHandler {
	return &GoalHandler{
		goal: goal,
		logger:    logger,
		msgbroker: msgbroker,
	}
}

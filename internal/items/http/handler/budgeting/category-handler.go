package budgeting

import (
	pb "gateway-service/genproto/category"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type CategoryHandler struct {
	category  pb.CategoryServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
}

func NewCategoryHandler(category pb.CategoryServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *CategoryHandler {
	return &CategoryHandler{
		category: category,
		logger:    logger,
		msgbroker: msgbroker,
	}
}

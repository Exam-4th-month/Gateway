package budgeting

import (
	pb "gateway-service/genproto/account"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type AccountHandler struct {
	account   pb.AccountServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
}

func NewAccountHandler(account pb.AccountServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *AccountHandler {
	return &AccountHandler{
		account:   account,
		logger:    logger,
		msgbroker: msgbroker,
	}
}

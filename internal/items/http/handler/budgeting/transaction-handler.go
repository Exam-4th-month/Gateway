package budgeting

import (
	pb "gateway-service/genproto/transaction"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type TransactionHandler struct {
	transaction pb.TransactionServiceClient
	logger      *slog.Logger
	msgbroker   *msgbroker.MsgBroker
}

func NewTransactionHandler(transaction pb.TransactionServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *TransactionHandler {
	return &TransactionHandler{
		transaction: transaction,
		logger:      logger,
		msgbroker:   msgbroker,
	}
}

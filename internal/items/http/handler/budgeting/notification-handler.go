package budgeting

import (
	pb "gateway-service/genproto/notification"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type NotificationHandler struct {
	notification pb.NotificationServiceClient
	logger       *slog.Logger
	msgbroker    *msgbroker.MsgBroker
}

func NewNotificationHandler(notification pb.NotificationServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *NotificationHandler {
	return &NotificationHandler{
		notification: notification,
		logger:       logger,
		msgbroker:    msgbroker,
	}
}

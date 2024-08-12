package budgeting

import (
	pb "gateway-service/genproto/report"
	"gateway-service/internal/items/msgbroker"
	"log/slog"
)

type ReportHandler struct {
	report    pb.ReportServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
}

func NewReportHandler(report pb.ReportServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker) *ReportHandler {
	return &ReportHandler{
		report:    report,
		logger:    logger,
		msgbroker: msgbroker,
	}
}

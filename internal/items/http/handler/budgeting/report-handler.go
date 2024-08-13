package budgeting

import (
	pb "gateway-service/genproto/report"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	report    pb.ReportServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
	config    *config.Config
}

func NewReportHandler(report pb.ReportServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *ReportHandler {
	return &ReportHandler{
		report:    report,
		logger:    logger,
		msgbroker: msgbroker,
		config:    config,
	}
}

// GetSpendingReportHandler godoc
// @Summary      Get spending report
// @Security     BearerAuth
// @Description  Retrieve a spending report for a user between the specified start and end dates
// @Tags         User Reports
// @Produce      json
// @Param        request  body  models.GetSpendingReportRequest  true  "Get Spending Report Request"
// @Success      200     {object}  pb.SpendingReportResponse
// @Failure      400     {object}  gin.H "Invalid request payload"
// @Failure      500     {object}  gin.H "Failed to retrieve spending report"
// @Router       /user/report/spending [get]
func (h *ReportHandler) GetSpendingReportHandler(c *gin.Context) {
	h.logger.Info("GetSpendingReportHandler called")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.GetSpendingReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.report.GetSpendingReport(c.Request.Context(), &pb.GetSpendingReportRequest{
		UserId:    userId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetIncomeReportHandler godoc
// @Summary      Get income report
// @Security     BearerAuth
// @Description  Retrieve an income report for a user between the specified start and end dates
// @Tags         User Reports
// @Produce      json
// @Param        request  body  models.GetIncomeReportRequest  true  "Get Income Report Request"
// @Success      200     {object}  pb.IncomeReportResponse
// @Failure      400     {object}  gin.H "Invalid request payload"
// @Failure      500     {object}  gin.H "Failed to retrieve income report"
// @Router       /user/report/incoming [get]
func (h *ReportHandler) GetIncomeReportHandler(c *gin.Context) {
	h.logger.Info("GetIncomeReportHandler called")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.GetIncomeReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.report.GetIncomeReport(c.Request.Context(), &pb.GetIncomeReportRequest{
		UserId:    userId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, resp)
}

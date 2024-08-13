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

// GetBudgetPerformanceReportHandler godoc
// @Summary      Get budget performance report
// @Security     BearerAuth
// @Description  Retrieve a budget performance report for a specific budget by its ID
// @Tags         User Reports
// @Produce      json
// @Param        id    path      string  true  "Budget ID"
// @Success      200    {object}  pb.BudgetPerformanceReportResponse
// @Failure      400    {object}  gin.H "Invalid request payload"
// @Failure      401    {object}  gin.H "User not authenticated"
// @Failure      500    {object}  gin.H "Failed to retrieve budget performance report"
// @Router       /user/report/bugdet [get]
func (h *ReportHandler) GetBudgetPerformanceReportHandler(c *gin.Context) {
	h.logger.Info("GetBudgetPerformanceReportHandler called")

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.report.GetBudgetPerformanceReport(c.Request.Context(), &pb.GetBudgetPerformanceReportRequest{
		UserId:   userId,
		BudgetId: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetGoalProgressReportHandler godoc
// @Summary      Get goal progress report
// @Security     BearerAuth
// @Description  Retrieve a goal progress report for a specific goal by its ID
// @Tags         User Reports
// @Produce      json
// @Param        id    path      string  true  "Goal ID"
// @Success      200    {object}  pb.GoalProgressReportResponse
// @Failure      400    {object}  gin.H "Invalid request payload"
// @Failure      401    {object}  gin.H "User not authenticated"
// @Failure      500    {object}  gin.H "Failed to retrieve goal progress report"
// @Router       /user/report/goal [get]
func (h *ReportHandler) GetGoalProgressReportHandler(c *gin.Context) {
	h.logger.Info("GetGoalProgressReportHandler called")

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.report.GetGoalProgressReport(c.Request.Context(), &pb.GetGoalProgressReportRequest{
		UserId: userId,
		GoalId: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, resp)
}

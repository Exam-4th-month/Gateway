package budgeting

import (
	pb "gateway-service/genproto/budget"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type BudgetHandler struct {
	budget    pb.BudgetServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
	config    *config.Config
}

func NewBudgetHandler(budget pb.BudgetServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *BudgetHandler {
	return &BudgetHandler{
		budget:    budget,
		logger:    logger,
		msgbroker: msgbroker,
		config:    config,
	}
}

// CreateBudgetHandler godoc
// @Summary      Create a budget
// @Security     BearerAuth
// @Description  Create a new budget for the authenticated user
// @Tags         User Budgets
// @Accept       json
// @Produce      json
// @Param        CreateBudgetRequest  body      models.CreateBudgetRequest  true  "Budget details"
// @Success      201                   {object}  pb.BudgetResponse
// @Failure      401                   {object}  gin.H "User not authenticated"
// @Failure      400                   {object}  gin.H "Invalid request body"
// @Failure      500                   {object}  gin.H "Failed to create budget"
// @Router       /user/budget [post]
func (h *BudgetHandler) CreateBudgetHandler(c *gin.Context) {
	h.logger.Info("CreateBudgetHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.budget.CreateBudget(c.Request.Context(), &pb.CreateBudgetRequest{
		UserId:     userId,
		CategoryId: req.CategoryID,
		Amount:     req.Amount,
		Period:     req.Period,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to create budget"})
		return
	}

	c.IndentedJSON(201, resp)
}

// GetBudgetsHandler godoc
// @Summary      Get budgets
// @Security     BearerAuth
// @Description  Get all budgets for the authenticated user
// @Tags         User Budgets
// @Produce      json
// @Success      200  {object}  pb.BudgetsResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to get budgets"
// @Router       /user/budget [get]
func (h *BudgetHandler) GetBudgetsHandler(c *gin.Context) {
	h.logger.Info("GetBudgetsHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.budget.GetBudgets(c.Request.Context(), &pb.GetBudgetsRequest{
		UserId: userId,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get budgets"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetBudgetByIdHandler godoc
// @Summary      Get budget by ID
// @Security     BearerAuth
// @Description  Get budget details by budget ID
// @Tags         User Budgets
// @Produce      json
// @Param        id   path      string  true  "Budget ID"
// @Success      200  {object}  pb.BudgetResponse
// @Failure      400  {object}  gin.H "Budget ID is required"
// @Failure      500  {object}  gin.H "Failed to retrieve budget"
// @Router       /user/budget/{id} [get]
func (h *BudgetHandler) GetBudgetByIdHandler(c *gin.Context) {
	h.logger.Info("GetBudgetByIdHandler")

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Budget ID is required"})
		return
	}

	resp, err := h.budget.GetBudgetById(c.Request.Context(), &pb.GetBudgetByIdRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to retrieve budget"})
		return
	}

	c.IndentedJSON(200, resp)
}

// UpdateBudgetHandler godoc
// @Summary      Update budget
// @Security     BearerAuth
// @Description  Update budget details by budget ID
// @Tags         User Budgets
// @Accept       json
// @Produce      json
// @Param        UpdateBudgetRequest  body      pb.UpdateBudgetRequest  true  "Updated budget details"
// @Success      200                   {object}  gin.H
// @Failure      400                   {object}  gin.H "Invalid request body"
// @Failure      500                   {object}  gin.H "Failed to update budget"
// @Router       /user/budget [put]
func (h *BudgetHandler) UpdateBudgetHandler(c *gin.Context) {
	h.logger.Info("UpdateBudgetHandler")

	var req pb.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while marshaling request"})
		return
	}

	err = h.msgbroker.BudgetUpdated(body)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while updating budjet"})
	}

	c.IndentedJSON(200, gin.H{"message": "Budget updated successfully!"})
}

// DeleteBudgetHandler godoc
// @Summary      Delete budget
// @Security     BearerAuth
// @Description  Delete budget by budget ID
// @Tags         User Budgets
// @Produce      json
// @Param        id   path      string  true  "Budget ID"
// @Success      200  {object}  gin.H "message: Budget deleted successfully"
// @Failure      400  {object}  gin.H "Budget ID is required"
// @Failure      500  {object}  gin.H "Failed to delete budget"
// @Router       /user/budget/{id} [delete]
func (h *BudgetHandler) DeleteBudgetHandler(c *gin.Context) {
	h.logger.Info("DeleteBudgetHandler")
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Budget ID is required"})
		return
	}

	_, err := h.budget.DeleteBudget(c.Request.Context(), &pb.DeleteBudgetRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to delete budget"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Budget deleted successfully"})
}

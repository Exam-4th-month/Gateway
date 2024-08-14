package budgeting

import (
	pb "gateway-service/genproto/goal"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type GoalHandler struct {
	goal      pb.GoalServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
	config    *config.Config
}

func NewGoalHandler(goal pb.GoalServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *GoalHandler {
	return &GoalHandler{
		goal:      goal,
		logger:    logger,
		msgbroker: msgbroker,
		config:    config,
	}
}

// CreateGoalHandler godoc
// @Summary      Create a goal
// @Security     BearerAuth
// @Description  Create a new financial goal for the authenticated user
// @Tags         User Goals
// @Accept       json
// @Produce      json
// @Param        CreateGoalRequest  body      models.CreateGoalRequest  true  "Goal details"
// @Success      200                {object}  pb.GoalResponse
// @Failure      401                {object}  gin.H "User not authenticated"
// @Failure      400                {object}  gin.H "Invalid request body"
// @Failure      500                {object}  gin.H "Failed to create goal"
// @Router       /user/goal [post]
func (h *GoalHandler) CreateGoalHandler(c *gin.Context) {
	h.logger.Info("CreateGoalHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.goal.CreateGoal(c.Request.Context(), &pb.CreateGoalRequest{
		UserId:        userId,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
		Deadline:      req.Deadline,
		Status:        req.Status,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to create goal"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetGoalsHandler godoc
// @Summary      Get goals
// @Security     BearerAuth
// @Description  Get all financial goals for the authenticated user
// @Tags         User Goals
// @Produce      json
// @Success      200  {object}  pb.GoalsResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to get goals"
// @Router       /user/goal [get]
func (h *GoalHandler) GetGoalsHandler(c *gin.Context) {
	h.logger.Info("GetGoalsHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.goal.GetGoals(c.Request.Context(), &pb.GetGoalsRequest{
		UserId: userId,
	})

	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get goals"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetGoalByIdHandler godoc
// @Summary      Get goal by ID
// @Security     BearerAuth
// @Description  Get financial goal details by goal ID
// @Tags         User Goals
// @Produce      json
// @Param        id   path      string  true  "Goal ID"
// @Success      200  {object}  pb.GoalResponse
// @Failure      400  {object}  gin.H "Goal ID is required"
// @Failure      500  {object}  gin.H "Failed to get goal"
// @Router       /user/goal/{id} [get]
func (h *GoalHandler) GetGoalByIdHandler(c *gin.Context) {
	h.logger.Info("GetGoalByIdHandler")
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Category ID is required"})
		return
	}

	resp, err := h.goal.GetGoalById(c.Request.Context(), &pb.GetGoalByIdRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get goal"})
		return
	}

	c.IndentedJSON(200, resp)
}

// UpdateGoalHandler godoc
// @Summary      Update goal
// @Security     BearerAuth
// @Description  Update financial goal details by goal ID
// @Tags         User Goals
// @Accept       json
// @Produce      json
// @Param        UpdateGoalRequest  body      pb.UpdateGoalRequest  true  "Updated goal details"
// @Success      200                {object}  gin.H
// @Failure      400                {object}  gin.H "Invalid request body"
// @Failure      500                {object}  gin.H "Failed to update goal"
// @Router       /user/goal [put]
func (h *GoalHandler) UpdateGoalHandler(c *gin.Context) {
	h.logger.Info("UpdateGoalHandler")

	var req pb.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	body, err := protojson.Marshal(&req)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while marshaling request"})
		return
	}

	err = h.msgbroker.GoalProgressUpdated(c.Request.Context(), body)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while updating goal"})
	}

	c.IndentedJSON(200, gin.H{"message": "Goal updated successfully!"})
}

// DeleteGoalHandler godoc
// @Summary      Delete goal
// @Security     BearerAuth
// @Description  Delete financial goal by goal ID
// @Tags         User Goals
// @Produce      json
// @Param        id   path      string  true  "Goal ID"
// @Success      200  {object}  gin.H "message: Goal deleted successfully"
// @Failure      400  {object}  gin.H "Goal ID is required"
// @Failure      500  {object}  gin.H "Failed to delete goal"
// @Router       /user/goal/{id} [delete]
func (h *GoalHandler) DeleteGoalHandler(c *gin.Context) {
	h.logger.Info("DeleteGoalHandler")
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Category ID is required"})
		return
	}

	_, err := h.goal.DeleteGoal(c.Request.Context(), &pb.DeleteGoalRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to delete goal"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Goal deleted successfully"})
}

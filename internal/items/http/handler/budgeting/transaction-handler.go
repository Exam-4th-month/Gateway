package budgeting

import (
	"fmt"
	not_pb "gateway-service/genproto/notification"
	pb "gateway-service/genproto/transaction"
	"time"

	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type TransactionHandler struct {
	transaction  pb.TransactionServiceClient
	notification not_pb.NotificationServiceClient
	logger       *slog.Logger
	msgbroker    *msgbroker.MsgBroker
	config       *config.Config
}

func NewTransactionHandler(notification not_pb.NotificationServiceClient, transaction pb.TransactionServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *TransactionHandler {
	return &TransactionHandler{
		transaction:  transaction,
		notification: notification,
		logger:       logger,
		msgbroker:    msgbroker,
		config:       config,
	}
}

// CreateTransactionHandler godoc
// @Summary      Create a transaction
// @Security     BearerAuth
// @Description  Create a new financial transaction for the authenticated user
// @Tags         User Transactions
// @Accept       json
// @Produce      json
// @Param        CreateTransactionRequest  body      models.CreateTransactionRequest  true  "Transaction details"
// @Success      201                       {object}  gin.H
// @Failure      401                       {object}  gin.H "User not authenticated"
// @Failure      400                       {object}  gin.H "Invalid request body"
// @Failure      500                       {object}  gin.H "Failed to create transaction"
// @Router       /user/transaction [post]
func (h *TransactionHandler) CreateTransactionHandler(c *gin.Context) {
	h.logger.Info("CreateTransactionHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var request = pb.CreateTransactionRequest{
		UserId:      userId,
		AccountId:   req.AccountID,
		CategoryId:  req.CategoryID,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
		Date:        time.Now().Format("2006-01-02"),
	}

	body, err := protojson.Marshal(&request)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while marshaling request"})
		return
	}

	err = h.msgbroker.TransactionCreated(c.Request.Context(), body)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while creating transaction"})
	}

	notification := not_pb.CreateNotificationRequest{
		UserId:  userId,
		Message: fmt.Sprintf("Transaction of %.2f has been created for account ID: %s", req.Amount, req.AccountID),
	}

	body, err = protojson.Marshal(&notification)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while marshaling request"})
		return
	}

	err = h.msgbroker.NotificationCreated(c.Request.Context(), body)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Error while creating notification"})
	}

	c.IndentedJSON(201, gin.H{"message": "Transaction created successfully!"})
}

// GetTransactionsHandler godoc
// @Summary      Get transactions
// @Security     BearerAuth
// @Description  Get all financial transactions for the authenticated user
// @Tags         User Transactions
// @Produce      json
// @Success      200  {object}  pb.TransactionsResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to get transactions"
// @Router       /user/transaction [get]
func (h *TransactionHandler) GetTransactionsHandler(c *gin.Context) {
	h.logger.Info("GetTransactionsHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.transaction.GetTransactions(c.Request.Context(), &pb.GetTransactionsRequest{
		UserId: userId,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetTransactionByIdHandler godoc
// @Summary      Get transaction by ID
// @Security     BearerAuth
// @Description  Get financial transaction details by transaction ID
// @Tags         User Transactions
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  pb.TransactionResponse
// @Failure      400  {object}  gin.H "Invalid transaction ID"
// @Failure      500  {object}  gin.H "Failed to get transaction"
// @Router       /user/transaction/{id} [get]
func (h *TransactionHandler) GetTransactionByIdHandler(c *gin.Context) {
	h.logger.Info("GetTransactionByIdHandler")
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Invalid transaction ID"})
		return
	}

	resp, err := h.transaction.GetTransactionById(c.Request.Context(), &pb.GetTransactionByIdRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get transaction"})
		return
	}

	c.IndentedJSON(200, resp)
}

// UpdateTransactionHandler godoc
// @Summary      Update transaction
// @Security     BearerAuth
// @Description  Update financial transaction details by transaction ID
// @Tags         User Transactions
// @Accept       json
// @Produce      json
// @Param        UpdateTransactionRequest  body      pb.UpdateTransactionRequest  true  "Updated transaction details"
// @Success      200                       {object}  pb.TransactionResponse
// @Failure      400                       {object}  gin.H "Invalid request body"
// @Failure      500                       {object}  gin.H "Failed to update transaction"
// @Router       /user/transaction [put]
func (h *TransactionHandler) UpdateTransactionHandler(c *gin.Context) {
	h.logger.Info("UpdateTransactionHandler")

	var req pb.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.transaction.UpdateTransaction(c.Request.Context(), &req)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.IndentedJSON(200, resp)
}

// DeleteTransactionHandler godoc
// @Summary      Delete transaction
// @Security     BearerAuth
// @Description  Delete financial transaction by transaction ID
// @Tags         User Transactions
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  gin.H "message: Transaction deleted successfully"
// @Failure      400  {object}  gin.H "Invalid transaction ID"
// @Failure      500  {object}  gin.H "Failed to delete transaction"
// @Router       /user/transaction/{id} [delete]
func (h *TransactionHandler) DeleteTransactionHandler(c *gin.Context) {
	h.logger.Info("DeleteTransactionHandler")
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Invalid transaction ID"})
		return
	}

	_, err := h.transaction.DeleteTransaction(c.Request.Context(), &pb.DeleteTransactionRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Transaction deleted successfully"})
}

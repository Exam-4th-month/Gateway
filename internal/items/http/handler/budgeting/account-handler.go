package budgeting

import (
	pb "gateway-service/genproto/account"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/items/redisservice"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	redis     *redisservice.RedisService
	account   pb.AccountServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
	config    *config.Config
}

func NewAccountHandler(redis *redisservice.RedisService, account pb.AccountServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *AccountHandler {
	return &AccountHandler{
		redis:     redis,
		account:   account,
		logger:    logger,
		msgbroker: msgbroker,
		config:    config,
	}
}

// CreateAccountHandler godoc
// @Summary      Create an account
// @Security     BearerAuth
// @Description  Create a new account for the authenticated user
// @Tags         User Accounts
// @Accept       json
// @Produce      json
// @Param        CreateAccountRequest  body      models.CreateAccountRequest  true  "Account details"
// @Success      201                   {object}  pb.AccountResponse
// @Failure      401                   {object}  gin.H "User not authenticated"
// @Failure      400                   {object}  gin.H "Invalid request body"
// @Failure      500                   {object}  gin.H "Failed to create account"
// @Router       /user/account [post]
func (h *AccountHandler) CreateAccountHandler(c *gin.Context) {
	h.logger.Info("CreateAccountHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.account.CreateAccount(c.Request.Context(), &pb.CreateAccountRequest{
		UserId:   userId,
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: req.Currency,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to create account"})
		return
	}

	if _, err := h.redis.StoreAccountInRedis(c.Request.Context(), resp); err != nil {
		h.logger.Error("Error storing account in Redis:", slog.String("err: ", err.Error()))
	}

	c.IndentedJSON(201, resp)
}

// GetAccountsHandler godoc
// @Summary      Get accounts
// @Security     BearerAuth
// @Description  Get all accounts for the authenticated user
// @Tags         User Accounts
// @Produce      json
// @Success      200  {object}  pb.AccountsResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to get accounts"
// @Router       /user/account [get]
func (h *AccountHandler) GetAccountsHandler(c *gin.Context) {
	h.logger.Info("GetAccountsHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.account.GetAccounts(c.Request.Context(), &pb.GetAccountsRequest{
		UserId: userId,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get accounts"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetAccountByIdHandler godoc
// @Summary      Get account by ID
// @Security     BearerAuth
// @Description  Get account details by account ID
// @Tags         User Accounts
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  pb.AccountResponse
// @Failure      400  {object}  gin.H "Account ID is required"
// @Failure      500  {object}  gin.H "Failed to retrieve account"
// @Router       /user/account/{id} [get]
func (h *AccountHandler) GetAccountByIdHandler(c *gin.Context) {
	h.logger.Info("GetAccountByIdHandler")

	accountID := c.Param("id")
	if accountID == "" {
		c.IndentedJSON(400, gin.H{"error": "Account ID is required"})
		return
	}
	
	acc, err := h.redis.GetAccountFromRedis(c.Request.Context(), accountID)
	if err != nil {
		h.logger.Error("Error getting account from Redis:", slog.String("err: ", err.Error()))
	}
	if acc != nil {
		h.logger.Info("Account found in Redis")
		c.IndentedJSON(200, acc)
	}

	resp, err := h.account.GetAccountById(c.Request.Context(), &pb.GetAccountByIdRequest{
		Id: accountID,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to retrieve account"})
		return
	}

	c.IndentedJSON(200, resp)
}

// UpdateAccountHandler godoc
// @Summary      Update account
// @Security     BearerAuth
// @Description  Update account details by account ID
// @Tags         User Accounts
// @Accept       json
// @Produce      json
// @Param        UpdateAccountRequest  body      pb.UpdateAccountRequest  true  "Updated account details"
// @Success      200                   {object}  pb.AccountResponse
// @Failure      400                   {object}  gin.H "Invalid request body"
// @Failure      500                   {object}  gin.H "Failed to update account"
// @Router       /user/account [put]
func (h *AccountHandler) UpdateAccountHandler(c *gin.Context) {
	h.logger.Info("UpdateAccountHandler")

	var req pb.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.account.UpdateAccount(c.Request.Context(), &req)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to update account"})
		return
	}

	c.IndentedJSON(200, resp)
}

// DeleteAccountHandler godoc
// @Summary      Delete account
// @Security     BearerAuth
// @Description  Delete account by account ID
// @Tags         User Accounts
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  gin.H "message: Account deleted successfully"
// @Failure      400  {object}  gin.H "Account ID is required"
// @Failure      500  {object}  gin.H "Failed to delete account"
// @Router       /user/account/{id} [delete]
func (h *AccountHandler) DeleteAccountHandler(c *gin.Context) {
	h.logger.Info("DeleteAccountHandler")

	accountID := c.Param("id")
	if accountID == "" {
		c.IndentedJSON(400, gin.H{"error": "Account ID is required"})
		return
	}

	_, err := h.account.DeleteAccount(c.Request.Context(), &pb.DeleteAccountRequest{
		Id: accountID,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to delete account"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Account deleted successfully"})
}

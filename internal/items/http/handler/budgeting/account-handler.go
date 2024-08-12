package budgeting

import (
	pb "gateway-service/genproto/account"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
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

// CreateAccountHandler godoc
// @Summary      Create an account
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

	UserID, exists := c.Get("user_id")
	if !exists {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(500, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.account.CreateAccount(c.Request.Context(), &pb.CreateAccountRequest{
		UserId:   userIDStr,
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: req.Currency,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to create account"})
		return
	}

	c.IndentedJSON(201, resp)
}

// GetAccountsHandler godoc
// @Summary      Get accounts
// @Description  Get all accounts for the authenticated user
// @Tags         User Accounts
// @Produce      json
// @Success      200  {object}  pb.AccountsResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to get accounts"
// @Router       /user/account [get]
func (h *AccountHandler) GetAccountsHandler(c *gin.Context) {
	h.logger.Info("GetAccountsHandler")

	UserID, exists := c.Get("user_id")
	if !exists {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := UserID.(string)
	if !ok {
		c.IndentedJSON(500, gin.H{"error": "Invalid user ID format"})
		return
	}

	resp, err := h.account.GetAccounts(c.Request.Context(), &pb.GetAccountsRequest{
		UserId: userIDStr,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get accounts"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetAccountByIdHandler godoc
// @Summary      Get account by ID
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

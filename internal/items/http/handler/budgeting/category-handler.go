package budgeting

import (
	pb "gateway-service/genproto/category"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"gateway-service/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	category  pb.CategoryServiceClient
	logger    *slog.Logger
	msgbroker *msgbroker.MsgBroker
	config    *config.Config
}

func NewCategoryHandler(category pb.CategoryServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *CategoryHandler {
	return &CategoryHandler{
		category:  category,
		logger:    logger,
		msgbroker: msgbroker,
		config:    config,
	}
}

// CreateCategoryHandler godoc
// @Summary      Create a category
// @Security     BearerAuth
// @Description  Create a new category for the authenticated user
// @Tags         User Categories
// @Accept       json
// @Produce      json
// @Param        CreateCategoryRequest  body      models.CreateCategoryRequest  true  "Category details"
// @Success      201                     {object}  pb.CategoryResponse
// @Failure      401                     {object}  gin.H "User not authenticated"
// @Failure      400                     {object}  gin.H "Invalid request body"
// @Failure      500                     {object}  gin.H "Failed to create category"
// @Router       /user/category [post]
func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	h.logger.Info("CreateCategoryHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.category.CreateCategory(c.Request.Context(), &pb.CreateCategoryRequest{
		UserId: userId,
		Name:   req.Name,
		Type:   req.Type,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to create category"})
		return
	}

	c.IndentedJSON(201, resp)
}

// GetCategoriesHandler godoc
// @Summary      Get categories
// @Security     BearerAuth
// @Description  Get all categories for the authenticated user
// @Tags         User Categories
// @Produce      json
// @Success      200  {object}  pb.CategoriesResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to get categories"
// @Router       /user/category [get]
func (h *CategoryHandler) GetCategoriesHandler(c *gin.Context) {
	h.logger.Info("GetCategoriesHandler")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.category.GetCategories(c.Request.Context(), &pb.GetCategoriesRequest{
		UserId: userId,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to get categories"})
		return
	}

	c.IndentedJSON(200, resp)
}

// GetCategoryByIdHandler godoc
// @Summary      Get category by ID
// @Security     BearerAuth
// @Description  Get category details by category ID
// @Tags         User Categories
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  pb.CategoryResponse
// @Failure      400  {object}  gin.H "Category ID is required"
// @Failure      500  {object}  gin.H "Failed to retrieve category"
// @Router       /user/category/{id} [get]
func (h *CategoryHandler) GetCategoryByIdHandler(c *gin.Context) {
	h.logger.Info("GetCategoryByIdHandler")

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Category ID is required"})
		return
	}

	resp, err := h.category.GetCategoryById(c.Request.Context(), &pb.GetCategoryByIdRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to retrieve category"})
		return
	}

	c.IndentedJSON(200, resp)
}

// UpdateCategoryHandler godoc
// @Summary      Update category
// @Security     BearerAuth
// @Description  Update category details by category ID
// @Tags         User Categories
// @Accept       json
// @Produce      json
// @Param        UpdateCategoryRequest  body      pb.UpdateCategoryRequest  true  "Updated category details"
// @Success      200                     {object}  pb.CategoryResponse
// @Failure      400                     {object}  gin.H "Invalid request body"
// @Failure      500                     {object}  gin.H "Failed to update category"
// @Router       /user/category [put]
func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	h.logger.Info("UpdateCategoryHandler")

	var req pb.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.category.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to update category"})
		return
	}

	c.IndentedJSON(200, resp)
}

// DeleteCategoryHandler godoc
// @Summary      Delete category
// @Security     BearerAuth
// @Description  Delete category by category ID
// @Tags         User Categories
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  gin.H "message: Category deleted successfully"
// @Failure      400  {object}  gin.H "Category ID is required"
// @Failure      500  {object}  gin.H "Failed to delete category"
// @Router       /user/category/{id} [delete]
func (h *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {
	h.logger.Info("DeleteCategoryHandler")
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Category ID is required"})
		return
	}

	_, err := h.category.DeleteCategory(c.Request.Context(), &pb.DeleteCategoryRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to delete category"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Category deleted successfully"})
}

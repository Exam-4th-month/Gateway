package budgeting

import (
	pb "gateway-service/genproto/notification"
	"gateway-service/internal/items/config"
	"gateway-service/internal/items/middleware"
	"gateway-service/internal/items/msgbroker"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notification pb.NotificationServiceClient
	logger       *slog.Logger
	msgbroker    *msgbroker.MsgBroker
	config       *config.Config
}

func NewNotificationHandler(notification pb.NotificationServiceClient, logger *slog.Logger, msgbroker *msgbroker.MsgBroker, config *config.Config) *NotificationHandler {
	return &NotificationHandler{
		notification: notification,
		logger:       logger,
		msgbroker:    msgbroker,
		config:       config,
	}
}

// GetNotifications godoc
// @Summary      Get user notifications
// @Security     BearerAuth
// @Description  Retrieve all notifications for the authenticated user
// @Tags         User Notifications
// @Accept       json
// @Produce      json
// @Success      200  {object}  pb.NotificationsResponse
// @Failure      401  {object}  gin.H "User not authenticated"
// @Failure      500  {object}  gin.H "Failed to retrieve notifications"
// @Router       /user/notification/ [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	h.logger.Info("GetNotifications")

	userId := middleware.GetUser_id(c, h.config)
	if userId == "" {
		c.IndentedJSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.notification.GetNotifications(c, &pb.GetNotificationsRequest{
		UserId: userId,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	c.IndentedJSON(200, resp)
}

// MarkNotificationAsRead godoc
// @Summary      Mark a notification as read
// @Security     BearerAuth
// @Description  Mark the specified notification as read
// @Tags         User Notifications
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Notification ID"
// @Success      200  {object}  gin.H
// @Failure      400  {object}  gin.H "Notification ID is required"
// @Failure      500  {object}  gin.H "Failed to mark notification as read"
// @Router       /user/notification/{id} [put]
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	h.logger.Info("MarkNotificationAsRead")

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(400, gin.H{"error": "Notification ID is required"})
		return
	}

	_, err := h.notification.MarkNotificationAsRead(c, &pb.MarkNotificationAsReadRequest{
		Id: id,
	})
	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Notification marked as read"})
}

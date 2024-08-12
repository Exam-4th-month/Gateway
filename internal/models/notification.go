package models

type GetNotificationsRequest struct{}

type MarkNotificationAsReadRequest struct {
    ID string `json:"id"`
}

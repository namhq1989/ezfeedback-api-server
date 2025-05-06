package dto

type CountNotificationsRequest struct{}

type CountNotificationResponse struct {
	Total int64 `json:"total"`
}

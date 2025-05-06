package dto

type GetNotificationsRequest struct {
	Page int64 `query:"page"`
}

type GetNotificationResponse struct {
	Notifications []Notification `json:"notifications"`
	Limit         int64          `json:"limit"`
}

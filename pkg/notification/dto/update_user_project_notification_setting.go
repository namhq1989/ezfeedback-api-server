package dto

type UpdateUserProjectNotificationSettingRequest struct {
	ProjectID            string `json:"projectId"`
	ReceiveNewFeedback   bool   `json:"receiveNewFeedback"`
	ReceiveDailySummary  bool   `json:"receiveDailySummary"`
	ReceiveWeeklySummary bool   `json:"receiveWeeklySummary"`
}

type UpdateUserProjectNotificationSettingResponse struct{}

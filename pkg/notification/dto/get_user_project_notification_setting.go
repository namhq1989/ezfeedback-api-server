package dto

type GetUserProjectNotificationSettingRequest struct {
	ProjectID string `query:"projectId"`
}

type GetUserProjectNotificationSettingResponse struct {
	Setting UserProjectNotificationSetting `json:"setting"`
}

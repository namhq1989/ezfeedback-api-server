package dto

type UpdateProjectCampaignRequest struct {
	Name           string `json:"name" validate:"required" message:"invalid_name"`
	Description    string `json:"description"`
	WidgetPosition string `json:"widgetPosition"`
}

type UpdateProjectCampaignResponse struct{}

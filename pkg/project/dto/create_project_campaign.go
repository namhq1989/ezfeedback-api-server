package dto

type CreateProjectCampaignRequest struct {
	Name           string `json:"name" validate:"required" message:"invalid_name"`
	Description    string `json:"description"`
	CampaignType   string `json:"campaignType"`
	WidgetPosition string `json:"widgetPosition"`
}

type CreateProjectCampaignResponse struct {
	ID string `json:"id"`
}

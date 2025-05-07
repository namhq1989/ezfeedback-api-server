package dto

type UpdateProjectCampaignRequest struct {
	Name     string                 `json:"name" validate:"required" message:"invalid_name"`
	Settings ProjectCampaignSetting `json:"settings"`
}

type UpdateProjectCampaignResponse struct{}

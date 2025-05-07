package dto

type CreateProjectCampaignRequest struct {
	Name         string                 `json:"name" validate:"required" message:"invalid_name"`
	CampaignType string                 `json:"campaignType"`
	Settings     ProjectCampaignSetting `json:"settings"`
}

type CreateProjectCampaignResponse struct {
	ID string `json:"id"`
}

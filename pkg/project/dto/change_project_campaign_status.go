package dto

type ChangeProjectCampaignStatusRequest struct {
	Status string `json:"status" validate:"required" message:"invalid_status"`
}

type ChangeProjectCampaignStatusResponse struct{}

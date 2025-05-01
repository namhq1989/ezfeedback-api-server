package dto

type CreateFeedbackRequest struct {
	CampaignID string                       `json:"campaignId" validate:"required" message:"project_invalid_campaign"`
	Email      string                       `json:"email"`
	CategoryID string                       `json:"categoryId"`
	Content    string                       `json:"content"`
	Rating     int32                        `json:"rating"`
	Context    CreateFeedbackRequestContext `json:"context"`
}

type CreateFeedbackRequestContext struct {
	UserID string `json:"userId"`
}

type CreateFeedbackResponse struct{}

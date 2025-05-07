package dto

type UpdateFeedbackReplyRequest struct {
	Content string `json:"content" validate:"required" message:"feedback_invalid_content"`
}

type UpdateFeedbackReplyResponse struct{}

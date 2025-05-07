package dto

type CreateFeedbackReplyRequest struct {
	Content string `json:"content" validate:"required" message:"feedback_invalid_content"`
}

type CreateFeedbackReplyResponse struct{}

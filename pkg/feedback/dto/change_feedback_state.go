package dto

type ChangeFeedbackStateRequest struct {
	State string `json:"state" validate:"required" message:"feedback_invalid_state"`
}

type ChangeFeedbackStateResponse struct{}

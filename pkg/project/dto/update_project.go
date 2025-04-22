package dto

type UpdateProjectRequest struct {
	Title                  string `json:"title" validate:"required" message:"invalid_title"`
	Description            string `json:"description"`
	IsFeedbackPublic       bool   `json:"isFeedbackPublic"`
	AllowAnonymousFeedback bool   `json:"allowAnonymousFeedback"`
	EnableVoting           bool   `json:"enableVoting"`
}

type UpdateProjectResponse struct{}

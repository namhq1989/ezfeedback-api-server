package dto

type UpdateProjectRequest struct {
	Title        string `json:"title" validate:"required" message:"invalid_title"`
	Description  string `json:"description"`
	Domain       string `json:"domain"`
	PrimaryColor string `json:"primaryColor"`
}

type UpdateProjectResponse struct{}

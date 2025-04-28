package dto

type CreateProjectRequest struct {
	Title        string `json:"title" validate:"required" message:"invalid_title"`
	Description  string `json:"description"`
	Domain       string `json:"domain"`
	PrimaryColor string `json:"primaryColor"`
}

type CreateProjectResponse struct {
	ID string `json:"id"`
}

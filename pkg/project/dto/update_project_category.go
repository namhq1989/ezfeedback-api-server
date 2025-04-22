package dto

type UpdateProjectCategoryRequest struct {
	Name string `json:"name" validate:"required" message:"invalid_name"`
}

type UpdateProjectCategoryResponse struct{}

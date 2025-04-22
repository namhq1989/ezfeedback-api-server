package dto

type CreateProjectCategoryRequest struct {
	Name string `json:"name" validate:"required" message:"invalid_name"`
}

type CreateProjectCategoryResponse struct {
	ID string `json:"id"`
}

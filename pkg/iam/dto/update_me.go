package dto

type UpdateMeRequest struct {
	Name string `json:"name" validate:"required" message:"invalid_name"`
}

type UpdateMeResponse struct{}

package dto

type ChangeProjectCategoryStatusRequest struct {
	Status string `json:"status" validate:"required" message:"invalid_status"`
}

type ChangeProjectCategoryStatusResponse struct{}

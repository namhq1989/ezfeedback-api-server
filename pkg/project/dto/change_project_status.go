package dto

type ChangeProjectStatusRequest struct {
	Status string `json:"status" validate:"required" message:"invalid_status"`
}

type ChangeProjectStatusResponse struct{}

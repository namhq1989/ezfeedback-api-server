package dto

type RequestVerificationCodeRequest struct {
	Email string `json:"email" validate:"required" message:"invalid_email"`
}

type RequestVerificationCodeResponse struct{}

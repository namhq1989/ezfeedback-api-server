package dto

type VerifyVerificationCodeRequest struct {
	Email string `json:"email" validate:"required" message:"invalid_email"`
	Code  string `json:"code" validate:"required" message:"invalid_code"`
}

type VerifyVerificationCodeResponse struct {
	IsNewUser bool   `json:"isNewUser"`
	Token     string `json:"token"`
}

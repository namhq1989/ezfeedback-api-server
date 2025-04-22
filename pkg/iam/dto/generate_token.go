package dto

type GenerateTokenRequest struct {
	UserID string `json:"userId"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

package dto

type PingRequest struct{}

type PingResponse struct {
	Success bool `json:"success"`
}

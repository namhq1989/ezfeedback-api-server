package dto

type GetMeRequest struct{}

type GetMeResponse struct {
	Me User `json:"me"`
}

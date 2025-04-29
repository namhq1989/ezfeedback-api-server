package dto

type GetProjectByIDRequest struct{}

type GetProjectByIDResponse struct {
	Project Project `json:"project"`
}

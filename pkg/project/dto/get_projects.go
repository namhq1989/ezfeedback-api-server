package dto

type GetProjectsRequest struct{}

type GetProjectsResponse struct {
	Projects []Project `json:"projects"`
}

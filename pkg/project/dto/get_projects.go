package dto

type GetProjectsRequest struct{}

type GetProjectsResponse struct {
	Projects []ProjectBrief `json:"projects"`
}

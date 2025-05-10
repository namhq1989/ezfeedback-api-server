package dto

type GetProjectCollaboratorsRequest struct{}

type GetProjectCollaboratorsResponse struct {
	Collaborators []ProjectCollaborator `json:"collaborators"`
}

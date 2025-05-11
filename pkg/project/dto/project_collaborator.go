package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCollaborator struct {
	ID        string                    `json:"id"`
	User      User                      `json:"user"`
	Role      string                    `json:"role"`
	CreatedAt *httprespond.TimeResponse `json:"createdAt"`
}

func (ProjectCollaborator) FromDomain(collaborator domain.ProjectCollaborator, user domain.User) ProjectCollaborator {
	return ProjectCollaborator{
		ID:        collaborator.ID,
		User:      User{}.FromDomain(user),
		Role:      collaborator.Role.String(),
		CreatedAt: httprespond.NewTimeResponse(collaborator.CreatedAt),
	}
}

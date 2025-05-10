package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCollaboratorMapper struct{}

func (ProjectCollaboratorMapper) FromModelToDomain(collaborator model.ProjectCollaborators) (*domain.ProjectCollaborator, error) {
	var result = &domain.ProjectCollaborator{
		ID:        collaborator.ID,
		ProjectID: collaborator.ProjectID,
		UserID:    collaborator.UserID,
		Role:      domain.ToProjectRole(collaborator.Role.String()),
		CreatedAt: collaborator.CreatedAt,
		UpdatedAt: collaborator.UpdatedAt,
	}

	return result, nil
}

func (ProjectCollaboratorMapper) FromDomainToModel(collaborator domain.ProjectCollaborator) (*model.ProjectCollaborators, error) {
	var result = &model.ProjectCollaborators{
		ID:        collaborator.ID,
		ProjectID: collaborator.ProjectID,
		UserID:    collaborator.UserID,
		Role:      model.ProjectRole(collaborator.Role.String()),
		CreatedAt: collaborator.CreatedAt,
		UpdatedAt: collaborator.UpdatedAt,
	}

	return result, nil
}

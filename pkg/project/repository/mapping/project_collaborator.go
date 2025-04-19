package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCollaboratorMapper struct{}

func (ProjectCollaboratorMapper) FromModelToDomain(projectCollaborator model.ProjectCollaborators) (*domain.ProjectCollaborator, error) {
	var result = &domain.ProjectCollaborator{
		ID:        projectCollaborator.ID,
		ProjectID: projectCollaborator.ProjectID,
		UserID:    projectCollaborator.UserID,
		Role:      domain.ToProjectRole(projectCollaborator.Role.String()),
		CreatedAt: projectCollaborator.CreatedAt,
		UpdatedAt: projectCollaborator.UpdatedAt,
	}

	return result, nil
}

func (ProjectCollaboratorMapper) FromDomainToModel(projectCollaborator domain.ProjectCollaborator) (*model.ProjectCollaborators, error) {
	var result = &model.ProjectCollaborators{
		ID:        projectCollaborator.ID,
		ProjectID: projectCollaborator.ProjectID,
		UserID:    projectCollaborator.UserID,
		Role:      model.ProjectRole(projectCollaborator.Role.String()),
		CreatedAt: projectCollaborator.CreatedAt,
		UpdatedAt: projectCollaborator.UpdatedAt,
	}

	return result, nil
}

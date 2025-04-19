package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectMapper struct{}

func (ProjectMapper) FromModelToDomain(project model.Projects) (*domain.Project, error) {
	var result = &domain.Project{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Description: project.Description,
		Slug:        project.Slug,
		Status:      domain.ToStatus(project.Status.String()),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	return result, nil
}

func (ProjectMapper) FromDomainToModel(project domain.Project) (*model.Projects, error) {
	var result = &model.Projects{
		ID:          project.ID,
		UserID:      project.UserID,
		Title:       project.Title,
		Description: project.Description,
		Slug:        project.Slug,
		Status:      model.Status(project.Status.String()),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	return result, nil
}

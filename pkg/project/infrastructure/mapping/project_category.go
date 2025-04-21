package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCategoryMapper struct{}

func (ProjectCategoryMapper) FromModelToDomain(projectCategory model.ProjectCategories) (*domain.ProjectCategory, error) {
	var result = &domain.ProjectCategory{
		ID:        projectCategory.ID,
		ProjectID: projectCategory.ProjectID,
		Name:      projectCategory.Name,
		Slug:      projectCategory.Slug,
		Status:    domain.ToStatus(projectCategory.Status.String()),
		CreatedAt: projectCategory.CreatedAt,
		UpdatedAt: projectCategory.UpdatedAt,
	}

	return result, nil
}

func (ProjectCategoryMapper) FromDomainToModel(projectCategory domain.ProjectCategory) (*model.ProjectCategories, error) {
	var result = &model.ProjectCategories{
		ID:        projectCategory.ID,
		ProjectID: projectCategory.ProjectID,
		Name:      projectCategory.Name,
		Slug:      projectCategory.Slug,
		Status:    model.Status(projectCategory.Status.String()),
		CreatedAt: projectCategory.CreatedAt,
		UpdatedAt: projectCategory.UpdatedAt,
	}

	return result, nil
}

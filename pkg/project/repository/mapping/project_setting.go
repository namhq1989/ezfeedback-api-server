package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectSettingMapper struct{}

func (ProjectSettingMapper) FromModelToDomain(projectSetting model.ProjectSettings) (*domain.ProjectSetting, error) {
	var result = &domain.ProjectSetting{
		ID:                     projectSetting.ID,
		ProjectID:              projectSetting.ProjectID,
		IsFeedbackPublic:       projectSetting.IsFeedbackPublic,
		AllowAnonymousFeedback: projectSetting.AllowAnonymousFeedback,
		EnableVoting:           projectSetting.EnableVoting,
		CreatedAt:              projectSetting.CreatedAt,
		UpdatedAt:              projectSetting.UpdatedAt,
	}

	return result, nil
}

func (ProjectSettingMapper) FromDomainToModel(projectSetting domain.ProjectSetting) (*model.ProjectSettings, error) {
	var result = &model.ProjectSettings{
		ID:                     projectSetting.ID,
		ProjectID:              projectSetting.ProjectID,
		IsFeedbackPublic:       projectSetting.IsFeedbackPublic,
		AllowAnonymousFeedback: projectSetting.AllowAnonymousFeedback,
		EnableVoting:           projectSetting.EnableVoting,
		CreatedAt:              projectSetting.CreatedAt,
		UpdatedAt:              projectSetting.UpdatedAt,
	}

	return result, nil
}

package shared

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type Service struct {
	projectRepository             domain.ProjectRepository
	projectSettingRepository      domain.ProjectSettingRepository
	projectCategoryRepository     domain.ProjectCategoryRepository
	projectCampaignRepository     domain.ProjectCampaignRepository
	projectCollaboratorRepository domain.ProjectCollaboratorRepository
	cachingRepository             domain.CachingRepository
}

func NewService(
	projectRepository domain.ProjectRepository,
	projectSettingRepository domain.ProjectSettingRepository,
	projectCategoryRepository domain.ProjectCategoryRepository,
	projectCampaignRepository domain.ProjectCampaignRepository,
	projectCollaboratorRepository domain.ProjectCollaboratorRepository,
	cachingRepository domain.CachingRepository,
) Service {
	return Service{
		projectRepository:             projectRepository,
		projectSettingRepository:      projectSettingRepository,
		projectCategoryRepository:     projectCategoryRepository,
		projectCampaignRepository:     projectCampaignRepository,
		projectCollaboratorRepository: projectCollaboratorRepository,
		cachingRepository:             cachingRepository,
	}
}

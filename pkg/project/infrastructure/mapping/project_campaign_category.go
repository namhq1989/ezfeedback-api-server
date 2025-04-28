package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCampaignCategoryMapper struct{}

func (ProjectCampaignCategoryMapper) FromModelToDomain(cm model.ProjectCampaignCategories) (*domain.ProjectCampaignCategory, error) {
	var result = &domain.ProjectCampaignCategory{
		ID:         cm.ID,
		CampaignID: cm.CampaignID,
		CategoryID: cm.CategoryID,
		CreatedAt:  cm.CreatedAt,
	}

	return result, nil
}

func (ProjectCampaignCategoryMapper) FromDomainToModel(cm domain.ProjectCampaignCategory) (*model.ProjectCampaignCategories, error) {
	var result = &model.ProjectCampaignCategories{
		ID:         cm.ID,
		CampaignID: cm.CampaignID,
		CategoryID: cm.CategoryID,
		CreatedAt:  cm.CreatedAt,
	}

	return result, nil
}

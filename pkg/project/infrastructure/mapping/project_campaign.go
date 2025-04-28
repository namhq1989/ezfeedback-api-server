package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCampaignMapper struct{}

func (ProjectCampaignMapper) FromModelToDomain(campaign model.ProjectCampaigns) (*domain.ProjectCampaign, error) {
	var result = &domain.ProjectCampaign{
		ID:                    campaign.ID,
		ProjectID:             campaign.ProjectID,
		Name:                  campaign.Name,
		Description:           campaign.Description,
		CampaignType:          domain.ToProjectCampaignType(campaign.CampaignType.String()),
		Status:                domain.ToStatus(campaign.Status.String()),
		SettingWidgetPosition: campaign.SettingWidgetPosition,
		CreatedAt:             campaign.CreatedAt,
		UpdatedAt:             campaign.UpdatedAt,
	}

	return result, nil
}

func (ProjectCampaignMapper) FromDomainToModel(campaign domain.ProjectCampaign) (*model.ProjectCampaigns, error) {
	var result = &model.ProjectCampaigns{
		ID:                    campaign.ID,
		ProjectID:             campaign.ProjectID,
		Name:                  campaign.Name,
		Description:           campaign.Description,
		CampaignType:          model.CampaignType(campaign.CampaignType.String()),
		Status:                model.Status(campaign.Status.String()),
		SettingWidgetPosition: campaign.SettingWidgetPosition,
		CreatedAt:             campaign.CreatedAt,
		UpdatedAt:             campaign.UpdatedAt,
	}

	return result, nil
}

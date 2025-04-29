package mapping

import (
	"github.com/lib/pq"
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

type ProjectCampaignWithData struct {
	ProjectCampaign model.ProjectCampaigns `alias:"pc"`
	Categories      pq.StringArray         `alias:"pc.categories"`
}

type ProjectCampaignWithDataMapper struct{}

func (ProjectCampaignWithDataMapper) FromModelToDomain(doc ProjectCampaignWithData) (*domain.ProjectCampaign, error) {
	var (
		campaignMapper = ProjectCampaignMapper{}
	)

	campaign, err := campaignMapper.FromModelToDomain(doc.ProjectCampaign)
	if err != nil {
		return nil, err
	}

	campaign.Categories = doc.Categories

	return campaign, nil
}

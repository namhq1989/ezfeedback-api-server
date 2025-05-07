package mapping

import (
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCampaignMapper struct{}

func (ProjectCampaignMapper) FromModelToDomain(campaign model.ProjectCampaigns) (*domain.ProjectCampaign, error) {
	var result = &domain.ProjectCampaign{
		ID:                  campaign.ID,
		ProjectID:           campaign.ProjectID,
		Name:                campaign.Name,
		CampaignType:        domain.ToProjectCampaignType(campaign.CampaignType.String()),
		Status:              domain.ToStatus(campaign.Status.String()),
		StatsTotalFeedbacks: campaign.StatsTotalFeedbacks,
		CreatedAt:           campaign.CreatedAt,
		UpdatedAt:           campaign.UpdatedAt,
	}

	if campaign.Settings != "" {
		if err := json.Unmarshal([]byte(campaign.Settings), &result.Settings); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (ProjectCampaignMapper) FromDomainToModel(campaign domain.ProjectCampaign) (*model.ProjectCampaigns, error) {
	var result = &model.ProjectCampaigns{
		ID:                  campaign.ID,
		ProjectID:           campaign.ProjectID,
		Name:                campaign.Name,
		CampaignType:        model.CampaignType(campaign.CampaignType.String()),
		Status:              model.Status(campaign.Status.String()),
		Settings:            "",
		StatsTotalFeedbacks: campaign.StatsTotalFeedbacks,
		CreatedAt:           campaign.CreatedAt,
		UpdatedAt:           campaign.UpdatedAt,
	}

	if settingsBytes, err := json.Marshal(campaign.Settings); err != nil {
		return nil, err
	} else {
		result.Settings = string(settingsBytes)
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

type ProjectCampaignHubData struct {
	ProjectCampaign model.ProjectCampaigns `alias:"pc"`
	Project         model.Projects         `alias:"p"`
	Setting         model.ProjectSettings  `alias:"ps"`
}

type ProjectCampaignHubDataMapper struct{}

func (ProjectCampaignHubDataMapper) FromModelToDomain(doc ProjectCampaignHubData) (*domain.ProjectCampaignHubData, error) {
	var (
		campaignMapper       = ProjectCampaignMapper{}
		projectMapper        = ProjectMapper{}
		projectSettingMapper = ProjectSettingMapper{}
	)

	campaign, err := campaignMapper.FromModelToDomain(doc.ProjectCampaign)
	if err != nil {
		return nil, err
	}

	project, err := projectMapper.FromModelToDomain(doc.Project)
	if err != nil {
		return nil, err
	}

	projectSetting, err := projectSettingMapper.FromModelToDomain(doc.Setting)
	if err != nil {
		return nil, err
	}

	return &domain.ProjectCampaignHubData{
		ProjectCampaign: *campaign,
		Project:         *project,
		ProjectSetting:  *projectSetting,
	}, nil
}

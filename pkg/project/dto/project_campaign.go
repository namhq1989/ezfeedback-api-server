package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCampaign struct {
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
	Description           string                    `json:"description"`
	CampaignType          string                    `json:"campaignType"`
	Status                string                    `json:"status"`
	SettingWidgetPosition string                    `json:"settingWidgetPosition"`
	CreatedAt             *httprespond.TimeResponse `json:"createdAt"`
	UpdatedAt             *httprespond.TimeResponse `json:"updatedAt"`
	Stats                 ProjectCampaignStats      `json:"stats"`
	Categories            []ProjectCategory         `json:"categories"`
}

type ProjectCampaignStats struct {
	TotalFeedback int `json:"totalFeedback"`
}

func (ProjectCampaign) FromDomain(campaign domain.ProjectCampaign, categories []domain.ProjectCategory) ProjectCampaign {
	var cats = make([]ProjectCategory, len(categories))
	for i, category := range categories {
		cats[i] = ProjectCategory{}.FromDomain(category)
	}

	return ProjectCampaign{
		ID:                    campaign.ID,
		Name:                  campaign.Name,
		Description:           campaign.Description,
		CampaignType:          campaign.CampaignType.String(),
		Status:                campaign.Status.String(),
		SettingWidgetPosition: campaign.SettingWidgetPosition,
		CreatedAt:             httprespond.NewTimeResponse(campaign.CreatedAt),
		UpdatedAt:             httprespond.NewTimeResponse(campaign.UpdatedAt),
		Stats: ProjectCampaignStats{
			TotalFeedback: 10,
		},
		Categories: cats,
	}
}

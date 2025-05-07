package dto

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/httprespond"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
)

type ProjectCampaign struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
	CampaignType string                    `json:"campaignType"`
	Status       string                    `json:"status"`
	Settings     ProjectCampaignSetting    `json:"settings"`
	CreatedAt    *httprespond.TimeResponse `json:"createdAt"`
	UpdatedAt    *httprespond.TimeResponse `json:"updatedAt"`
	Stats        ProjectCampaignStats      `json:"stats"`
	Categories   []ProjectCategory         `json:"categories"`
}

type ProjectCampaignStats struct {
	TotalFeedbacks int32 `json:"totalFeedbacks"`
}

func (ProjectCampaign) FromDomain(campaign domain.ProjectCampaign, categories []domain.ProjectCategory) ProjectCampaign {
	var cats = make([]ProjectCategory, len(categories))
	for i, category := range categories {
		cats[i] = ProjectCategory{}.FromDomain(category)
	}

	return ProjectCampaign{
		ID:           campaign.ID,
		Name:         campaign.Name,
		CampaignType: campaign.CampaignType.String(),
		Status:       campaign.Status.String(),
		Settings:     ProjectCampaignSetting{}.FromDomain(campaign.Settings),
		CreatedAt:    httprespond.NewTimeResponse(campaign.CreatedAt),
		UpdatedAt:    httprespond.NewTimeResponse(campaign.UpdatedAt),
		Stats: ProjectCampaignStats{
			TotalFeedbacks: campaign.StatsTotalFeedbacks,
		},
		Categories: cats,
	}
}

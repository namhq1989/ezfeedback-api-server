package dto

import "github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"

type ProjectCampaign struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (ProjectCampaign) FromDomain(campaign domain.ProjectCampaign) ProjectCampaign {
	return ProjectCampaign{
		ID:   campaign.ID,
		Name: campaign.Name,
	}
}

type ProjectCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (ProjectCategory) FromDomain(category domain.ProjectCategory) ProjectCategory {
	return ProjectCategory{
		ID:   category.ID,
		Name: category.Name,
	}
}

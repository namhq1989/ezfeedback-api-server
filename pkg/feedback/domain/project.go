package domain

import (
	"github.com/namhq1989/go-utilities/appcontext"
)

type ProjectHub interface {
	GetProjectByID(ctx *appcontext.AppContext, projectID string) (*Project, error)
	GetProjectCampaignByID(ctx *appcontext.AppContext, campaignID string) (*ProjectCampaignHubData, error)
}

type ProjectCampaignHubData struct {
	Project         Project
	ProjectCampaign ProjectCampaign
}

type Project struct {
	ID      string
	UserID  string
	Title   string
	Status  Status
	Setting ProjectSetting
}

type ProjectCampaign struct {
	ID           string
	CampaignType ProjectCampaignType
	Status       Status
}

type ProjectSetting struct {
	Domain string
}

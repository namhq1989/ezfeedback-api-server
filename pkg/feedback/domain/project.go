package domain

import (
	"github.com/namhq1989/go-utilities/appcontext"
)

type ProjectHub interface {
	GetProjectByID(ctx *appcontext.AppContext, projectID, userID string) (*Project, error)
	GetProjectCampaignByID(ctx *appcontext.AppContext, campaignID string) (*ProjectCampaignHubData, error)
	GetProjectCategories(ctx *appcontext.AppContext, projectID string) ([]ProjectCategory, error)
	GetProjectCampaigns(ctx *appcontext.AppContext, projectID string) ([]ProjectCampaign, error)
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
	Name         string
	CampaignType ProjectCampaignType
	Status       Status
}

type ProjectCategory struct {
	ID   string
	Name string
}

type ProjectSetting struct {
	Domain string
}

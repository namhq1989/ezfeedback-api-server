package domain

import "github.com/namhq1989/go-utilities/appcontext"

type Service interface {
	GetProjectByID(ctx *appcontext.AppContext, projectID, userID string) (*Project, error)
	GetProjectSettingByProjectID(ctx *appcontext.AppContext, projectID string) (*ProjectSetting, error)
	GetProjectCategoriesByProjectID(ctx *appcontext.AppContext, projectID string, status Status) ([]ProjectCategory, error)
	GetProjectCampaignsByProjectID(ctx *appcontext.AppContext, projectID string) ([]ProjectCampaign, error)
	GetProjectCampaign(ctx *appcontext.AppContext, projectID, campaignID, userID string) (*ProjectCampaign, error)
	GetProjectCategory(ctx *appcontext.AppContext, projectID, categoryID, userID string) (*ProjectCategory, error)
}

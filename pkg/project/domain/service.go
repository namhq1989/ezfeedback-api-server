package domain

import "github.com/namhq1989/go-utilities/appcontext"

type Service interface {
	GetProjectByID(ctx *appcontext.AppContext, id string) (*Project, error)
	GetProjectSettingByProjectID(ctx *appcontext.AppContext, projectID string) (*ProjectSetting, error)
	GetProjectCategoriesByProjectID(ctx *appcontext.AppContext, projectID string, status Status) ([]ProjectCategory, error)
}

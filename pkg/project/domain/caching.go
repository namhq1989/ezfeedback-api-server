package domain

import "github.com/namhq1989/go-utilities/appcontext"

type CachingRepository interface {
	GetProjectByID(ctx *appcontext.AppContext, id string) (*Project, error)
	SetProjectByID(ctx *appcontext.AppContext, id string, project Project) error
	DeleteProjectByID(ctx *appcontext.AppContext, id string) error

	GetProjectSettingByProjectID(ctx *appcontext.AppContext, id string) (*ProjectSetting, error)
	SetProjectSettingByProjectID(ctx *appcontext.AppContext, id string, setting ProjectSetting) error
	DeleteProjectSettingByProjectID(ctx *appcontext.AppContext, id string) error
}

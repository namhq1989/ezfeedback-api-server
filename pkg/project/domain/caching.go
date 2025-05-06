package domain

import "github.com/namhq1989/go-utilities/appcontext"

type CachingRepository interface {
	GetProjectByID(ctx *appcontext.AppContext, id string) (*Project, error)
	SetProjectByID(ctx *appcontext.AppContext, id string, project Project) error
	DeleteProjectByID(ctx *appcontext.AppContext, id string) error

	GetProjectSettingByProjectID(ctx *appcontext.AppContext, id string) (*ProjectSetting, error)
	SetProjectSettingByProjectID(ctx *appcontext.AppContext, id string, setting ProjectSetting) error
	DeleteProjectSettingByProjectID(ctx *appcontext.AppContext, id string) error

	GetProjectCategoriesByProjectID(ctx *appcontext.AppContext, id string) ([]ProjectCategory, error)
	SetProjectCategoriesByProjectID(ctx *appcontext.AppContext, id string, categories []ProjectCategory) error
	DeleteProjectCategoriesByProjectID(ctx *appcontext.AppContext, id string) error

	GetProjectCampaignsByProjectID(ctx *appcontext.AppContext, id string) ([]ProjectCampaign, error)
	SetProjectCampaignsByProjectID(ctx *appcontext.AppContext, id string, campaigns []ProjectCampaign) error
	DeleteProjectCampaignsByProjectID(ctx *appcontext.AppContext, id string) error

	GetProjectCollaboratorsByProjectID(ctx *appcontext.AppContext, id string) ([]ProjectCollaborator, error)
	SetProjectCollaboratorsByProjectID(ctx *appcontext.AppContext, id string, collaborators []ProjectCollaborator) error
	DeleteProjectCollaboratorsByProjectID(ctx *appcontext.AppContext, id string) error

	GetApiGetProjectsByUserID(ctx *appcontext.AppContext, userID string) (*string, error)
	SetApiGetProjectsByUserID(ctx *appcontext.AppContext, userID string, data string) error
	DeleteApiGetProjectsByUserID(ctx *appcontext.AppContext, userID string) error

	GetApiGetProjectByID(ctx *appcontext.AppContext, id string) (*string, error)
	SetApiGetProjectByID(ctx *appcontext.AppContext, id string, data string) error
	DeleteApiGetProjectByID(ctx *appcontext.AppContext, id string) error
}

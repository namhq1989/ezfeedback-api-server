package domain

import "github.com/namhq1989/go-utilities/appcontext"

type CachingRepository interface {
	GetUserProjectNotificationSetting(ctx *appcontext.AppContext, userID, projectID string) (*UserProjectNotificationSetting, error)
	SetUserProjectNotificationSetting(ctx *appcontext.AppContext, userID, projectID string, setting UserProjectNotificationSetting) error
	DeleteProjectByID(ctx *appcontext.AppContext, userID, projectID string) error
}

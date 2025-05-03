package domain

import "github.com/namhq1989/go-utilities/appcontext"

type Service interface {
	GetUserProjectNotificationSetting(ctx *appcontext.AppContext, userID, projectID string) (*UserProjectNotificationSetting, error)
}

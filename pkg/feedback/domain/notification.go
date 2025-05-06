package domain

import "github.com/namhq1989/go-utilities/appcontext"

type NotificationHub interface {
	CreateNotificationReminder(ctx *appcontext.AppContext, projectID string) error
}

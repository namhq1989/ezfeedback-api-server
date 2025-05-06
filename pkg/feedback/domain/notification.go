package domain

import (
	"github.com/namhq1989/go-utilities/appcontext"
)

type NotificationHub interface {
	CreateNewFeedbackNotificationDocument(ctx *appcontext.AppContext, userID string, metadata NotificationMetadata) error
	CreateNotificationReminder(ctx *appcontext.AppContext, projectID string) error
}

type NotificationMetadata struct {
	ProjectID    string
	ProjectTitle string
}

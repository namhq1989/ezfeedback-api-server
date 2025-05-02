package domain

import (
	"time"

	"github.com/namhq1989/go-utilities/appcontext"
)

type NotificationReminderRepository interface {
	Create(ctx *appcontext.AppContext, reminder NotificationReminder) error
	Delete(ctx *appcontext.AppContext, reminder NotificationReminder) error
	FindAllExisting(ctx *appcontext.AppContext) ([]NotificationReminder, error)
	CleanupStale(ctx *appcontext.AppContext) error
}

type NotificationReminder struct {
	ID        string
	ProjectID string
	CreatedAt time.Time
}

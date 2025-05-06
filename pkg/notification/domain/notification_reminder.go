package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type NotificationReminderRepository interface {
	Create(ctx *appcontext.AppContext, reminder NotificationReminder) error
	Delete(ctx *appcontext.AppContext, reminder NotificationReminder) error
	FindByProjectID(ctx *appcontext.AppContext, projectID string) (*NotificationReminder, error)
	FindAllExisting(ctx *appcontext.AppContext) ([]NotificationReminder, error)
	CleanupStale(ctx *appcontext.AppContext) error
}

type NotificationReminder struct {
	ID        string
	ProjectID string
	CreatedAt time.Time
}

func NewNotificationReminder(projectID string) (*NotificationReminder, error) {
	if !uuid.IsValidID(projectID) {
		return nil, apperrors.Project.InvalidProjectID
	}

	var nr = &NotificationReminder{
		ID:        uuid.New(),
		ProjectID: projectID,
		CreatedAt: manipulation.NowUTC(),
	}
	return nr, nil
}

package domain

import (
	"time"

	"github.com/namhq1989/go-utilities/appcontext"
)

type FeedbackHub interface {
	GetProjectStatsForNotificationReminder(ctx *appcontext.AppContext, projectID string, timestamp time.Time, limit int64) (int64, []Feedback, error)
}

const (
	NumOfFeedbacksForTheReminderStats int64 = 2
)

type Feedback struct {
	ID           string
	Project      Project
	AppUserID    *string
	Email        *string
	Content      string
	Rating       int32
	CampaignType ProjectCampaignType
	CreatedAt    time.Time
}

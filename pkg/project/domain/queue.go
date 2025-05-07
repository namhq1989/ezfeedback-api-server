package domain

import "github.com/namhq1989/go-utilities/appcontext"

type QueueRepository interface {
	OnProjectCreated(ctx *appcontext.AppContext, payload QueueOnProjectCreatedPayload) error
	OnFeedbackCreated(ctx *appcontext.AppContext, payload QueueOnFeedbackCreatedPayload) error
}

type QueueOnProjectCreatedPayload struct {
	UserID    string
	ProjectID string
}

type QueueOnFeedbackCreatedPayload struct {
	ProjectID  string
	CampaignID string
}

package domain

import "github.com/namhq1989/go-utilities/appcontext"

type QueueRepository interface {
	OnFeedbackCreated(ctx *appcontext.AppContext, payload QueueOnFeedbackCreatedPayload) error
}

type QueueOnFeedbackCreatedPayload struct {
	ProjectID  string
	CampaignID string
}

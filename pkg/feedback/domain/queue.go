package domain

import "github.com/namhq1989/go-utilities/appcontext"

type QueueRepository interface {
	FeedbackCreated(ctx *appcontext.AppContext, payload QueueFeedbackCreatedPayload) error
}

type QueueFeedbackCreatedPayload struct {
	Feedback Feedback
}

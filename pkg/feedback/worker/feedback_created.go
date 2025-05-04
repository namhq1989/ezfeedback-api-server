package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type FeedbackCreatedHandler struct {
}

func NewFeedbackCreatedHandler() FeedbackCreatedHandler {
	return FeedbackCreatedHandler{}
}

func (h FeedbackCreatedHandler) FeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueFeedbackCreatedPayload) error {
	tracer := otel.Tracer("[tracer] feedback created")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] feedback created")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Print("payload", payload)

	// call Notification service to create a notification
	// call Notification service to create a notification reminder
	// check & create project user

	return nil
}

package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type FeedbackCreatedHandler struct {
	notificationHub domain.NotificationHub
}

func NewFeedbackCreatedHandler(notificationHub domain.NotificationHub) FeedbackCreatedHandler {
	return FeedbackCreatedHandler{
		notificationHub: notificationHub,
	}
}

func (h FeedbackCreatedHandler) FeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueFeedbackCreatedPayload) error {
	tracer := otel.Tracer("[tracer] feedback created")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] feedback created")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Text("create notification reminder")
	if err := h.notificationHub.CreateNotificationReminder(ctx, payload.Feedback.ProjectID); err != nil {
		ctx.Logger().Error("failed to create notification reminder", err, appcontext.Fields{})
	}

	// call Notification service to create a notification
	// check & create project user

	return nil
}

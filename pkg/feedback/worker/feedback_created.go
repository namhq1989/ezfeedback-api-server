package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type FeedbackCreatedHandler struct {
	projectHub      domain.ProjectHub
	notificationHub domain.NotificationHub
}

func NewFeedbackCreatedHandler(
	projectHub domain.ProjectHub,
	notificationHub domain.NotificationHub,
) FeedbackCreatedHandler {
	return FeedbackCreatedHandler{
		projectHub:      projectHub,
		notificationHub: notificationHub,
	}
}

func (h FeedbackCreatedHandler) FeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueFeedbackCreatedPayload) error {
	tracer := otel.Tracer("[tracer] feedback created")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] feedback created")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Text("find project data via grpc")
	project, err := h.projectHub.GetProjectByID(ctx, payload.Feedback.ProjectID)
	if err != nil {
		ctx.Logger().Error("failed to find project data via grpc", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("create notification reminder")
	if err = h.notificationHub.CreateNotificationReminder(ctx, payload.Feedback.ProjectID); err != nil {
		ctx.Logger().Error("failed to create notification reminder", err, appcontext.Fields{})
	}

	ctx.Logger().Text("create notification document")
	if err = h.notificationHub.CreateNewFeedbackNotificationDocument(ctx, project.UserID, domain.NotificationMetadata{
		ProjectID:    project.ID,
		ProjectTitle: project.Title,
	}); err != nil {
		ctx.Logger().Error("failed to create notification document", err, appcontext.Fields{})
	}

	return nil
}

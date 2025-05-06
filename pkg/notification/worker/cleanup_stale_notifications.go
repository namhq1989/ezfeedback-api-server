package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type CleanupStaleNotificationsHandler struct {
	notificationRepository domain.NotificationRepository
}

func NewCleanupStaleNotificationsHandler(notificationRepository domain.NotificationRepository) CleanupStaleNotificationsHandler {
	return CleanupStaleNotificationsHandler{
		notificationRepository: notificationRepository,
	}
}

func (h CleanupStaleNotificationsHandler) CleanupStaleNotifications(ctx *appcontext.AppContext, _ domain.QueueCleanupStaleNotificationsPayload) error {
	tracer := otel.Tracer("[tracer] cleanup stale notifications")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] cleanup stale notifications")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Text("delete in db")
	return h.notificationRepository.CleanupStale(ctx)
}

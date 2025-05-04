package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type ScanNotificationRemindersHandler struct {
	notificationReminderRepository domain.NotificationReminderRepository
	queueRepository                domain.QueueRepository
}

func NewScanNotificationRemindersHandler(
	notificationReminderRepository domain.NotificationReminderRepository,
	queueRepository domain.QueueRepository,
) ScanNotificationRemindersHandler {
	return ScanNotificationRemindersHandler{
		notificationReminderRepository: notificationReminderRepository,
		queueRepository:                queueRepository,
	}
}

func (h ScanNotificationRemindersHandler) ScanNotificationReminders(ctx *appcontext.AppContext, _ domain.QueueScanNotificationRemindersPayload) error {
	// ctx.Logger().Text("find reminders in db")
	reminders, err := h.notificationReminderRepository.FindAllExisting(ctx)
	if err != nil {
		ctx.Logger().Error("failed to find reminders in db", err, appcontext.Fields{})
		return err
	}
	if len(reminders) == 0 {
		// ctx.Logger().Text("no reminders found, skip")
		return nil
	}

	tracer := otel.Tracer("[tracer] scan notification reminder")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] scan notification reminder")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Info("found reminders", appcontext.Fields{"count": len(reminders)})
	for _, reminder := range reminders {
		ctx.Logger().Info("add reminder to queue", appcontext.Fields{"reminderId": reminder.ID, "projectId": reminder.ProjectID})
		if err = h.queueRepository.ProcessNotificationReminder(ctx, domain.QueueProcessNotificationReminderPayload{
			Reminder: reminder,
		}); err != nil {
			ctx.Logger().Error("failed to add reminder to queue", err, appcontext.Fields{"reminderId": reminder.ID})
			return err
		}
	}

	ctx.Logger().Text("done scan notification reminder")
	return nil
}

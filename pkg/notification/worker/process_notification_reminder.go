package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type ProcessNotificationReminderHandler struct {
	notificationReminderRepository domain.NotificationReminderRepository
	mailerRepository               domain.MailerRepository
	feedbackHub                    domain.FeedbackHub
	projectHub                     domain.ProjectHub
}

func NewProcessNotificationReminderHandler(
	notificationReminderRepository domain.NotificationReminderRepository,
	mailerRepository domain.MailerRepository,
	feedbackHub domain.FeedbackHub,
	projectHub domain.ProjectHub,
) ProcessNotificationReminderHandler {
	return ProcessNotificationReminderHandler{
		notificationReminderRepository: notificationReminderRepository,
		mailerRepository:               mailerRepository,
		feedbackHub:                    feedbackHub,
		projectHub:                     projectHub,
	}
}

func (h ProcessNotificationReminderHandler) ProcessNotificationReminder(ctx *appcontext.AppContext, payload domain.QueueProcessNotificationReminderPayload) error {
	tracer := otel.Tracer("[tracer] process notification reminder")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] process notification reminder")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Info("process notification reminder", appcontext.Fields{"reminderId": payload.Reminder.ID})

	ctx.Logger().Text("get project feedback stats")
	totalFeedbacks, feedbacks, err := h.feedbackHub.GetProjectStatsForNotificationReminder(ctx, payload.Reminder.ProjectID, payload.Reminder.CreatedAt, domain.NumOfFeedbacksForTheReminderStats)
	if err != nil {
		ctx.Logger().Error("failed to get project feedback stats", err, appcontext.Fields{})
		return err
	}
	if totalFeedbacks == 0 || len(feedbacks) == 0 {
		ctx.Logger().ErrorText("total feedbacks is 0, or no feedbacks found, this reminder is invalid, delete then respond")
		if err = h.notificationReminderRepository.Delete(ctx, payload.Reminder); err != nil {
			ctx.Logger().Error("failed to delete notification reminder", err, appcontext.Fields{})
			return err
		}
		return nil
	}

	ctx.Logger().Text("find project collaborators")
	collaborators, err := h.projectHub.GetProjectCollaborators(ctx, payload.Reminder.ProjectID)
	if err != nil {
		ctx.Logger().Error("failed to find project collaborators", err, appcontext.Fields{})
		return err
	}
	if len(collaborators) == 0 {
		ctx.Logger().ErrorText("no collaborators found, this reminder is invalid, delete then respond")
		if err = h.notificationReminderRepository.Delete(ctx, payload.Reminder); err != nil {
			ctx.Logger().Error("failed to delete notification reminder", err, appcontext.Fields{})
			return err
		}
		return nil
	}

	ctx.Logger().Text("send email to collaborators")
	for _, collaborator := range collaborators {
		var user = collaborator.User
		if err = h.mailerRepository.SendNewFeedbackEmail(ctx, user.Email, feedbacks); err != nil {
			ctx.Logger().Error("failed to send email", err, appcontext.Fields{})
			continue
		}
	}

	ctx.Logger().Text("delete notification reminder")
	if err = h.notificationReminderRepository.Delete(ctx, payload.Reminder); err != nil {
		ctx.Logger().Error("failed to delete notification reminder", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done process notification reminder")
	return nil
}

package infrastructure

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type NotificationHub struct {
	client notificationpb.NotificationServiceClient
}

func NewNotificationHub(client notificationpb.NotificationServiceClient) NotificationHub {
	return NotificationHub{
		client: client,
	}
}

func (r NotificationHub) CreateNewFeedbackNotificationDocument(ctx *appcontext.AppContext, userID string, metadata domain.NotificationMetadata) error {
	_, err := r.client.CreateNewFeedbackNotificationDocument(ctx.Context(), &notificationpb.CreateNewFeedbackNotificationDocumentRequest{
		TraceId: ctx.GetTraceID(),
		UserId:  userID,
		Metadata: &notificationpb.CreateNewFeedbackNotificationDocumentMetadata{
			ProjectId:    metadata.ProjectID,
			ProjectTitle: metadata.ProjectTitle,
		},
	})
	return err
}

func (r NotificationHub) CreateNotificationReminder(ctx *appcontext.AppContext, projectID string) error {
	_, err := r.client.CreateNotificationReminder(ctx.Context(), &notificationpb.CreateNotificationReminderRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
	})
	return err
}

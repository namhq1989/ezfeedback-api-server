package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		CreateNewFeedbackNotificationDocument(ctx *appcontext.AppContext, req *notificationpb.CreateNewFeedbackNotificationDocumentRequest) (*notificationpb.CreateNewFeedbackNotificationDocumentResponse, error)
		CreateNotificationReminder(ctx *appcontext.AppContext, req *notificationpb.CreateNotificationReminderRequest) (*notificationpb.CreateNotificationReminderResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		CreateNewFeedbackNotificationDocumentHandler
		CreateNotificationReminderHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	notificationRepository domain.NotificationRepository,
	notificationReminderRepository domain.NotificationReminderRepository,
) *Application {
	return &Application{
		appHubHandler: appHubHandler{
			CreateNewFeedbackNotificationDocumentHandler: NewCreateNewFeedbackNotificationDocumentHandler(notificationRepository),
			CreateNotificationReminderHandler:            NewCreateNotificationReminderHandler(notificationReminderRepository),
		},
	}
}

package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		CreateNotificationReminder(ctx *appcontext.AppContext, req *notificationpb.CreateNotificationReminderRequest) (*notificationpb.CreateNotificationReminderResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		CreateNotificationReminderHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	notificationReminderRepository domain.NotificationReminderRepository,
) *Application {
	return &Application{
		appHubHandler: appHubHandler{
			CreateNotificationReminderHandler: NewCreateNotificationReminderHandler(notificationReminderRepository),
		},
	}
}

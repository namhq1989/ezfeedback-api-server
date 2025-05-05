package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/feedbackpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		GetProjectStatsForNotificationReminder(ctx *appcontext.AppContext, req *feedbackpb.GetProjectStatsForNotificationReminderRequest) (*feedbackpb.GetProjectStatsForNotificationReminderResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		GetProjectStatsForNotificationReminderHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	feedbackRepository domain.FeedbackRepository,
	projectHub domain.ProjectHub,
) *Application {
	return &Application{
		appHubHandler: appHubHandler{
			GetProjectStatsForNotificationReminderHandler: NewGetProjectStatsForNotificationReminderHandler(feedbackRepository, projectHub),
		},
	}
}

package application

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Commands interface {
		UpdateUserProjectNotificationSetting(ctx *appcontext.AppContext, performerID string, req dto.UpdateUserProjectNotificationSettingRequest) (*dto.UpdateUserProjectNotificationSettingResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)

		GetNotifications(ctx *appcontext.AppContext, performerID string, req dto.GetNotificationsRequest) (*dto.GetNotificationResponse, error)
		GetUserProjectNotificationSetting(ctx *appcontext.AppContext, performerID string, req dto.GetUserProjectNotificationSettingRequest) (*dto.GetUserProjectNotificationSettingResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.UpdateUserProjectNotificationSettingHandler
	}
	queryHandlers struct {
		query.PingHandler

		query.GetNotificationsHandler
		query.GetUserProjectNotificationSettingHandler
	}
	Application struct {
		commandHandlers
		queryHandlers
	}
)

var _ Instance = (*Application)(nil)

func New(
	notificationRepository domain.NotificationRepository,
	userProjectNotificationSettingRepository domain.UserProjectNotificationSettingRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) *Application {
	return &Application{
		commandHandlers: commandHandlers{
			UpdateUserProjectNotificationSettingHandler: command.NewUpdateUserProjectNotificationSettingHandler(
				userProjectNotificationSettingRepository,
				cachingRepository,
				service,
			),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),

			GetNotificationsHandler:                  query.NewGetNotificationsHandler(notificationRepository),
			GetUserProjectNotificationSettingHandler: query.NewGetUserProjectNotificationSettingHandler(service),
		},
	}
}

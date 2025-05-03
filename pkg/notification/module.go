package notification

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/infrastructure"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/rest"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/shared"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Module struct{}

func (Module) Name() string {
	return "NOTIFICATION"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		userProjectNotificationSettingRepository = infrastructure.NewUserProjectNotificationSettingRepository(mono.Database())
		cachingRepository                        = infrastructure.NewCachingRepository(mono.Caching(), mono.Config().IsEnvRelease)

		service = shared.NewService(
			userProjectNotificationSettingRepository,
			cachingRepository,
		)

		app = application.New(
			userProjectNotificationSettingRepository,
			cachingRepository,
			service,
		)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT(), mono.Config().IsEnvRelease); err != nil {
		return err
	}

	return nil
}

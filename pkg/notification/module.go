package notification

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/grpcclient"
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/application"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/grpc"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/infrastructure"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/rest"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/shared"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/worker"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Module struct{}

func (Module) Name() string {
	return "NOTIFICATION"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	projectGRPCClient, err := grpcclient.NewProjectClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	feedbackGRPCClient, err := grpcclient.NewFeedbackClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	var (
		notificationRepository                   = infrastructure.NewNotificationRepository(mono.Database())
		notificationReminderRepository           = infrastructure.NewNotificationReminderRepository(mono.Database())
		userProjectNotificationSettingRepository = infrastructure.NewUserProjectNotificationSettingRepository(mono.Database())
		cachingRepository                        = infrastructure.NewCachingRepository(mono.Caching(), mono.Config().IsEnvRelease)
		queueRepository                          = infrastructure.NewQueueRepository(mono.Queue())
		mailerRepository                         = infrastructure.NewMailerRepository(mono.Mailer())
		projectHub                               = infrastructure.NewProjectHub(projectGRPCClient)
		feedbackHub                              = infrastructure.NewFeedbackHub(feedbackGRPCClient)

		service = shared.NewService(
			userProjectNotificationSettingRepository,
			cachingRepository,
		)

		app = application.New(
			userProjectNotificationSettingRepository,
			cachingRepository,
			service,
		)

		hub = grpc.New(
			notificationRepository,
			notificationReminderRepository,
		)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT(), mono.Config().IsEnvRelease); err != nil {
		return err
	}

	// grpc
	if err := grpc.RegisterServer(ctx, mono.RPC(), hub); err != nil {
		return err
	}

	// worker
	w := worker.New(
		mono.Queue(),
		notificationReminderRepository,
		queueRepository,
		mailerRepository,
		feedbackHub,
		projectHub,
	)
	w.Start()

	return nil
}

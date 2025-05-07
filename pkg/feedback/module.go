package feedback

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/grpcclient"
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/application"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/grpc"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/infrastructure"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/rest"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/shared"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/worker"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Module struct{}

func (Module) Name() string {
	return "FEEDBACK"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	billingGRPCClient, err := grpcclient.NewBillingClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	projectGRPCClient, err := grpcclient.NewProjectClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	notificationGRPCClient, err := grpcclient.NewNotificationClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	iamGRPCClient, err := grpcclient.NewIAMClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	var (
		feedbackRepository      = infrastructure.NewFeedbackRepository(mono.Database())
		feedbackReplyRepository = infrastructure.NewFeedbackReplyRepository(mono.Database())
		queueRepository         = infrastructure.NewQueueRepository(mono.Queue())
		cachingRepository       = infrastructure.NewCachingRepository(mono.Caching(), mono.Config().IsEnvRelease)
		externalAPIRepository   = infrastructure.NewExternalAPIRepository(mono.ExternalAPI())
		feedbackHub             = infrastructure.NewFeedbackHub(mono.Database())
		billingHub              = infrastructure.NewBillingHub(billingGRPCClient)
		projectHub              = infrastructure.NewProjectHub(projectGRPCClient)
		notificationHub         = infrastructure.NewNotificationHub(notificationGRPCClient)
		iamHub                  = infrastructure.NewIAMHub(iamGRPCClient)

		service = shared.NewService(
			feedbackRepository,
			cachingRepository,
			externalAPIRepository,
		)

		app = application.New(
			feedbackRepository,
			feedbackReplyRepository,
			cachingRepository,
			queueRepository,
			billingHub,
			projectHub,
			iamHub,
			service,
		)

		hub = grpc.New(
			feedbackHub,
			projectHub,
		)
	)

	// rest server
	if err = rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT(), mono.Config().IsEnvRelease); err != nil {
		return err
	}

	// grpc
	if err = grpc.RegisterServer(ctx, mono.RPC(), hub); err != nil {
		return err
	}

	// worker
	w := worker.New(
		mono.Queue(),
		projectHub,
		notificationHub,
	)
	w.Start()

	return nil
}

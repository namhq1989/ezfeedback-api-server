package project

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/grpcclient"
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/grpc"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/infrastructure"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/rest"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/shared"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/worker"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Module struct{}

func (Module) Name() string {
	return "PROJECT"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	billingGRPCClient, err := grpcclient.NewBillingClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	iamGRPCClient, err := grpcclient.NewIAMClient(ctx, mono.Config().GRPCPort)
	if err != nil {
		return err
	}

	var (
		projectRepository                 = infrastructure.NewProjectRepository(mono.Database())
		projectSettingRepository          = infrastructure.NewProjectSettingRepository(mono.Database())
		projectCategoryRepository         = infrastructure.NewProjectCategoryRepository(mono.Database())
		projectCampaignRepository         = infrastructure.NewProjectCampaignRepository(mono.Database())
		projectCampaignCategoryRepository = infrastructure.NewProjectCampaignCategoryRepository(mono.Database())
		projectCollaboratorRepository     = infrastructure.NewProjectCollaboratorRepository(mono.Database())
		cachingRepository                 = infrastructure.NewCachingRepository(mono.Caching(), mono.Config().IsEnvRelease)
		queueRepository                   = infrastructure.NewQueueRepository(mono.Queue())
		billingHub                        = infrastructure.NewBillingHub(billingGRPCClient)
		iamHub                            = infrastructure.NewIAMHub(iamGRPCClient)
		projectCampaignHub                = infrastructure.NewProjectCampaignHub(mono.Database())

		service = shared.NewService(
			projectRepository,
			projectSettingRepository,
			projectCategoryRepository,
			projectCampaignRepository,
			projectCollaboratorRepository,
			cachingRepository,
		)

		app = application.New(
			projectRepository,
			projectSettingRepository,
			projectCategoryRepository,
			projectCampaignRepository,
			projectCampaignCategoryRepository,
			cachingRepository,
			queueRepository,
			billingHub,
			iamHub,
			service,
		)

		hub = grpc.New(
			queueRepository,
			projectCampaignHub,
			iamHub,
			service,
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
		projectRepository,
		projectCampaignRepository,
		projectCollaboratorRepository,
		cachingRepository,
	)
	w.Start()

	return nil
}

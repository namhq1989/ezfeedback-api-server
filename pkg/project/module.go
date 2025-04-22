package project

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/grpcclient"
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/infrastructure"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/rest"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/shared"
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

	var (
		projectRepository         = infrastructure.NewProjectRepository(mono.Database())
		projectSettingRepository  = infrastructure.NewProjectSettingRepository(mono.Database())
		projectCategoryRepository = infrastructure.NewProjectCategoryRepository(mono.Database())
		cachingRepository         = infrastructure.NewCachingRepository(mono.Caching(), mono.Config().IsEnvRelease)
		billingHub                = infrastructure.NewBillingHub(billingGRPCClient)

		service = shared.NewService(
			projectRepository,
			projectSettingRepository,
			projectCategoryRepository,
			cachingRepository,
		)

		app = application.New(
			projectRepository,
			projectSettingRepository,
			projectCategoryRepository,
			cachingRepository,
			billingHub,
			service,
		)
	)

	// rest server
	if err = rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT(), mono.Config().IsEnvRelease); err != nil {
		return err
	}

	return nil
}

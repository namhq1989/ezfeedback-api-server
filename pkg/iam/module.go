package iam

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/infrastructure"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/rest"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/worker"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Module struct{}

func (Module) Name() string {
	return "IAM"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		verificationCodeRepository = infrastructure.NewVerificationCodeRepository(mono.Database())
		queueRepository            = infrastructure.NewQueueRepository(mono.Queue())

		app = application.New(
			verificationCodeRepository,
			queueRepository,
		)
	)

	// rest server
	if err := rest.RegisterServer(ctx, app, mono.Rest(), mono.JWT(), mono.Config().IsEnvRelease); err != nil {
		return err
	}

	// worker
	w := worker.New(
		mono.Queue(),
		verificationCodeRepository,
	)
	w.Start()

	return nil
}

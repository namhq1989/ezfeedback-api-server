package billing

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/pkg/billing/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
)

type Module struct{}

func (Module) Name() string {
	return "BILLING"
}

func (Module) Startup(ctx *appcontext.AppContext, mono monolith.Monolith) error {
	var (
		hub = grpc.New()
	)

	// grpc
	if err := grpc.RegisterServer(ctx, mono.RPC(), hub); err != nil {
		return err
	}

	return nil
}

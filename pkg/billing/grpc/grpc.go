package grpc

import (
	"context"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/billingpb"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/grpc"
)

type server struct {
	hub Application
	billingpb.UnimplementedBillingServiceServer
}

var _ billingpb.BillingServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, registrar grpc.ServiceRegistrar, hub *Application) error {
	billingpb.RegisterBillingServiceServer(registrar, server{hub: *hub})
	return nil
}

func (s server) CanCreateProject(bgCtx context.Context, req *billingpb.CanCreateProjectRequest) (*billingpb.CanCreateProjectResponse, error) {
	return s.hub.CanCreateProject(appcontext.NewGRPC(bgCtx), req)
}

func (s server) CanAcceptFeedback(bgCtx context.Context, req *billingpb.CanAcceptFeedbackRequest) (*billingpb.CanAcceptFeedbackResponse, error) {
	return s.hub.CanAcceptFeedback(appcontext.NewGRPC(bgCtx), req)
}

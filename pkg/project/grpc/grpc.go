package grpc

import (
	"context"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/grpc"
)

type server struct {
	hub Application
	projectpb.UnimplementedProjectServiceServer
}

var _ projectpb.ProjectServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, registrar grpc.ServiceRegistrar, hub *Application) error {
	projectpb.RegisterProjectServiceServer(registrar, server{hub: *hub})
	return nil
}

func (s server) GetProjectById(bgCtx context.Context, req *projectpb.GetProjectByIdRequest) (*projectpb.GetProjectByIdResponse, error) {
	return s.hub.GetProjectById(appcontext.NewGRPC(bgCtx), req)
}

func (s server) GetProjectCampaignById(bgCtx context.Context, req *projectpb.GetProjectCampaignByIdRequest) (*projectpb.GetProjectCampaignByIdResponse, error) {
	return s.hub.GetProjectCampaignByID(appcontext.NewGRPC(bgCtx), req)
}

func (s server) GetProjectCollaborators(bgCtx context.Context, req *projectpb.GetProjectCollaboratorsRequest) (*projectpb.GetProjectCollaboratorsResponse, error) {
	return s.hub.GetProjectCollaborators(appcontext.NewGRPC(bgCtx), req)
}

func (s server) GetProjectCategories(bgCtx context.Context, req *projectpb.GetProjectCategoriesRequest) (*projectpb.GetProjectCategoriesResponse, error) {
	return s.hub.GetProjectCategories(appcontext.NewGRPC(bgCtx), req)
}

func (s server) GetProjectCampaigns(bgCtx context.Context, req *projectpb.GetProjectCampaignsRequest) (*projectpb.GetProjectCampaignsResponse, error) {
	return s.hub.GetProjectCampaigns(appcontext.NewGRPC(bgCtx), req)
}

func (s server) OnFeedbackCreated(bgCtx context.Context, req *projectpb.OnFeedbackCreatedRequest) (*projectpb.OnFeedbackCreatedResponse, error) {
	return s.hub.OnFeedbackCreated(appcontext.NewGRPC(bgCtx), req)
}

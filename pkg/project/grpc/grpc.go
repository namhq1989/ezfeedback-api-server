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

func (s server) GetProjectCampaignById(bgCtx context.Context, req *projectpb.GetProjectCampaignByIdRequest) (*projectpb.GetProjectCampaignByIdResponse, error) {
	return s.hub.GetProjectCampaignByID(appcontext.NewGRPC(bgCtx), req)
}

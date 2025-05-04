package grpc

import (
	"context"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/grpc"
)

type server struct {
	hub Application
	iampb.UnimplementedIAMServiceServer
}

var _ iampb.IAMServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, registrar grpc.ServiceRegistrar, hub *Application) error {
	iampb.RegisterIAMServiceServer(registrar, server{hub: *hub})
	return nil
}

func (s server) GetUsersByIds(bgCtx context.Context, req *iampb.GetUsersByIdsRequest) (*iampb.GetUsersByIdsResponse, error) {
	return s.hub.GetUsersByIds(appcontext.NewGRPC(bgCtx), req)
}

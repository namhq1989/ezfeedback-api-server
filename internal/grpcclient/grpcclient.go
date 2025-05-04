package grpcclient

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/billingpb"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
}

func NewBillingClient(_ *appcontext.AppContext, addr string) (billingpb.BillingServiceClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}

	return billingpb.NewBillingServiceClient(conn), nil
}

func NewProjectClient(_ *appcontext.AppContext, addr string) (projectpb.ProjectServiceClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}

	return projectpb.NewProjectServiceClient(conn), nil
}

func NewIAMClient(_ *appcontext.AppContext, addr string) (iampb.IAMServiceClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}

	return iampb.NewIAMServiceClient(conn), nil
}

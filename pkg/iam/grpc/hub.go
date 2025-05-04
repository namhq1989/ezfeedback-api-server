package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/iampb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		GetUsersByIds(ctx *appcontext.AppContext, req *iampb.GetUsersByIdsRequest) (*iampb.GetUsersByIdsResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		GetUsersByIDsHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	service domain.Service,
) *Application {
	return &Application{
		appHubHandler: appHubHandler{
			GetUsersByIDsHandler: NewGetUsersByIDsHandler(service),
		},
	}
}

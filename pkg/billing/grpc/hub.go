package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/billingpb"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		CanCreateProject(ctx *appcontext.AppContext, req *billingpb.CanCreateProjectRequest) (*billingpb.CanCreateProjectResponse, error)
		CanAcceptFeedback(ctx *appcontext.AppContext, req *billingpb.CanAcceptFeedbackRequest) (*billingpb.CanAcceptFeedbackResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		CanCreateProjectHandler
		CanAcceptFeedbackHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New() *Application {
	return &Application{
		appHubHandler: appHubHandler{
			CanCreateProjectHandler:  NewCanCreateProjectHandler(),
			CanAcceptFeedbackHandler: NewCanAcceptFeedbackHandler(),
		},
	}
}

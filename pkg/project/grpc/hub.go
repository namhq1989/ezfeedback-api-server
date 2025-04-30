package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		GetProjectCampaignByID(ctx *appcontext.AppContext, req *projectpb.GetProjectCampaignByIdRequest) (*projectpb.GetProjectCampaignByIdResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		GetProjectCampaignByIDHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	projectCampaignHub domain.ProjectCampaignHub,
) *Application {
	return &Application{
		appHubHandler: appHubHandler{
			GetProjectCampaignByIDHandler: NewGetProjectCampaignByIDHandler(projectCampaignHub),
		},
	}
}

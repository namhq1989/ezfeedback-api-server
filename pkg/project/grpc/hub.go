package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Hubs interface {
		GetProjectById(ctx *appcontext.AppContext, req *projectpb.GetProjectByIdRequest) (*projectpb.GetProjectByIdResponse, error)
		GetProjectCampaignByID(ctx *appcontext.AppContext, req *projectpb.GetProjectCampaignByIdRequest) (*projectpb.GetProjectCampaignByIdResponse, error)
		GetProjectCollaborators(ctx *appcontext.AppContext, req *projectpb.GetProjectCollaboratorsRequest) (*projectpb.GetProjectCollaboratorsResponse, error)
	}
	App interface {
		Hubs
	}

	appHubHandler struct {
		GetProjectByIDHandler
		GetProjectCampaignByIDHandler
		GetProjectCollaboratorsHandler
	}
	Application struct {
		appHubHandler
	}
)

var _ App = (*Application)(nil)

func New(
	projectCampaignHub domain.ProjectCampaignHub,
	service domain.Service,
) *Application {
	return &Application{
		appHubHandler: appHubHandler{
			GetProjectByIDHandler:          NewGetProjectByIDHandler(service),
			GetProjectCampaignByIDHandler:  NewGetProjectCampaignByIDHandler(projectCampaignHub),
			GetProjectCollaboratorsHandler: NewGetProjectCollaboratorsHandler(service),
		},
	}
}

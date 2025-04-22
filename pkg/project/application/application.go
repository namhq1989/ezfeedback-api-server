package application

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Commands interface {
		CreateProject(ctx *appcontext.AppContext, performerID string, req dto.CreateProjectRequest) (*dto.CreateProjectResponse, error)
		UpdateProject(ctx *appcontext.AppContext, performerID, projectID string, req dto.UpdateProjectRequest) (*dto.UpdateProjectResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.CreateProjectHandler
		command.UpdateProjectHandler
	}
	queryHandlers struct {
		query.PingHandler
	}
	Application struct {
		commandHandlers
		queryHandlers
	}
)

var _ Instance = (*Application)(nil)

func New(
	projectRepository domain.ProjectRepository,
	projectSettingRepository domain.ProjectSettingRepository,
	cachingRepository domain.CachingRepository,
	billingHub domain.BillingHub,
) *Application {
	return &Application{
		commandHandlers: commandHandlers{
			CreateProjectHandler: command.NewCreateProjectHandler(projectRepository, projectSettingRepository, billingHub),
			UpdateProjectHandler: command.NewUpdateProjectHandler(projectRepository, projectSettingRepository, cachingRepository),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),
		},
	}
}

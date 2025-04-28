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
		ChangeProjectStatus(ctx *appcontext.AppContext, performerID, projectID string, req dto.ChangeProjectStatusRequest) (*dto.ChangeProjectStatusResponse, error)

		CreateProjectCategory(ctx *appcontext.AppContext, performerID, projectID string, req dto.CreateProjectCategoryRequest) (*dto.CreateProjectCategoryResponse, error)
		UpdateProjectCategory(ctx *appcontext.AppContext, performerID, projectID, categoryID string, req dto.UpdateProjectCategoryRequest) (*dto.UpdateProjectCategoryResponse, error)
		ChangeProjectCategoryStatus(ctx *appcontext.AppContext, performerID, projectID, categoryID string, req dto.ChangeProjectCategoryStatusRequest) (*dto.ChangeProjectCategoryStatusResponse, error)

		CreateProjectCampaign(ctx *appcontext.AppContext, performerID, projectID string, req dto.CreateProjectCampaignRequest) (*dto.CreateProjectCampaignResponse, error)
		UpdateProjectCampaign(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.UpdateProjectCampaignRequest) (*dto.UpdateProjectCampaignResponse, error)
		ChangeProjectCampaignStatus(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.ChangeProjectCampaignStatusRequest) (*dto.ChangeProjectCampaignStatusResponse, error)
		AddProjectCampaignCategory(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.AddProjectCampaignCategoryRequest) (*dto.AddProjectCampaignCategoryResponse, error)
		RemoveProjectCampaignCategory(ctx *appcontext.AppContext, performerID, projectID, campaignID string, req dto.RemoveProjectCampaignCategoryRequest) (*dto.RemoveProjectCampaignCategoryResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)

		GetProjects(ctx *appcontext.AppContext, performerID string, _ dto.GetProjectsRequest) (*dto.GetProjectsResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.CreateProjectHandler
		command.UpdateProjectHandler
		command.ChangeProjectStatusHandler

		command.CreateProjectCategoryHandler
		command.UpdateProjectCategoryHandler
		command.ChangeProjectCategoryStatusHandler

		command.CreateProjectCampaignHandler
		command.UpdateProjectCampaignHandler
		command.ChangeProjectCampaignStatusHandler
		command.AddProjectCampaignCategoryHandler
		command.RemoveProjectCampaignCategoryHandler
	}
	queryHandlers struct {
		query.PingHandler

		query.GetProjectsHandler
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
	projectCategoryRepository domain.ProjectCategoryRepository,
	projectCampaignRepository domain.ProjectCampaignRepository,
	projectCampaignCategoryRepository domain.ProjectCampaignCategoryRepository,
	cachingRepository domain.CachingRepository,
	billingHub domain.BillingHub,
	service domain.Service,
) *Application {
	return &Application{
		commandHandlers: commandHandlers{
			CreateProjectHandler: command.NewCreateProjectHandler(projectRepository, projectSettingRepository, cachingRepository, billingHub),
			UpdateProjectHandler: command.NewUpdateProjectHandler(projectRepository, projectSettingRepository, cachingRepository),
			ChangeProjectStatusHandler: command.NewChangeProjectStatusHandler(
				projectRepository,
				cachingRepository,
			),

			CreateProjectCategoryHandler: command.NewCreateProjectCategoryHandler(
				projectRepository,
				projectCategoryRepository,
				cachingRepository,
			),
			UpdateProjectCategoryHandler: command.NewUpdateProjectCategoryHandler(
				projectCategoryRepository,
				cachingRepository,
				service,
			),
			ChangeProjectCategoryStatusHandler: command.NewChangeProjectCategoryStatusHandler(
				projectCategoryRepository,
				cachingRepository,
				service,
			),

			CreateProjectCampaignHandler: command.NewCreateProjectCampaignHandler(
				projectRepository,
				projectCampaignRepository,
				cachingRepository,
			),
			UpdateProjectCampaignHandler: command.NewUpdateProjectCampaignHandler(
				projectCampaignRepository,
				cachingRepository,
				service,
			),
			ChangeProjectCampaignStatusHandler: command.NewChangeProjectCampaignStatusHandler(
				projectCampaignRepository,
				cachingRepository,
				service,
			),
			AddProjectCampaignCategoryHandler: command.NewAddProjectCampaignCategoryHandler(
				projectCampaignCategoryRepository,
				cachingRepository,
				service,
			),
			RemoveProjectCampaignCategoryHandler: command.NewRemoveProjectCampaignCategoryHandler(
				projectCampaignCategoryRepository,
				cachingRepository,
				service,
			),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),

			GetProjectsHandler: query.NewGetProjectsHandler(projectRepository, cachingRepository, service),
		},
	}
}

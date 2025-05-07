package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateProjectHandler struct {
	projectRepository        domain.ProjectRepository
	projectSettingRepository domain.ProjectSettingRepository
	cachingRepository        domain.CachingRepository
	queueRepository          domain.QueueRepository
	billingHub               domain.BillingHub
}

func NewCreateProjectHandler(
	projectRepository domain.ProjectRepository,
	projectSettingRepository domain.ProjectSettingRepository,
	cachingRepository domain.CachingRepository,
	queueRepository domain.QueueRepository,
	billingHub domain.BillingHub,
) CreateProjectHandler {
	return CreateProjectHandler{
		projectRepository:        projectRepository,
		projectSettingRepository: projectSettingRepository,
		cachingRepository:        cachingRepository,
		queueRepository:          queueRepository,
		billingHub:               billingHub,
	}
}

// CreateProject godoc
// @tags     Project
// @summary  Create project
// @id       project-create
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload body    dto.CreateProjectRequest true "Body"
// @success  200     {object} dto.CreateProjectResponse
// @router   /api/project [post]
func (h CreateProjectHandler) CreateProject(ctx *appcontext.AppContext, performerID string, req dto.CreateProjectRequest) (*dto.CreateProjectResponse, error) {
	ctx.Logger().Info("new create project request", appcontext.Fields{
		"performerID": performerID, "title": req.Title, "description": req.Description,
		"domain": req.Domain, "primaryColor": req.PrimaryColor,
	})

	ctx.Logger().Text("check user usage limit")
	canCreateProject, err := h.billingHub.CanCreateProject(ctx, performerID, 5)
	if err != nil {
		ctx.Logger().Error("failed to check user usage limit", err, appcontext.Fields{})
		return nil, apperrors.Common.BadRequest
	}
	if !canCreateProject {
		ctx.Logger().ErrorText("user usage limit exceeded")
		return nil, apperrors.Billing.ProjectLimitExceeded
	}

	ctx.Logger().Text("create project model")
	project, err := domain.NewProject(performerID, req.Title, req.Description)
	if err != nil {
		ctx.Logger().Error("failed to create project model", err, appcontext.Fields{})
		return nil, err
	}
	ctx.Logger().Text("persist project in db")
	if err = h.projectRepository.Create(ctx, *project); err != nil {
		ctx.Logger().Error("failed to persist project in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("create setting model")
	setting, err := domain.NewProjectSetting(project.ID, req.Domain, req.PrimaryColor)
	if err != nil {
		ctx.Logger().Error("failed to create setting model", err, appcontext.Fields{})
		return nil, err
	}
	ctx.Logger().Text("persist setting in db")
	if err = h.projectSettingRepository.Create(ctx, *setting); err != nil {
		ctx.Logger().Error("failed to persist setting in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteApiGetProjectsByUserID(ctx, performerID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("add task to queue")
	if err = h.queueRepository.OnProjectCreated(ctx, domain.QueueOnProjectCreatedPayload{
		UserID:    performerID,
		ProjectID: project.ID,
	}); err != nil {
		ctx.Logger().Error("failed to add task to queue", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done create project request")
	return &dto.CreateProjectResponse{ID: project.ID}, nil
}

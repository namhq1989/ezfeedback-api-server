package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type UpdateProjectHandler struct {
	projectRepository        domain.ProjectRepository
	projectSettingRepository domain.ProjectSettingRepository
	cachingRepository        domain.CachingRepository
}

func NewUpdateProjectHandler(
	projectRepository domain.ProjectRepository,
	projectSettingRepository domain.ProjectSettingRepository,
	cachingRepository domain.CachingRepository,
) UpdateProjectHandler {
	return UpdateProjectHandler{
		projectRepository:        projectRepository,
		projectSettingRepository: projectSettingRepository,
		cachingRepository:        cachingRepository,
	}
}

// UpdateProject godoc
// @tags     Project
// @summary  Update project
// @id       project-update
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Project id"
// @param    payload body    dto.UpdateProjectRequest true "Body"
// @success  200     {object} dto.UpdateProjectResponse
// @router   /api/project/{id} [put]
func (h UpdateProjectHandler) UpdateProject(ctx *appcontext.AppContext, performerID, projectID string, req dto.UpdateProjectRequest) (*dto.UpdateProjectResponse, error) {
	ctx.Logger().Info("new update project request", appcontext.Fields{
		"performerID": performerID, "title": req.Title, "description": req.Description,
		"isFeedbackPublic": req.IsFeedbackPublic, "allowAnonymousFeedback": req.AllowAnonymousFeedback,
		"enableVoting": req.EnableVoting,
	})

	ctx.Logger().Text("find project in db")
	project, err := h.projectRepository.FindByID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project in db", err, appcontext.Fields{})
		return nil, err
	}
	if project == nil {
		ctx.Logger().ErrorText("project not found")
		return nil, apperrors.Project.ProjectNotFound
	}

	ctx.Logger().Text("update project data")
	if err = project.SetTitle(req.Title); err != nil {
		ctx.Logger().Error("failed to update project title", err, appcontext.Fields{})
		return nil, err
	}
	if err = project.SetDescription(req.Description); err != nil {
		ctx.Logger().Error("failed to update project description", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update project in db")
	if err = h.projectRepository.Update(ctx, *project); err != nil {
		ctx.Logger().Error("failed to update project in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("find setting in db")
	setting, err := h.projectSettingRepository.FindByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find setting in db", err, appcontext.Fields{})
		return nil, err
	}
	if setting == nil {
		ctx.Logger().Text("setting not found, create new")
		setting, _ = domain.NewProjectSetting(project.ID, req.IsFeedbackPublic, req.AllowAnonymousFeedback, req.EnableVoting)
	}

	ctx.Logger().Text("update setting in db")
	if err = h.projectSettingRepository.Update(ctx, *setting); err != nil {
		ctx.Logger().Error("failed to update setting in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete caching data")
	if err = h.cachingRepository.DeleteProjectByID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	if err = h.cachingRepository.DeleteProjectSettingByProjectID(ctx, projectID); err != nil {
		ctx.Logger().Error("failed to delete caching data", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done update project request")
	return &dto.UpdateProjectResponse{}, nil
}
